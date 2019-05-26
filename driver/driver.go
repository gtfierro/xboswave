package driver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"reflect"
	"sync"
	"time"

	//"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/immesys/wavemq/mqpb"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
)

type URI struct {
	Namespace string
	Value     string
}

type Triple struct {
	Subject   URI
	Predicate URI
	Object    URI
}

type Config struct {
	Namespace  string
	EntityFile string
	SiteRouter string
}

type Driver struct {
	ctx          context.Context
	brickContext map[Triple]time.Time
	namespace    []byte
	namespaceStr string
	perspective  *mqpb.Perspective
	client       mqpb.WAVEMQClient

	sync.RWMutex
}

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
	}

	conn, err := grpc.Dial(cfg.SiteRouter, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	// Create the WAVEMQ client
	driver.client = mqpb.NewWAVEMQClient(conn)

	return driver, nil
}

func (driver *Driver) addToContext(instance string, msg *xbospb.XBOSIoTDeviceState) {

	dyn, err := dynamic.AsDynamicMessage(msg)
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	var triples []Triple
	for _, field := range dyn.GetKnownFields() {
		if isNilReflect(dyn.GetField(field)) {
			continue
		}
		asmsg := field.GetMessageType()
		if asmsg == nil {
			continue
		}
		opts := asmsg.GetOptions()
		c, e := proto.GetExtension(opts, xbospb.E_BrickEquipClass)
		var equipURI *xbospb.URI
		if e == nil {
			equipURI = c.(*xbospb.URI)
		} else {
			continue
		}

		triples = append(triples, Triple{
			Subject:   URI{Namespace: driver.namespaceStr, Value: instance},
			Predicate: URI{Namespace: "rdf", Value: "type"},
			Object:    URI{Namespace: equipURI.Namespace, Value: equipURI.Value},
		})

		value := dyn.GetField(field)
		t, e := dynamic.AsDynamicMessage(value.(proto.Message))
		if e != nil {
			fmt.Println("ERROR", err)
			return
		}

		for _, field := range asmsg.GetFields() {
			if isNilReflect(t.GetField(field)) {
				continue
			}
			opts := field.GetOptions()
			f, e := proto.GetExtension(opts, xbospb.E_BrickPointClass)
			if e == nil {
				uri := f.(*xbospb.URI)

				//TODO: add the URI to the context
				triples = append(triples, Triple{
					Subject:   URI{Namespace: driver.namespaceStr, Value: fmt.Sprintf("%s%s", instance, field.GetJSONName())},
					Predicate: URI{Namespace: "rdf", Value: "type"},
					Object:    URI{Namespace: uri.Namespace, Value: uri.Value},
				}, Triple{
					Subject:   URI{Namespace: driver.namespaceStr, Value: instance},
					Predicate: URI{Namespace: "brickframe", Value: "hasPoint"},
					Object:    URI{Namespace: driver.namespaceStr, Value: fmt.Sprintf("%s%s", instance, field.GetJSONName())},
				})
			}
		}
	}

	driver.Lock()
	defer driver.Unlock()
	for _, triple := range triples {
		driver.brickContext[triple] = time.Now()
	}
}

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

func (driver *Driver) Respond(service, instance string, requestid uint64, msg *xbospb.XBOSIoTDeviceState) error {
	return driver.report(service, instance, requestid, msg)
}

func (driver *Driver) Report(service, instance string, msg *xbospb.XBOSIoTDeviceState) error {
	return driver.report(service, instance, 0, msg)
}

// This method is called by device drivers to publish a reading, encapsulated in an XBOSIoTDeviceState message.
// This is not called automatically. The device driver must choose when to call Report(). This is likely
// either on a regular timer or in response to the receipt of an actuation message
// If a time is not provided in msg, Report will add the current timestamp
func (driver *Driver) report(service, instance string, requestid uint64, msg *xbospb.XBOSIoTDeviceState) error {
	// add the timestamp if it doesn't exist
	if msg.Time == 0 {
		msg.Time = uint64(time.Now().UnixNano())
	}

	driver.addToContext(instance, msg)
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
				log.Println(err)
				continue
			}
			if m.Error != nil {
				log.Println(m.Error)
				continue
			}
			for _, po := range m.Message.Tbs.Payload {
				var msg xbospb.XBOS
				err := proto.Unmarshal(po.Content, &msg)
				if err != nil {
					log.Println(err)
					continue
				}
				received := time.Unix(0, m.Message.Timestamps[len(m.Message.Timestamps)-1])
				if msg.XBOSIoTDeviceActuation != nil {
					if err := cb(msg.XBOSIoTDeviceActuation, received); err != nil {
						log.Println(err)
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
