package main

import (
	"context"
	"log"
	"time"

	//"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/proto"
	"github.com/gtfierro/xboswave/grpcserver"
	xbospb "github.com/gtfierro/xboswave/proto"
	grpc "google.golang.org/grpc"
)

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

type stream struct {
	c   chan proto.Message
	ctx context.Context
	grpc.ServerStream
}

func newStream(duration time.Duration) *stream {
	ctx, _ := context.WithTimeout(context.Background(), duration)

	return &stream{
		c:   make(chan proto.Message),
		ctx: ctx,
	}
}

func (s *stream) Send(msg *xbospb.TestResponse) error {
	s.c <- msg
	return nil
}
func (s *stream) Context() context.Context {
	return s.ctx
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

	var backend testserver

	srv.OnUnary("TestUnary", func(call *xbospb.UnaryCall) (*xbospb.UnaryResponse, error) {
		//log.Printf("%+v", call)
		//log.Println("query id", call.QueryId)
		//log.Println("type", call.Payload.GetTypeUrl())
		//log.Println("len pay", len(call.Payload.GetValue()))

		var msg xbospb.TestParams
		err := grpcserver.GetUnaryPayload(call, &msg)
		if err != nil {
			return nil, err
		}

		reply, err := backend.TestUnary(context.Background(), &msg)

		log.Printf("Reply: %+v", reply)

		resp, err := grpcserver.MakeUnaryResponse(call, reply, err)
		return resp, err
	})

	srv.OnStream("TestStream", func(call *xbospb.StreamingCall) (chan *xbospb.StreamingResponse, error) {
		var msg xbospb.TestParams

		err := grpcserver.GetStreamingPayload(call, &msg)
		if err != nil {
			return nil, err
		}

		s := newStream(30 * time.Second)

		finished := make(chan bool)
		c := make(chan *xbospb.StreamingResponse)

		ctx, cancel := context.WithTimeout(s.Context(), 3*time.Second)
		go func() {
			err := backend.TestStream(ctx, &msg, s)
			if err != nil {
				log.Println("Error?", err)
				//TODO: transmit error
			}
			cancel()
			finished <- true
		}()

		go func() {
		replyloop:
			for {
				select {
				case reply := <-s.c:
					resp, err := grpcserver.MakeStreamingResponse(call, reply, err)
					if err != nil {
						log.Println("error make stream", err)
					}
					if resp == nil {
						log.Println("FINISH")
						continue
					}
					c <- resp
				case <-finished:
					log.Println("DONE2")
					break replyloop
				case <-s.Context().Done():
					log.Println("DONE")
					break replyloop
				}
			}
			log.Println("closing up")
			close(c)
			close(s.c)
		}()
		return c, nil
	})

	srv.Serve()
	select {}
}
