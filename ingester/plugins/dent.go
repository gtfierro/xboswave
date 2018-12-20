package main

import (
	"fmt"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
	"strings"
)

var lookup = map[string]func(msg xbospb.XBOS, idx int) float64{
	"apparent_energy": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].ApparentEnergy)
	},
	"apparent_pf": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].ApparentPf)
	},
	"apparent_power": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].ApparentPower)
	},
	"current": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].Current)
	},
	"displacement_pf": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].DisplacementPf)
	},
	"line_frequency": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].LineFrequency)
	},
	"phase_neutral_voltage": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].PhaseNeutralVoltage)
	},
	"reactive_energy": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].ReactiveEnergy)
	},
	"reactive_power": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].ReactivePower)
	},
	"true_energy": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].TrueEnergy)
	},
	"true_power": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].TruePower)
	},
	"volts": func(msg xbospb.XBOS, idx int) float64 {
		return float64(msg.DentMeterState.Phases[idx].Volts)
	},
}

var units = map[string]string{
	"apparent_energy":       "kVAh",
	"apparent_pf":           "PF",
	"apparent_power":        "kVA",
	"current":               "A",
	"displacement_pf":       "PF",
	"line_frequency":        "HZ",
	"phase_neutral_voltage": "V",
	"reactive_energy":       "kVARh",
	"reactive_power":        "kVAR",
	"true_energy":           "kWh",
	"true_power":            "kW",
	"volts":                 "V",
}

func build(uri types.SubscriptionURI, name string, msg xbospb.XBOS, idx int) types.ExtractedTimeseries {

	value := lookup[name](msg, idx)
	phasename := msg.DentMeterState.Phases[idx].Phase
	eltname := msg.DentMeterState.Phases[idx].Annotation

	var extracted types.ExtractedTimeseries
	time := int64(msg.DentMeterState.Time)
	extracted.Values = append(extracted.Values, value)
	extracted.Times = append(extracted.Times, time)

	// incomming uri is dentmeter/<bldg>/<meterid>
	parts := strings.Split(uri.Resource, "/")
	// remove initial part
	baseuri := strings.Join(parts[1:], "/") // now just <bldg>/<meterid>
	extracted.UUID = types.GenerateUUID(uri, []byte(uri.Resource+name+phasename+eltname+"dent_meter_archive_1"))
	extracted.Collection = fmt.Sprintf("dent_meter/berkeley/%s/%s/%s", baseuri, eltname, phasename) //uri.Resource
	//fmt.Println(extracted.Collection, extracted.UUID.String())
	extracted.Tags = map[string]string{
		"unit":  units[name],
		"name":  name,
		"phase": phasename,
		"elt":   eltname,
	}
	return extracted
}

func Extract(uri types.SubscriptionURI, msg xbospb.XBOS, add func(types.ExtractedTimeseries) error) error {
	if msg.DentMeterState != nil {
		for idx := range msg.DentMeterState.Phases {
			for name := range lookup {
				extracted := build(uri, name, msg, idx)
				//fmt.Println(idx, name, uri, len(extracted.Values))
				if err := add(extracted); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
