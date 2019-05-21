package driver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (driver *Driver) addToContext(msg *xbospb.XBOSIoTDeviceState) {

	dyn, err := dynamic.AsDynamicMessage(msg)
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	var triples []Triple
	for _, field := range dyn.GetKnownFields() {
		if isNilReflect(dyn.GetField(field)) {
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

				//TODO: fix this so that it actually creates instances
				triples = append(triples, Triple{
					Subject:   URI{Namespace: equipURI.Namespace, Value: equipURI.Value},
					Predicate: URI{Namespace: "brickframe", Value: "hasPoint"},
					Object:    URI{Namespace: uri.Namespace, Value: uri.Value},
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

func (driver *Driver) ReportContext() error {
	driver.RLock()
	defer driver.RUnlock()

	var triples []*xbospb.Triple

	for triple, announced := range driver.brickContext {
		fmt.Printf("Triple %v announced %s\n", triple, announced)
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
		Uri:         "context",
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

// This method is called by device drivers to publish a reading, encapsulated in an XBOSIoTDeviceState message.
// This is not called automatically. The device driver must choose when to call Report(). This is likely
// either on a regular timer or in response to the receipt of an actuation message
// If a time is not provided in msg, Report will add the current timestamp
func (driver *Driver) Report(resource string, msg *xbospb.XBOSIoTDeviceState) error {
	// add the timestamp if it doesn't exist
	if msg.Time == 0 {
		msg.Time = uint64(time.Now().UnixNano())
	}

	driver.addToContext(msg)

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
		Namespace:   driver.namespace,
		Uri:         resource,
		Content: []*mqpb.PayloadObject{
			{Schema: "xbosproto/XBOS", Content: po},
		},
	})
	if resp.Error != nil {
		return fmt.Errorf("Error publishing: %s", resp.Error.Message)
	}
	return err
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
