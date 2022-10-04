package main
import (
	"fmt"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
)
func has_device(msg xbospb.XBOS) bool {
	return msg.FlexstatState!= nil
}
var device_units = map[string]string{
	"time": "seconds",
	"space_temp_sensor": "F",
	"minimum_proportional": "F",
	"active_cooling_setpt": "F",
	"active_heating_setpt": "F",
	"unocc_cooling_setpt": "F",
	"unocc_heating_setpt": "F",
	"occ_min_clg_setpt": "F",
	"occ_max_htg_setpt": "F",
	"override_timer": "hours",
	"occ_cooling_setpt": "F",
	"occ_heating_setpt": "F",
	"current_mode_setpt": "F",
	"ui_setpt": "F",
	"cooling_need": "F",
	"heating_need": "F",
	"unocc_min_clg_setpt": "F",
	"unocc_max_htg_setpt": "F",
	"min_setpt_diff": "F",
	"min_setpt_limit": "F",
	"space_temp": "F",
	"cooling_prop": "F",
	"heating_prop": "F",
	"cooling_intg": "per hour",
	"heating_intg": "per hour",
	"fan": "T/F",
	"occupancy_mode": "T/F",
	"setpt_override_mode": "T/F",
	"fan_alarm": "T/F",
	"fan_need": "T/F",
	"heating_cooling_mode": "T/F",
	"occ_fan_auto_on": "T/F",
	"unocc_fan_auto_on": "T/F",
	"fan_status": "T/F",
}

var device_lookup = map[string]func(msg xbospb.XBOS) (float64, bool){
	"space_temp_sensor": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.SpaceTempSensor != nil {
			return float64(msg.FlexstatState.SpaceTempSensor.Value), true
		}
		return 0, false
	},
	"minimum_proportional": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.MinimumProportional != nil {
			return float64(msg.FlexstatState.MinimumProportional.Value), true
		}
		return 0, false
	},
	"active_cooling_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.ActiveCoolingSetpt != nil {
			return float64(msg.FlexstatState.ActiveCoolingSetpt.Value), true
		}
		return 0, false
	},
	"active_heating_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.ActiveHeatingSetpt != nil {
			return float64(msg.FlexstatState.ActiveHeatingSetpt.Value), true
		}
		return 0, false
	},
	"unocc_cooling_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.UnoccCoolingSetpt != nil {
			return float64(msg.FlexstatState.UnoccCoolingSetpt.Value), true
		}
		return 0, false
	},
	"unocc_heating_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.UnoccHeatingSetpt != nil {
			return float64(msg.FlexstatState.UnoccHeatingSetpt.Value), true
		}
		return 0, false
	},
	"occ_min_clg_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.OccMinClgSetpt != nil {
			return float64(msg.FlexstatState.OccMinClgSetpt.Value), true
		}
		return 0, false
	},
	"occ_max_htg_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.OccMaxHtgSetpt != nil {
			return float64(msg.FlexstatState.OccMaxHtgSetpt.Value), true
		}
		return 0, false
	},
	"override_timer": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.OverrideTimer != nil {
			return float64(msg.FlexstatState.OverrideTimer.Value), true
		}
		return 0, false
	},
	"occ_cooling_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.OccCoolingSetpt != nil {
			return float64(msg.FlexstatState.OccCoolingSetpt.Value), true
		}
		return 0, false
	},
	"occ_heating_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.OccHeatingSetpt != nil {
			return float64(msg.FlexstatState.OccHeatingSetpt.Value), true
		}
		return 0, false
	},
	"current_mode_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.CurrentModeSetpt != nil {
			return float64(msg.FlexstatState.CurrentModeSetpt.Value), true
		}
		return 0, false
	},
	"ui_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.UiSetpt != nil {
			return float64(msg.FlexstatState.UiSetpt.Value), true
		}
		return 0, false
	},
	"cooling_need": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.CoolingNeed != nil {
			return float64(msg.FlexstatState.CoolingNeed.Value), true
		}
		return 0, false
	},
	"heating_need": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.HeatingNeed != nil {
			return float64(msg.FlexstatState.HeatingNeed.Value), true
		}
		return 0, false
	},
	"unocc_min_clg_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.UnoccMinClgSetpt != nil {
			return float64(msg.FlexstatState.UnoccMinClgSetpt.Value), true
		}
		return 0, false
	},
	"unocc_max_htg_setpt": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.UnoccMaxHtgSetpt != nil {
			return float64(msg.FlexstatState.UnoccMaxHtgSetpt.Value), true
		}
		return 0, false
	},
	"min_setpt_diff": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.MinSetptDiff != nil {
			return float64(msg.FlexstatState.MinSetptDiff.Value), true
		}
		return 0, false
	},
	"min_setpt_limit": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.MinSetptLimit != nil {
			return float64(msg.FlexstatState.MinSetptLimit.Value), true
		}
		return 0, false
	},
	"space_temp": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.SpaceTemp != nil {
			return float64(msg.FlexstatState.SpaceTemp.Value), true
		}
		return 0, false
	},
	"cooling_prop": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.CoolingProp != nil {
			return float64(msg.FlexstatState.CoolingProp.Value), true
		}
		return 0, false
	},
	"heating_prop": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.HeatingProp != nil {
			return float64(msg.FlexstatState.HeatingProp.Value), true
		}
		return 0, false
	},
	"cooling_intg": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.CoolingIntg != nil {
			return float64(msg.FlexstatState.CoolingIntg.Value), true
		}
		return 0, false
	},
	"heating_intg": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.HeatingIntg != nil {
			return float64(msg.FlexstatState.HeatingIntg.Value), true
		}
		return 0, false
	},
	"fan": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.Fan != nil {
			return float64(msg.FlexstatState.Fan.Value), true
		}
		return 0, false
	},
	"occupancy_mode": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.OccupancyMode != nil {
			return float64(msg.FlexstatState.OccupancyMode.Value), true
		}
		return 0, false
	},
	"setpt_override_mode": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.SetptOverrideMode != nil {
			return float64(msg.FlexstatState.SetptOverrideMode.Value), true
		}
		return 0, false
	},
	"fan_alarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.FanAlarm != nil {
			return float64(msg.FlexstatState.FanAlarm.Value), true
		}
		return 0, false
	},
	"fan_need": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.FanNeed != nil {
			return float64(msg.FlexstatState.FanNeed.Value), true
		}
		return 0, false
	},
	"heating_cooling_mode": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.HeatingCoolingMode != nil {
			return float64(msg.FlexstatState.HeatingCoolingMode.Value), true
		}
		return 0, false
	},
	"occ_fan_auto_on": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.HeatingCoolingMode != nil {
			return float64(msg.FlexstatState.HeatingCoolingMode.Value), true
		}
		return 0, false
	},
	"unocc_fan_auto_on": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.HeatingCoolingMode != nil {
			return float64(msg.FlexstatState.HeatingCoolingMode.Value), true
		}
		return 0, false
	},
	"fan_status": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.FlexstatState.FanStatus != nil {
			return float64(msg.FlexstatState.FanStatus.Value), true
		}
		return 0, false
	},
}
func build_device(uri types.SubscriptionURI, name string, msg xbospb.XBOS) types.ExtractedTimeseries {
	
	if extractfunc, found := device_lookup[name]; found {
		if value, found := extractfunc(msg); found {
			var extracted types.ExtractedTimeseries
			time := int64(msg.FlexstatState.Time)
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
	if msg.FlexstatState != nil {
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
