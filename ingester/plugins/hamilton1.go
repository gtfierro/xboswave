package main

import (
	"fmt"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
	"strings"
)

var lookup = map[string]func(msg xbospb.XBOS) float64{
	"uptime": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.Uptime)
	},
	"acc_x": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.AccX)
	},
	"acc_y": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.AccY)
	},
	"acc_z": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.AccZ)
	},
	"mag_x": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.MagX)
	},
	"mag_y": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.MagY)
	},
	"mag_z": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.MagZ)
	},
	"air_temp": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.AirTemp)
	},
	"air_hum": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.AirHum)
	},
	"air_rh": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.AirRh)
	},
	"light_lux": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.LightLux)
	},
	"buttons": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.Buttons)
	},
	"occupancy": func(msg xbospb.XBOS) float64 {
		return float64(msg.HamiltonData.H3C.Occupancy)
	},
}

func build(uri types.SubscriptionURI, name string, msg xbospb.XBOS) types.ExtractedTimeseries {

	value := lookup[name](msg)

	var extracted types.ExtractedTimeseries
	time := int64(msg.HamiltonData.Time)
	extracted.Values = append(extracted.Values, value)
	extracted.Times = append(extracted.Times, time)
	extracted.UUID = types.GenerateUUID(uri, []byte(name))
	parts := strings.Split(uri.Resource, "/")
	extracted.Collection = fmt.Sprintf("hamilton/%s/%s", name, parts[2]) //uri.Resource
	extracted.Tags = map[string]string{
		"unit": "seconds",
		"name": name,
	}
	return extracted
}

func Extract(uri types.SubscriptionURI, msg xbospb.XBOS, add func(types.ExtractedTimeseries) error) error {
	if msg.HamiltonData != nil {
		for name := range lookup {
			extracted := build(uri, name, msg)
			if err := add(extracted); err != nil {
				return err
			}
		}
	}
	return nil
}
