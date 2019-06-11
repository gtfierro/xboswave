package driver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/immesys/wavemq/mqpb"
	logrus "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, ForceColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
}

func ParseResource(namespace, uri string) *xbospb.Resource {
	parts := strings.Split(uri, "/")
	return &xbospb.Resource{
		Transport: xbospb.Transport_WAVEMQ,
		Namespace: namespace,
		Service:   parts[0],
		Instance:  parts[1],
	}
}

func MakeResource(namespace, service, instance string) *xbospb.Resource {
	return &xbospb.Resource{
		Transport: xbospb.Transport_WAVEMQ,
		Namespace: namespace,
		Service:   service,
		Instance:  instance,
	}
}

func ResourceToString(res *xbospb.Resource) string {
	if res.Transport == xbospb.Transport_WAVEMQ {
		return fmt.Sprintf("%s/%s/%s", res.Namespace, res.Service, res.Instance)
	} else {
		return fmt.Sprintf("%s/%s/%s/%s", res.Location, res.Namespace, res.Service, res.Instance)
	}
}

type URI struct {
	Namespace string
	Value     string
}

type Triple struct {
	Subject   URI
	Predicate URI
	Object    URI
}

type XBOSDriver interface {
	Init(Config) error
	Start() error
}

// A Driver is a persistent process that handles the reporting+actuation to and from
// one or more devices
type Driver struct {
	ctx          context.Context
	brickContext map[Triple]time.Time
	namespace    []byte
	namespaceStr string
	report_rate  time.Duration
	perspective  *mqpb.Perspective
	client       mqpb.WAVEMQClient
	errors       chan error

	sync.RWMutex
}

