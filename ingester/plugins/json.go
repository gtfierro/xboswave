package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
)

var units = map[string]string{
	"discharge_air_temperature_setpoint":               "F",
	"mixed_air_temperature_setpoint":                   "F",
	"economizer_differential_air_temperature_setpoint": "F",
	"zone_temperature_setpoint":                        "F",
	"thermostat_adjust_setpoint":                       "F",

	"supply_air_flow_setpoint":                      "cfm",
	"occupied_heating_min_supply_air_flow_setpoint": "cfm",

	"discharge_air_temperature_sensor": "F",
	"outside_air_temperature_sensor":   "F",
	"return_air_temperature_sensor":    "F",
	"zone_temperature_sensor":          "F",
	"mixed_air_temperature_sensor":     "F",

	"heating_valve_command": "%",
	"cooling_valve_command": "%",
	"occupancy_command":     "t/f",
	"shutdown_command":      "t/f",

	"supply_air_flow_sensor": "cfm",

	"building_static_pressure_setpoint":    "wg",
	"discharge_air_static_pressure_sensor": "wg",
	"building_static_pressure_sensor":      "wg",
	"supply_air_velocity_pressure_sensor":  "wg",

	"cooling_demand": "%",
	"heating_demand": "%",

	"supply_air_damper_min_position_setpoint":  "%",
	"cooling_max_supply_air_flow_setpoint":     "cfm",
	"mixed_air_temperature_low_limit_setpoint": "F",

	"damper_position_command": "%",
	"damper_position_sensor":  "%",
	"on_off_command":          "t/f",
	"fan_reset_command":       "t/f",

	"fan_overload_alarm": "t/f",
	"vfd_alarm":          "t/f",

	"fan_speed_setpoint": "%",

	"fan_status":           "unknown",
	"box_mode":             "unknown",
	"filter_status":        "unknown",
	"occupied_mode_status": "unknown",
}

func Extract(uri types.SubscriptionURI, msg xbospb.XBOS, add func(types.ExtractedTimeseries) error) error {
	marshaler := jsonpb.Marshaler{
		EnumsAsInts:  true,
		EmitDefaults: true,
		Indent:       "   ",
		OrigName:     true,
	}

	if msg.XBOSIoTDeviceState == nil {
		return nil
	}
	state := msg.XBOSIoTDeviceState

	b := new(bytes.Buffer)
	jsonmsg := make(map[string]interface{})
	dec := json.NewDecoder(b)
	if err := marshaler.Marshal(b, state); err != nil {
		return err
	}
	if err := dec.Decode(&jsonmsg); err != nil {
		return err
	}

	timestamp := int64(state.Time)
	delete(jsonmsg, "time")
	delete(jsonmsg, "requestid")
	delete(jsonmsg, "error")

	for equipname, _fields := range jsonmsg {
		fields, ok := _fields.(map[string]interface{})
		if !ok {
			continue
		}
		for fieldname, fieldvalue := range fields {
			if fieldvalue == nil {
				continue
			}
			var ex types.ExtractedTimeseries
			fv, ok := getValue(fieldvalue)
			if !ok {
				fmt.Println("skipping", equipname, fieldname, fieldvalue)
				continue
			}
			ex.Values = append(ex.Values, fv)
			ex.Times = append(ex.Times, timestamp)
			ex.UUID = types.GenerateUUID(uri, []byte(equipname+fieldname))
			ex.Collection = fmt.Sprintf("xbos/%s", uri.Resource)
			ex.Tags = map[string]string{
				"unit": units[fieldname],
				"name": fieldname,
			}
			if err := add(ex); err != nil {
				return err
			}
		}
	}

	return nil
}

func build(uri types.SubscriptionURI, name string, dict map[string]interface{}) types.ExtractedTimeseries {
	return types.ExtractedTimeseries{}
}

func getValue(i interface{}) (float64, bool) {
	if m, ok := i.(map[string]interface{}); ok {
		val, vok := m["value"].(float64)
		return val, vok
	}
	return -1, false
}
