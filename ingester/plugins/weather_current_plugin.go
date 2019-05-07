package main
import (
	"fmt"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
)
func has_device(msg xbospb.XBOS) bool {
	return msg.XBOSIoTDeviceState.WeatherCurrent!= nil
}
var device_units = map[string]string{
	"time":	"seconds",
	"icon":	"",
	"nearestStormDistance":	"miles",
	"nearestStormBearing":	"degrees",
	"precipIntensity":	"inches per hour",
	"precipIntensityError":	"",
	"precipProbability":	"",
	"precipType":	"",
	"temperature":	"F",
	"apparentTemperature":	"F",
	"dewPoint":	"F",
	"humidity":	"",
	"pressure":	"millibars",
	"windSpeed":	"miles per hour",
	"windGust":	"miles per hour",
	"windBearing":	"degree",
	"cloudCover":	"",
	"uvIndex":	"miles",
	"visibility":	"",
	"ozone":	"Dobson",
}
var device_lookup = map[string]func(msg xbospb.XBOS) (float64, bool){

	"time": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.Time != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.Time.Value), true
		}
		return 0, false
	},
	"nearestStormDistance": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.NearestStormDistance != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.NearestStormDistance.Value), true
		}
		return 0, false
	},
	"nearestStormBearing": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.NearestStormBearing != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.NearestStormBearing.Value), true
		}
		return 0, false
	},
	"precipIntensity": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.PrecipIntensity != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.PrecipIntensity.Value), true
		}
		return 0, false
	},
	"precipIntensityError": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.PrecipIntensityError != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.PrecipIntensityError.Value), true
		}
		return 0, false
	},
	"precipProbability": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.PrecipProbability != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.PrecipProbability.Value), true
		}
		return 0, false
	},
	"temperature": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.Temperature != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.Temperature.Value), true
		}
		return 0, false
	},
	"apparentTemperature": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.ApparentTemperature != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.ApparentTemperature.Value), true
		}
		return 0, false
	},
	"dewPoint": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.DewPoint != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.DewPoint.Value), true
		}
		return 0, false
	},
	"humidity": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.Humidity != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.Humidity.Value), true
		}
		return 0, false
	},
	"pressure": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.Pressure != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.Pressure.Value), true
		}
		return 0, false
	},
	"windSpeed": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.WindSpeed != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.WindSpeed.Value), true
		}
		return 0, false
	},
	"windGust": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.WindGust != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.WindGust.Value), true
		}
		return 0, false
	},
	"windBearing": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.WindBearing != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.WindBearing.Value), true
		}
		return 0, false
	},
	"cloudCover": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.CloudCover != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.CloudCover.Value), true
		}
		return 0, false
	},
	"uvIndex": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.UvIndex != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.UvIndex.Value), true
		}
		return 0, false
	},
	"visibility": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.Visibility != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.Visibility.Value), true
		}
		return 0, false
	},
	"ozone": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.XBOSIoTDeviceState.WeatherCurrent.Ozone != nil {
			return float64(msg.XBOSIoTDeviceState.WeatherCurrent.Ozone.Value), true
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