// creates a new driver with the given configuration
func NewDriver(cfg Config) (*Driver, error) {
	namespace, err := base64.URLEncoding.DecodeString(cfg.Namespace)
	if err != nil {
		return nil, err
	}

	entity, err := ioutil.ReadFile(cfg.EntityFile)
	if err != nil {
		return nil, err
	}

	driver := &Driver{
		ctx: context.Background(),
		perspective: &mqpb.Perspective{
			EntitySecret: &mqpb.EntitySecret{
				DER: entity,
			},
		},
		namespace:    namespace,
		namespaceStr: cfg.Namespace,
		brickContext: make(map[Triple]time.Time),
		report_rate:  cfg.ReportRate,
		errors:       make(chan error),
	}

	conn, err := grpc.Dial(cfg.SiteRouter, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	// Create the WAVEMQ client
	driver.client = mqpb.NewWAVEMQClient(conn)

	return driver, nil
}

func (driver *Driver) BlockUntilError() error {
	return <-driver.errors
}

func (driver *Driver) ExitWithError(err error) {
	driver.errors <- err
}

// add triples to the driver's context. These will be deduped automatically so
// adding triples is idempotent. This method does *not* report any triples; use ReportContext()
// to do this.
func (driver *Driver) AddToContext(triples []Triple) {
	driver.Lock()
	defer driver.Unlock()
	for _, triple := range triples {
		driver.brickContext[triple] = time.Now()
	}
}

// publishes the driver's device context
// TODO: where is good URI to publish this?
func (driver *Driver) ReportContext(instance string) error {
	driver.RLock()
	defer driver.RUnlock()

	var triples []*xbospb.Triple

	for triple := range driver.brickContext {
		triples = append(triples, &xbospb.Triple{
			Subject:   &xbospb.URI{Namespace: triple.Subject.Namespace, Value: triple.Subject.Value},
			Predicate: &xbospb.URI{Namespace: triple.Predicate.Namespace, Value: triple.Predicate.Value},
			Object:    &xbospb.URI{Namespace: triple.Object.Namespace, Value: triple.Object.Value},
		})
	}

	msg := &xbospb.XBOS{
		XBOSIoTContext: &xbospb.XBOSIoTContext{
			Time:    uint64(time.Now().UnixNano()),
			Context: triples,
		},
	}
	po, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(driver.ctx, 1*time.Second)
	defer cancel()

	resp, err := driver.client.Publish(ctx, &mqpb.PublishParams{
		Perspective: driver.perspective,
		Namespace:   driver.namespace,
		Uri:         fmt.Sprintf("context/%s/signal/state", instance),
		Content: []*mqpb.PayloadObject{
			{Schema: "xbosproto/XBOS", Content: po},
		},
	})
	if resp.Error != nil {
		return fmt.Errorf("Error publishing: %s", resp.Error.Message)
	}
	return err

	return nil
}

// publishes a device's state in response to a Request
func (driver *Driver) Respond(service, instance string, requestid uint64, msg *xbospb.XBOSIoTDeviceState) error {
	return driver.report(service, instance, requestid, msg)
}

// This method is called by device drivers to publish a reading, encapsulated in an XBOSIoTDeviceState message.
// This is called automatically by AddReport(), which uses a regular timer.
// Respond() is used when publishing in response to the receipt of an actuation message
// If a time is not provided in msg, Report will add the current timestamp.
func (driver *Driver) Report(service, instance string, msg *xbospb.XBOSIoTDeviceState) error {
	return driver.report(service, instance, 0, msg)
}

func (driver *Driver) report(service, instance string, requestid uint64, msg *xbospb.XBOSIoTDeviceState) error {
	// add the timestamp if it doesn't exist
	if msg.Time == 0 {
		msg.Time = uint64(time.Now().UnixNano())
	}

	msg.Requestid = int64(requestid)

	xbosmsg := &xbospb.XBOS{
		XBOSIoTDeviceState: msg,
	}
	po, err := proto.Marshal(xbosmsg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(driver.ctx, 1*time.Second)
	defer cancel()

	resp, err := driver.client.Publish(ctx, &mqpb.PublishParams{
		Perspective: driver.perspective,
		Persist:     true,
		Namespace:   driver.namespace,
		Uri:         fmt.Sprintf("%s/%s/signal/state", service, instance),
		Content: []*mqpb.PayloadObject{
			{Schema: "xbosproto/XBOS", Content: po},
		},
	})
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return fmt.Errorf("Error publishing: %s", resp.Error.Message)
	}
	return nil
}

// Call the given function at config.ReportRate and publish the returned device state
func (driver *Driver) AddReport(service, instance string, cb func() (*xbospb.XBOSIoTDeviceState, error)) error {
	doread := func() {
		if msg, err := cb(); err != nil {
			log.Error("Report err", err)
		} else if err := driver.Report(service, instance, msg); err != nil {
			driver.ExitWithError(err)
		} else {
			log.Info("Reporting", service, instance, msg)
		}
	}
	go func() {
		doread()
		for range time.Tick(driver.report_rate) {
			doread()
		}
	}()
	return nil
}

// Call the given function whenever the driver receives an actuation request for a given device
func (driver *Driver) AddActuationCallback(service, instance string, cb func(msg *xbospb.XBOSIoTDeviceActuation, received time.Time) error) error {
	//TODO: handle mutiple callbacks?
	fmt.Println("Subscribing to", fmt.Sprintf("%s/%s/slot/cmd", service, instance))
	sub, err := driver.client.Subscribe(context.Background(), &mqpb.SubscribeParams{
		Perspective: driver.perspective,
		Namespace:   driver.namespace,
		Uri:         fmt.Sprintf("%s/%s/slot/cmd", service, instance),
		Identifier:  fmt.Sprintf("%s|%s", service, instance),
		Expiry:      60,
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			m, err := sub.Recv()
			if err != nil && err != io.EOF {
				driver.ExitWithError(err)
				continue
			}
			if m.Error != nil {
				driver.ExitWithError(fmt.Errorf("%s", m.Error))
				continue
			}
			for _, po := range m.Message.Tbs.Payload {
				var msg xbospb.XBOS
				err := proto.Unmarshal(po.Content, &msg)
				if err != nil {
					driver.ExitWithError(err)
					continue
				}
				received := time.Unix(0, m.Message.Timestamps[len(m.Message.Timestamps)-1])
				if msg.XBOSIoTDeviceActuation != nil {
					if err := cb(msg.XBOSIoTDeviceActuation, received); err != nil {
						log.Error(err)
						continue
					}
				}
			}
		}
	}()
	return nil
}

func msg2json(msg proto.Message) (map[string]interface{}, error) {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       "   ",
		OrigName:     true,
	}
	s, err := marshaler.MarshalToString(msg)
	if err != nil {
		return nil, err
	}
	var m = make(map[string]interface{})
	err = json.Unmarshal([]byte(s), &m)
	return m, err
}

func isNilReflect(v interface{}) bool {
	if v == nil {
		return true
	}
	value := reflect.ValueOf(v)
	return (value.Kind() == reflect.Ptr && value.IsNil())
}

func Manage(cfg Config, drivers ...XBOSDriver) error {
	for _, driver := range drivers {
		driver := driver
		go func() {
			for {
				if err := driver.Init(cfg); err != nil {
					log.Errorf("Driver %v init failed with config %v and error %v\nRestarting in 30 seconds...", driver, cfg, err)
				} else if err := driver.Start(); err != nil {
					log.Errorf("Driver %v exited with config %v and error %v\nRestarting in 30 seconds...", driver, cfg, err)
				} else {
					log.Warningf("Driver %v with config %v exited without an error\nRestarting in 30 seconds...")
				}
				time.Sleep(30 * time.Second)
			}
		}()
	}
	select {}
}
