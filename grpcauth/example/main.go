package main

import (
	"context"
	"encoding/pem"
	"io/ioutil"
	"net"
	"time"

	"github.com/cloudflare/cfssl/log"
	"github.com/gtfierro/xboswave/grpcauth"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/immesys/wave/eapi"
	pb "github.com/immesys/wave/eapi/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// This is a sample GRPC server that supports both a unary call and a streaming call
// the proto file for this server is in xboswave/proto/grpcserver.proto
type testserver struct{}

// SayHello implements helloworld.GreeterServer
func (s testserver) TestUnary(ctx context.Context, in *xbospb.TestParams) (*xbospb.TestResponse, error) {
	log.Debugf("Received: %v", in.X)
	return &xbospb.TestResponse{X: in.X}, nil
}

func (s testserver) TestStream(in *xbospb.TestParams, server xbospb.Test_TestStreamServer) error {
	log.Debugf("Received stream: %v", in.X)
	for i := 0; i < 10; i++ {
		if err := server.Send(&xbospb.TestResponse{X: in.X}); err != nil {
			return err
		}
	}
	return nil
}

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

func main() {
	server_perspective := loadPerspective("service2.ent")

	//wavecreds, err := grpcauth.NewWaveCredentials(perspective, "localhost:410", "serviceproof.pem")
	serverwavecreds, err := grpcauth.NewWaveCredentials(server_perspective, "localhost:410", "proof1.pem", "")
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not create wave creds"))
	}

	l, err := net.Listen("tcp", "localhost:7373")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(serverwavecreds))
	xbospb.RegisterTestServer(grpcServer, testserver{})
	go grpcServer.Serve(l)

	client_perspective := loadPerspective("client.ent")
	clientwavecreds, err := grpcauth.NewWaveCredentials(client_perspective, "localhost:410", "", "GyDFHWM5l8TzGz2h7_f_okOKWV67H_wRWHM-heDl8WLjXQ==")
	if err != nil {
		log.Fatal(err)
	}

	//setup client
	clientconn, err := grpc.Dial("localhost:7373", grpc.WithTransportCredentials(clientwavecreds), grpc.FailOnNonTempDialError(true), grpc.WithBlock(), grpc.WithTimeout(30*time.Second))
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
