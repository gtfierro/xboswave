package main

import (
	"fmt"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
	logrus "github.com/sirupsen/logrus"
	"os"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, ForceColors: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
}

func has_meter(msg xbospb.XBOS) bool {
	return msg.XBOSIoTDeviceState.Meter != nil
}

var lookup = map[string]func(msg xbospb.XBOS) (float64, bool){
	// XBOSIoTDeviceState.Meter
	"apparent_power": func(msg xbospb.XBOS) (float64, bool) {
		if has_meter(msg) && msg.XBOSIoTDeviceState.Meter.ApparentPower != nil {
			return float64(msg.XBOSIoTDeviceState.Meter.ApparentPower.Value), true
		}
		return 0, false
	},
	"power": func(msg xbospb.XBOS) (float64, bool) {
		if has_meter(msg) && msg.XBOSIoTDeviceState.Meter.Power != nil {
			return float64(msg.XBOSIoTDeviceState.Meter.Power.Value), true
		}
		return 0, false
	},
	"voltage": func(msg xbospb.XBOS) (float64, bool) {
		if has_meter(msg) && msg.XBOSIoTDeviceState.Meter.Voltage != nil {
			return float64(msg.XBOSIoTDeviceState.Meter.Voltage.Value), true
		}
		return 0, false
	},
	"energy": func(msg xbospb.XBOS) (float64, bool) {
		if has_meter(msg) && msg.XBOSIoTDeviceState.Meter.Energy != nil {
			return float64(msg.XBOSIoTDeviceState.Meter.Energy.Value), true
		}
		return 0, false
	},
}

var units = map[string]string{
	"apparent_power": "kVA",
	"current":        "A",
	"power":          "kW",
	"voltage":        "V",
	"energy":         "KWh",
}

func build(uri types.SubscriptionURI, name string, msg xbospb.XBOS) types.ExtractedTimeseries {

	if extractfunc, found := lookup[name]; found {
		if value, found := extractfunc(msg); found {
			var extracted types.ExtractedTimeseries
			time := int64(msg.XBOSIoTDeviceState.Time)
			extracted.Values = append(extracted.Values, value)
			extracted.Times = append(extracted.Times, time)
			extracted.UUID = types.GenerateUUID(uri, []byte(name))
			extracted.Collection = fmt.Sprintf("xbos/%s/meter", uri.Resource)
			extracted.Tags = map[string]string{
				"unit": units[name],
				"name": name,
			}
			return extracted
		}
	}

	return types.ExtractedTimeseries{}
}

func Extract(uri types.SubscriptionURI, msg xbospb.XBOS, add func(types.ExtractedTimeseries) error) error {
	if msg.XBOSIoTDeviceState != nil {
		if has_meter(msg) {
			for name := range lookup {
				extracted := build(uri, name, msg)
				if extracted.Empty() {
					continue
				}
				log.Debugf("Adding %s", extracted)
				if err := add(extracted); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
