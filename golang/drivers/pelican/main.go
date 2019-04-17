package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/john-b-yang/xboswave/golang/drivers/pelican/storage"
	pb "github.com/john-b-yang/xboswave/golang/drivers/protos"
	pb2 "github.com/john-b-yang/xboswave/proto"
	"google.golang.org/grpc"
)

type setpointsMsg struct {
	HeatingSetpoint *float64 `msgpack:"heating_setpoint"`
	CoolingSetpoint *float64 `msgpack:"cooling_setpoint"`
}

type stateMsg struct {
	HeatingSetpoint *float64 `msgpack:"heating_setpoint"`
	CoolingSetpoint *float64 `msgpack:"cooling_setpoint"`
	Override        *bool    `msgpack:"override"`
	Mode            *int     `msgpack:"mode"`
	Fan             *bool    `msgpack:"fan"`
}

type stageMsg struct {
	HeatingStages *int32 `msgpack:"enabled_heat_stages"`
	CoolingStages *int32 `msgpack:"enabled_cool_stages"`
}

type occupancyMsg struct {
	Occupancy bool  `msgpack:"occupancy"`
	Time      int64 `msgpack:"time"`
}

// TODO(john-b-yang): Figure out what other fields are necessary
// Required Fields for Publishing and Subscribing to WAVEMQ site router
const SiteRouter = "127.0.0.1:4516"

