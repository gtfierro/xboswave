package main

import (
	"fmt"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
)

func has_device(msg xbospb.XBOS) bool {
	return msg.XBOSIoTDeviceState.WeatherStation != nil
}

var device_units = map[string]string{
	"time":                 "seconds",
	"icon":                 "",
	"nearestStormDistance": "miles",
	"nearestStormBearing":  "degrees",
	"precipIntensity":      "inches per hour",
	"precipIntensityError": "",
	"precipProbability":    "",
	"precipType":           "",
	"temperature":          "F",
	"apparentTemperature":  "F",
	"dewPoint":             "F",
	"humidity":             "",
	"pressure":             "millibars",
	"windSpeed":            "miles per hour",
	"windGust":             "miles per hour",
	"windBearing":          "degree",
	"cloudCover":           "",
	"uvIndex":              "miles",
	"visibility":           "",
	"ozone":                "Dobson",
}
var device_lookup = map[string]func(msg xbospb.XBOS) (float64, bool){

	"time": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.Time != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.Time.Value), true
		}
		return 0, false
	},
	"nearestStormDistance": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.NearestStormDistance != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.NearestStormDistance.Value), true
		}
		return 0, false
	},
	"nearestStormBearing": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.NearestStormBearing != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.NearestStormBearing.Value), true
		}
		return 0, false
	},
	"precipIntensity": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.PrecipIntensity != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.PrecipIntensity.Value), true
		}
		return 0, false
	},
	"precipIntensityError": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.PrecipIntensityError != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.PrecipIntensityError.Value), true
		}
		return 0, false
	},
	"precipProbability": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.PrecipProbability != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.PrecipProbability.Value), true
		}
		return 0, false
	},
	"temperature": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.Temperature != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.Temperature.Value), true
		}
		return 0, false
	},
	"apparentTemperature": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.ApparentTemperature != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.ApparentTemperature.Value), true
		}
		return 0, false
	},
	"dewPoint": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.DewPoint != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.DewPoint.Value), true
		}
		return 0, false
	},
	"humidity": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.Humidity != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.Humidity.Value), true
		}
		return 0, false
	},
	"pressure": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.Pressure != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.Pressure.Value), true
		}
		return 0, false
	},
	"windSpeed": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.WindSpeed != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.WindSpeed.Value), true
		}
		return 0, false
	},
	"windGust": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.WindGust != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.WindGust.Value), true
		}
		return 0, false
	},
	"windBearing": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.WindBearing != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.WindBearing.Value), true
		}
		return 0, false
	},
	"cloudCover": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.CloudCover != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.CloudCover.Value), true
		}
		return 0, false
	},
	"uvIndex": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.UvIndex != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.UvIndex.Value), true
		}
		return 0, false
	},
	"visibility": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.Visibility != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.Visibility.Value), true
		}
		return 0, false
	},
	"ozone": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherStation.Ozone != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherStation.Ozone.Value), true
		}
		return 0, false
	},
}

func build_device(uri types.SubscriptionURI, name string, msg xbospb.XBOS) types.ExtractedTimeseries {

	if extractfunc, found := device_lookup[name]; found {
		if value, found := extractfunc(msg); found {
			var extracted types.ExtractedTimeseries
			time := int64(msg.XBOSIoTDeviceState.Time)
			extracted.Values = append(extracted.Values, value)
			extracted.Times = append(extracted.Times, time)
			extracted.UUID = types.GenerateUUID(uri, []byte(name))
			extracted.Collection = fmt.Sprintf("xbos/%s", uri.Resource)
			extracted.Tags = map[string]string{
				"unit": device_units[name],
				"name": name,
			}
			return extracted
		}
	}
	return types.ExtractedTimeseries{}
}

func Extract(uri types.SubscriptionURI, msg xbospb.XBOS, add func(types.ExtractedTimeseries) error) error {
	if msg.XBOSIoTDeviceState != nil {
		if has_device(msg) {
			for name := range device_lookup {
				extracted := build_device(uri, name, msg)
				if err := add(extracted); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
