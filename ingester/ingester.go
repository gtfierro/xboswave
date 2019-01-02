package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/immesys/wavemq/mqpb"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	logrus "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"plugin"
	"sync"
	"time"
)

type ExponentialBackoff struct {
	delay int
}

func NewExponentialBackoff() *ExponentialBackoff {
	return &ExponentialBackoff{
		delay: 2,
	}
}

func (exp *ExponentialBackoff) Wait() {
	time.Sleep(time.Duration(exp.delay) * time.Second)
	exp.delay *= 4
	if time.Duration(exp.delay)*time.Second > 10*time.Minute {
		exp.delay = 600
	}
}
func (exp *ExponentialBackoff) Reset() {
	exp.delay = 2
}

//TODO: need to stagger commits when there are really full buffers (such as when ingester restarts)

type Ingester struct {
	plugin_mapping map[string]pluginlist
	client         mqpb.WAVEMQClient
	perspective    *mqpb.Perspective
	//btrdbClient    *btrdbClient
	tsdbClient timeseriesInterface
	cfgmgr     *ConfigManager
	ctx        context.Context

	// pending subscription changes
	pendingSubs chan subscriptionChange

	subs     map[types.SubscriptionURI]*subscription
	subsLock sync.RWMutex
}

func NewIngester(client mqpb.WAVEMQClient, persp *mqpb.Perspective, cfg Config, cfgmgr *ConfigManager, ctx context.Context) *Ingester {
	ingest := &Ingester{
		plugin_mapping: make(map[string]pluginlist),
		perspective:    persp,
		client:         client,
		subs:           make(map[types.SubscriptionURI]*subscription),
		pendingSubs:    make(chan subscriptionChange),
		cfgmgr:         cfgmgr,
		ctx:            ctx,
	}

	// setup prometheus endpoint
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		port := fmt.Sprintf(":%d", cfg.Monitoring.PrometheusPort)
		logrus.Infof("Prometheus endpoint at %s", port)
		if err := http.ListenAndServe(port, nil); err != nil {
			logrus.Fatal(err)
		}
	}()

	// instantiate the timeseries database
	if cfg.Database.BTrDB != nil {
		ingest.tsdbClient = newBTrDBv4(cfg.Database.BTrDB)
	} else if cfg.Database.InfluxDB != nil {
		ingest.tsdbClient = newInfluxDB(cfg.Database.InfluxDB)
	}

	// monitor pendingSubs channel for changes to the subscriptions
	ingest.handleSubscriptionChanges()

	// add existing archive requests that are enabled
	existingRequests, err := ingest.cfgmgr.List(&RequestFilter{Enabled: &_TRUE})
	if err != nil {
		logrus.Fatal(err)
	}
	for _, req := range existingRequests {
		if err := ingest.addArchiveRequest(&req); err != nil {
			logrus.Fatal(err)
		}
	}

	go ingest.shell()

	return ingest
}

func (ingest *Ingester) Finish() error {
	ingest.subsLock.Lock()
	defer ingest.subsLock.Unlock()

	for uri, sub := range ingest.subs {
		logrus.Info("stopping ", uri)
		sub.stop <- struct{}{}
	}
	logrus.Info("Flushing buffered timeseries values")
	return ingest.tsdbClient.Flush()
}

// registers a .so with a schema
func (ingest *Ingester) addPlugin(schema, plugin_filename string) error {
	logrus.Info("Adding plugin=", plugin_filename, " schema=", schema)

	loadedPlugin, err := plugin.Open(plugin_filename)
	if err != nil {
		return errors.Wrapf(err, "Could not open plugin %s", plugin_filename)
	}
	extractedSymbol, err := loadedPlugin.Lookup("Extract")
	if err != nil {
		return errors.Wrapf(err, "Could not lookup Extract symbol in plugin %s", plugin_filename)
	}
	extractFunc, ok := extractedSymbol.(func(uri types.SubscriptionURI, msg xbospb.XBOS, add func(types.ExtractedTimeseries) error) error)
	if !ok {
		return fmt.Errorf("Could not pull Extract symbol from plugin %s (%T) should be (%T)", plugin_filename, extractedSymbol)
	}

	if _, found := ingest.plugin_mapping[schema]; !found {
		ingest.plugin_mapping[schema] = newPluginlist()
	}
	ingest.plugin_mapping[schema].add(plugin_filename, extractFunc)

	return nil
}

func (ingest *Ingester) runPlugins(uri types.SubscriptionURI, msg *mqpb.Message) error {
	for _, po := range msg.Tbs.Payload {
		msgUri := types.SubscriptionURI{
			Namespace: uri.Namespace,
			Resource:  msg.Tbs.Uri,
		}
		var msg xbospb.XBOS
		err := proto.Unmarshal(po.Content, &msg)
		if err != nil {
			logrus.Error(errors.Wrap(err, "Could not unmarshal message into xbospb.XBOS"))
			continue
		}
		//fmt.Printf("    schema=%q type=%+v\n", po.Schema, msg)

		list, found := ingest.plugin_mapping[po.Schema]
		if !found {
			logrus.Errorf("No plugins found for %s", po.Schema)
			continue
		}
		for _, extractFunc := range list.mapping {
			err := extractFunc(msgUri, msg, func(extracted types.ExtractedTimeseries) error {
				if len(extracted.Times) == 0 {
					return nil
				}
				err = ingest.tsdbClient.write(extracted)
				if err != nil {
					logrus.Error(errors.Wrap(err, "Could not write to btrdb buffer"))
				}
				return err
			})
			if err == types.NoMatch {
				continue
			} else if err != nil {
				logrus.Error(errors.Wrap(err, "Could not run extractfunc"))
			}
		}
	}
	return nil
}

