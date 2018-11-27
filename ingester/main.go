package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gtfierro/xboswave/ingester/types"
	"github.com/immesys/wavemq/mqpb"
	logrus "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, ForceColors: true})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

// This is an example that shows how to publish and subscribe to a WAVEMQ site router
// Fill these fields in:
const EntityFile = "wavemqingester.ent"
const Namespace = "GyAlyQyfJuai4MCyg6Rx9KkxnZZXWyDaIo0EXGY9-WEq6w=="
const SiteRouter = "127.0.0.1:4516"

var IngesterName = "testingester2"
var IngestSubscriptionExpiry = int64(48 * 60 * 60) // 48 hours
var MaxInMemoryTimeseriesBuffer = 1000             // # of time/reading pairs
var TimeseriesOperationTimeout = 1 * time.Minute

var namespaceBytes []byte

func main() {
	var err error
	ctx := context.Background()

	namespaceBytes, err = base64.URLEncoding.DecodeString(Namespace)
	if err != nil {
		fmt.Printf("failed to decode namespace: %v\n", err)
		os.Exit(1)
	}
	// Load the WAVE3 entity that will be used
	perspective, err := ioutil.ReadFile(EntityFile)
	if err != nil {
		fmt.Printf("could not load entity %q, you might need to create one and grant it permissions\n", EntityFile)
		os.Exit(1)
	}

	// Establish a GRPC connection to the site router.
	conn, err := grpc.DialContext(ctx, SiteRouter, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		fmt.Printf("could not connect to the site router: %v\n", err)
		os.Exit(1)
	}

	// Create the WAVEMQ client
	client := mqpb.NewWAVEMQClient(conn)
	//subscribe(client, perspective)

	// setup kill
	interruptSignal := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		killSignal := <-interruptSignal
		switch killSignal {
		case os.Interrupt, syscall.SIGINT:
			logrus.Warning("Caught SIGINT; closing...")
		case syscall.SIGTERM:
			logrus.Warning("Caught SIGTERM; closing...")
		default:
			logrus.Warning(killSignal)
		}
		conn.Close()
		close(done)
	}()

	persp := &mqpb.Perspective{
		EntitySecret: &mqpb.EntitySecret{
			DER: perspective,
		},
	}
	//btrdbCfg := &btrdbConfig{
	//	addresses: []string{"127.0.0.1:4410"},
	//}
	influxCfg := &influxdbConfig{
		address: "http://127.0.0.1:8086",
	}
	ingest := NewIngester(client, persp, nil, influxCfg, ctx)

	//store := NewArchiveRequestStore(client, persp, extract)
	req := &ArchiveRequest{
		Schema: "xbosproto/XBOS",
		Plugin: "plugins/hamilton1.so",
		URI: types.SubscriptionURI{
			Namespace: "GyAlyQyfJuai4MCyg6Rx9KkxnZZXWyDaIo0EXGY9-WEq6w==", // XBOS
			Resource:  "*",
		},
	}
	if err := ingest.addArchiveRequest(req); err != nil {
		logrus.Fatal(err)
	}

	req2 := &ArchiveRequest{
		Schema: "xbosproto/XBOS",
		Plugin: "plugins/iot_plugin.so",
		URI: types.SubscriptionURI{
			Namespace: "GyAlyQyfJuai4MCyg6Rx9KkxnZZXWyDaIo0EXGY9-WEq6w==", // XBOS
			Resource:  "*",
		},
	}
	if err := ingest.addArchiveRequest(req2); err != nil {
		logrus.Fatal(err)
	}

	<-done
	logrus.Info(ingest.Finish())
}

//func subscribe(client mqpb.WAVEMQClient, perspective []byte) {
//	sub, err := client.Subscribe(context.Background(), &mqpb.SubscribeParams{
//		Perspective: &mqpb.Perspective{
//			EntitySecret: &mqpb.EntitySecret{
//				DER: perspective,
//			},
//		},
//		Namespace: namespaceBytes,
//		Uri:       "*",
//		//If you want a persistent subscription between different runs of this program,
//		//specify this to be something constant (but unique)
//		Identifier: uuid.NewRandom().String(),
//		//This subscription will automatically unsubscribe one minute after this
//		//program ends
//		Expiry: 60,
//	})
//	if err != nil {
//		fmt.Printf("subscribe error: %v\n", err)
//		os.Exit(1)
//	}
//
//	req := &ArchiveRequest{
//		Schema: "xbosproto/XBOS",
//		Plugin: "plugins/hamilton1.so",
//		URI: SubscriptionURI{
//			Namespace: "GyAlyQyfJuai4MCyg6Rx9KkxnZZXWyDaIo0EXGY9-WEq6w==", // XBOS
//			Resource:  "*",
//		},
//	}
//	err = extract.addPlugin(req.Schema, req.Plugin)
//	if err != nil {
//		logrus.Fatal(err)
//	}
//
//	for {
//		m, err := sub.Recv()
//		if err != nil {
//			fmt.Printf("subscribe error: %v\n", err)
//			os.Exit(1)
//		}
//		if m.Error != nil {
//			fmt.Printf("subscribe error: %v\n", m.Error.Message)
//			os.Exit(1)
//		}
//		fmt.Printf("received message on URI: %s\n", m.Message.Tbs.Uri)
//		fmt.Printf("  contents:\n")
//		for _, po := range m.Message.Tbs.Payload {
//			var msg xbospb.XBOS
//			err := proto.Unmarshal(po.Content, &msg)
//			if err != nil {
//				fmt.Println(err)
//			}
//			//fmt.Printf("    schema=%q type=%+v\n", po.Schema, msg)
//			extract.matchSchema(po.Schema, msg)
//
//		}
//	}
//}