func main() {
	var confBytes []byte
	confBytes, confReadErr := ioutil.ReadFile("config.json")
	if confReadErr != nil {
		fmt.Printf("Failed to read config.json file properly: %s\n", confReadErr)
	}
	var confData map[string]interface{}
	if unmarshalErr := json.Unmarshal(confBytes, &confData); unmarshalErr != nil {
		fmt.Printf("Failed to unmarshal config.json file properly: %s\n", unmarshalErr)
	}
	username := confData["username"].(string)
	password := confData["password"].(string)
	sitename := confData["sitename"].(string)
	namespace := confData["namespace"].(string)
	namespaceBytes, namespaceErr := base64.URLEncoding.DecodeString(namespace)
	if namespaceErr != nil {
		fmt.Printf("Failed to convert namespace to bytes: %v", namespaceErr)
		os.Exit(1)
	}

	pelicans, err := storage.ReadPelicans(username, password, sitename)
	if err != nil {
		fmt.Printf("Failed to read thermostat info: %v\n", err)
		os.Exit(1)
	}

	var paramBytes []byte
	paramBytes, paramReadErr := ioutil.ReadFile("params.json")
	if paramReadErr != nil {
		fmt.Printf("Failed to read params.json file properly: %s\n", paramReadErr)
		os.Exit(1)
	}
	var paramData map[string]interface{}
	if unmarshalErr := json.Unmarshal(paramBytes, &paramData); unmarshalErr != nil {
		fmt.Printf("Failed to unmarshal params.json file properly %s\n", unmarshalErr)
	}

	pollInt, pollIntErr := time.ParseDuration(paramData["poll_interval"].(string))
	pollDr, pollDrErr := time.ParseDuration(paramData["poll_interval_dr"].(string))
	pollSched, pollSchedErr := time.ParseDuration(paramData["poll_interval_sched"].(string))
	if pollIntErr != nil {
		fmt.Printf("Failed to parse duration of poll interval properly: %v", pollIntErr)
	}
	if pollDrErr != nil {
		fmt.Printf("Failed to parse duration of demond response (DR) poll interval properly: %v", pollDrErr)
	}
	if pollSchedErr != nil {
		fmt.Printf("Failed to parse duration of schedule poll interval properly: %v", pollDrErr)
	}

	// Establish a GRPC connection to the site router.
	conn, err := grpc.Dial(SiteRouter, grpc.WithInsecure(), grpc.FailOnNonTempDialError(true), grpc.WithBlock())
	if err != nil {
		fmt.Printf("could not connect to the site router: %v\n", err)
		os.Exit(1)
	}

	client := pb.NewWAVEMQClient(conn)

	done := make(chan bool)
	for _, pelican := range pelicans {
		currentPelican := pelican

		go func() {
			for {
				if status, err := currentPelican.GetStatus(); err != nil {
					fmt.Printf("Failed to retrieve Pelican status: %v\n", err)
					done <- true
				} else if status != nil {
					fmt.Printf("%s %+v\n", currentPelican.Name, status)
					statusMsg := &pb2.XBOSIoTDeviceState{
						Time: uint64(status.Time),
						Thermostat: &pb2.Thermostat{
							Temperature:      &pb2.Double{Value: status.Temperature},
							RelativeHumidity: &pb2.Double{Value: status.RelHumidity},
							Override:         &pb2.Bool{Value: status.Override},
							FanState:         &pb2.Bool{Value: status.Fan},
							Mode:             pb2.HVACMode(status.Mode),
							State:            pb2.HVACState(status.State),
						},
					}
					statusBytes, statusErr := proto.Marshal(statusMsg)
					if statusErr != nil {
						fmt.Printf("Failed to serialized Pelican status message: %v", statusErr)
					}
					payload := &pb.PayloadObject{
						Schema:  "PelicanStatus", // TODO(john-b-yang) Check what schema supposed to be
						Content: statusBytes,
					}
					publishParams := &pb.PublishParams{
						Content:   []*pb.PayloadObject{payload},
						Namespace: namespaceBytes,
					}

					client.Publish(context.Background(), publishParams)
				}
				time.Sleep(pollInt)
			}
		}()

		go func() {
			for {
				if drStatus, drErr := currentPelican.TrackDREvent(); drErr != nil {
					fmt.Printf("Failed to retrieve Pelican's DR status: %v\n", drErr)
				} else if drStatus != nil {
					fmt.Printf("%s DR Status: %+v\n", currentPelican.Name, drStatus)
				}
				time.Sleep(pollDr)
			}
		}()

		go func() {
			for {
				if schedStatus, schedErr := currentPelican.GetSchedule(); schedErr != nil {
					fmt.Printf("Failed to retrieve Pelican's Schedule: %v\n", schedErr)
				} else {
					fmt.Printf("%s Schedule: %+v\n", currentPelican.Name, schedStatus)

					// Convert Go Struct to Proto Message
					schedule := &pb2.ThermostatSchedule{}
					for day, blockSchedules := range schedStatus.DaySchedules {
						var blockList []*pb2.ThermostatScheduleBlock
						for _, block := range blockSchedules {
							blockMsg := &pb2.ThermostatScheduleBlock{
								HeatingSetpoint: &pb2.Double{Value: block.HeatSetting},
								CoolingSetpoint: &pb2.Double{Value: block.CoolSetting},
								Time:            block.Time,
							}
							blockList = append(blockList, blockMsg)
						}
						schedule.ScheduleMap[day] = &pb2.ThermostatScheduleDay{
							Blocks: blockList,
						}
					}

					scheduleBytes, scheduleErr := proto.Marshal(schedule)
					if scheduleErr != nil {
						fmt.Printf("Failed to serialized Pelican schedule message: %v", scheduleErr)
					}
					payload := &pb.PayloadObject{
						Schema:  "PelicanScheduleStatus", // TODO(john-b-yang) Check what schema supposed to be
						Content: scheduleBytes,
					}
					publishParams := &pb.PublishParams{
						Content:   []*pb.PayloadObject{payload},
						Namespace: namespaceBytes,
					}

					client.Publish(context.Background(), publishParams)
				}
				time.Sleep(pollSched)
			}
		}()
	}
	<-done

	/*
		for i, pelican := range pelicans {
			pelican := pelican
			name := strings.Replace(pelican.Name, " ", "_", -1)
			name = strings.Replace(name, "&", "_and_", -1)
			name = strings.Replace(name, "'", "", -1)
			fmt.Println("Transforming", pelican.Name, "=>", name)

			// Ensure thermostat is running with correct number of stages
			if err := pelican.ModifyStages(&types.PelicanStageParams{
				HeatingStages: &pelican.HeatingStages,
				CoolingStages: &pelican.CoolingStages,
			}); err != nil {
				fmt.Printf("Failed to configure heating/cooling stages for pelican %s: %s\n",
					pelican.Name, err)
				//os.Exit(1)
			}

			tstatIfaces[i].SubscribeSlot("setpoints", func(msg *bw2.SimpleMessage) {
				po := msg.GetOnePODF(TSTAT_PO_DF)
				if po == nil {
					fmt.Println("Received message on setpoints slot without required PO. Droping.")
					return
				}

				var setpoints setpointsMsg
				if err := po.(bw2.MsgPackPayloadObject).ValueInto(&setpoints); err != nil {
					fmt.Println("Received malformed PO on setpoints slot. Dropping.", err)
					return
				}

				params := types.PelicanSetpointParams{
					HeatingSetpoint: setpoints.HeatingSetpoint,
					CoolingSetpoint: setpoints.CoolingSetpoint,
				}
				if err := pelican.ModifySetpoints(&params); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("Set heating setpoint to %v and cooling setpoint to %v\n",
						setpoints.HeatingSetpoint, setpoints.CoolingSetpoint)
				}
			})

			tstatIfaces[i].SubscribeSlot("state", func(msg *bw2.SimpleMessage) {
				po := msg.GetOnePODF(TSTAT_PO_DF)
				if po == nil {
					fmt.Println("Received message on state slot without required PO. Dropping.")
					return
				}

				var state stateMsg
				if err := po.(bw2.MsgPackPayloadObject).ValueInto(&state); err != nil {
					fmt.Println("Received malformed PO on state slot. Dropping.", err)
					return
				}

				params := types.PelicanStateParams{
					HeatingSetpoint: state.HeatingSetpoint,
					CoolingSetpoint: state.CoolingSetpoint,
				}
				fmt.Printf("%+v", state)
				if state.Mode != nil {
					m := float64(*state.Mode)
					params.Mode = &m
				}

				if state.Override != nil && *state.Override {
					f := float64(1)
					params.Override = &f
				} else {
					f := float64(0)
					params.Override = &f
				}

				if state.Fan != nil && *state.Fan {
					f := float64(1)
					params.Fan = &f
				} else {
					f := float64(0)
					params.Fan = &f
				}

				if err := pelican.ModifyState(&params); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("Set Pelican state to: %+v\n", params)
				}
			})

			tstatIfaces[i].SubscribeSlot("stages", func(msg *bw2.SimpleMessage) {
				po := msg.GetOnePODF(TSTAT_PO_DF)
				if po == nil {
					fmt.Println("Received message on stage slot without required PO. Dropping.")
					return
				}

				var stages stageMsg
				if err := po.(bw2.MsgPackPayloadObject).ValueInto(&stages); err != nil {
					fmt.Println("Received malformed PO on stage slot. Dropping.", err)
					return
				}
				if stages.HeatingStages == nil && stages.CoolingStages == nil {
					fmt.Println("Received message on stage slot with no content. Dropping.")
					return
				}

				params := types.PelicanStageParams{
					HeatingStages: stages.HeatingStages,
					CoolingStages: stages.CoolingStages,
				}
				if err := pelican.ModifyStages(&params); err != nil {
					fmt.Println(err)
				} else {
					if stages.HeatingStages != nil {
						fmt.Printf("Set pelican heating stages to: %d\n", *stages.HeatingStages)
					}
					if stages.CoolingStages != nil {
						fmt.Printf("Set pelican cooling stages to: %d\n", *stages.CoolingStages)
					}
				}
			})

			schedstatIfaces[i].SubscribeSlot("schedule", func(msg *bw2.SimpleMessage) {
				po := msg.GetOnePODF(SCHED_PO_DF)
				if po == nil {
					fmt.Println("Received message on stage slot without required PO. Dropping.")
					return
				}

				var schedule types.ThermostatSchedule
				if err := po.(bw2.MsgPackPayloadObject).ValueInto(&schedule); err != nil {
					fmt.Println("Received malformed PO on stage slot. Dropping.", err)
					return
				}
				if schedule.DaySchedules == nil {
					fmt.Println("Received message on stage slot with no content. Dropping.")
					return
				}

				if err := pelican.SetSchedule(&schedule); err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("Set pelican schedule to: %v", schedule)
				}
			})
		}

		done := make(chan bool)
		for i, pelican := range pelicans {

			occupancy, err := currentPelican.GetOccupancy()
			if err != nil {
				fmt.Printf("Failed to retrieve initial occupancy reading: %s\n", err)
				return
			}
			// Start occupancy tracking loop only if thermostat has the necessary sensor
			if occupancy != types.OCCUPANCY_UNKNOWN {
				go func() {
					for {
						occupancy, err := currentPelican.GetOccupancy()
						if err != nil {
							fmt.Printf("Failed to read thermostat occupancy: %s\n", err)
						} else {
							occupancyStatus := occupancyMsg{
								Occupancy: (occupancy == types.OCCUPANCY_OCCUPIED),
								Time:      time.Now().UnixNano(),
							}
							po, err := bw2.CreateMsgPackPayloadObject(bw2.FromDotForm(OCCUPANCY_PO_DF), occupancyStatus)
							if err != nil {
								fmt.Printf("Failed to create occupancy msgpack PO: %s\n", err)
							} else {
								currentOccupancyIface.PublishSignal("info", po)
							}
						}
						time.Sleep(pollInt)
					}
				}()
			}
		}
		<-done
	*/
}