func (ingest *Ingester) delArchiveRequest(req *ArchiveRequest) error {
	// remove the request from the config manager
	remainingSubs, err := ingest.cfgmgr.Delete(*req)
	if err != nil {
		return err
	}

	// request an unsubscription, but don't unsub if someone else is using it
	if !remainingSubs {
		ingest.pendingSubs <- subscriptionChange{
			add: false,
			uri: req.URI,
		}
	}

	return nil
}

func (ingest *Ingester) enableArchiveRequest(req *ArchiveRequest) error {
	ingest.pendingSubs <- subscriptionChange{
		add: true,
		uri: req.URI,
	}

	return ingest.cfgmgr.Enable(*req)
}

func (ingest *Ingester) disableArchiveRequest(req *ArchiveRequest) error {
	remainingSubs, err := ingest.cfgmgr.Disable(*req)
	if err != nil {
		return err
	}
	// request an unsubscription, but don't unsub if someone else is using it
	if !remainingSubs {
		ingest.pendingSubs <- subscriptionChange{
			add: false,
			uri: req.URI,
		}
	}
	return nil
}

func (ingest *Ingester) addArchiveRequest(req *ArchiveRequest) error {
	// register plugin
	err := ingest.addPlugin(req.Schema, req.Plugin)
	if err != nil {
		return err
	}

	// request a subscription
	ingest.pendingSubs <- subscriptionChange{
		add: true,
		uri: req.URI,
	}

	if err := ingest.cfgmgr.Add(*req); err != nil {
		return err
	}

	return nil
}

func (ingest *Ingester) newSubscription(uri types.SubscriptionURI) (*subscription, error) {
	logrus.Info("New Subscription ns=", uri.Namespace, "uri=", uri.Resource)
	nsbytes, err := base64.URLEncoding.DecodeString(uri.Namespace)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not decode namespace %s", uri.Namespace)
	}
	var sub = new(subscription)
	sub.S, err = ingest.client.Subscribe(context.Background(), &mqpb.SubscribeParams{
		Perspective: ingest.perspective,
		Namespace:   nsbytes,
		Uri:         uri.Resource,
		Identifier:  IngesterName,
		Expiry:      IngestSubscriptionExpiry,
	})
	sub.timer = NewExponentialBackoff()
	sub.stop = make(chan struct{}, 1)
	sub.uri = uri
	if err != nil {
		return nil, errors.Wrapf(err, "could not subscribe to namespace %s", uri.Namespace)
	}
	// increase # of active subs
	activeSubscriptions.Inc()
	if serr := ingest.cfgmgr.ClearErrorURI(uri); serr != nil {
		logrus.Warning(serr)
	}
	sub.timer.Reset()

	resub := func() error {
		sub.S, err = ingest.client.Subscribe(context.Background(), &mqpb.SubscribeParams{
			Perspective: ingest.perspective,
			Namespace:   nsbytes,
			Uri:         uri.Resource,
			Identifier:  IngesterName,
			Expiry:      IngestSubscriptionExpiry,
		})
		return err
	}

	// this handles the state machine for the subscriptions
	resubtrigger := false
	go func() {
		for {
			select {
			case <-sub.stop:
				logrus.Warning("Stopping subscription to ", sub.uri)
				activeSubscriptions.Dec()
				return
			default:
			}

			// if subscription is currently errored, then make sure we reconnect later
			if sub.S == nil {
				resubtrigger = true
			} else {
				// otherwise, process the message
				m, err := sub.S.Recv()
				if err != nil {
					logrus.Error("subscribe err1:", err)
					activeSubscriptions.Dec()
					sub.timer.Wait()
					if err != io.EOF {
						resubtrigger = true
					}
				} else {
					// check for broker-reported error
					if m.Error != nil {
						err := errors.New(m.Error.Message)
						logrus.Errorf("Error resubscribing to %s from broker: %s. Retrying...", uri, err)
						if serr := ingest.cfgmgr.MarkErrorURI(uri, err.Error()); serr != nil {
							logrus.Warning(serr)
						}
						resubtrigger = true
						continue
					}
					// get uri
					msgsProcessed.Inc()
					err = ingest.runPlugins(uri, m.Message)
					if err != nil {
						logrus.Error("plugins err:", err)
						continue
					}
				}
			}

			// process resub
			if resubtrigger {
				if err := resub(); err != nil {
					logrus.Errorf("Error resubscribing to %s (%s). Retrying...", uri, err)
					if serr := ingest.cfgmgr.MarkErrorURI(uri, err.Error()); serr != nil {
						logrus.Warning(serr)
					}
					sub.timer.Wait()
					continue
				} else {
					logrus.Info("Restablished subscription to", uri.Resource)
					if serr := ingest.cfgmgr.ClearErrorURI(uri); serr != nil {
						logrus.Warning(serr)
					}
					sub.timer.Reset()
					activeSubscriptions.Inc()
					resubtrigger = false
				}
			}
		}
	}()

	return sub, nil
}
