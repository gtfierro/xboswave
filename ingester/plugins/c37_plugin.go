package main

import "github.com/gtfierro/xboswave/ingester/types"
import xbospb "github.com/gtfierro/xboswave/proto"
import "fmt"

func has_frame(msg xbospb.XBOS) bool {
	return msg.C37DataFrame != nil
}

func Extract(uri types.SubscriptionURI, msg xbospb.XBOS, add func(types.ExtractedTimeseries) error) error {
	if !has_frame(msg) {
		return nil
	}

	// go through phasors
	for _, phasor_channel := range msg.C37DataFrame.PhasorChannels {
		var extracted types.ExtractedTimeseries
		extracted.Tags = map[string]string{
			"station_name": msg.C37DataFrame.StationName,
			"id_code":      fmt.Sprintf("%d", msg.C37DataFrame.IdCode),
		}
		// handle angles
		for _, phasor := range phasor_channel.Data {
			extracted.Tags["unit"] = "degrees"
			extracted.Tags["channel_name"] = phasor_channel.ChannelName
			extracted.Collection = fmt.Sprintf("xbos/%s/%s/%s/%s", uri.Resource, extracted.Tags["station_name"], extracted.Tags["id_code"], extracted.Tags["channel_name"])
			extracted.Tags["name"] = "angle"

			time := phasor.Time
			angle := phasor.Angle
			extracted.Values = append(extracted.Values, angle)
			extracted.Times = append(extracted.Times, time)
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"angle"))
		}

		if !extracted.Empty() {
			if err := add(extracted); err != nil {
				fmt.Println(err)
				return err
			}
		}

		extracted = types.ExtractedTimeseries{}
		extracted.Tags = map[string]string{
			"station_name": msg.C37DataFrame.StationName,
			"id_code":      fmt.Sprintf("%d", msg.C37DataFrame.IdCode),
		}

		// handle magnitudes
		for _, phasor := range phasor_channel.Data {
			extracted.Tags["unit"] = phasor_channel.Unit
			extracted.Tags["channel_name"] = phasor_channel.ChannelName
			extracted.Collection = fmt.Sprintf("xbos/%s/%s/%s/%s", uri.Resource, extracted.Tags["station_name"], extracted.Tags["id_code"], extracted.Tags["channel_name"])
			extracted.Tags["name"] = "magnitude"

			time := phasor.Time
			magnitude := phasor.Magnitude
			extracted.Values = append(extracted.Values, magnitude)
			extracted.Times = append(extracted.Times, time)
			extracted.UUID = types.GenerateUUID(uri, []byte(extracted.Collection+"magnitude"))
		}

		if !extracted.Empty() {
			if err := add(extracted); err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}
