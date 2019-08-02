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
			extracted.Collection = fmt.Sprintf("xbos/%s/%s", uri.Resource, phasor_target.ChannelName)

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
			fmt.Printf("%+v\n", extracted)

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
			fmt.Printf("%+v\n", extracted)
		}

	}

	return nil
}
