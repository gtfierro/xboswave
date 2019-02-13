package main

import (
	"fmt"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
)

func has_weather_station(msg xbospb.XBOS) bool {
	return msg.XBOSIoTDeviceState.WeatherStation != nil
}

var weather_units = map[string]string{
	"temperature":            "celsius",
	"precip_intensity":       "unknown",
	"nearest_storm_distance": "km",
	"nearest_storm_bearing":  "degrees",
	"humidity":               "unknown",
}

var weather_lookup = map[string]func(msg xbospb.XBOS) (float64, bool){
	// XBOSIoTDeviceState.WeatherStation
	"temperature": func(msg xbospb.XBOS) (float64, bool) {
		if has_weather_station(msg) && msg.XBOSIoTDeviceState.WeatherStation.Temperature != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.Temperature.Value), true
		}
		return 0, false
	},
	"humidity": func(msg xbospb.XBOS) (float64, bool) {
		if has_weather_station(msg) && msg.XBOSIoTDeviceState.WeatherStation.Humidity != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.Humidity.Value), true
		}
		return 0, false
	},
	"nearest_storm_distance": func(msg xbospb.XBOS) (float64, bool) {
		if has_weather_station(msg) && msg.XBOSIoTDeviceState.WeatherStation.NearestStormDistance != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.NearestStormDistance.Value), true
		}
		return 0, false
	},
	"nearest_storm_bearing": func(msg xbospb.XBOS) (float64, bool) {
		if has_weather_station(msg) && msg.XBOSIoTDeviceState.WeatherStation.NearestStormBearing != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.NearestStormBearing.Value), true
		}
		return 0, false
	},
	"precip_intensity": func(msg xbospb.XBOS) (float64, bool) {
		if has_weather_station(msg) && msg.XBOSIoTDeviceState.WeatherStation.PrecipIntensity != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.PrecipIntensity.Value), true
		}
		return 0, false
	},
}

func build_weather(uri types.SubscriptionURI, name string, msg xbospb.XBOS) types.ExtractedTimeseries {
	if extractfunc, found := weather_lookup[name]; found {
		if value, found := extractfunc(msg); found {
			var extracted types.ExtractedTimeseries
			time := int64(msg.XBOSIoTDeviceState.Time)
			extracted.Values = append(extracted.Values, value)
			extracted.Times = append(extracted.Times, time)
			extracted.UUID = types.GenerateUUID(uri, []byte(name))
			extracted.Collection = fmt.Sprintf("xbos/%s", uri.Resource)
			extracted.Tags = map[string]string{
				"unit": weather_units[name],
				"name": name,
			}
			return extracted
		}
	}

	return types.ExtractedTimeseries{}
}

func Extract(uri types.SubscriptionURI, msg xbospb.XBOS, add func(types.ExtractedTimeseries) error) error {
	if msg.XBOSIoTDeviceState != nil {
		if has_weather_station(msg) {
			for name := range weather_lookup {
				extracted := build_weather(uri, name, msg)
				if err := add(extracted); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
