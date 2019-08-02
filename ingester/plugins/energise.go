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
		extracted.Tags = make(map[string]string)
		for _, phasor_target := range spbc.PhasorTargets {
			extracted.Tags["nodeID"] = phasor_target.NodeID
			extracted.Tags["channel_name"] = phasor_target.ChannelName
			if phasor_target.Kvbase != nil {
				extracted.Tags["kvbase"] = fmt.Sprintf("%f", phasor_target.Kvbase.Value)
			}
			if phasor_target.KVAbase != nil {
				extracted.Tags["KVAbase"] = fmt.Sprintf("%f", phasor_target.KVAbase.Value)
			}
			extracted.Collection = fmt.Sprintf("energise/%s/%s", uri.Resource, phasor_target.ChannelName)

			// archive angle
			extracted.Tags["name"] = "angle"
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
			extracted.Tags["name"] = "magnitude"
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
		extracted.Tags = make(map[string]string)
		for _, channel_status := range lpbc.Statuses {
			extracted.Tags["nodeID"] = channel_status.NodeID
			extracted.Tags["channel_name"] = channel_status.ChannelName
			extracted.Collection = fmt.Sprintf("energise/%s/%s", uri.Resource, channel_status.ChannelName)

			// archive phasor_errors
			// V == magnitude
			extracted.Tags["name"] = "V"
			extracted.Values = []float64{channel_status.PhasorErrors.Magnitude}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"V"))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// delta == Angle
			extracted.Tags["name"] = "delta"
			extracted.Values = []float64{channel_status.PhasorErrors.Angle}
			extracted.Times = []int64{timestamp}
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"delta"))
			if !extracted.Empty() {
				if err := add(extracted); err != nil {
					fmt.Println(err)
					return err
				}
			}

			// archive p_saturated
			extracted.Tags["name"] = "p_saturated"
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
			extracted.Tags["name"] = "q_saturated"
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
				extracted.Tags["name"] = "p_max"
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
				extracted.Tags["name"] = "q_max"
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

	return nil
}
