package main

import (
	"fmt"
	"github.com/gtfierro/xboswave/ingester/types"
	xbospb "github.com/gtfierro/xboswave/proto"
)

func has_device(msg xbospb.XBOS) bool {
	return msg.ParkerState != nil
}

var device_units = map[string]string{
	"compressor_working_hours":               "hours",
	"on_standby_status":                      "",
	"light_status":                           "",
	"aux_output_status":                      "",
	"next_defrost_counter":                   "seconds",
	"door_switch_input_status":               "",
	"multipurpose_input_status":              "",
	"compressor_status":                      "",
	"output_defrost_status":                  "",
	"fans_status":                            "",
	"output_k4_status":                       "",
	"cabinet_temperature":                    "C",
	"evaporator_temperature":                 "C",
	"auxiliary_temperature":                  "C",
	"probe1_failure_alarm":                   "",
	"probe2_failure_alarm":                   "",
	"probe3_failure_alarm":                   "",
	"minimum_temperature_alarm":              "",
	"maximum_temperture_alarm":               "",
	"condensor_temperature_failure_alarm":    "",
	"condensor_pre_alarm":                    "",
	"door_alarm":                             "",
	"multipurpose_input_alarm":               "",
	"compressor_blocked_alarm":               "",
	"power_failure_alarm":                    "",
	"rtc_error_alarm":                        "",
	"energy_saving_regulator_flag":           "",
	"energy_saving_real_time_regulator_flag": "",
	"service_request_regulator_flag":         "",
	"on_standby_regulator_flag":              "",
	"new_alarm_to_read_regulator_flag":       "",
	"defrost_status_regulator_flag":          "",
	"active_setpoint":                        "C",
	"time_until_defrost":                     "seconds",
	"current_defrost_counter":                "seconds",
	"compressor_delay":                       "seconds",
	"num_alarms_in_history":                  "",
	"energy_saving_status":                   "",
	"service_request_status":                 "",
	"resistors_activated_by_aux_key_status":  "",
	"evaporator_valve_state":                 "",
	"output_defrost_state":                   "",
	"output_lux_state":                       "",
	"output_aux_state":                       "",
	"resistors_state":                        "",
	"output_alarm_state":                     "",
	"second_compressor_state":                "",
	"setpoint":                               "",
	"r1":                                     "C",
	"r2":                                     "C",
	"r4":                                     "",
	"C0":                                     "minutes",
	"C1":                                     "minutes",
	"d0":                                     "hours",
	"d3":                                     "minutes",
	"d5":                                     "minutes",
	"d7":                                     "minutes",
	"d8":                                     "",
	"A0":                                     "",
	"A1":                                     "C",
	"A2":                                     "",
	"A3":                                     "",
	"A4":                                     "C",
	"A5":                                     "",
	"A6":                                     "minutes",
	"A7":                                     "minutes",
	"A8":                                     "minutes",
	"A9":                                     "minutes",
	"F0":                                     "",
	"F1":                                     "C",
	"F2":                                     "",
	"F3":                                     "minutes",
	"Hd1":                                    "hh:mm",
	"Hd2":                                    "hh:mm",
	"Hd3":                                    "hh:mm",
	"Hd4":                                    "hh:mm",
	"Hd5":                                    "hh:mm",
	"Hd6":                                    "hh:mm",
}
var device_lookup = map[string]func(msg xbospb.XBOS) (float64, bool){

	// Mapping of string to extract the value from part of Xbos message
	"CompressorWorkingHours": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.CompressorWorkingHours != nil {
			return float64(msg.ParkerState.CompressorWorkingHours.Value), true
		}
		return 0, false
	},
	"OnStandbyStatus": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.OnStandbyStatus != nil {
			return float64(msg.ParkerState.OnStandbyStatus.Value), true
		}
		return 0, false
	},
	"LightStatus": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.LightStatus != nil {
			return float64(msg.ParkerState.LightStatus.Value), true
		}
		return 0, false
	},
	"AuxOutputStatus": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.AuxOutputStatus != nil {
			return float64(msg.ParkerState.AuxOutputStatus.Value), true
		}
		return 0, false
	},
	"NextDefrostCounter": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.NextDefrostCounter != nil {
			return float64(msg.ParkerState.NextDefrostCounter.Value), true
		}
		return 0, false
	},
	"DoorSwitchInputStatus": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.DoorSwitchInputStatus != nil {
			return float64(msg.ParkerState.DoorSwitchInputStatus.Value), true
		}
		return 0, false
	},
	"MultipurposeInputStatus": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.MultipurposeInputStatus != nil {
			return float64(msg.ParkerState.MultipurposeInputStatus.Value), true
		}
		return 0, false
	},
	"CompressorStatus": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.CompressorStatus != nil {
			return float64(msg.ParkerState.CompressorStatus.Value), true
		}
		return 0, false
	},
	"OutputDefrostStatus": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.OutputDefrostStatus != nil {
			return float64(msg.ParkerState.OutputDefrostStatus.Value), true
		}
		return 0, false
	},
	"FansStatus": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.FansStatus != nil {
			return float64(msg.ParkerState.FansStatus.Value), true
		}
		return 0, false
	},
	"OutputK4Status": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.OutputK4Status != nil {
			return float64(msg.ParkerState.OutputK4Status.Value), true
		}
		return 0, false
	},
	"CabinetTemperature": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.CabinetTemperature != nil {
			return float64(msg.ParkerState.CabinetTemperature.Value), true
		}
		return 0, false
	},
	"EvaporatorTemperature": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.EvaporatorTemperature != nil {
			return float64(msg.ParkerState.EvaporatorTemperature.Value), true
		}
		return 0, false
	},
	"AuxiliaryTemperature": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.AuxiliaryTemperature != nil {
			return float64(msg.ParkerState.AuxiliaryTemperature.Value), true
		}
		return 0, false
	},
	"Probe1FailureAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.Probe1FailureAlarm != nil {
			return float64(msg.ParkerState.Probe1FailureAlarm.Value), true
		}
		return 0, false
	},
	"Probe2FailureAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.Probe2FailureAlarm != nil {
			return float64(msg.ParkerState.Probe2FailureAlarm.Value), true
		}
		return 0, false
	},
	"Probe3FailureAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.Probe3FailureAlarm != nil {
			return float64(msg.ParkerState.Probe3FailureAlarm.Value), true
		}
		return 0, false
	},
	"MinimumTemperatureAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.MinimumTemperatureAlarm != nil {
			return float64(msg.ParkerState.MinimumTemperatureAlarm.Value), true
		}
		return 0, false
	},
	"MaximumTempertureAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.MaximumTempertureAlarm != nil {
			return float64(msg.ParkerState.MaximumTempertureAlarm.Value), true
		}
		return 0, false
	},
	"CondensorTemperatureFailureAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.CondensorTemperatureFailureAlarm != nil {
			return float64(msg.ParkerState.CondensorTemperatureFailureAlarm.Value), true
		}
		return 0, false
	},
	"CondensorPreAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.CondensorPreAlarm != nil {
			return float64(msg.ParkerState.CondensorPreAlarm.Value), true
		}
		return 0, false
	},
	"DoorAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.DoorAlarm != nil {
			return float64(msg.ParkerState.DoorAlarm.Value), true
		}
		return 0, false
	},
	"MultipurposeInputAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.MultipurposeInputAlarm != nil {
			return float64(msg.ParkerState.MultipurposeInputAlarm.Value), true
		}
		return 0, false
	},
	"CompressorBlockedAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.CompressorBlockedAlarm != nil {
			return float64(msg.ParkerState.CompressorBlockedAlarm.Value), true
		}
		return 0, false
	},
	"PowerFailureAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.PowerFailureAlarm != nil {
			return float64(msg.ParkerState.PowerFailureAlarm.Value), true
		}
		return 0, false
	},
	"RtcErrorAlarm": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.RtcErrorAlarm != nil {
			return float64(msg.ParkerState.RtcErrorAlarm.Value), true
		}
		return 0, false
	},
	"EnergySavingRegulatorFlag": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.EnergySavingRegulatorFlag != nil {
			return float64(msg.ParkerState.EnergySavingRegulatorFlag.Value), true
		}
		return 0, false
	},
	"EnergySavingRealTimeRegulatorFlag": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.EnergySavingRealTimeRegulatorFlag != nil {
			return float64(msg.ParkerState.EnergySavingRealTimeRegulatorFlag.Value), true
		}
		return 0, false
	},
	"ServiceRequestRegulatorFlag": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.ServiceRequestRegulatorFlag != nil {
			return float64(msg.ParkerState.ServiceRequestRegulatorFlag.Value), true
		}
		return 0, false
	},
	"OnStandbyRegulatorFlag": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.OnStandbyRegulatorFlag != nil {
			return float64(msg.ParkerState.OnStandbyRegulatorFlag.Value), true
		}
		return 0, false
	},
	"NewAlarmToReadRegulatorFlag": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.NewAlarmToReadRegulatorFlag != nil {
			return float64(msg.ParkerState.NewAlarmToReadRegulatorFlag.Value), true
		}
		return 0, false
	},
	"DefrostStatusRegulatorFlag": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.DefrostStatusRegulatorFlag != nil {
			return float64(msg.ParkerState.DefrostStatusRegulatorFlag.Value), true
		}
		return 0, false
	},
	"ActiveSetpoint": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.ActiveSetpoint != nil {
			return float64(msg.ParkerState.ActiveSetpoint.Value), true
		}
		return 0, false
	},
	"TimeUntilDefrost": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.TimeUntilDefrost != nil {
			return float64(msg.ParkerState.TimeUntilDefrost.Value), true
		}
		return 0, false
	},
	"CurrentDefrostCounter": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.CurrentDefrostCounter != nil {
			return float64(msg.ParkerState.CurrentDefrostCounter.Value), true
		}
		return 0, false
	},
	"CompressorDelay": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.CompressorDelay != nil {
			return float64(msg.ParkerState.CompressorDelay.Value), true
		}
		return 0, false
	},
	"NumAlarmsInHistory": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.NumAlarmsInHistory != nil {
			return float64(msg.ParkerState.NumAlarmsInHistory.Value), true
		}
		return 0, false
	},
	"EnergySavingStatus": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.EnergySavingStatus != nil {
			return float64(msg.ParkerState.EnergySavingStatus.Value), true
		}
		return 0, false
	},
	"ServiceRequestStatus": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.ServiceRequestStatus != nil {
			return float64(msg.ParkerState.ServiceRequestStatus.Value), true
		}
		return 0, false
	},
	"ResistorsActivatedByAuxKeyStatus": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.ResistorsActivatedByAuxKeyStatus != nil {
			return float64(msg.ParkerState.ResistorsActivatedByAuxKeyStatus.Value), true
		}
		return 0, false
	},
	"EvaporatorValveState": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.EvaporatorValveState != nil {
			return float64(msg.ParkerState.EvaporatorValveState.Value), true
		}
		return 0, false
	},
	"OutputDefrostState": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.OutputDefrostState != nil {
			return float64(msg.ParkerState.OutputDefrostState.Value), true
		}
		return 0, false
	},
	"OutputLuxState": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.OutputLuxState != nil {
			return float64(msg.ParkerState.OutputLuxState.Value), true
		}
		return 0, false
	},
	"OutputAuxState": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.OutputAuxState != nil {
			return float64(msg.ParkerState.OutputAuxState.Value), true
		}
		return 0, false
	},
	"ResistorsState": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.ResistorsState != nil {
			return float64(msg.ParkerState.ResistorsState.Value), true
		}
		return 0, false
	},
	"OutputAlarmState": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.OutputAlarmState != nil {
			return float64(msg.ParkerState.OutputAlarmState.Value), true
		}
		return 0, false
	},
	"SecondCompressorState": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.SecondCompressorState != nil {
			return float64(msg.ParkerState.SecondCompressorState.Value), true
		}
		return 0, false
	},
	"Setpoint": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.Setpoint != nil {
			return float64(msg.ParkerState.Setpoint.Value), true
		}
		return 0, false
	},
	"R1": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.R1 != nil {
			return float64(msg.ParkerState.R1.Value), true
		}
		return 0, false
	},
	"R2": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.R2 != nil {
			return float64(msg.ParkerState.R2.Value), true
		}
		return 0, false
	},
	"R4": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.R4 != nil {
			return float64(msg.ParkerState.R4.Value), true
		}
		return 0, false
	},
	"C0": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.C0 != nil {
			return float64(msg.ParkerState.C0.Value), true
		}
		return 0, false
	},
	"C1": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.C1 != nil {
			return float64(msg.ParkerState.C1.Value), true
		}
		return 0, false
	},
	"D0": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.D0 != nil {
			return float64(msg.ParkerState.D0.Value), true
		}
		return 0, false
	},
	"D3": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.D3 != nil {
			return float64(msg.ParkerState.D3.Value), true
		}
		return 0, false
	},
	"D5": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.D5 != nil {
			return float64(msg.ParkerState.D5.Value), true
		}
		return 0, false
	},
	"D7": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.D7 != nil {
			return float64(msg.ParkerState.D7.Value), true
		}
		return 0, false
	},
	"D8": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.D8 != nil {
			return float64(msg.ParkerState.D8.Value), true
		}
		return 0, false
	},
	"A0": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.A0 != nil {
			return float64(msg.ParkerState.A0.Value), true
		}
		return 0, false
	},
	"A1": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.A1 != nil {
			return float64(msg.ParkerState.A1.Value), true
		}
		return 0, false
	},
	"A2": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.A2 != nil {
			return float64(msg.ParkerState.A2.Value), true
		}
		return 0, false
	},
	"A3": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.A3 != nil {
			return float64(msg.ParkerState.A3.Value), true
		}
		return 0, false
	},
	"A4": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.A4 != nil {
			return float64(msg.ParkerState.A4.Value), true
		}
		return 0, false
	},
	"A5": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.A5 != nil {
			return float64(msg.ParkerState.A5.Value), true
		}
		return 0, false
	},
	"A6": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.A6 != nil {
			return float64(msg.ParkerState.A6.Value), true
		}
		return 0, false
	},
	"A7": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.A7 != nil {
			return float64(msg.ParkerState.A7.Value), true
		}
		return 0, false
	},
	"A8": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.A8 != nil {
			return float64(msg.ParkerState.A8.Value), true
		}
		return 0, false
	},
	"A9": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.A9 != nil {
			return float64(msg.ParkerState.A9.Value), true
		}
		return 0, false
	},
	"F0": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.F0 != nil {
			return float64(msg.ParkerState.F0.Value), true
		}
		return 0, false
	},
	"F1": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.F1 != nil {
			return float64(msg.ParkerState.F1.Value), true
		}
		return 0, false
	},
	"F2": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.F2 != nil {
			return float64(msg.ParkerState.F2.Value), true
		}
		return 0, false
	},
	"F3": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.F3 != nil {
			return float64(msg.ParkerState.F3.Value), true
		}
		return 0, false
	},
	"Hd1": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.Hd1 != nil {
			return float64(msg.ParkerState.Hd1.Value), true
		}
		return 0, false
	},
	"Hd2": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.Hd2 != nil {
			return float64(msg.ParkerState.Hd2.Value), true
		}
		return 0, false
	},
	"Hd3": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.Hd3 != nil {
			return float64(msg.ParkerState.Hd3.Value), true
		}
		return 0, false
	},
	"Hd4": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.Hd4 != nil {
			return float64(msg.ParkerState.Hd4.Value), true
		}
		return 0, false
	},
	"Hd5": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.Hd5 != nil {
			return float64(msg.ParkerState.Hd5.Value), true
		}
		return 0, false
	},
	"Hd6": func(msg xbospb.XBOS) (float64, bool) {
		if has_device(msg) && msg.ParkerState.Hd6 != nil {
			return float64(msg.ParkerState.Hd6.Value), true
		}
		return 0, false
	},
}

func build_device(uri types.SubscriptionURI, name string, msg xbospb.XBOS) types.ExtractedTimeseries {

	if extractfunc, found := device_lookup[name]; found {
		if value, found := extractfunc(msg); found {
			var extracted types.ExtractedTimeseries
			time := int64(msg.ParkerState.Time)
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
			// Go through each Field in the Xbos Message
			for name := range device_lookup {
				extracted := build_device(uri, name, msg)
				//add function takes in an extracted timeseries and adds it to database
				if err := add(extracted); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
