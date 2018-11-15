// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hamilton.proto

package xbospb

/*
This is designed to be included by the main xbos proto file and includes the
definitions for the Hamilton project

Version 1.0
*/

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Data emitted from Hamilton Sensors
// Maintainer: Michael Andersen
type HamiltonData struct {
	Serial               uint32      `protobuf:"varint,1,opt,name=serial,proto3" json:"serial,omitempty"`
	Model                string      `protobuf:"bytes,2,opt,name=model,proto3" json:"model,omitempty"`
	Time                 uint64      `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	H3C                  *Hamilton3C `protobuf:"bytes,4,opt,name=h3c,proto3" json:"h3c,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *HamiltonData) Reset()         { *m = HamiltonData{} }
func (m *HamiltonData) String() string { return proto.CompactTextString(m) }
func (*HamiltonData) ProtoMessage()    {}
func (*HamiltonData) Descriptor() ([]byte, []int) {
	return fileDescriptor_hamilton_b14faa22e79b32e8, []int{0}
}
func (m *HamiltonData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HamiltonData.Unmarshal(m, b)
}
func (m *HamiltonData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HamiltonData.Marshal(b, m, deterministic)
}
func (dst *HamiltonData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HamiltonData.Merge(dst, src)
}
func (m *HamiltonData) XXX_Size() int {
	return xxx_messageInfo_HamiltonData.Size(m)
}
func (m *HamiltonData) XXX_DiscardUnknown() {
	xxx_messageInfo_HamiltonData.DiscardUnknown(m)
}

var xxx_messageInfo_HamiltonData proto.InternalMessageInfo

func (m *HamiltonData) GetSerial() uint32 {
	if m != nil {
		return m.Serial
	}
	return 0
}

func (m *HamiltonData) GetModel() string {
	if m != nil {
		return m.Model
	}
	return ""
}

func (m *HamiltonData) GetTime() uint64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *HamiltonData) GetH3C() *Hamilton3C {
	if m != nil {
		return m.H3C
	}
	return nil
}

