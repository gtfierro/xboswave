package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/immesys/wavemq/mqpb"
	"github.com/pborman/uuid"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("subscriber <base64 namespace> <entity> <resource>")
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

	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  false,
		EmitDefaults: false,
		Indent:       "   ",
		OrigName:     true,
	}

	// Create the WAVEMQ client
	client := mqpb.NewWAVEMQClient(conn)
	perspective := &mqpb.Perspective{
		EntitySecret: &mqpb.EntitySecret{
			DER: entity,
		},
	}

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
			s, err := marshaler.MarshalToString(&msg)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(base64.URLEncoding.EncodeToString(m.Message.Tbs.Namespace), m.Message.Tbs.Uri)
			log.Println(s)
		}
	}
}
