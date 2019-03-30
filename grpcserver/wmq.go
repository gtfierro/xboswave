package grpcserver

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/immesys/wavemq/mqpb"
	"github.com/pkg/errors"
	logrus "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, ForceColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
}

// Implements WAVEMQ frontend to a GRPC backend

type Config struct {
	SiteRouter string
	EntityFile string
	Namespace  string
	BaseURI    string
	ServerName string
}

// call URI format
// <base uri> / s.grpcserver / <server name> / i.grpc / slot / <call name>

type UnaryCallback func(*xbospb.UnaryCall) (*xbospb.UnaryResponse, error)
type StreamCallback func(*xbospb.StreamingCall) (chan *xbospb.StreamingResponse, error)

type WaveMQServer struct {
	client         mqpb.WAVEMQClient
	unaryHandlers  map[string]UnaryCallback
	streamHandlers map[string]StreamCallback
	perspective    *mqpb.Perspective
	namespace      []byte
	baseURI        string
	name           string
}

func NewWaveMQServer(cfg *Config) (*WaveMQServer, error) {
	ctx := context.Background()

	//setup namespace
	namespaceBytes, err := base64.URLEncoding.DecodeString(cfg.Namespace)
	if err != nil {
		log.Fatalf("failed to decode namespace: %v", err)
	}

	// load perspective
	perspectivefile, err := ioutil.ReadFile(cfg.EntityFile)
	if err != nil {
		log.Fatalf("could not load entity (%v) you might need to create one and grant it permissions\n", err)
	}
	perspective := &mqpb.Perspective{
		EntitySecret: &mqpb.EntitySecret{
			DER: perspectivefile,
		},
	}

	conn, err := grpc.DialContext(ctx, cfg.SiteRouter, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect to site router %v", err)
	}
	// Create the WAVEMQ client
	server := &WaveMQServer{
		client:         mqpb.NewWAVEMQClient(conn),
		unaryHandlers:  make(map[string]UnaryCallback),
		streamHandlers: make(map[string]StreamCallback),
		baseURI:        cfg.BaseURI,
		name:           cfg.ServerName,
		perspective:    perspective,
		namespace:      namespaceBytes,
	}
	return server, nil
}

func (wmq *WaveMQServer) OnUnary(method string, cb UnaryCallback) {
	wmq.unaryHandlers[method] = cb
}

func (wmq *WaveMQServer) OnStream(method string, cb StreamCallback) {
	wmq.streamHandlers[method] = cb
}

func (wmq *WaveMQServer) Serve() error {
	// listen on URIs for each method
	interfaceURI := fmt.Sprintf("%s/s.grpcserver/%s/i.grpc/slot/call", wmq.baseURI, wmq.name)
	incoming := wmq.subscribe(interfaceURI)

	for methodName, callback := range wmq.unaryHandlers {
		methodName := methodName
		callback := callback
		go func() {
			log.Println("Listening")
			for msg := range incoming {
				msg := msg

				go func() {
					log.Printf("incoming call")
					unarycall := getUnaryCallByMethod(msg, methodName)
					if unarycall == nil {
						return
					}

					resp, err := callback(unarycall)
					if err != nil {
						log.Error(err)
						return
						//TODO: handle this
					}

					// TODO: send response back
					log.Printf("%+v", resp)

					respuri := fmt.Sprintf("%s/s.grpcserver/%s/i.grpc/signal/response", wmq.baseURI, wmq.name)
					b, err := proto.Marshal(resp)
					if err != nil {
						log.Error(err)
						return
					}

					log.Println("Publish on", respuri)
					x, err := wmq.client.Publish(context.Background(), &mqpb.PublishParams{
						Perspective: wmq.perspective,
						Namespace:   wmq.namespace,
						Uri:         respuri,
						Content:     []*mqpb.PayloadObject{{Schema: "xbosproto/GRPCServer", Content: b}},
						Persist:     true,
					})
					if err != nil {
						log.Error(err)
						return
					}
					log.Println("result", x)
				}()
			}
		}()
	}

	//	for methodName, callback := range wmq.streamHandlers {
	//		methodName := methodName
	//		callback := callback
	//		go func() {
	//			for msg := range incoming {
	//				if err := callback(msg); err != nil {
	//					log.Error(err)
	//					//TODO: handle this
	//				}
	//			}
	//		}()
	//	}
	return nil
}

func (wmq *WaveMQServer) subscribe(uri string) chan *mqpb.Message {
	log.Infof("Subscribe to %s", uri)
	var msgs = make(chan *mqpb.Message)
	go func() {
		for {
			subscription, err := wmq.client.Subscribe(context.Background(), &mqpb.SubscribeParams{
				Perspective: wmq.perspective,
				Namespace:   wmq.namespace,
				Uri:         uri,
				Identifier:  "grpc server frontend test",
				Expiry:      int64(10 * 60), // 10 minutes
			})
			if err != nil {
				log.Error(errors.Wrapf(err, "Could not subscribe to %s. Retrying in 30 sec...", uri))
				time.Sleep(30 * time.Second)
				continue
			}
			for {
				m, err := subscription.Recv()
				if err != nil {
					log.Error(errors.Wrapf(err, "Could not subscribe to %s. Retrying in 30 sec...", uri))
					time.Sleep(30 * time.Second)
					break
				}
				if m.Error != nil {
					log.Error(errors.Wrapf(err, "Error in message (%s). Retrying in 30 sec...", m.Error.Message, uri))
					time.Sleep(30 * time.Second)
					break
				}
				msgs <- m.Message
			}

		}
	}()
	return msgs
}

func getUnaryCallByMethod(m *mqpb.Message, method string) *xbospb.UnaryCall {

	pos := m.Tbs.Payload
	if len(pos) == 0 {
		return nil
	}

	if pos[0].Schema != "xbosproto/GRPCServer" {
		return nil
	}

	var msg xbospb.UnaryCall
	err := proto.Unmarshal(pos[0].Content, &msg)
	if err != nil {
		log.Error(err)
		return nil
	}
	return &msg
}

func GetPayload(call *xbospb.UnaryCall, msg proto.Message) error {
	return ptypes.UnmarshalAny(call.Payload, msg)
}

func MakeUnaryResponse(call *xbospb.UnaryCall, msg proto.Message, err error) (*xbospb.UnaryResponse, error) {
	packed, err := ptypes.MarshalAny(msg)
	if err != nil {
		return nil, err
	}
	var errstr string
	if err != nil {
		errstr = err.Error()
	}
	return &xbospb.UnaryResponse{
		QueryId: call.QueryId,
		Error:   errstr,
		Payload: packed,
	}, nil

}