// Data specific to a Hamilton 3C/7C sensor
// Maintainer: Michael Andersen
type Hamilton3C struct {
	// unit:seconds
	Uptime     uint64  `protobuf:"varint,1,opt,name=uptime,proto3" json:"uptime,omitempty"`
	Flags      uint32  `protobuf:"varint,2,opt,name=flags,proto3" json:"flags,omitempty"`
	AccX       float64 `protobuf:"fixed64,3,opt,name=acc_x,json=accX,proto3" json:"acc_x,omitempty"`
	AccY       float64 `protobuf:"fixed64,4,opt,name=acc_y,json=accY,proto3" json:"acc_y,omitempty"`
	AccZ       float64 `protobuf:"fixed64,5,opt,name=acc_z,json=accZ,proto3" json:"acc_z,omitempty"`
	MagX       float64 `protobuf:"fixed64,6,opt,name=mag_x,json=magX,proto3" json:"mag_x,omitempty"`
	MagY       float64 `protobuf:"fixed64,7,opt,name=mag_y,json=magY,proto3" json:"mag_y,omitempty"`
	MagZ       float64 `protobuf:"fixed64,8,opt,name=mag_z,json=magZ,proto3" json:"mag_z,omitempty"`
	TmpDie     float64 `protobuf:"fixed64,9,opt,name=tmp_die,json=tmpDie,proto3" json:"tmp_die,omitempty"`
	TmpVoltage float64 `protobuf:"fixed64,10,opt,name=tmp_voltage,json=tmpVoltage,proto3" json:"tmp_voltage,omitempty"`
	// unit:celsius
	AirTemp float64 `protobuf:"fixed64,11,opt,name=air_temp,json=airTemp,proto3" json:"air_temp,omitempty"`
	// unit:humidity
	AirHum float64 `protobuf:"fixed64,12,opt,name=air_hum,json=airHum,proto3" json:"air_hum,omitempty"`
	// unit:%rh
	AirRh float64 `protobuf:"fixed64,13,opt,name=air_rh,json=airRh,proto3" json:"air_rh,omitempty"`
	// unit:lux
	LightLux float64 `protobuf:"fixed64,14,opt,name=light_lux,json=lightLux,proto3" json:"light_lux,omitempty"`
	// unit:# pushes
	Buttons uint32 `protobuf:"varint,15,opt,name=buttons,proto3" json:"buttons,omitempty"`
	// unit:% occupied
	Occupancy            float64  `protobuf:"fixed64,16,opt,name=occupancy,proto3" json:"occupancy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Hamilton3C) Reset()         { *m = Hamilton3C{} }
func (m *Hamilton3C) String() string { return proto.CompactTextString(m) }
func (*Hamilton3C) ProtoMessage()    {}
func (*Hamilton3C) Descriptor() ([]byte, []int) {
	return fileDescriptor_hamilton_b14faa22e79b32e8, []int{1}
}
func (m *Hamilton3C) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hamilton3C.Unmarshal(m, b)
}
func (m *Hamilton3C) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hamilton3C.Marshal(b, m, deterministic)
}
func (dst *Hamilton3C) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hamilton3C.Merge(dst, src)
}
func (m *Hamilton3C) XXX_Size() int {
	return xxx_messageInfo_Hamilton3C.Size(m)
}
func (m *Hamilton3C) XXX_DiscardUnknown() {
	xxx_messageInfo_Hamilton3C.DiscardUnknown(m)
}

var xxx_messageInfo_Hamilton3C proto.InternalMessageInfo

func (m *Hamilton3C) GetUptime() uint64 {
	if m != nil {
		return m.Uptime
	}
	return 0
}

func (m *Hamilton3C) GetFlags() uint32 {
	if m != nil {
		return m.Flags
	}
	return 0
}

func (m *Hamilton3C) GetAccX() float64 {
	if m != nil {
		return m.AccX
	}
	return 0
}

func (m *Hamilton3C) GetAccY() float64 {
	if m != nil {
		return m.AccY
	}
	return 0
}

func (m *Hamilton3C) GetAccZ() float64 {
	if m != nil {
		return m.AccZ
	}
	return 0
}

func (m *Hamilton3C) GetMagX() float64 {
	if m != nil {
		return m.MagX
	}
	return 0
}

func (m *Hamilton3C) GetMagY() float64 {
	if m != nil {
		return m.MagY
	}
	return 0
}

func (m *Hamilton3C) GetMagZ() float64 {
	if m != nil {
		return m.MagZ
	}
	return 0
}

func (m *Hamilton3C) GetTmpDie() float64 {
	if m != nil {
		return m.TmpDie
	}
	return 0
}

func (m *Hamilton3C) GetTmpVoltage() float64 {
	if m != nil {
		return m.TmpVoltage
	}
	return 0
}

func (m *Hamilton3C) GetAirTemp() float64 {
	if m != nil {
		return m.AirTemp
	}
	return 0
}

func (m *Hamilton3C) GetAirHum() float64 {
	if m != nil {
		return m.AirHum
	}
	return 0
}

func (m *Hamilton3C) GetAirRh() float64 {
	if m != nil {
		return m.AirRh
	}
	return 0
}

func (m *Hamilton3C) GetLightLux() float64 {
	if m != nil {
		return m.LightLux
	}
	return 0
}

func (m *Hamilton3C) GetButtons() uint32 {
	if m != nil {
		return m.Buttons
	}
	return 0
}

func (m *Hamilton3C) GetOccupancy() float64 {
	if m != nil {
		return m.Occupancy
	}
	return 0
}

// Data specific to a Hamilton 330/370 sensor
// Maintainer: Michael Andersen
type Hamilton330 struct {
	Uptime     uint64  `protobuf:"varint,1,opt,name=uptime,proto3" json:"uptime,omitempty"`
	Flags      uint32  `protobuf:"varint,2,opt,name=flags,proto3" json:"flags,omitempty"`
	AccX       float64 `protobuf:"fixed64,3,opt,name=acc_x,json=accX,proto3" json:"acc_x,omitempty"`
	AccY       float64 `protobuf:"fixed64,4,opt,name=acc_y,json=accY,proto3" json:"acc_y,omitempty"`
	AccZ       float64 `protobuf:"fixed64,5,opt,name=acc_z,json=accZ,proto3" json:"acc_z,omitempty"`
	MagX       float64 `protobuf:"fixed64,6,opt,name=mag_x,json=magX,proto3" json:"mag_x,omitempty"`
	MagY       float64 `protobuf:"fixed64,7,opt,name=mag_y,json=magY,proto3" json:"mag_y,omitempty"`
	MagZ       float64 `protobuf:"fixed64,8,opt,name=mag_z,json=magZ,proto3" json:"mag_z,omitempty"`
	TmpDie     float64 `protobuf:"fixed64,9,opt,name=tmp_die,json=tmpDie,proto3" json:"tmp_die,omitempty"`
	TmpVoltage float64 `protobuf:"fixed64,10,opt,name=tmp_voltage,json=tmpVoltage,proto3" json:"tmp_voltage,omitempty"`
	// unit:celsius
	AirTemp float64 `protobuf:"fixed64,11,opt,name=air_temp,json=airTemp,proto3" json:"air_temp,omitempty"`
	// unit:humidity
	AirHum float64 `protobuf:"fixed64,12,opt,name=air_hum,json=airHum,proto3" json:"air_hum,omitempty"`
	// unit:%rh
	AirRh float64 `protobuf:"fixed64,13,opt,name=air_rh,json=airRh,proto3" json:"air_rh,omitempty"`
	// unit:lux
	LightLux float64 `protobuf:"fixed64,14,opt,name=light_lux,json=lightLux,proto3" json:"light_lux,omitempty"`
	// unit:# pushes
	Buttons uint32 `protobuf:"varint,15,opt,name=buttons,proto3" json:"buttons,omitempty"`
	// unit:% occupied
	Occupancy            float64  `protobuf:"fixed64,16,opt,name=occupancy,proto3" json:"occupancy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Hamilton330) Reset()         { *m = Hamilton330{} }
func (m *Hamilton330) String() string { return proto.CompactTextString(m) }
func (*Hamilton330) ProtoMessage()    {}
func (*Hamilton330) Descriptor() ([]byte, []int) {
	return fileDescriptor_hamilton_b14faa22e79b32e8, []int{2}
}
func (m *Hamilton330) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hamilton330.Unmarshal(m, b)
}
func (m *Hamilton330) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hamilton330.Marshal(b, m, deterministic)
}
func (dst *Hamilton330) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hamilton330.Merge(dst, src)
}
func (m *Hamilton330) XXX_Size() int {
	return xxx_messageInfo_Hamilton330.Size(m)
}
func (m *Hamilton330) XXX_DiscardUnknown() {
	xxx_messageInfo_Hamilton330.DiscardUnknown(m)
}

