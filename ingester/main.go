package main

import (
	"context"
	"encoding/base64"
	//"fmt"
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

	cfg, err := ReadConfig("ingester.yml")
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Printf("%+v\n", cfg)

	ctx := context.Background()

	namespaceBytes, err = base64.URLEncoding.DecodeString(Namespace)
	if err != nil {
		logrus.Fatalf("failed to decode namespace: %v", err)
	}

	// Establish a GRPC connection to the site router.
	conn, err := grpc.DialContext(ctx, cfg.WAVEMQ.SiteRouter, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		logrus.Fatalf("Could not connect to site router %v", err)
	}

	// Create the WAVEMQ client
	client := mqpb.NewWAVEMQClient(conn)

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

	// Load the WAVE3 entity that will be used
	perspective, err := ioutil.ReadFile(cfg.WAVEMQ.EntityFile)
	if err != nil {
		logrus.Fatalf("could not load entity (%v) you might need to create one and grant it permissions\n", err)
	}
	persp := &mqpb.Perspective{
		EntitySecret: &mqpb.EntitySecret{
			DER: perspective,
		},
	}
	//btrdbCfg := &btrdbConfig{
	//	addresses: []string{"127.0.0.1:4410"},
	//}
	//influxCfg := &influxdbConfig{
	//	address: "http://127.0.0.1:8086",
	//}
	ingest := NewIngester(client, persp, cfg.Database, ctx)

	//store := NewArchiveRequestStore(client, persp, extract)
	//req := &ArchiveRequest{
	//	Schema: "xbosproto/XBOS",
	//	Plugin: "plugins/hamilton1.so",
	//	URI: types.SubscriptionURI{
	//		Namespace: "GyAlyQyfJuai4MCyg6Rx9KkxnZZXWyDaIo0EXGY9-WEq6w==", // XBOS
	//		Resource:  "*",
	//	},
	//}
	//if err := ingest.addArchiveRequest(req); err != nil {
	//	logrus.Fatal(err)
	//}

	//req2 := &ArchiveRequest{
	//	Schema: "xbosproto/XBOS",
	//	Plugin: "plugins/iot_plugin.so",
	//	URI: types.SubscriptionURI{
	//		Namespace: "GyAlyQyfJuai4MCyg6Rx9KkxnZZXWyDaIo0EXGY9-WEq6w==", // XBOS
	//		Resource:  "*",
	//	},
	//}
	//if err := ingest.addArchiveRequest(req2); err != nil {
	//	logrus.Fatal(err)
	//}

	req3 := &ArchiveRequest{
		Schema: "xbosproto/XBOS",
		Plugin: "plugins/dent.so",
		URI: types.SubscriptionURI{
			Namespace: "GyAlyQyfJuai4MCyg6Rx9KkxnZZXWyDaIo0EXGY9-WEq6w==", // cory-hall
			Resource:  "dentmeter/*",
		},
	}
	if err := ingest.addArchiveRequest(req3); err != nil {
		logrus.Fatal(err)
	}

	<-done
	logrus.Info(ingest.Finish())
}
