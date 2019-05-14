// Code generated by protoc-gen-go. DO NOT EDIT.
// source: parker.proto

package xbospb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ParkerState struct {
	//unit: hours
	CompressorWorkingHours *Double `protobuf:"bytes,1,opt,name=compressor_working_hours,json=compressorWorkingHours,proto3" json:"compressor_working_hours,omitempty"`
	//on/standby
	OnStandbyStatus *Int64 `protobuf:"bytes,2,opt,name=on_standby_status,json=onStandbyStatus,proto3" json:"on_standby_status,omitempty"`
	LightStatus     *Int64 `protobuf:"bytes,3,opt,name=light_status,json=lightStatus,proto3" json:"light_status,omitempty"`
	AuxOutputStatus *Int64 `protobuf:"bytes,4,opt,name=aux_output_status,json=auxOutputStatus,proto3" json:"aux_output_status,omitempty"`
	//counter reduces periodically (in 1/4 of a minute)
	//unit: seconds
	NextDefrostCounter *Double `protobuf:"bytes,5,opt,name=next_defrost_counter,json=nextDefrostCounter,proto3" json:"next_defrost_counter,omitempty"`
	//digital_io_status & 0x0001
	DoorSwitchInputStatus *Int64 `protobuf:"bytes,6,opt,name=door_switch_input_status,json=doorSwitchInputStatus,proto3" json:"door_switch_input_status,omitempty"`
	//digital_io_status & 0x0002
	MultipurposeInputStatus *Int64 `protobuf:"bytes,7,opt,name=multipurpose_input_status,json=multipurposeInputStatus,proto3" json:"multipurpose_input_status,omitempty"`
	//digital_io_status & 0x0100
	CompressorStatus *Int64 `protobuf:"bytes,8,opt,name=compressor_status,json=compressorStatus,proto3" json:"compressor_status,omitempty"`
	//digital_io_status & 0x0200
	OutputDefrostStatus *Int64 `protobuf:"bytes,9,opt,name=output_defrost_status,json=outputDefrostStatus,proto3" json:"output_defrost_status,omitempty"`
	//digital_io_status & 0x0400
	FansStatus *Int64 `protobuf:"bytes,10,opt,name=fans_status,json=fansStatus,proto3" json:"fans_status,omitempty"`
	//digital_io_status & 0x0800
	OutputK4Status *Int64 `protobuf:"bytes,11,opt,name=output_k4_status,json=outputK4Status,proto3" json:"output_k4_status,omitempty"`
	//temperature measured by cabinet probe (in 10x actual value)
	//unit: C
	CabinetTemperature *Double `protobuf:"bytes,12,opt,name=cabinet_temperature,json=cabinetTemperature,proto3" json:"cabinet_temperature,omitempty"`
	//temperature measured by evaporator probe (in 10x actual value)
	//unit: C
	EvaporatorTemperature *Double `protobuf:"bytes,13,opt,name=evaporator_temperature,json=evaporatorTemperature,proto3" json:"evaporator_temperature,omitempty"`
	//temperature measured by auxiliary probe (if present) (in 10x actual value)
	//unit: C
	AuxiliaryTemperature *Double `protobuf:"bytes,14,opt,name=auxiliary_temperature,json=auxiliaryTemperature,proto3" json:"auxiliary_temperature,omitempty"`
	//alarm_status & 0x0100
	Probe1FailureAlarm *Int64 `protobuf:"bytes,15,opt,name=probe1_failure_alarm,json=probe1FailureAlarm,proto3" json:"probe1_failure_alarm,omitempty"`
	//alarm_status & 0x0200
	Probe2FailureAlarm *Int64 `protobuf:"bytes,16,opt,name=probe2_failure_alarm,json=probe2FailureAlarm,proto3" json:"probe2_failure_alarm,omitempty"`
	//alarm_status & 0x0400
	Probe3FailureAlarm *Int64 `protobuf:"bytes,17,opt,name=probe3_failure_alarm,json=probe3FailureAlarm,proto3" json:"probe3_failure_alarm,omitempty"`
	//alarm_status & 0x1000
	MinimumTemperatureAlarm *Int64 `protobuf:"bytes,18,opt,name=minimum_temperature_alarm,json=minimumTemperatureAlarm,proto3" json:"minimum_temperature_alarm,omitempty"`
	//alarm_status & 0x2000
	MaximumTempertureAlarm *Int64 `protobuf:"bytes,19,opt,name=maximum_temperture_alarm,json=maximumTempertureAlarm,proto3" json:"maximum_temperture_alarm,omitempty"`
	//alarm_status & 0x4000
	CondensorTemperatureFailureAlarm *Int64 `protobuf:"bytes,20,opt,name=condensor_temperature_failure_alarm,json=condensorTemperatureFailureAlarm,proto3" json:"condensor_temperature_failure_alarm,omitempty"`
	//alarm_status & 0x8000
	CondensorPreAlarm *Int64 `protobuf:"bytes,21,opt,name=condensor_pre_alarm,json=condensorPreAlarm,proto3" json:"condensor_pre_alarm,omitempty"`
	//alarm_status & 0x0004
	DoorAlarm *Int64 `protobuf:"bytes,22,opt,name=door_alarm,json=doorAlarm,proto3" json:"door_alarm,omitempty"`
	//alarm_status & 0x0008
	MultipurposeInputAlarm *Int64 `protobuf:"bytes,23,opt,name=multipurpose_input_alarm,json=multipurposeInputAlarm,proto3" json:"multipurpose_input_alarm,omitempty"`
	//alarm_status & 0x0010
	CompressorBlockedAlarm *Int64 `protobuf:"bytes,24,opt,name=compressor_blocked_alarm,json=compressorBlockedAlarm,proto3" json:"compressor_blocked_alarm,omitempty"`
	//alarm_status & 0x0020
	PowerFailureAlarm *Int64 `protobuf:"bytes,25,opt,name=power_failure_alarm,json=powerFailureAlarm,proto3" json:"power_failure_alarm,omitempty"`
	//alarm_status & 0x0080
	RtcErrorAlarm *Int64 `protobuf:"bytes,26,opt,name=rtc_error_alarm,json=rtcErrorAlarm,proto3" json:"rtc_error_alarm,omitempty"`
	//regulator_flag_1 & 0x0100
	EnergySavingRegulatorFlag *Int64 `protobuf:"bytes,27,opt,name=energy_saving_regulator_flag,json=energySavingRegulatorFlag,proto3" json:"energy_saving_regulator_flag,omitempty"`
	//regulator_flag_1 & 0x0200
	EnergySavingRealTimeRegulatorFlag *Int64 `protobuf:"bytes,28,opt,name=energy_saving_real_time_regulator_flag,json=energySavingRealTimeRegulatorFlag,proto3" json:"energy_saving_real_time_regulator_flag,omitempty"`
	//regulator_flag_1 & 0x0400
	ServiceRequestRegulatorFlag *Int64 `protobuf:"bytes,29,opt,name=service_request_regulator_flag,json=serviceRequestRegulatorFlag,proto3" json:"service_request_regulator_flag,omitempty"`
	//regulator_flag_2 & 0x0001; 1=standby
	OnStandbyRegulatorFlag *Int64 `protobuf:"bytes,30,opt,name=on_standby_regulator_flag,json=onStandbyRegulatorFlag,proto3" json:"on_standby_regulator_flag,omitempty"`
	//regulator_flag_2 & 0x0080
	NewAlarmToReadRegulatorFlag *Int64 `protobuf:"bytes,31,opt,name=new_alarm_to_read_regulator_flag,json=newAlarmToReadRegulatorFlag,proto3" json:"new_alarm_to_read_regulator_flag,omitempty"`
	//regulator_flag_2 & 0x0700; 0/1/2/3 = no defrost active/defrost running/dripping/fans stop
	DefrostStatusRegulatorFlag *Int64 `protobuf:"bytes,32,opt,name=defrost_status_regulator_flag,json=defrostStatusRegulatorFlag,proto3" json:"defrost_status_regulator_flag,omitempty"`
	//active_setpoint=setpoint(when no energy saving); else=setpoint+r4
	//unit: C
	ActiveSetpoint *Int64 `protobuf:"bytes,33,opt,name=active_setpoint,json=activeSetpoint,proto3" json:"active_setpoint,omitempty"`
	//time remaining to next defrost
	//unit: seconds
	TimeUntilDefrost *Int64 `protobuf:"bytes,34,opt,name=time_until_defrost,json=timeUntilDefrost,proto3" json:"time_until_defrost,omitempty"`
	//current defrost counter countdown (in 1/4 of a minute)
	//unit: seconds
	CurrentDefrostCounter *Int64 `protobuf:"bytes,35,opt,name=current_defrost_counter,json=currentDefrostCounter,proto3" json:"current_defrost_counter,omitempty"`
	//compressor delay in seconds
	//unit: seconds
	CompressorDelay *Int64 `protobuf:"bytes,36,opt,name=compressor_delay,json=compressorDelay,proto3" json:"compressor_delay,omitempty"`
	//number of HACCP alarms in history (max of last 9 stored)
	NumAlarmsInHistory *Int64 `protobuf:"bytes,37,opt,name=num_alarms_in_history,json=numAlarmsInHistory,proto3" json:"num_alarms_in_history,omitempty"`
	//is energy saving mode active or not; digital_output_flags & 0x0100
	EnergySavingStatus *Int64 `protobuf:"bytes,38,opt,name=energy_saving_status,json=energySavingStatus,proto3" json:"energy_saving_status,omitempty"`
	//digital_output_flags & 0x0200
	ServiceRequestStatus *Int64 `protobuf:"bytes,39,opt,name=service_request_status,json=serviceRequestStatus,proto3" json:"service_request_status,omitempty"`
	//digital_output_flags & 0x001
	ResistorsActivatedByAuxKeyStatus *Int64 `protobuf:"bytes,40,opt,name=resistors_activated_by_aux_key_status,json=resistorsActivatedByAuxKeyStatus,proto3" json:"resistors_activated_by_aux_key_status,omitempty"`
	//digital_output_flags & 0x002
	EvaporatorValveState *Int64 `protobuf:"bytes,41,opt,name=evaporator_valve_state,json=evaporatorValveState,proto3" json:"evaporator_valve_state,omitempty"`
	//digital_output_flags & 0x004
	OutputDefrostState *Int64 `protobuf:"bytes,42,opt,name=output_defrost_state,json=outputDefrostState,proto3" json:"output_defrost_state,omitempty"`
	//digital_output_flags & 0x008
	OutputLuxState *Int64 `protobuf:"bytes,43,opt,name=output_lux_state,json=outputLuxState,proto3" json:"output_lux_state,omitempty"`
	//digital_output_flags & 0x0010
	OutputAuxState *Int64 `protobuf:"bytes,44,opt,name=output_aux_state,json=outputAuxState,proto3" json:"output_aux_state,omitempty"`
	//activated by cabinet probe; digital_output_flags & 0x0020
	ResistorsState *Int64 `protobuf:"bytes,45,opt,name=resistors_state,json=resistorsState,proto3" json:"resistors_state,omitempty"`
	//digital_output_flags & 0x0040
	OutputAlarmState *Int64 `protobuf:"bytes,46,opt,name=output_alarm_state,json=outputAlarmState,proto3" json:"output_alarm_state,omitempty"`
	//digital_output_flags & 0x0080
	SecondCompressorState *Int64 `protobuf:"bytes,47,opt,name=second_compressor_state,json=secondCompressorState,proto3" json:"second_compressor_state,omitempty"`
	//setpoint
	Setpoint *Double `protobuf:"bytes,48,opt,name=setpoint,proto3" json:"setpoint,omitempty"`
	//min working setpoint
	//unit: C
	R1 *Double `protobuf:"bytes,49,opt,name=r1,proto3" json:"r1,omitempty"`
	//max working setpoint
	//unit: C
	R2 *Double `protobuf:"bytes,50,opt,name=r2,proto3" json:"r2,omitempty"`
	//used for active_set_point calculation in energy saving mode; adds to active setpoint
	R4 *Double `protobuf:"bytes,51,opt,name=r4,proto3" json:"r4,omitempty"`
	//compressor delay after turning on controller
	//unit: minutes
	C0 *Double `protobuf:"bytes,52,opt,name=C0,proto3" json:"C0,omitempty"`
	//min time between 2 activations in succession of compressor
	//unit: minutes
	C1 *Double `protobuf:"bytes,53,opt,name=C1,proto3" json:"C1,omitempty"`
	//defrost interval (only if d8 = 0/1/2); 0 = the defrost at intervals will never be activated
	//unit: hours
	D0 *Double `protobuf:"bytes,54,opt,name=d0,proto3" json:"d0,omitempty"`
	//defrost duration if P3=0 or 2; max duration if P3=1
	//unit: minutes
	D3 *Double `protobuf:"bytes,55,opt,name=d3,proto3" json:"d3,omitempty"`
	//defrost delay when you turn on controller; only if d4=1
	//unit: minutes
	D5 *Double `protobuf:"bytes,56,opt,name=d5,proto3" json:"d5,omitempty"`
	//drip delay
	//unit: minutes
	D7 *Double `protobuf:"bytes,57,opt,name=d7,proto3" json:"d7,omitempty"`
	//kind of defrost interval; 0/1/2/3=defrost on when controller/compressor/evaporator temperature is below d9  is on for d0 hours/realtime
	D8 *Int64 `protobuf:"bytes,58,opt,name=d8,proto3" json:"d8,omitempty"`
	//measured input for low temp alarm; 0/1/2=cab/evap/aux (only if P4=1/2)
	A0 *Int64 `protobuf:"bytes,59,opt,name=A0,proto3" json:"A0,omitempty"`
	//temperature below which low temperature alarm is activated
	//unit: C
	A1 *Double `protobuf:"bytes,60,opt,name=A1,proto3" json:"A1,omitempty"`
	//kind of lower temp alarm; 0/1/2=disabled/working setpoint-A1/absolute (or A1)
	A2 *Int64 `protobuf:"bytes,61,opt,name=A2,proto3" json:"A2,omitempty"`
	//measured input for high temp alarm; 0/1/2=cab/evap/aux (only if P4=1/2)
	A3 *Int64 `protobuf:"bytes,62,opt,name=A3,proto3" json:"A3,omitempty"`
	//temperature above which high temperature alarm is activated
	//unit: C
	A4 *Double `protobuf:"bytes,63,opt,name=A4,proto3" json:"A4,omitempty"`
	//kind of high temp alarm; 0/1/2=disabled/working setpoint+A4/absolute (or A4)
	A5 *Int64 `protobuf:"bytes,64,opt,name=A5,proto3" json:"A5,omitempty"`
	//high temperature alarm delay after turning on controller; only if A3=0
	//unit: minutes
	A6 *Double `protobuf:"bytes,65,opt,name=A6,proto3" json:"A6,omitempty"`
	//temperature alarm delay
	//unit: minutes
	A7 *Double `protobuf:"bytes,66,opt,name=A7,proto3" json:"A7,omitempty"`
	//high temperature alarm delay after end of defrost; only if A3=0
	//unit: minutes
	A8 *Double `protobuf:"bytes,67,opt,name=A8,proto3" json:"A8,omitempty"`
	//high temperature alarm delay after deactivation of microport input only if A3=0
	//unit: minutes
	A9 *Double `protobuf:"bytes,68,opt,name=A9,proto3" json:"A9,omitempty"`
	//evap fan activity during normal operation; 0/1/2/3/4=off/on/in parallel with compressor/dependent on F1/off if compressor is off and depedent on F1 if compressor is on
	F0 *Int64 `protobuf:"bytes,69,opt,name=F0,proto3" json:"F0,omitempty"`
	//evap temperature above which evap fan is turned off; only if F0=3/4
	//unit: C
	F1 *Double `protobuf:"bytes,70,opt,name=F1,proto3" json:"F1,omitempty"`
	//evap fan activity during defrost and drip delay; 0/1/2 = off/on/dependent on F0
	F2 *Int64 `protobuf:"bytes,71,opt,name=F2,proto3" json:"F2,omitempty"`
	//fan delay after evap drip completes
	//unit: minutes
	F3 *Double `protobuf:"bytes,72,opt,name=F3,proto3" json:"F3,omitempty"`
	//first real time defrost activation time; only if d8=3
	//unit: hh:mm
	Hd1 *Double `protobuf:"bytes,73,opt,name=Hd1,proto3" json:"Hd1,omitempty"`
	//second real time defrost activation time; only if d8=3
	//unit: hh:mm
	Hd2 *Double `protobuf:"bytes,74,opt,name=Hd2,proto3" json:"Hd2,omitempty"`
	//third real time defrost activation time; only if d8=3
	//unit: hh:mm
	Hd3 *Double `protobuf:"bytes,75,opt,name=Hd3,proto3" json:"Hd3,omitempty"`
	//fourth real time defrost activation time; only if d8=3
	//unit: hh:mm
	Hd4 *Double `protobuf:"bytes,76,opt,name=Hd4,proto3" json:"Hd4,omitempty"`
	//fifth real time defrost activation time; only if d8=3
	//unit: hh:mm
	Hd5 *Double `protobuf:"bytes,77,opt,name=Hd5,proto3" json:"Hd5,omitempty"`
	//sixth real time defrost activation time; only if d8=3
	//unit: hh:mm
	Hd6 *Double `protobuf:"bytes,78,opt,name=Hd6,proto3" json:"Hd6,omitempty"`
	// current UNIX epoch time
	//unit:ns
	Time                 uint64   `protobuf:"varint,79,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ParkerState) Reset()         { *m = ParkerState{} }
func (m *ParkerState) String() string { return proto.CompactTextString(m) }
func (*ParkerState) ProtoMessage()    {}
func (*ParkerState) Descriptor() ([]byte, []int) {
	return fileDescriptor_48357db32399db41, []int{0}
}

func (m *ParkerState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ParkerState.Unmarshal(m, b)
}
func (m *ParkerState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ParkerState.Marshal(b, m, deterministic)
}
func (m *ParkerState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ParkerState.Merge(m, src)
}
func (m *ParkerState) XXX_Size() int {
	return xxx_messageInfo_ParkerState.Size(m)
}
func (m *ParkerState) XXX_DiscardUnknown() {
	xxx_messageInfo_ParkerState.DiscardUnknown(m)
}

var xxx_messageInfo_ParkerState proto.InternalMessageInfo

func (m *ParkerState) GetCompressorWorkingHours() *Double {
	if m != nil {
		return m.CompressorWorkingHours
	}
	return nil
}

func (m *ParkerState) GetOnStandbyStatus() *Int64 {
	if m != nil {
		return m.OnStandbyStatus
	}
	return nil
}

func (m *ParkerState) GetLightStatus() *Int64 {
	if m != nil {
		return m.LightStatus
	}
	return nil
}

func (m *ParkerState) GetAuxOutputStatus() *Int64 {
	if m != nil {
		return m.AuxOutputStatus
	}
	return nil
}

func (m *ParkerState) GetNextDefrostCounter() *Double {
	if m != nil {
		return m.NextDefrostCounter
	}
	return nil
}

func (m *ParkerState) GetDoorSwitchInputStatus() *Int64 {
	if m != nil {
		return m.DoorSwitchInputStatus
	}
	return nil
}

func (m *ParkerState) GetMultipurposeInputStatus() *Int64 {
	if m != nil {
		return m.MultipurposeInputStatus
	}
	return nil
}

func (m *ParkerState) GetCompressorStatus() *Int64 {
	if m != nil {
		return m.CompressorStatus
	}
	return nil
}

func (m *ParkerState) GetOutputDefrostStatus() *Int64 {
	if m != nil {
		return m.OutputDefrostStatus
	}
	return nil
}

func (m *ParkerState) GetFansStatus() *Int64 {
	if m != nil {
		return m.FansStatus
	}
	return nil
}

func (m *ParkerState) GetOutputK4Status() *Int64 {
	if m != nil {
		return m.OutputK4Status
	}
	return nil
}

func (m *ParkerState) GetCabinetTemperature() *Double {
	if m != nil {
		return m.CabinetTemperature
	}
	return nil
}

func (m *ParkerState) GetEvaporatorTemperature() *Double {
	if m != nil {
		return m.EvaporatorTemperature
	}
	return nil
}

func (m *ParkerState) GetAuxiliaryTemperature() *Double {
	if m != nil {
		return m.AuxiliaryTemperature
	}
	return nil
}

func (m *ParkerState) GetProbe1FailureAlarm() *Int64 {
	if m != nil {
		return m.Probe1FailureAlarm
	}
	return nil
}

func (m *ParkerState) GetProbe2FailureAlarm() *Int64 {
	if m != nil {
		return m.Probe2FailureAlarm
	}
	return nil
}

func (m *ParkerState) GetProbe3FailureAlarm() *Int64 {
	if m != nil {
		return m.Probe3FailureAlarm
	}
	return nil
}

func (m *ParkerState) GetMinimumTemperatureAlarm() *Int64 {
	if m != nil {
		return m.MinimumTemperatureAlarm
	}
	return nil
}

func (m *ParkerState) GetMaximumTempertureAlarm() *Int64 {
	if m != nil {
		return m.MaximumTempertureAlarm
	}
	return nil
}

func (m *ParkerState) GetCondensorTemperatureFailureAlarm() *Int64 {
	if m != nil {
		return m.CondensorTemperatureFailureAlarm
	}
	return nil
}

func (m *ParkerState) GetCondensorPreAlarm() *Int64 {
	if m != nil {
		return m.CondensorPreAlarm
	}
	return nil
}

func (m *ParkerState) GetDoorAlarm() *Int64 {
	if m != nil {
		return m.DoorAlarm
	}
	return nil
}

func (m *ParkerState) GetMultipurposeInputAlarm() *Int64 {
	if m != nil {
		return m.MultipurposeInputAlarm
	}
	return nil
}

func (m *ParkerState) GetCompressorBlockedAlarm() *Int64 {
	if m != nil {
		return m.CompressorBlockedAlarm
	}
	return nil
}

func (m *ParkerState) GetPowerFailureAlarm() *Int64 {
	if m != nil {
		return m.PowerFailureAlarm
	}
	return nil
}

func (m *ParkerState) GetRtcErrorAlarm() *Int64 {
	if m != nil {
		return m.RtcErrorAlarm
	}
	return nil
}

func (m *ParkerState) GetEnergySavingRegulatorFlag() *Int64 {
	if m != nil {
		return m.EnergySavingRegulatorFlag
	}
	return nil
}

func (m *ParkerState) GetEnergySavingRealTimeRegulatorFlag() *Int64 {
	if m != nil {
		return m.EnergySavingRealTimeRegulatorFlag
	}
	return nil
}

func (m *ParkerState) GetServiceRequestRegulatorFlag() *Int64 {
	if m != nil {
		return m.ServiceRequestRegulatorFlag
	}
	return nil
}

func (m *ParkerState) GetOnStandbyRegulatorFlag() *Int64 {
	if m != nil {
		return m.OnStandbyRegulatorFlag
	}
	return nil
}

func (m *ParkerState) GetNewAlarmToReadRegulatorFlag() *Int64 {
	if m != nil {
		return m.NewAlarmToReadRegulatorFlag
	}
	return nil
}

func (m *ParkerState) GetDefrostStatusRegulatorFlag() *Int64 {
	if m != nil {
		return m.DefrostStatusRegulatorFlag
	}
	return nil
}

func (m *ParkerState) GetActiveSetpoint() *Int64 {
	if m != nil {
		return m.ActiveSetpoint
	}
	return nil
}

func (m *ParkerState) GetTimeUntilDefrost() *Int64 {
	if m != nil {
		return m.TimeUntilDefrost
	}
	return nil
}

func (m *ParkerState) GetCurrentDefrostCounter() *Int64 {
	if m != nil {
		return m.CurrentDefrostCounter
	}
	return nil
}

func (m *ParkerState) GetCompressorDelay() *Int64 {
	if m != nil {
		return m.CompressorDelay
	}
	return nil
}

func (m *ParkerState) GetNumAlarmsInHistory() *Int64 {
	if m != nil {
		return m.NumAlarmsInHistory
	}
	return nil
}

func (m *ParkerState) GetEnergySavingStatus() *Int64 {
	if m != nil {
		return m.EnergySavingStatus
	}
	return nil
}

func (m *ParkerState) GetServiceRequestStatus() *Int64 {
	if m != nil {
		return m.ServiceRequestStatus
	}
	return nil
}

func (m *ParkerState) GetResistorsActivatedByAuxKeyStatus() *Int64 {
	if m != nil {
		return m.ResistorsActivatedByAuxKeyStatus
	}
	return nil
}

func (m *ParkerState) GetEvaporatorValveState() *Int64 {
	if m != nil {
		return m.EvaporatorValveState
	}
	return nil
}

func (m *ParkerState) GetOutputDefrostState() *Int64 {
	if m != nil {
		return m.OutputDefrostState
	}
	return nil
}

func (m *ParkerState) GetOutputLuxState() *Int64 {
	if m != nil {
		return m.OutputLuxState
	}
	return nil
}

func (m *ParkerState) GetOutputAuxState() *Int64 {
	if m != nil {
		return m.OutputAuxState
	}
	return nil
}

func (m *ParkerState) GetResistorsState() *Int64 {
	if m != nil {
		return m.ResistorsState
	}
	return nil
}

func (m *ParkerState) GetOutputAlarmState() *Int64 {
	if m != nil {
		return m.OutputAlarmState
	}
	return nil
}

func (m *ParkerState) GetSecondCompressorState() *Int64 {
	if m != nil {
		return m.SecondCompressorState
	}
	return nil
}

func (m *ParkerState) GetSetpoint() *Double {
	if m != nil {
		return m.Setpoint
	}
	return nil
}

func (m *ParkerState) GetR1() *Double {
	if m != nil {
		return m.R1
	}
	return nil
}

func (m *ParkerState) GetR2() *Double {
	if m != nil {
		return m.R2
	}
	return nil
}

func (m *ParkerState) GetR4() *Double {
	if m != nil {
		return m.R4
	}
	return nil
}

func (m *ParkerState) GetC0() *Double {
	if m != nil {
		return m.C0
	}
	return nil
}

func (m *ParkerState) GetC1() *Double {
	if m != nil {
		return m.C1
	}
	return nil
}

func (m *ParkerState) GetD0() *Double {
	if m != nil {
		return m.D0
	}
	return nil
}

func (m *ParkerState) GetD3() *Double {
	if m != nil {
		return m.D3
	}
	return nil
}

func (m *ParkerState) GetD5() *Double {
	if m != nil {
		return m.D5
	}
	return nil
}

func (m *ParkerState) GetD7() *Double {
	if m != nil {
		return m.D7
	}
	return nil
}

func (m *ParkerState) GetD8() *Int64 {
	if m != nil {
		return m.D8
	}
	return nil
}

func (m *ParkerState) GetA0() *Int64 {
	if m != nil {
		return m.A0
	}
	return nil
}

func (m *ParkerState) GetA1() *Double {
	if m != nil {
		return m.A1
	}
	return nil
}

func (m *ParkerState) GetA2() *Int64 {
	if m != nil {
		return m.A2
	}
	return nil
}

func (m *ParkerState) GetA3() *Int64 {
	if m != nil {
		return m.A3
	}
	return nil
}

func (m *ParkerState) GetA4() *Double {
	if m != nil {
		return m.A4
	}
	return nil
}

func (m *ParkerState) GetA5() *Int64 {
	if m != nil {
		return m.A5
	}
	return nil
}

func (m *ParkerState) GetA6() *Double {
	if m != nil {
		return m.A6
	}
	return nil
}

func (m *ParkerState) GetA7() *Double {
	if m != nil {
		return m.A7
	}
	return nil
}

func (m *ParkerState) GetA8() *Double {
	if m != nil {
		return m.A8
	}
	return nil
}

func (m *ParkerState) GetA9() *Double {
	if m != nil {
		return m.A9
	}
	return nil
}

func (m *ParkerState) GetF0() *Int64 {
	if m != nil {
		return m.F0
	}
	return nil
}

func (m *ParkerState) GetF1() *Double {
	if m != nil {
		return m.F1
	}
	return nil
}

func (m *ParkerState) GetF2() *Int64 {
	if m != nil {
		return m.F2
	}
	return nil
}

func (m *ParkerState) GetF3() *Double {
	if m != nil {
		return m.F3
	}
	return nil
}

func (m *ParkerState) GetHd1() *Double {
	if m != nil {
		return m.Hd1
	}
	return nil
}

func (m *ParkerState) GetHd2() *Double {
	if m != nil {
		return m.Hd2
	}
	return nil
}

func (m *ParkerState) GetHd3() *Double {
	if m != nil {
		return m.Hd3
	}
	return nil
}

func (m *ParkerState) GetHd4() *Double {
	if m != nil {
		return m.Hd4
	}
	return nil
}

func (m *ParkerState) GetHd5() *Double {
	if m != nil {
		return m.Hd5
	}
	return nil
}

func (m *ParkerState) GetHd6() *Double {
	if m != nil {
		return m.Hd6
	}
	return nil
}

func (m *ParkerState) GetTime() uint64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func init() {
	proto.RegisterType((*ParkerState)(nil), "xbospb.ParkerState")
}

func init() { proto.RegisterFile("parker.proto", fileDescriptor_48357db32399db41) }

var fileDescriptor_48357db32399db41 = []byte{
	// 1323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x98, 0xdb, 0x73, 0xd4, 0x36,
	0x14, 0x87, 0x87, 0x40, 0x29, 0x28, 0x40, 0x88, 0x73, 0x53, 0xb8, 0x35, 0x40, 0xa1, 0x94, 0xd2,
	0x74, 0xef, 0x49, 0xa0, 0x34, 0x2c, 0x9b, 0x2c, 0x49, 0xa1, 0xc0, 0x6c, 0xd2, 0xf6, 0xa5, 0x53,
	0x57, 0xbb, 0x56, 0x12, 0x4f, 0xbc, 0x92, 0x2b, 0x4b, 0xc9, 0xee, 0x6b, 0xff, 0xf2, 0x8e, 0x2e,
	0xeb, 0xb5, 0x9d, 0x63, 0x9e, 0x60, 0xe6, 0x7c, 0xbf, 0xcf, 0xf2, 0x91, 0xd6, 0x47, 0x13, 0x74,
	0x23, 0x26, 0xe2, 0x94, 0x8a, 0xf5, 0x58, 0x70, 0xc9, 0xbd, 0xab, 0xa3, 0x3e, 0x4f, 0xe2, 0xfe,
	0x9d, 0x05, 0xa6, 0xa2, 0x88, 0xf4, 0x23, 0x2a, 0xc7, 0x31, 0x4d, 0x6c, 0xf1, 0xd1, 0x7f, 0x0f,
	0xd1, 0xec, 0x67, 0x43, 0x1f, 0x48, 0x22, 0xa9, 0xb7, 0x87, 0xf0, 0x80, 0x0f, 0x63, 0x41, 0x93,
	0x84, 0x0b, 0xff, 0x9c, 0x8b, 0xd3, 0x90, 0x1d, 0xfb, 0x27, 0x5c, 0x89, 0x04, 0x5f, 0x5a, 0xbb,
	0xf4, 0x6c, 0xb6, 0x76, 0x6b, 0xdd, 0xfa, 0xd6, 0x77, 0xb8, 0xea, 0x47, 0xb4, 0xb7, 0x3c, 0xe5,
	0xff, 0xb4, 0xf8, 0x9e, 0xa6, 0xbd, 0x2d, 0x34, 0xcf, 0x99, 0x9f, 0x48, 0xc2, 0x82, 0xfe, 0x58,
	0xff, 0x2b, 0x55, 0x82, 0x67, 0x8c, 0xe2, 0xe6, 0x44, 0xb1, 0xcf, 0x64, 0xab, 0xd1, 0x9b, 0xe3,
	0xec, 0xc0, 0x62, 0x07, 0x86, 0xf2, 0x2a, 0xe8, 0x46, 0x14, 0x1e, 0x9f, 0xc8, 0x49, 0xea, 0x32,
	0x94, 0x9a, 0x35, 0x88, 0x4b, 0x6c, 0xa1, 0x79, 0xa2, 0x46, 0x3e, 0x57, 0x32, 0x56, 0x69, 0xec,
	0x0a, 0xf8, 0x30, 0xa2, 0x46, 0x9f, 0x0c, 0xe6, 0xa2, 0x6f, 0xd0, 0x22, 0xa3, 0x23, 0xe9, 0x07,
	0xf4, 0x48, 0xf0, 0x44, 0xfa, 0x03, 0xae, 0x98, 0xa4, 0x02, 0x7f, 0x05, 0xbe, 0xad, 0xa7, 0xd9,
	0x1d, 0x8b, 0x76, 0x2c, 0xe9, 0x75, 0x11, 0x0e, 0x38, 0x17, 0x7e, 0x72, 0x1e, 0xca, 0xc1, 0x89,
	0x1f, 0xb2, 0xcc, 0x1a, 0xae, 0x42, 0x6b, 0x58, 0xd2, 0xf8, 0x81, 0xa1, 0xf7, 0xd9, 0x74, 0x25,
	0xfb, 0x68, 0x75, 0xa8, 0x22, 0x19, 0xc6, 0x4a, 0xc4, 0x3c, 0xa1, 0x79, 0xd1, 0xd7, 0x90, 0x68,
	0x25, 0xcb, 0x67, 0x55, 0x2f, 0xd1, 0x7c, 0x66, 0x1b, 0x9d, 0xe2, 0x1a, 0xa4, 0xb8, 0x3d, 0xe5,
	0x5c, 0xb6, 0x8d, 0x96, 0x5c, 0x1f, 0x27, 0x2d, 0x71, 0xf9, 0xeb, 0x50, 0x7e, 0xc1, 0xb2, 0xae,
	0x25, 0x4e, 0xb1, 0x8e, 0x66, 0x8f, 0x08, 0x4b, 0x26, 0x41, 0x04, 0x05, 0x91, 0x26, 0x1c, 0xbf,
	0x81, 0x6e, 0xbb, 0x47, 0x9e, 0x36, 0x26, 0xa1, 0x59, 0x28, 0x74, 0xcb, 0x62, 0xef, 0x1b, 0x2e,
	0xb8, 0x8d, 0x16, 0x06, 0xa4, 0x1f, 0x32, 0x2a, 0x7d, 0x49, 0x87, 0x31, 0x15, 0x44, 0x2a, 0x41,
	0xf1, 0x0d, 0x78, 0xef, 0x1c, 0x7a, 0x38, 0x25, 0xbd, 0x5d, 0xb4, 0x4c, 0xcf, 0x48, 0xcc, 0x05,
	0x91, 0x5c, 0xe4, 0x1c, 0x37, 0x41, 0xc7, 0xd2, 0x94, 0xce, 0x6a, 0x3a, 0x68, 0x89, 0xa8, 0x51,
	0x18, 0x85, 0x44, 0x8c, 0x73, 0x96, 0x5b, 0xa0, 0x65, 0x31, 0x85, 0xb3, 0x92, 0x6d, 0xb4, 0x18,
	0x0b, 0xde, 0xa7, 0x55, 0xff, 0x88, 0x84, 0x91, 0x12, 0xd4, 0x27, 0x11, 0x11, 0x43, 0x3c, 0x07,
	0x75, 0xc2, 0xb3, 0x68, 0xd7, 0x92, 0x6d, 0x0d, 0xa6, 0x82, 0x5a, 0x41, 0x70, 0xbb, 0x5c, 0x50,
	0x03, 0x05, 0xf5, 0x82, 0x60, 0xbe, 0x5c, 0x50, 0xcf, 0x09, 0xf4, 0x11, 0x0e, 0x59, 0x38, 0x54,
	0xc3, 0x6c, 0x17, 0x9c, 0xc5, 0x83, 0x8f, 0xb0, 0xe5, 0x33, 0x8d, 0xb0, 0xaa, 0x77, 0x08, 0x0f,
	0xc9, 0x28, 0xa3, 0xca, 0x98, 0x16, 0x20, 0xd3, 0xb2, 0xc3, 0x0f, 0x53, 0xda, 0x8a, 0xfe, 0x42,
	0x8f, 0x07, 0x9c, 0x05, 0x94, 0x25, 0xf9, 0x1d, 0x2e, 0xbc, 0xe3, 0x22, 0xe4, 0x5c, 0x4b, 0x93,
	0x99, 0xf5, 0xe5, 0xde, 0xf8, 0x35, 0x5a, 0x98, 0xda, 0xe3, 0xd4, 0xb6, 0x04, 0xd9, 0xe6, 0x53,
	0xf2, 0xf3, 0x24, 0xfe, 0x02, 0x21, 0xf3, 0xed, 0xb0, 0xa9, 0x65, 0x28, 0x75, 0x5d, 0x03, 0xd3,
	0x9e, 0x5c, 0xfc, 0x42, 0xd8, 0xec, 0x0a, 0xdc, 0x93, 0xe2, 0x07, 0x22, 0x15, 0x65, 0xbe, 0x0f,
	0xfd, 0x88, 0x0f, 0x4e, 0x69, 0xe0, 0x44, 0x18, 0x14, 0x4d, 0xf1, 0xb7, 0x96, 0x4e, 0x5f, 0x3f,
	0xe6, 0xe7, 0x54, 0x14, 0x9a, 0xb9, 0x0a, 0xbe, 0xbe, 0x21, 0x73, 0xdd, 0x6b, 0xa2, 0x39, 0x21,
	0x07, 0x3e, 0x15, 0x22, 0xed, 0xc1, 0x1d, 0x28, 0x7a, 0x53, 0xc8, 0xc1, 0xae, 0x86, 0x6c, 0xec,
	0x23, 0xba, 0x47, 0x19, 0x15, 0xc7, 0x63, 0x3f, 0x21, 0x67, 0x7a, 0x40, 0x09, 0x7a, 0xac, 0x22,
	0xf3, 0x13, 0x3e, 0x8a, 0xc8, 0x31, 0xbe, 0x0b, 0x39, 0x56, 0x6d, 0xe4, 0xc0, 0x24, 0x7a, 0x93,
	0x40, 0x37, 0x22, 0xc7, 0xde, 0x3f, 0xe8, 0x69, 0xd1, 0x47, 0x22, 0x5f, 0x86, 0x43, 0x5a, 0x34,
	0xdf, 0x83, 0xcc, 0x0f, 0xf3, 0x66, 0x12, 0x1d, 0x86, 0x43, 0x9a, 0x7f, 0x42, 0x0f, 0x3d, 0x48,
	0xa8, 0x38, 0x0b, 0x07, 0xda, 0xf8, 0xaf, 0xa2, 0x89, 0x2c, 0x9a, 0xef, 0x43, 0xe6, 0xbb, 0x2e,
	0xd4, 0xb3, 0x99, 0xbc, 0x73, 0x0f, 0xad, 0x66, 0x26, 0x6c, 0x41, 0xf7, 0x00, 0xdc, 0xc5, 0x74,
	0xd2, 0xe6, 0x4d, 0x87, 0x68, 0x8d, 0xd1, 0x73, 0xbb, 0x01, 0xbe, 0xe4, 0xfa, 0xf5, 0x83, 0xa2,
	0xf0, 0x1b, 0x70, 0x7d, 0x8c, 0x9e, 0x9b, 0x2d, 0x39, 0xe4, 0x3d, 0x4a, 0x82, 0xbc, 0xf5, 0x33,
	0xba, 0x9f, 0x9f, 0x20, 0x45, 0xe5, 0x1a, 0xa4, 0xbc, 0x13, 0x64, 0x47, 0x49, 0xde, 0xd8, 0x42,
	0x73, 0x64, 0x20, 0xc3, 0x33, 0xea, 0x27, 0x54, 0xc6, 0x3c, 0x64, 0x12, 0x3f, 0x04, 0xc7, 0x84,
	0xa5, 0x0e, 0x1c, 0xe4, 0xbd, 0x42, 0x9e, 0xd9, 0x4c, 0xc5, 0x64, 0x18, 0x4d, 0xc6, 0x1a, 0x7e,
	0x04, 0xce, 0x43, 0x0d, 0xfe, 0xae, 0x39, 0x37, 0xd2, 0xbc, 0x5d, 0xb4, 0x32, 0x50, 0x42, 0x50,
	0x76, 0xf1, 0x8e, 0xf0, 0x18, 0x9c, 0xee, 0x8e, 0x2e, 0xdc, 0x12, 0x36, 0x51, 0x66, 0xd4, 0xfa,
	0x01, 0x8d, 0xc8, 0x18, 0x7f, 0x0b, 0xde, 0x50, 0xa6, 0xd8, 0x8e, 0xa6, 0xbc, 0x37, 0x68, 0x89,
	0xa9, 0xa1, 0xdd, 0x9d, 0xc4, 0x0f, 0x99, 0x7f, 0x12, 0x26, 0x92, 0x8b, 0x31, 0x7e, 0x02, 0x7e,
	0x96, 0x99, 0x1a, 0x9a, 0x2d, 0x49, 0xf6, 0xd9, 0x9e, 0x05, 0xf5, 0x77, 0x3d, 0x7f, 0xbe, 0xdd,
	0x8c, 0x7d, 0x0a, 0x0a, 0xb2, 0xa7, 0xd9, 0xcd, 0xd9, 0x0e, 0x5a, 0x2e, 0x1e, 0x5f, 0xa7, 0xf8,
	0x0e, 0x52, 0x2c, 0xe6, 0x8f, 0xad, 0x93, 0xfc, 0x8d, 0x9e, 0x08, 0x9a, 0x98, 0x25, 0x25, 0xbe,
	0xd9, 0x21, 0x22, 0x69, 0xe0, 0xf7, 0xc7, 0xbe, 0xbe, 0xbb, 0x9d, 0xd2, 0xf4, 0x96, 0xf8, 0x0c,
	0xfc, 0x14, 0xa7, 0xd9, 0xf6, 0x24, 0xfa, 0x76, 0xdc, 0x56, 0xa3, 0xf7, 0x74, 0x3c, 0x5d, 0x64,
	0x66, 0x96, 0x9f, 0x91, 0x48, 0x9f, 0x13, 0x7d, 0xab, 0xc5, 0xdf, 0x83, 0x8b, 0x9c, 0xc2, 0x7f,
	0x68, 0xd6, 0x5e, 0x80, 0xb7, 0xd1, 0x22, 0x70, 0xfb, 0xa1, 0xf8, 0x39, 0xd8, 0xaa, 0x0b, 0x97,
	0x1f, 0x9a, 0xb9, 0xcb, 0x44, 0x6a, 0xe4, 0xc2, 0x3f, 0x7c, 0xe1, 0x2e, 0xf3, 0x41, 0x8d, 0x8a,
	0x41, 0x92, 0x06, 0x5f, 0x7c, 0x21, 0xd8, 0x9e, 0x04, 0x5b, 0x68, 0x6e, 0xda, 0x57, 0x9b, 0xfb,
	0x11, 0xcc, 0xa5, 0x94, 0xcd, 0xbd, 0x42, 0xde, 0xe4, 0x81, 0xe6, 0x87, 0x6f, 0xa3, 0xeb, 0xe0,
	0xaf, 0xc2, 0x3d, 0x52, 0x73, 0x36, 0xbc, 0x8b, 0x56, 0x12, 0xaa, 0xe7, 0x99, 0x5f, 0xb8, 0x68,
	0x52, 0xfc, 0x13, 0xf8, 0xab, 0xb0, 0x74, 0x27, 0x77, 0xdb, 0xa4, 0xde, 0x73, 0x74, 0x2d, 0xfd,
	0x29, 0x57, 0xc0, 0xbb, 0x52, 0x5a, 0xf7, 0x1e, 0xa0, 0x19, 0x51, 0xc5, 0x55, 0x90, 0x9a, 0x11,
	0x55, 0x53, 0xaf, 0xe1, 0x5a, 0x49, 0xbd, 0x66, 0xea, 0x0d, 0x5c, 0x2f, 0xa9, 0x37, 0x74, 0xbd,
	0x53, 0xc1, 0x0d, 0xb8, 0xde, 0xa9, 0x98, 0x7a, 0x15, 0x37, 0x4b, 0xea, 0xe6, 0xf9, 0x41, 0x05,
	0xb7, 0xe0, 0x7a, 0x60, 0xf2, 0x41, 0x1d, 0x6f, 0x94, 0xd4, 0xeb, 0xa6, 0xde, 0xc4, 0x9b, 0x25,
	0xf5, 0xa6, 0xa9, 0x6f, 0xe0, 0xad, 0x92, 0xfa, 0x86, 0x77, 0x1f, 0xcd, 0x04, 0x9b, 0xf8, 0x25,
	0xd4, 0xfd, 0x99, 0x60, 0x53, 0x97, 0xdb, 0x15, 0xfc, 0x0a, 0x2c, 0xb7, 0xcd, 0xea, 0xda, 0x55,
	0xfc, 0x33, 0x6c, 0x6f, 0x57, 0x4d, 0xbc, 0x86, 0x5f, 0xc3, 0xf1, 0x9a, 0x29, 0xd7, 0xf1, 0x2f,
	0x70, 0xd9, 0xbc, 0x5b, 0xbb, 0x81, 0xb7, 0x4b, 0xec, 0x0d, 0x13, 0x6f, 0xe2, 0x37, 0x70, 0xdc,
	0xbc, 0x7a, 0xbb, 0x85, 0xdb, 0x25, 0xf1, 0x96, 0xa9, 0x6f, 0xe0, 0xb7, 0x25, 0xf5, 0x0d, 0x53,
	0xdf, 0xc4, 0x9d, 0x92, 0xfa, 0xa6, 0xa9, 0x6f, 0xe1, 0x9d, 0x92, 0xfa, 0x96, 0x5e, 0x5e, 0xb7,
	0x82, 0x77, 0xc1, 0xe5, 0x75, 0x4d, 0xef, 0xba, 0x55, 0xdc, 0x85, 0xe3, 0x5d, 0xd3, 0xbb, 0x6e,
	0x0d, 0xbf, 0x83, 0xe3, 0xe6, 0x60, 0x76, 0xeb, 0x78, 0xaf, 0x24, 0x5e, 0xf7, 0xd6, 0xd0, 0xe5,
	0xbd, 0xa0, 0x8a, 0xf7, 0x41, 0x40, 0x97, 0x2c, 0x51, 0xc3, 0xbf, 0x96, 0x11, 0x35, 0x4b, 0xd4,
	0xf1, 0xfb, 0x32, 0xc2, 0x3d, 0xa5, 0x81, 0x3f, 0x94, 0x11, 0x0d, 0x4b, 0x34, 0xf1, 0x6f, 0x65,
	0x44, 0xd3, 0x12, 0x2d, 0xfc, 0xb1, 0x8c, 0x68, 0x79, 0x1e, 0xba, 0xa2, 0x27, 0x2c, 0xfe, 0xb4,
	0x76, 0xe9, 0xd9, 0x95, 0x9e, 0xf9, 0x7f, 0xff, 0xaa, 0xf9, 0x5b, 0x44, 0xfd, 0xff, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xd5, 0x88, 0x5e, 0xa0, 0xb8, 0x10, 0x00, 0x00,
}