var xxx_messageInfo_Hamilton330 proto.InternalMessageInfo

func (m *Hamilton330) GetUptime() uint64 {
	if m != nil {
		return m.Uptime
	}
	return 0
}

func (m *Hamilton330) GetFlags() uint32 {
	if m != nil {
		return m.Flags
	}
	return 0
}

func (m *Hamilton330) GetAccX() float64 {
	if m != nil {
		return m.AccX
	}
	return 0
}

func (m *Hamilton330) GetAccY() float64 {
	if m != nil {
		return m.AccY
	}
	return 0
}

func (m *Hamilton330) GetAccZ() float64 {
	if m != nil {
		return m.AccZ
	}
	return 0
}

func (m *Hamilton330) GetMagX() float64 {
	if m != nil {
		return m.MagX
	}
	return 0
}

func (m *Hamilton330) GetMagY() float64 {
	if m != nil {
		return m.MagY
	}
	return 0
}

func (m *Hamilton330) GetMagZ() float64 {
	if m != nil {
		return m.MagZ
	}
	return 0
}

func (m *Hamilton330) GetTmpDie() float64 {
	if m != nil {
		return m.TmpDie
	}
	return 0
}

func (m *Hamilton330) GetTmpVoltage() float64 {
	if m != nil {
		return m.TmpVoltage
	}
	return 0
}

