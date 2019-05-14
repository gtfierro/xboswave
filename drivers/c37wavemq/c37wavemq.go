package main

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	xbospb "github.com/gtfierro/xboswave/proto"
	"github.com/immesys/wavemq/mqpb"
	logging "github.com/op/go-logging"
	"google.golang.org/grpc"
)

var lg *logging.Logger

func init() {
	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
	logging.SetFormatter(logging.MustStringFormatter("[%{level:-8s}]%{time:2006-01-02T15:04:05.000000} %{shortfile:18s} > %{message}"))
	lg = logging.MustGetLogger("log")
}

type ProtocolAdapterConfig struct {
	C37TargetAddress string `yaml:"c37TargetAddress"`
	C37ID            uint16 `yaml:"c37id"`

	SiteRouter string                  `yaml:"siteRouter"`
	Namespace  string                  `yaml:"namespace"`
	EntityFile string                  `yaml:"entityFile"`
	Outputs    []ProtocolAdapterOutput `yaml:"outputs"`
}

type ProtocolAdapterOutput struct {
	URI      string   `yaml:"uri"`
	Channels []string `yaml:"channels"`
}

func StartProtocolAdapter(cfg *ProtocolAdapterConfig) {
	//Initiate waveMQ connection
	entityFile, err := ioutil.ReadFile(cfg.EntityFile)
	if err != nil {
		lg.Fatalf("could not read entity file: %v", err)
	}
	namespaceBytes, err := base64.URLEncoding.DecodeString(cfg.Namespace)
	if err != nil {
		lg.Fatalf("failed to decode namespace: %v", err)
	}

	// Establish a GRPC connection to the site router.
	conn, err := grpc.Dial(cfg.SiteRouter, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		lg.Fatalf("could not connect to the site router: %v", err)
	}

	// Create the WAVEMQ client
	client := mqpb.NewWAVEMQClient(conn)

	// Create downstream channels
	downstream := makeDownstreams(entityFile, namespaceBytes, client, cfg.Outputs)

	// Start connection to upstream device
	output := HandleDevice(context.Background(), "upstream", cfg.C37TargetAddress, cfg.C37ID)
	for {
		o := <-output
		for _, d := range downstream {
			d <- o
		}
	}
}

const BatchSize = 120
const Proportion = .1 // 10% of readings
var _batchstep = int(Proportion * BatchSize)

func makeDownstreams(entity []byte, namespace []byte, client mqpb.WAVEMQClient, outputs []ProtocolAdapterOutput) []chan *DataFrame {
	perspective := &mqpb.Perspective{
		EntitySecret: &mqpb.EntitySecret{
			DER: entity,
		},
	}
	rv := make([]chan *DataFrame, len(outputs))
	for idx, o := range outputs {
		rv[idx] = make(chan *DataFrame, 100)
		go func(ch chan *DataFrame, cfg ProtocolAdapterOutput) {
			include := make(map[string]bool)
			for _, c := range cfg.Channels {
				include[c] = true
			}
			published := 0
			go func() {
				for {
					time.Sleep(10 * time.Second)
					lg.Infof("published %d wavemq frames", published)
				}
			}()
			batch := make([]*DataFrame, 0, BatchSize)
			for df := range ch {
				batch = append(batch, df)
				if len(batch) == BatchSize {
					out := &xbospb.C37DataFrame{
						StationName: batch[0].Data[0].STN,
						IdCode:      uint32(batch[0].Data[0].IDCODE),
					}
					for pi, pn := range batch[0].Data[0].PHASOR_NAMES {
						if include[pn] {
							ch := &xbospb.PhasorChannel{
								ChannelName: pn,
							}
							if batch[0].Data[0].PHASOR_ISVOLT[pi] {
								ch.Unit = "Volt"
							} else {
								ch.Unit = "Amp"
							}
							for i := 0; i < len(batch); i += _batchstep {
								ch.Data = append(ch.Data, &xbospb.Phasor{
									Time:      batch[i].UTCUnixNanos,
									Angle:     batch[i].Data[0].PHASOR_ANG[pi],
									Magnitude: batch[i].Data[0].PHASOR_MAG[pi],
								})
							}
							out.PhasorChannels = append(out.PhasorChannels, ch)
						}
					} //loop over phasor names
					for ai, an := range batch[0].Data[0].ANALOG_NAMES {
						if include[an] {
							ch := &xbospb.ScalarChannel{
								ChannelName: an,
								Unit:        "Analog",
							}
							for i := 0; i < len(batch); i += _batchstep {
								ch.Data = append(ch.Data, &xbospb.Scalar{
									Time:  batch[i].UTCUnixNanos,
									Value: batch[i].Data[0].ANALOG[ai],
								})
							}
							out.ScalarChannels = append(out.ScalarChannels, ch)
						}
					}
					for di, dn := range batch[0].Data[0].DIGITAL_NAMES {
						if include[dn] {
							ch := &xbospb.ScalarChannel{
								ChannelName: dn,
								Unit:        "Digital",
							}
							for i := 0; i < len(batch); i += _batchstep {
								ch.Data = append(ch.Data, &xbospb.Scalar{
									Time:  batch[i].UTCUnixNanos,
									Value: float64(batch[i].Data[0].DIGITAL[di]),
								})
							}
							out.ScalarChannels = append(out.ScalarChannels, ch)
						}
					}
					if include["FREQ"] {
						ch := &xbospb.ScalarChannel{
							ChannelName: "FREQ",
							Unit:        "Hz",
						}
						for i := 0; i < len(batch); i += _batchstep {
							ch.Data = append(ch.Data, &xbospb.Scalar{
								Time:  batch[i].UTCUnixNanos,
								Value: batch[i].Data[0].FREQ,
							})
						}
						out.ScalarChannels = append(out.ScalarChannels, ch)
					}
					if include["DFREQ"] {
						ch := &xbospb.ScalarChannel{
							ChannelName: "DFREQ",
							Unit:        "Hz/s",
						}
						for i := 0; i < len(batch); i += _batchstep {
							ch.Data = append(ch.Data, &xbospb.Scalar{
								Time:  batch[i].UTCUnixNanos,
								Value: batch[i].Data[0].DFREQ,
							})
						}
						out.ScalarChannels = append(out.ScalarChannels, ch)
					}
					msg := &xbospb.XBOS{
						C37DataFrame: out,
					}

					po, err := proto.Marshal(msg)
					if err != nil {
						lg.Fatalf("could not marshal output proto: %v", err)
					}
					batch = batch[:0]
					stat, err := client.Publish(context.Background(), &mqpb.PublishParams{
						Perspective: perspective,
						Namespace:   namespace,
						Uri:         cfg.URI,
						Content: []*mqpb.PayloadObject{
							{Schema: "xbosproto/XBOS", Content: po},
						},
					})
					if err != nil {
						lg.Warningf("could not publish, grpc error: %v", err)
					} else if stat.Error != nil {
						lg.Warningf("could not publish: wavemq error: %v", stat.Error.Message)
					} else {
						published += 1
					}
					// spew.Dump(out)
					// _ = perspective
					// _ = po
				} //if batch size
			} //loop over data frames
		}(rv[idx], o)
	}
	return rv
}
