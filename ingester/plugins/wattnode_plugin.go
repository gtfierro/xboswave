package main
import (
	"fmt"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
)
func has_device(msg xbospb.XBOS) bool {
	return msg.WattnodeState!= nil
}
var device_units = map[string]string{
	"EnergySum":	"kWh",
	"EnergyPosSum":	"kWh",
	"EnergySumNR":	"kWh",
	"EnergyPosSumNr":	"kWh",
	"PowerSum":	"W",
	"PowerA":	"W",
	"PowerB":	"W",
	"PowerC":	"W",
	"VoltAvgLN":	"V",
	"VoltA":	"V",
	"VoltB":	"V",
	"VoltC":	"V",
	"VoltAvgLL":	"V",
	"VoltAB":	"V",
	"VoltBC":	"V",
	"VoltAC":	"V",
	"Freq":	"Hz",
	"EnergyA":	"kWh",
	"EnergyB":	"kWh",
	"EnergyC":	"kWh",
	"EnergyPosA":	"kWh",
	"EnergyPosB":	"kWh",
	"EnergyPosC":	"kWh",
	"EnergyNegSum":	"kWh",
	"EnergyNegSumNR":	"kWh",
	"EnergyNegA":	"kWh",
	"EnergyNegB":	"kWh",
	"EnergyNegC":	"kWh",
	"EnergyReacSum":	"kVARh",
	"EnergyReacA":	"kVARh",
	"EnergyReacB":	"kVARh",
	"EnergyReacC":	"kVARh",
	"EnergyAppSum":	"kVAh",
	"EnergyAppA":	"kVAh",
	"EnergyAppB":	"kVAh",
	"EnergyAppC":	"kVAh",
	"PowerFactorAvg":	"",
	"PowerFactorA":	"",
	"PowerFactorB":	"",
	"PowerFactorC":	"",
	"PowerReacSum":	"VAR",
	"PowerReacA":	"VAR",
	"PowerReacB":	"VAR",
	"PowerReacC":	"VAR",
	"PowerAppSum":	"VA",
	"PowerAppA":	"VA",
	"PowerAppB":	"VA",
	"PowerAppC":	"VA",
	"CurrentA":	"A",
	"CurrentB":	"A",
	"CurrentC":	"A",
	"Demand":	"W",
	"DemandMin":	"W",
	"DemandMax":	"W",
	"DemandApp":	"W",
	"DemandA":	"W",
	"DemandB":	"W",
	"DemandC":	"W",
}
var device_lookup = map[string]func(msg xbospb.XBOS) (float64, bool){

	"EnergySum": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergySum != nil {
			return float64(msg.WattnodeState.EnergySum.Value), true
		}
		return 0, false
	},
	"EnergyPosSum": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyPosSum != nil {
			return float64(msg.WattnodeState.EnergyPosSum.Value), true
		}
		return 0, false
	},
	"EnergySumNR": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergySumNR != nil {
			return float64(msg.WattnodeState.EnergySumNR.Value), true
		}
		return 0, false
	},
	"EnergyPosSumNr": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyPosSumNr != nil {
			return float64(msg.WattnodeState.EnergyPosSumNr.Value), true
		}
		return 0, false
	},
	"PowerSum": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerSum != nil {
			return float64(msg.WattnodeState.PowerSum.Value), true
		}
		return 0, false
	},
	"PowerA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerA != nil {
			return float64(msg.WattnodeState.PowerA.Value), true
		}
		return 0, false
	},
	"PowerB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerB != nil {
			return float64(msg.WattnodeState.PowerB.Value), true
		}
		return 0, false
	},
	"PowerC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerC != nil {
			return float64(msg.WattnodeState.PowerC.Value), true
		}
		return 0, false
	},
	"VoltAvgLN": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.VoltAvgLN != nil {
			return float64(msg.WattnodeState.VoltAvgLN.Value), true
		}
		return 0, false
	},
	"VoltA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.VoltA != nil {
			return float64(msg.WattnodeState.VoltA.Value), true
		}
		return 0, false
	},
	"VoltB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.VoltB != nil {
			return float64(msg.WattnodeState.VoltB.Value), true
		}
		return 0, false
	},
	"VoltC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.VoltC != nil {
			return float64(msg.WattnodeState.VoltC.Value), true
		}
		return 0, false
	},
	"VoltAvgLL": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.VoltAvgLL != nil {
			return float64(msg.WattnodeState.VoltAvgLL.Value), true
		}
		return 0, false
	},
	"VoltAB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.VoltAB != nil {
			return float64(msg.WattnodeState.VoltAB.Value), true
		}
		return 0, false
	},
	"VoltBC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.VoltBC != nil {
			return float64(msg.WattnodeState.VoltBC.Value), true
		}
		return 0, false
	},
	"VoltAC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.VoltAC != nil {
			return float64(msg.WattnodeState.VoltAC.Value), true
		}
		return 0, false
	},
	"Freq": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.Freq != nil {
			return float64(msg.WattnodeState.Freq.Value), true
		}
		return 0, false
	},
	"EnergyA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyA != nil {
			return float64(msg.WattnodeState.EnergyA.Value), true
		}
		return 0, false
	},
	"EnergyB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyB != nil {
			return float64(msg.WattnodeState.EnergyB.Value), true
		}
		return 0, false
	},
	"EnergyC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyC != nil {
			return float64(msg.WattnodeState.EnergyC.Value), true
		}
		return 0, false
	},
	"EnergyPosA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyPosA != nil {
			return float64(msg.WattnodeState.EnergyPosA.Value), true
		}
		return 0, false
	},
	"EnergyPosB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyPosB != nil {
			return float64(msg.WattnodeState.EnergyPosB.Value), true
		}
		return 0, false
	},
	"EnergyPosC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyPosC != nil {
			return float64(msg.WattnodeState.EnergyPosC.Value), true
		}
		return 0, false
	},
	"EnergyNegSum": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyNegSum != nil {
			return float64(msg.WattnodeState.EnergyNegSum.Value), true
		}
		return 0, false
	},
	"EnergyNegSumNR": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyNegSumNR != nil {
			return float64(msg.WattnodeState.EnergyNegSumNR.Value), true
		}
		return 0, false
	},
	"EnergyNegA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyNegA != nil {
			return float64(msg.WattnodeState.EnergyNegA.Value), true
		}
		return 0, false
	},
	"EnergyNegB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyNegB != nil {
			return float64(msg.WattnodeState.EnergyNegB.Value), true
		}
		return 0, false
	},
	"EnergyNegC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyNegC != nil {
			return float64(msg.WattnodeState.EnergyNegC.Value), true
		}
		return 0, false
	},
	"EnergyReacSum": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyReacSum != nil {
			return float64(msg.WattnodeState.EnergyReacSum.Value), true
		}
		return 0, false
	},
	"EnergyReacA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyReacA != nil {
			return float64(msg.WattnodeState.EnergyReacA.Value), true
		}
		return 0, false
	},
	"EnergyReacB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyReacB != nil {
			return float64(msg.WattnodeState.EnergyReacB.Value), true
		}
		return 0, false
	},
	"EnergyReacC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyReacC != nil {
			return float64(msg.WattnodeState.EnergyReacC.Value), true
		}
		return 0, false
	},
	"EnergyAppSum": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyAppSum != nil {
			return float64(msg.WattnodeState.EnergyAppSum.Value), true
		}
		return 0, false
	},
	"EnergyAppA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyAppA != nil {
			return float64(msg.WattnodeState.EnergyAppA.Value), true
		}
		return 0, false
	},
	"EnergyAppB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyAppB != nil {
			return float64(msg.WattnodeState.EnergyAppB.Value), true
		}
		return 0, false
	},
	"EnergyAppC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.EnergyAppC != nil {
			return float64(msg.WattnodeState.EnergyAppC.Value), true
		}
		return 0, false
	},
	"PowerFactorAvg": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerFactorAvg != nil {
			return float64(msg.WattnodeState.PowerFactorAvg.Value), true
		}
		return 0, false
	},
	"PowerFactorA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerFactorA != nil {
			return float64(msg.WattnodeState.PowerFactorA.Value), true
		}
		return 0, false
	},
	"PowerFactorB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerFactorB != nil {
			return float64(msg.WattnodeState.PowerFactorB.Value), true
		}
		return 0, false
	},
	"PowerFactorC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerFactorC != nil {
			return float64(msg.WattnodeState.PowerFactorC.Value), true
		}
		return 0, false
	},
	"PowerReacSum": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerReacSum != nil {
			return float64(msg.WattnodeState.PowerReacSum.Value), true
		}
		return 0, false
	},
	"PowerReacA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerReacA != nil {
			return float64(msg.WattnodeState.PowerReacA.Value), true
		}
		return 0, false
	},
	"PowerReacB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerReacB != nil {
			return float64(msg.WattnodeState.PowerReacB.Value), true
		}
		return 0, false
	},
	"PowerReacC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerReacC != nil {
			return float64(msg.WattnodeState.PowerReacC.Value), true
		}
		return 0, false
	},
	"PowerAppSum": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerAppSum != nil {
			return float64(msg.WattnodeState.PowerAppSum.Value), true
		}
		return 0, false
	},
	"PowerAppA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerAppA != nil {
			return float64(msg.WattnodeState.PowerAppA.Value), true
		}
		return 0, false
	},
	"PowerAppB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerAppB != nil {
			return float64(msg.WattnodeState.PowerAppB.Value), true
		}
		return 0, false
	},
	"PowerAppC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.PowerAppC != nil {
			return float64(msg.WattnodeState.PowerAppC.Value), true
		}
		return 0, false
	},
	"CurrentA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.CurrentA != nil {
			return float64(msg.WattnodeState.CurrentA.Value), true
		}
		return 0, false
	},
	"CurrentB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.CurrentB != nil {
			return float64(msg.WattnodeState.CurrentB.Value), true
		}
		return 0, false
	},
	"CurrentC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.CurrentC != nil {
			return float64(msg.WattnodeState.CurrentC.Value), true
		}
		return 0, false
	},
	"Demand": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.Demand != nil {
			return float64(msg.WattnodeState.Demand.Value), true
		}
		return 0, false
	},
	"DemandMin": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.DemandMin != nil {
			return float64(msg.WattnodeState.DemandMin.Value), true
		}
		return 0, false
	},
	"DemandMax": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.DemandMax != nil {
			return float64(msg.WattnodeState.DemandMax.Value), true
		}
		return 0, false
	},
	"DemandApp": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.DemandApp != nil {
			return float64(msg.WattnodeState.DemandApp.Value), true
		}
		return 0, false
	},
	"DemandA": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.DemandA != nil {
			return float64(msg.WattnodeState.DemandA.Value), true
		}
		return 0, false
	},
	"DemandB": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.DemandB != nil {
			return float64(msg.WattnodeState.DemandB.Value), true
		}
		return 0, false
	},
	"DemandC": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.WattnodeState.DemandC != nil {
			return float64(msg.WattnodeState.DemandC.Value), true
		}
		return 0, false
	},
}
func build_device(uri types.SubscriptionURI, name string, msg xbospb.XBOS) types.ExtractedTimeseries {
	
	if extractfunc, found := device_lookup[name]; found {
		if value, found := extractfunc(msg); found {
			var extracted types.ExtractedTimeseries
			time := int64(msg.WattnodeState.Time)
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
	if msg.WattnodeState != nil {
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