func (m *Hamilton330) GetAirTemp() float64 {
	if m != nil {
		return m.AirTemp
	}
	return 0
}

func (m *Hamilton330) GetAirHum() float64 {
	if m != nil {
		return m.AirHum
	}
	return 0
}

func (m *Hamilton330) GetAirRh() float64 {
	if m != nil {
		return m.AirRh
	}
	return 0
}

func (m *Hamilton330) GetLightLux() float64 {
	if m != nil {
		return m.LightLux
	}
	return 0
}

func (m *Hamilton330) GetButtons() uint32 {
	if m != nil {
		return m.Buttons
	}
	return 0
}

func (m *Hamilton330) GetOccupancy() float64 {
	if m != nil {
		return m.Occupancy
	}
	return 0
}

// Published by Hamilton Border routers periodically
type HamiltonBRLinkStats struct {
	BadFrames            uint64   `protobuf:"varint,1,opt,name=BadFrames,json=badFrames,proto3" json:"BadFrames,omitempty"`
	LostFrames           uint64   `protobuf:"varint,2,opt,name=LostFrames,json=lostFrames,proto3" json:"LostFrames,omitempty"`
	DropNotConnected     uint64   `protobuf:"varint,3,opt,name=DropNotConnected,json=dropNotConnected,proto3" json:"DropNotConnected,omitempty"`
	SumSerialReceived    uint64   `protobuf:"varint,4,opt,name=SumSerialReceived,json=sumSerialReceived,proto3" json:"SumSerialReceived,omitempty"`
	SumDomainForwarded   uint64   `protobuf:"varint,5,opt,name=SumDomainForwarded,json=sumDomainForwarded,proto3" json:"SumDomainForwarded,omitempty"`
	SumDropNotConnected  uint64   `protobuf:"varint,6,opt,name=SumDropNotConnected,json=sumDropNotConnected,proto3" json:"SumDropNotConnected,omitempty"`
	SumDomainReceived    uint64   `protobuf:"varint,7,opt,name=SumDomainReceived,json=sumDomainReceived,proto3" json:"SumDomainReceived,omitempty"`
	SumSerialForwarded   uint64   `protobuf:"varint,8,opt,name=SumSerialForwarded,json=sumSerialForwarded,proto3" json:"SumSerialForwarded,omitempty"`
	PublishOkay          uint64   `protobuf:"varint,9,opt,name=PublishOkay,json=publishOkay,proto3" json:"PublishOkay,omitempty"`
	PublishError         uint64   `protobuf:"varint,10,opt,name=PublishError,json=publishError,proto3" json:"PublishError,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HamiltonBRLinkStats) Reset()         { *m = HamiltonBRLinkStats{} }
func (m *HamiltonBRLinkStats) String() string { return proto.CompactTextString(m) }
func (*HamiltonBRLinkStats) ProtoMessage()    {}
func (*HamiltonBRLinkStats) Descriptor() ([]byte, []int) {
	return fileDescriptor_hamilton_b14faa22e79b32e8, []int{3}
}
func (m *HamiltonBRLinkStats) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HamiltonBRLinkStats.Unmarshal(m, b)
}
func (m *HamiltonBRLinkStats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HamiltonBRLinkStats.Marshal(b, m, deterministic)
}
func (dst *HamiltonBRLinkStats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HamiltonBRLinkStats.Merge(dst, src)
}
func (m *HamiltonBRLinkStats) XXX_Size() int {
	return xxx_messageInfo_HamiltonBRLinkStats.Size(m)
}
func (m *HamiltonBRLinkStats) XXX_DiscardUnknown() {
	xxx_messageInfo_HamiltonBRLinkStats.DiscardUnknown(m)
}

var xxx_messageInfo_HamiltonBRLinkStats proto.InternalMessageInfo

func (m *HamiltonBRLinkStats) GetBadFrames() uint64 {
	if m != nil {
		return m.BadFrames
	}
	return 0
}

func (m *HamiltonBRLinkStats) GetLostFrames() uint64 {
	if m != nil {
		return m.LostFrames
	}
	return 0
}

func (m *HamiltonBRLinkStats) GetDropNotConnected() uint64 {
	if m != nil {
		return m.DropNotConnected
	}
	return 0
}

func (m *HamiltonBRLinkStats) GetSumSerialReceived() uint64 {
	if m != nil {
		return m.SumSerialReceived
	}
	return 0
}

func (m *HamiltonBRLinkStats) GetSumDomainForwarded() uint64 {
	if m != nil {
		return m.SumDomainForwarded
	}
	return 0
}

func (m *HamiltonBRLinkStats) GetSumDropNotConnected() uint64 {
	if m != nil {
		return m.SumDropNotConnected
	}
	return 0
}

func (m *HamiltonBRLinkStats) GetSumDomainReceived() uint64 {
	if m != nil {
		return m.SumDomainReceived
	}
	return 0
}

func (m *HamiltonBRLinkStats) GetSumSerialForwarded() uint64 {
	if m != nil {
		return m.SumSerialForwarded
	}
	return 0
}

func (m *HamiltonBRLinkStats) GetPublishOkay() uint64 {
	if m != nil {
		return m.PublishOkay
	}
	return 0
}

func (m *HamiltonBRLinkStats) GetPublishError() uint64 {
	if m != nil {
		return m.PublishError
	}
	return 0
}

// Published by Hamilton Border routers for each message
type HamiltonBRMessage struct {
	SrcMAC               string   `protobuf:"bytes,1,opt,name=SrcMAC,json=srcMAC,proto3" json:"SrcMAC,omitempty"`
	SrcIP                string   `protobuf:"bytes,2,opt,name=SrcIP,json=srcIP,proto3" json:"SrcIP,omitempty"`
	PopID                string   `protobuf:"bytes,3,opt,name=PopID,json=popID,proto3" json:"PopID,omitempty"`
	PopTime              int64    `protobuf:"varint,4,opt,name=PopTime,json=popTime,proto3" json:"PopTime,omitempty"`
	BRTime               int64    `protobuf:"varint,5,opt,name=BRTime,json=bRTime,proto3" json:"BRTime,omitempty"`
	RSSI                 int32    `protobuf:"varint,6,opt,name=RSSI,json=rSSI,proto3" json:"RSSI,omitempty"`
	LQI                  int32    `protobuf:"varint,7,opt,name=LQI,json=lQI,proto3" json:"LQI,omitempty"`
	Payload              []byte   `protobuf:"bytes,8,opt,name=Payload,json=payload,proto3" json:"Payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HamiltonBRMessage) Reset()         { *m = HamiltonBRMessage{} }
func (m *HamiltonBRMessage) String() string { return proto.CompactTextString(m) }
func (*HamiltonBRMessage) ProtoMessage()    {}
func (*HamiltonBRMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_hamilton_b14faa22e79b32e8, []int{4}
}
func (m *HamiltonBRMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HamiltonBRMessage.Unmarshal(m, b)
}
func (m *HamiltonBRMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HamiltonBRMessage.Marshal(b, m, deterministic)
}
func (dst *HamiltonBRMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HamiltonBRMessage.Merge(dst, src)
}
func (m *HamiltonBRMessage) XXX_Size() int {
	return xxx_messageInfo_HamiltonBRMessage.Size(m)
}
func (m *HamiltonBRMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_HamiltonBRMessage.DiscardUnknown(m)
}

var xxx_messageInfo_HamiltonBRMessage proto.InternalMessageInfo

func (m *HamiltonBRMessage) GetSrcMAC() string {
	if m != nil {
		return m.SrcMAC
	}
	return ""
}

func (m *HamiltonBRMessage) GetSrcIP() string {
	if m != nil {
		return m.SrcIP
	}
	return ""
}

func (m *HamiltonBRMessage) GetPopID() string {
	if m != nil {
		return m.PopID
	}
	return ""
}

func (m *HamiltonBRMessage) GetPopTime() int64 {
	if m != nil {
		return m.PopTime
	}
	return 0
}

func (m *HamiltonBRMessage) GetBRTime() int64 {
	if m != nil {
		return m.BRTime
	}
	return 0
}

func (m *HamiltonBRMessage) GetRSSI() int32 {
	if m != nil {
		return m.RSSI
	}
	return 0
}

func (m *HamiltonBRMessage) GetLQI() int32 {
	if m != nil {
		return m.LQI
	}
	return 0
}

func (m *HamiltonBRMessage) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func init() {
	proto.RegisterType((*HamiltonData)(nil), "xbospb.HamiltonData")
	proto.RegisterType((*Hamilton3C)(nil), "xbospb.Hamilton3C")
	proto.RegisterType((*Hamilton330)(nil), "xbospb.Hamilton330")
	proto.RegisterType((*HamiltonBRLinkStats)(nil), "xbospb.HamiltonBRLinkStats")
	proto.RegisterType((*HamiltonBRMessage)(nil), "xbospb.HamiltonBRMessage")
}

func init() { proto.RegisterFile("hamilton.proto", fileDescriptor_hamilton_b14faa22e79b32e8) }

var fileDescriptor_hamilton_b14faa22e79b32e8 = []byte{
	// 681 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x95, 0xdf, 0x8e, 0x1b, 0x35,
	0x14, 0xc6, 0x35, 0x9d, 0x3f, 0x49, 0x4e, 0xb2, 0x25, 0xeb, 0x40, 0x31, 0x02, 0x41, 0x14, 0x71,
	0x11, 0x21, 0x14, 0x55, 0xe4, 0x09, 0xd8, 0x0d, 0x55, 0x23, 0xa5, 0x90, 0x7a, 0x2a, 0xd4, 0xee,
	0x4d, 0xe4, 0xcc, 0x98, 0x64, 0xb4, 0xe3, 0xd8, 0xf2, 0x78, 0x96, 0x64, 0x1f, 0x03, 0x1e, 0x85,
	0xa7, 0xe0, 0xad, 0x90, 0xcf, 0xfc, 0x09, 0x64, 0x5f, 0x80, 0x8b, 0xde, 0xe5, 0xfb, 0x7d, 0xe3,
	0x73, 0x8e, 0xfc, 0xd9, 0x0e, 0x3c, 0xdf, 0x73, 0x99, 0xe5, 0x56, 0x1d, 0x66, 0xda, 0x28, 0xab,
	0x48, 0x74, 0xdc, 0xaa, 0x42, 0x6f, 0x27, 0x0f, 0x30, 0x78, 0x5d, 0x3b, 0x0b, 0x6e, 0x39, 0x79,
	0x01, 0x51, 0x21, 0x4c, 0xc6, 0x73, 0xea, 0x8d, 0xbd, 0xe9, 0x15, 0xab, 0x15, 0xf9, 0x14, 0x42,
	0xa9, 0x52, 0x91, 0xd3, 0x67, 0x63, 0x6f, 0xda, 0x63, 0x95, 0x20, 0x04, 0x02, 0x9b, 0x49, 0x41,
	0xfd, 0xb1, 0x37, 0x0d, 0x18, 0xfe, 0x26, 0xdf, 0x82, 0xbf, 0x9f, 0x27, 0x34, 0x18, 0x7b, 0xd3,
	0xfe, 0x0f, 0x64, 0x56, 0xf5, 0x99, 0x35, 0x4d, 0xe6, 0xb7, 0xcc, 0xd9, 0x93, 0x3f, 0x7c, 0x80,
	0x33, 0x73, 0x6d, 0x4b, 0x8d, 0xa5, 0x3c, 0x2c, 0x55, 0x2b, 0xd7, 0xf6, 0xb7, 0x9c, 0xef, 0x0a,
	0x6c, 0x7b, 0xc5, 0x2a, 0x41, 0x46, 0x10, 0xf2, 0x24, 0xd9, 0x1c, 0xb1, 0xaf, 0xc7, 0x02, 0x9e,
	0x24, 0xef, 0x1b, 0x78, 0xc2, 0xce, 0x15, 0xfc, 0xd0, 0xc0, 0x47, 0x1a, 0xb6, 0xf0, 0xce, 0x41,
	0xc9, 0x77, 0x9b, 0x23, 0x8d, 0x2a, 0x28, 0xf9, 0xee, 0x7d, 0x03, 0x4f, 0xb4, 0xd3, 0xc2, 0x0f,
	0x0d, 0x7c, 0xa4, 0xdd, 0x16, 0xde, 0x91, 0xcf, 0xa1, 0x63, 0xa5, 0xde, 0xa4, 0x99, 0xa0, 0x3d,
	0xc4, 0x91, 0x95, 0x7a, 0x91, 0x09, 0xf2, 0x0d, 0xf4, 0x9d, 0xf1, 0xa0, 0x72, 0xcb, 0x77, 0x82,
	0x02, 0x9a, 0x60, 0xa5, 0xfe, 0xb5, 0x22, 0xe4, 0x0b, 0xe8, 0xf2, 0xcc, 0x6c, 0xac, 0x90, 0x9a,
	0xf6, 0xd1, 0xed, 0xf0, 0xcc, 0xbc, 0x13, 0x52, 0xbb, 0xa2, 0xce, 0xda, 0x97, 0x92, 0x0e, 0xaa,
	0xa2, 0x3c, 0x33, 0xaf, 0x4b, 0x49, 0x3e, 0x03, 0xf7, 0x6b, 0x63, 0xf6, 0xf4, 0x0a, 0x79, 0xc8,
	0x33, 0xc3, 0xf6, 0xe4, 0x4b, 0xe8, 0xe5, 0xd9, 0x6e, 0x6f, 0x37, 0x79, 0x79, 0xa4, 0xcf, 0xd1,
	0xe9, 0x22, 0x58, 0x95, 0x47, 0x42, 0xa1, 0xb3, 0x2d, 0xad, 0x55, 0x87, 0x82, 0x7e, 0x82, 0xfb,
	0xd6, 0x48, 0xf2, 0x15, 0xf4, 0x54, 0x92, 0x94, 0x9a, 0x1f, 0x92, 0x13, 0x1d, 0xe2, 0xb2, 0x33,
	0x98, 0xfc, 0xe9, 0x43, 0xbf, 0x0d, 0x65, 0xfe, 0xf2, 0x63, 0x2a, 0xff, 0x8b, 0x54, 0xfe, 0xf2,
	0x61, 0xd4, 0xa4, 0x72, 0xc3, 0x56, 0xd9, 0xe1, 0x3e, 0xb6, 0xdc, 0xe2, 0xaa, 0x1b, 0x9e, 0xbe,
	0x32, 0x5c, 0x8a, 0xa2, 0x0e, 0xa8, 0xb7, 0x6d, 0x00, 0xf9, 0x1a, 0x60, 0xa5, 0x0a, 0x5b, 0xdb,
	0xcf, 0xd0, 0x86, 0xbc, 0x25, 0xe4, 0x3b, 0x18, 0x2e, 0x8c, 0xd2, 0x3f, 0x2b, 0x7b, 0xab, 0x0e,
	0x07, 0x91, 0x58, 0x91, 0xd6, 0xd7, 0x78, 0x98, 0x5e, 0x70, 0xf2, 0x3d, 0x5c, 0xc7, 0xa5, 0x8c,
	0xf1, 0x25, 0x60, 0x22, 0x11, 0xd9, 0x83, 0x48, 0x31, 0xd0, 0x80, 0x5d, 0x17, 0x97, 0x06, 0x99,
	0x01, 0x89, 0x4b, 0xb9, 0x50, 0x92, 0x67, 0x87, 0x57, 0xca, 0xfc, 0xce, 0x4d, 0x2a, 0x52, 0x8c,
	0x3a, 0x60, 0xa4, 0x78, 0xe2, 0x90, 0x97, 0x30, 0x72, 0xdf, 0x5f, 0x0e, 0x13, 0xe1, 0x82, 0x51,
	0xf1, 0xd4, 0xaa, 0xe7, 0xa9, 0xea, 0xb4, 0xf3, 0x74, 0xda, 0x79, 0xfe, 0x6b, 0xd4, 0xf3, 0x54,
	0x43, 0x9e, 0xe7, 0xe9, 0xb6, 0xf3, 0x5c, 0x38, 0x64, 0x0c, 0xfd, 0x75, 0xb9, 0xcd, 0xb3, 0x62,
	0xff, 0xcb, 0x3d, 0x3f, 0xe1, 0x69, 0x0a, 0x58, 0x5f, 0x9f, 0x11, 0x99, 0xc0, 0xa0, 0xfe, 0xe2,
	0x27, 0x63, 0x94, 0xc1, 0x33, 0x15, 0xb0, 0x81, 0xfe, 0x17, 0x9b, 0xfc, 0xed, 0xc1, 0xf5, 0x39,
	0xb5, 0x37, 0xa2, 0x28, 0xdc, 0x59, 0x7b, 0x01, 0x51, 0x6c, 0x92, 0x37, 0x3f, 0xde, 0x62, 0x60,
	0x3d, 0x16, 0x15, 0xa8, 0xdc, 0x8d, 0x8a, 0x4d, 0xb2, 0x5c, 0x37, 0xcf, 0x6b, 0xe1, 0x84, 0xa3,
	0x6b, 0xa5, 0x97, 0x0b, 0x0c, 0xa6, 0xc7, 0x42, 0xed, 0x84, 0x3b, 0x47, 0x6b, 0xa5, 0xdf, 0xb9,
	0x6b, 0xe9, 0x32, 0xf0, 0x59, 0x47, 0x57, 0xd2, 0x55, 0xbf, 0x61, 0x68, 0x84, 0x68, 0x44, 0x5b,
	0x54, 0xee, 0x99, 0x66, 0x71, 0xbc, 0xc4, 0x2d, 0x0d, 0x59, 0x60, 0xe2, 0x78, 0x49, 0x86, 0xe0,
	0xaf, 0xde, 0x2e, 0x71, 0xd7, 0x42, 0xe6, 0xe7, 0x6f, 0x97, 0x58, 0x97, 0x9f, 0x72, 0xc5, 0xab,
	0xcd, 0x19, 0xb0, 0x8e, 0xae, 0xe4, 0x36, 0xc2, 0xff, 0x8c, 0xf9, 0x3f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x81, 0x31, 0x12, 0x9c, 0x45, 0x06, 0x00, 0x00,
}
