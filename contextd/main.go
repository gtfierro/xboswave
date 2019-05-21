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
	"github.com/pborman/uuid"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

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
