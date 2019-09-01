package main

import "github.com/gtfierro/xboswave/ingester/types"
import xbospb "github.com/gtfierro/xboswave/proto"
import "fmt"

func has_frame(msg xbospb.XBOS) bool {
	return msg.EnergiseMessage != nil
}

func Extract(uri types.SubscriptionURI, msg xbospb.XBOS, add func(types.ExtractedTimeseries) error) error {
	if !has_frame(msg) {
		return nil
	}

	// archive SPBC phasor targets
	if msg.EnergiseMessage.SPBC != nil {
		spbc := msg.EnergiseMessage.SPBC
		timestamp := spbc.Time

		var extracted types.ExtractedTimeseries
		for _, phasor_target := range spbc.PhasorTargets {

			// archive angle
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = phasor_target.NodeID
			extracted.Tags["channel_name"] = phasor_target.ChannelName
			if phasor_target.Kvbase != nil {
				extracted.Tags["kvbase"] = fmt.Sprintf("%f", phasor_target.Kvbase.Value)
			}
			if phasor_target.KVAbase != nil {
				extracted.Tags["KVAbase"] = fmt.Sprintf("%f", phasor_target.KVAbase.Value)
			}
			extracted.Tags["name"] = "angle"
			extracted.Tags["unit"] = "degrees"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/angle", uri.Resource, phasor_target.ChannelName)
			extracted.Values = []float64{phasor_target.Angle}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"angle"))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// archive magnitude
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = phasor_target.NodeID
			extracted.Tags["channel_name"] = phasor_target.ChannelName
			if phasor_target.Kvbase != nil {
				extracted.Tags["kvbase"] = fmt.Sprintf("%f", phasor_target.Kvbase.Value)
			}
			if phasor_target.KVAbase != nil {
				extracted.Tags["KVAbase"] = fmt.Sprintf("%f", phasor_target.KVAbase.Value)
			}
			extracted.Tags["name"] = "magnitude"
			extracted.Tags["unit"] = "per unit"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/magnitude", uri.Resource, phasor_target.ChannelName)
			extracted.Values = []float64{phasor_target.Magnitude}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"magnitude"))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}
		}
	}

	// archive LPBC Status
	if msg.EnergiseMessage.LPBCStatus != nil {
		lpbc := msg.EnergiseMessage.LPBCStatus
		timestamp := lpbc.Time

		var extracted types.ExtractedTimeseries
		for _, channel_status := range lpbc.Statuses {

			// archive phasor_errors
			// V == magnitude
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = channel_status.NodeID
			extracted.Tags["channel_name"] = channel_status.ChannelName
			extracted.Tags["name"] = "magnitude_error"
			extracted.Tags["unit"] = "per unit"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/magnitude_error", uri.Resource, channel_status.ChannelName)
			extracted.Values = []float64{channel_status.PhasorErrors.Magnitude}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"magnitude_error"))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// delta == Angle
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = channel_status.NodeID
			extracted.Tags["channel_name"] = channel_status.ChannelName
			extracted.Tags["name"] = "angle_error"
			extracted.Tags["unit"] = "degrees"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/angle_error", uri.Resource, channel_status.ChannelName)
			extracted.Values = []float64{channel_status.PhasorErrors.Angle}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"angle_error"))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// archive p_saturated
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = channel_status.NodeID
			extracted.Tags["channel_name"] = channel_status.ChannelName
			extracted.Tags["name"] = "p_saturated"
			extracted.Tags["unit"] = "true/false"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/p_saturated", uri.Resource, channel_status.ChannelName)
			var pSatVal = 0
			if channel_status.PSaturated {
				pSatVal = 1
			}
			extracted.Values = []float64{float64(pSatVal)}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"p_saturated"))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// archive q_saturated
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = channel_status.NodeID
			extracted.Tags["channel_name"] = channel_status.ChannelName
			extracted.Tags["name"] = "q_saturated"
			extracted.Tags["unit"] = "true/false"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/q_saturated", uri.Resource, channel_status.ChannelName)
			var qSatVal = 0
			if channel_status.QSaturated {
				qSatVal = 1
			}
			extracted.Values = []float64{float64(qSatVal)}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"q_saturated"))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// archive p_max
			if channel_status.PMax != nil {
				extracted = types.ExtractedTimeseries{}
				extracted.Tags = make(map[string]string)
				extracted.Tags["node_id"] = channel_status.NodeID
				extracted.Tags["channel_name"] = channel_status.ChannelName
				extracted.Tags["name"] = "p_max"
				extracted.Tags["unit"] = "kW"
				extracted.Collection = fmt.Sprintf("energise/%s/%s/p_max", uri.Resource, channel_status.ChannelName)
				extracted.Values = []float64{channel_status.PMax.Value}
				extracted.Times = []int64{timestamp}
				extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"p_max"))
				if !extracted.Empty() {
					if err := add(extracted); err != nil {
						fmt.Println(err)
						return err
					}
				}
			}

			// archive q_max
			if channel_status.QMax != nil {
				extracted = types.ExtractedTimeseries{}
				extracted.Tags = make(map[string]string)
				extracted.Tags["node_id"] = channel_status.NodeID
				extracted.Tags["channel_name"] = channel_status.ChannelName
				extracted.Tags["name"] = "q_max"
				extracted.Tags["unit"] = "kVAR"
				extracted.Collection = fmt.Sprintf("energise/%s/%s/q_max", uri.Resource, channel_status.ChannelName)
				extracted.Values = []float64{channel_status.QMax.Value}
				extracted.Times = []int64{timestamp}
				extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"q_max"))
				if !extracted.Empty() {
					if err := add(extracted); err != nil {
						fmt.Println(err)
						return err
					}
				}
			}
		}

	}

	// archive LPBC Status
	if msg.EnergiseMessage.ActuatorCommand != nil {
		var extracted types.ExtractedTimeseries
		_msg := msg.EnergiseMessage.ActuatorCommand
		timestamp := _msg.Time
		for idx, phase := range _msg.Phases {

			// pcmd
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = _msg.NodeID
			extracted.Tags["channel_name"] = phase
			extracted.Tags["name"] = "p_cmd"
			extracted.Tags["unit"] = "kVA"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/actuation/%s", uri.Resource, phase, extracted.Tags["name"])
			extracted.Values = []float64{_msg.PCmd[idx]}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+extracted.Tags["name"]))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// qcmd
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = _msg.NodeID
			extracted.Tags["channel_name"] = phase
			extracted.Tags["name"] = "q_cmd"
			extracted.Tags["unit"] = "kVA"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/actuation/%s", uri.Resource, phase, extracted.Tags["name"])
			extracted.Values = []float64{_msg.QCmd[idx]}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+extracted.Tags["name"]))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// pact
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = _msg.NodeID
			extracted.Tags["channel_name"] = phase
			extracted.Tags["name"] = "p_act"
			extracted.Tags["unit"] = "kVA"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/actuation/%s", uri.Resource, phase, extracted.Tags["name"])
			extracted.Values = []float64{_msg.PAct[idx]}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+extracted.Tags["name"]))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// qact
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = _msg.NodeID
			extracted.Tags["channel_name"] = phase
			extracted.Tags["name"] = "q_act"
			extracted.Tags["unit"] = "kVA"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/actuation/%s", uri.Resource, phase, extracted.Tags["name"])
			extracted.Values = []float64{_msg.QAct[idx]}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+extracted.Tags["name"]))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// P_pv
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = _msg.NodeID
			extracted.Tags["channel_name"] = phase
			extracted.Tags["name"] = "p_pv"
			extracted.Tags["unit"] = "kW"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/actuation/%s", uri.Resource, phase, extracted.Tags["name"])
			extracted.Values = []float64{_msg.P_PV[idx]}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+extracted.Tags["name"]))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// BattCmd
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = _msg.NodeID
			extracted.Tags["channel_name"] = phase
			extracted.Tags["name"] = "batt_cmd"
			extracted.Tags["unit"] = "W"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/actuation/%s", uri.Resource, phase, extracted.Tags["name"])
			extracted.Values = []float64{_msg.BattCmd[idx]}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+extracted.Tags["name"]))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// pfctrl
			extracted = types.ExtractedTimeseries{}
			extracted.Tags = make(map[string]string)
			extracted.Tags["node_id"] = _msg.NodeID
			extracted.Tags["channel_name"] = phase
			extracted.Tags["name"] = "pf_ctrl"
			extracted.Tags["unit"] = "pf"
			extracted.Collection = fmt.Sprintf("energise/%s/%s/actuation/%s", uri.Resource, phase, extracted.Tags["name"])
			extracted.Values = []float64{_msg.PfCtrl[idx]}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+extracted.Tags["name"]))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

		}
	}

	return nil
}
