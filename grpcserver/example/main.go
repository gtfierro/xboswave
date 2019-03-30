package main

import (
	"context"
	"log"

	//"github.com/golang/protobuf/ptypes"
	"github.com/gtfierro/xboswave/grpcserver"
	xbospb "github.com/gtfierro/xboswave/proto"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type helloserver struct{}

// SayHello implements helloworld.GreeterServer
func (s *helloserver) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	srv, err := grpcserver.NewWaveMQServer(&grpcserver.Config{
		SiteRouter: "localhost:4516",
		EntityFile: "defaultentity.ent",
		Namespace:  "GyBnl_UdduxPIcOwkrnZfqJGQiztUWKyHj9m5zHiFHS1uQ==",
		BaseURI:    "testgrpc",
		ServerName: "helloworld",
	})
	if err != nil {
		log.Fatal(err)
	}

	var backend helloserver

	srv.OnUnary("SayHello", func(call *xbospb.UnaryCall) (*xbospb.UnaryResponse, error) {
		log.Printf("%+v", call)
		log.Println("query id", call.QueryId)
		log.Println("type", call.Payload.GetTypeUrl())
		log.Println("len pay", len(call.Payload.GetValue()))

		var msg pb.HelloRequest

		err := grpcserver.GetPayload(call, &msg)
		if err != nil {
			return nil, err
		}

		reply, err := backend.SayHello(context.Background(), &msg)

		log.Printf("Reply: %+v", reply)

		resp, err := grpcserver.MakeUnaryResponse(call, reply, err)
		return resp, err
	})

	srv.Serve()
	select {}
}
