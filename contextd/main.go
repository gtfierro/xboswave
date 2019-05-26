package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"git.sr.ht/~gabe/hod/turtle"
	rdf "git.sr.ht/~gabe/hod/turtle/rdfparser"
	"github.com/golang/protobuf/proto"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/immesys/wavemq/mqpb"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/pborman/uuid"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sync"
	"time"
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

type contextd struct {
	brickContext map[Triple]time.Time
	sync.RWMutex
}

func newContextd() *contextd {
	return &contextd{
		brickContext: make(map[Triple]time.Time),
	}
}

func (cd *contextd) addtoContext(namespace string, resource string, msg *xbospb.XBOSIoTDeviceState) {
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

		fmt.Println(resource)
		//TODO: extract instance from URI
		instance := "TODO"
		triples = append(triples, Triple{
			Subject:   URI{Namespace: namespace, Value: instance},
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
					Subject:   URI{Namespace: namespace, Value: fmt.Sprintf("%s%s", instance, field.GetJSONName())},
					Predicate: URI{Namespace: "rdf", Value: "type"},
					Object:    URI{Namespace: uri.Namespace, Value: uri.Value},
				}, Triple{
					Subject:   URI{Namespace: namespace, Value: instance},
					Predicate: URI{Namespace: "brickframe", Value: "hasPoint"},
					Object:    URI{Namespace: namespace, Value: fmt.Sprintf("%s%s", instance, field.GetJSONName())},
				})
			}
		}
	}

	cd.Lock()
	defer cd.Unlock()
	for _, triple := range triples {
		cd.brickContext[triple] = time.Now()
	}
}

func main() {

	if len(os.Args) < 4 {
		fmt.Println("contextd <base64 namespace> <entity> <resource>")
		os.Exit(1)
	}

	namespace, err := base64.URLEncoding.DecodeString(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	entity, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.Dial("localhost:4516", grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}

	cd := newContextd()

	// Create the WAVEMQ client
	client := mqpb.NewWAVEMQClient(conn)
	perspective := &mqpb.Perspective{
		EntitySecret: &mqpb.EntitySecret{
			DER: entity,
		},
	}

	var l sync.Mutex
	triples := make(map[turtle.Triple]time.Time)

	go func() {
		for range time.Tick(30 * time.Second) {
			f, err := os.Create("dump.ttl")
			if err != nil {
				log.Println(err)
			}
			enc := rdf.NewTripleEncoder(f, rdf.Turtle)
			l.Lock()
			for triple := range triples {
				newt, err := convertTriple(triple)
				if err != nil {
					log.Println(err)
					break
				}
				if err := enc.Encode(newt); err != nil {
					log.Println(err)
					break
				}
			}
			l.Unlock()
			if err := enc.Close(); err != nil {
				log.Println(err)
			}
			if err := f.Close(); err != nil {
				log.Println(err)
			}
			log.Println("Serialized to dump.ttl")
		}
	}()

	sub, err := client.Subscribe(context.Background(), &mqpb.SubscribeParams{
		Perspective: perspective,
		Namespace:   namespace,
		Uri:         os.Args[3],
		Identifier:  uuid.NewRandom().String(),
		Expiry:      10,
	})
	if err != nil {
		log.Fatal(err)
	}
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
			if msg.XBOSIoTContext != nil {
				log.Println(base64.URLEncoding.EncodeToString(m.Message.Tbs.Namespace), m.Message.Tbs.Uri)
				l.Lock()
				for _, triple := range msg.XBOSIoTContext.Context {
					t := turtle.Triple{
						Subject:   turtle.URI{triple.Subject.Namespace, triple.Subject.Value},
						Predicate: turtle.URI{triple.Predicate.Namespace, triple.Predicate.Value},
						Object:    turtle.URI{triple.Object.Namespace, triple.Object.Value},
					}
					triples[t] = time.Now()
				}
				l.Unlock()
			} else if msg.XBOSIoTDeviceState != nil {
				cd.addtoContext(base64.URLEncoding.EncodeToString(m.Message.Tbs.Namespace), m.Message.Tbs.Uri, msg.XBOSIoTDeviceState)
			}
		}
	}
}

func convertTriple(t turtle.Triple) (rdf.Triple, error) {
	new := rdf.Triple{}
	subj, err := rdf.NewIRI(t.Subject.String())
	if err != nil {
		return new, err
	}
	pred, err := rdf.NewIRI(t.Predicate.String())
	if err != nil {
		return new, err
	}
	obj, err := rdf.NewIRI(t.Object.String())
	if err != nil {
		return new, err
	}
	return rdf.Triple{
		Subj: subj,
		Pred: pred,
		Obj:  obj,
	}, nil
}

func isNilReflect(v interface{}) bool {
	if v == nil {
		return true
	}
	value := reflect.ValueOf(v)
	return (value.Kind() == reflect.Ptr && value.IsNil())
}
