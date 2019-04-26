package main

import (
	"context"
	"encoding/pem"
	"io/ioutil"
	"time"

	"github.com/cloudflare/cfssl/log"
	"github.com/gtfierro/xboswave/grpcauth"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/immesys/wave/eapi"
	pb "github.com/immesys/wave/eapi/pb"
	"google.golang.org/grpc"
)

func loadPerspective(filename string) *pb.Perspective {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("could not read file %q: %v\n", filename, err)
	}
	block, _ := pem.Decode(contents)
	if block == nil {
		log.Fatalf("file %q is not a PEM file\n", filename)
	}
	if block.Type != eapi.PEM_ENTITY_SECRET {
		log.Fatalf("PEM is not an entity secret\n")
	}

	return &pb.Perspective{
		EntitySecret: &pb.EntitySecret{
			DER: block.Bytes,
		},
	}
}

func init() {
	log.Level = log.LevelDebug
}

func main() {
	client_perspective := loadPerspective("client.ent")
	clientcred, err := grpcauth.NewClientCredentials(client_perspective, "localhost:410", "GyBHxjkpzmGxXk9qgJW6AJHCXleNifvhgusCs0v1MLFWJg==", "xbospb/Test/*")
	if err != nil {
		log.Fatal(err)
	}
	clientcred.AddGRPCProofFile("clientproof.pem")

	//setup client
	clientconn, err := grpc.Dial("localhost:7373", grpc.WithTransportCredentials(clientcred), grpc.FailOnNonTempDialError(true), grpc.WithBlock(), grpc.WithTimeout(30*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	testclient := xbospb.NewTestClient(clientconn)
	resp, err := testclient.TestUnary(context.Background(), &xbospb.TestParams{
		X: "hello1",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("%+v\n", resp)
}
