package main

import (
	"context"
	"log"

	"github.com/gtfierro/xboswave/grpcserver"
	xbospb "github.com/gtfierro/xboswave/proto"
)

// This is a sample GRPC server that supports both a unary call and a streaming call
// the proto file for this server is in xboswave/proto/grpcserver.proto
type testserver struct{}

// SayHello implements helloworld.GreeterServer
func (s *testserver) TestUnary(ctx context.Context, in *xbospb.TestParams) (*xbospb.TestResponse, error) {
	log.Printf("Received: %v", in.X)
	return &xbospb.TestResponse{X: in.X}, nil
}

func (s *testserver) TestStream(ctx context.Context, in *xbospb.TestParams, server xbospb.Test_TestStreamServer) error {
	log.Printf("Received stream: %v", in.X)
	for i := 0; i < 10; i++ {
		if err := server.Send(&xbospb.TestResponse{X: in.X}); err != nil {
			return err
		}
	}
	return nil
}

func main() {

	// create a new WAVEMQ frontend for a GRPC server
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

	// instantiate GRPC backend
	var backend testserver

	// on the "TestUnary" unary call
	srv.OnUnary("TestUnary", func(call *xbospb.UnaryCall) (*xbospb.UnaryResponse, error) {

		// unpack the generic argument inside "call" into the type
		// expected by the GRPC service. GetUnaryPayload does this for us
		var msg xbospb.TestParams
		err := grpcserver.GetUnaryPayload(call, &msg)
		if err != nil {
			return nil, err
		}

		// make the call to the GRPC backend by calling the method
		reply, err := backend.TestUnary(context.Background(), &msg)
		log.Printf("Reply: %+v", reply)

		// wrap the response object and return
		resp, err := grpcserver.MakeUnaryResponse(call, reply, err)
		return resp, err
	})

	srv.OnStream("TestStream", func(call *xbospb.StreamingCall, stream *grpcserver.StreamContext) error {

		// unpack the generic argument inside "call" into the type expected by the
		// GRPC service. GetStreamingpayload does this for us
		var msg xbospb.TestParams
		err := grpcserver.GetStreamingPayload(call, &msg)
		if err != nil {
			return err
		}

		// dispatch the call to the GRPC backend. The "stream" object will handle delivery
		// return the error when finished
		err = backend.TestStream(stream.Context(), &msg, stream)
		stream.Finish(call, err)
		return err
	})

	srv.Serve()
	select {}
}
