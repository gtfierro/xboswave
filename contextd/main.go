package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"git.sr.ht/~gabe/hod/turtle"
	rdf "git.sr.ht/~gabe/hod/turtle/rdfparser"
	"github.com/golang/protobuf/proto"
	"github.com/gtfierro/xboswave/driver"
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

type contextd struct {
	brickContext map[turtle.Triple]time.Time
	sync.RWMutex
}

func newContextd() *contextd {
	return &contextd{
		brickContext: make(map[turtle.Triple]time.Time),
	}
}

func (cd *contextd) readContextMessage(msg *xbospb.XBOS) {
	var triples []turtle.Triple
	if msg.XBOSIoTContext == nil {
		return
	}
	for _, triple := range msg.XBOSIoTContext.Context {
		t := turtle.Triple{
			Subject:   turtle.URI{triple.Subject.Namespace, triple.Subject.Value},
			Predicate: turtle.URI{triple.Predicate.Namespace, triple.Predicate.Value},
			Object:    turtle.URI{triple.Object.Namespace, triple.Object.Value},
		}
		triples = append(triples, t)
	}
	cd.addtoContext(triples...)
}

func (cd *contextd) extractContextFromState(namespace string, resource string, m *xbospb.XBOS) {

	if m.XBOSIoTDeviceState == nil {
		return
	}
	msg := m.XBOSIoTDeviceState

	dyn, err := dynamic.AsDynamicMessage(msg)
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	var triples []turtle.Triple
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
		res := driver.ParseResource(namespace, resource)

		triples = append(triples, turtle.Triple{
			Subject:   turtle.URI{Namespace: namespace, Value: res.Instance},
			Predicate: turtle.URI{Namespace: turtle.RDF_NAMESPACE, Value: "type"},
			Object:    turtle.URI{Namespace: equipURI.Namespace, Value: equipURI.Value},
		}, turtle.Triple{
			Subject:   turtle.URI{Namespace: namespace, Value: res.Instance},
			Predicate: turtle.URI{Namespace: turtle.BRICK_NAMESPACE, Value: "uri"},
			Object:    turtle.URI{Value: fmt.Sprintf("%s/%s", namespace, resource)},
		})

		value := dyn.GetField(field)
		t, e := dynamic.AsDynamicMessage(value.(proto.Message))
		if e != nil {
			fmt.Println("ERROR", err)
			continue
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
				triples = append(triples, turtle.Triple{
					Subject:   turtle.URI{Namespace: namespace, Value: fmt.Sprintf("%s%s", res.Instance, field.GetJSONName())},
					Predicate: turtle.URI{Namespace: turtle.RDF_NAMESPACE, Value: "type"},
					Object:    turtle.URI{Namespace: uri.Namespace, Value: uri.Value},
				}, turtle.Triple{
					Subject:   turtle.URI{Namespace: namespace, Value: res.Instance},
					Predicate: turtle.URI{Namespace: turtle.BRICK_NAMESPACE, Value: "hasPoint"},
					Object:    turtle.URI{Namespace: namespace, Value: fmt.Sprintf("%s%s", res.Instance, field.GetJSONName())},
				})
			}
		}
	}
	cd.addtoContext(triples...)
}

func (cd *contextd) addtoContext(triples ...turtle.Triple) {
	cd.Lock()
	defer cd.Unlock()
	for _, triple := range triples {
		cd.brickContext[triple] = time.Now()
	}
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("contextd <base64 namespace> <entity>")
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

	go func() {
		for range time.Tick(30 * time.Second) {
			f, err := os.Create("dump.ttl")
			if err != nil {
				log.Println(err)
			}
			enc := rdf.NewTripleEncoder(f, rdf.Turtle)
			cd.RLock()
			for triple := range cd.brickContext {
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
			cd.RUnlock()
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
		Uri:         "*",
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
			namespace := base64.URLEncoding.EncodeToString(m.Message.Tbs.Namespace)
			resource := m.Message.Tbs.Uri
			cd.readContextMessage(&msg)
			cd.extractContextFromState(namespace, resource, &msg)
		}
	}
}

func convertTriple(t turtle.Triple) (rdf.Triple, error) {
	var err error
	newtriple := rdf.Triple{}
	newtriple.Subj, err = rdf.NewIRI(t.Subject.String())
	if err != nil {
		return newtriple, err
	}

	newtriple.Pred, err = rdf.NewIRI(t.Predicate.String())
	if err != nil {
		return newtriple, err
	}

	if len(t.Object.Namespace) == 0 {
		newtriple.Obj, err = rdf.NewLiteral(t.Object.Value)
	} else {
		newtriple.Obj, err = rdf.NewIRI(t.Object.String())
	}
	return newtriple, err

}

func isNilReflect(v interface{}) bool {
	if v == nil {
		return true
	}
	value := reflect.ValueOf(v)
	return (value.Kind() == reflect.Ptr && value.IsNil())
}
