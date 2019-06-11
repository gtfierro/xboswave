package main

import (
	"fmt"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
)

type add_fn func(types.ExtractedTimeseries) error

func has_device(msg xbospb.XBOS) bool {
	return msg.XBOSIoTDeviceState.WeatherStationPrediction != nil
}

// This contains the mapping of each field's value to the unit
var device_units = map[string]string{
	"time":                 "seconds",
	"icon":                 "",
	"precipintensity":      "inches per hour",
	"precipintensityError": "",
	"precipprobability":    "",
	"temperature":          "F",
	"apparenttemperature":  "F",
	"dewpoint":             "F",
	"humidity":             "",
	"pressure":             "millibars",
	"windspeed":            "miles per hour",
	"windgust":             "miles per hour",
	"windbearing":          "degree",
	"cloudcover":           "",
	"uvindex":              "",
	"visibility":           "miles",
	"ozone":                "Dobson",
}

func ingest_time_series(value float64, name string, toInflux types.ExtractedTimeseries,
	pass_add add_fn, prediction_time int64, step int, uri types.SubscriptionURI) error {
	toInflux.Values = append(toInflux.Values, value)

	//This UUID is unique to each field in the message
	toInflux.UUID = types.GenerateUUID(uri, []byte(name))
	//The collection comes from the resource name of the driver
	toInflux.Collection = fmt.Sprintf("xbos/%s", uri.Resource)
	//These are the tags that will be used when the point is written
	toInflux.Tags = map[string]string{
		"unit":            device_units[name],
		"name":            name,
		"prediction_time": fmt.Sprintf("%d", prediction_time/1e9),
		"prediction_step": fmt.Sprintf("%d", step),
	}
	//This add function is passed in from the ingester and when it is executed
	//a point is written into influx
	if err := pass_add(toInflux); err != nil {
		return err
	}
	return nil
}

func Extract(uri types.SubscriptionURI, msg xbospb.XBOS, add func(types.ExtractedTimeseries) error) error {
	if msg.XBOSIoTDeviceState != nil {

		if has_device(msg) {
			step := 1

			//Iterate through each hour of prediction from current to 48 hours from current
			for _, _prediction := range msg.XBOSIoTDeviceState.WeatherStationPrediction.Predictions {
				//This prediction contains all of the fields that were present in WeatherCurrent message
				//There is one for each hour that is retrieved from the DarkSky API
				prediction := _prediction.Prediction

				//This will contain all the information necessary to send one prediction for one hour out of 0-48
				var extracted types.ExtractedTimeseries
				prediction_time := int64(_prediction.PredictionTime)

				//This is the xbos message time
				time := int64(msg.XBOSIoTDeviceState.Time)

				//This is the time that is being put into influx as the timestamp
				extracted.Times = append(extracted.Times, time)

				if prediction.Time != nil {
					err := ingest_time_series(float64(prediction.Time.Value),
						"time", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.PrecipIntensity != nil {
					err := ingest_time_series(float64(prediction.PrecipIntensity.Value),
						"precipintensity", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.PrecipIntensityError != nil {
					err := ingest_time_series(float64(prediction.PrecipIntensityError.Value),
						"precipintensityerror", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.PrecipProbability != nil {
					err := ingest_time_series(float64(prediction.PrecipProbability.Value),
						"precipprobability", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.Temperature != nil {
					err := ingest_time_series(float64(prediction.Temperature.Value),
						"temperature", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.ApparentTemperature != nil {
					err := ingest_time_series(float64(prediction.ApparentTemperature.Value),
						"apparenttemperature", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.DewPoint != nil {
					err := ingest_time_series(float64(prediction.DewPoint.Value),
						"dewpoint", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.Humidity != nil {
					err := ingest_time_series(float64(prediction.Humidity.Value),
						"humidity", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.Pressure != nil {
					err := ingest_time_series(float64(prediction.Pressure.Value),
						"pressure", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.WindSpeed != nil {
					err := ingest_time_series(float64(prediction.WindSpeed.Value),
						"windspeed", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.WindGust != nil {
					err := ingest_time_series(float64(prediction.WindGust.Value),
						"windgust", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.WindBearing != nil {
					err := ingest_time_series(float64(prediction.WindBearing.Value),
						"windbearing", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.CloudCover != nil {
					err := ingest_time_series(float64(prediction.CloudCover.Value),
						"cloudcover", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.UvIndex != nil {
					err := ingest_time_series(float64(prediction.UvIndex.Value),
						"uvindex", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.Visibility != nil {
					err := ingest_time_series(float64(prediction.Visibility.Value),
						"visibility", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}
				if prediction.Ozone != nil {
					err := ingest_time_series(float64(prediction.Ozone.Value),
						"ozone", extracted, add, prediction_time, step, uri)
					if err != nil {
						return err
					}
				}

				step++
			}
		}
	}
	return nil
}
