// Code generated by protoc-gen-go. DO NOT EDIT.
// source: c37.proto

// from https://github.com/PingThingsIO/c37-wavemq-adapter

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

type C37DataFrame struct {
	StationName          string           `protobuf:"bytes,1,opt,name=stationName,proto3" json:"stationName,omitempty"`
	IdCode               uint32           `protobuf:"varint,2,opt,name=idCode,proto3" json:"idCode,omitempty"`
	PhasorChannels       []*PhasorChannel `protobuf:"bytes,3,rep,name=phasorChannels,proto3" json:"phasorChannels,omitempty"`
	ScalarChannels       []*ScalarChannel `protobuf:"bytes,4,rep,name=scalarChannels,proto3" json:"scalarChannels,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *C37DataFrame) Reset()         { *m = C37DataFrame{} }
func (m *C37DataFrame) String() string { return proto.CompactTextString(m) }
func (*C37DataFrame) ProtoMessage()    {}
func (*C37DataFrame) Descriptor() ([]byte, []int) {
	return fileDescriptor_eaf04cc48428f74f, []int{0}
}

func (m *C37DataFrame) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_C37DataFrame.Unmarshal(m, b)
}
func (m *C37DataFrame) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_C37DataFrame.Marshal(b, m, deterministic)
}
func (m *C37DataFrame) XXX_Merge(src proto.Message) {
	xxx_messageInfo_C37DataFrame.Merge(m, src)
}
func (m *C37DataFrame) XXX_Size() int {
	return xxx_messageInfo_C37DataFrame.Size(m)
}
func (m *C37DataFrame) XXX_DiscardUnknown() {
	xxx_messageInfo_C37DataFrame.DiscardUnknown(m)
}

var xxx_messageInfo_C37DataFrame proto.InternalMessageInfo

func (m *C37DataFrame) GetStationName() string {
	if m != nil {
		return m.StationName
	}
	return ""
}

func (m *C37DataFrame) GetIdCode() uint32 {
	if m != nil {
		return m.IdCode
	}
	return 0
}

func (m *C37DataFrame) GetPhasorChannels() []*PhasorChannel {
	if m != nil {
		return m.PhasorChannels
	}
	return nil
}

func (m *C37DataFrame) GetScalarChannels() []*ScalarChannel {
	if m != nil {
		return m.ScalarChannels
	}
	return nil
}

type PhasorChannel struct {
	ChannelName          string    `protobuf:"bytes,1,opt,name=channelName,proto3" json:"channelName,omitempty"`
	Unit                 string    `protobuf:"bytes,2,opt,name=unit,proto3" json:"unit,omitempty"`
	Data                 []*Phasor `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *PhasorChannel) Reset()         { *m = PhasorChannel{} }
func (m *PhasorChannel) String() string { return proto.CompactTextString(m) }
func (*PhasorChannel) ProtoMessage()    {}
func (*PhasorChannel) Descriptor() ([]byte, []int) {
	return fileDescriptor_eaf04cc48428f74f, []int{1}
}

func (m *PhasorChannel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PhasorChannel.Unmarshal(m, b)
}
func (m *PhasorChannel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PhasorChannel.Marshal(b, m, deterministic)
}
func (m *PhasorChannel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PhasorChannel.Merge(m, src)
}
func (m *PhasorChannel) XXX_Size() int {
	return xxx_messageInfo_PhasorChannel.Size(m)
}
func (m *PhasorChannel) XXX_DiscardUnknown() {
	xxx_messageInfo_PhasorChannel.DiscardUnknown(m)
}

var xxx_messageInfo_PhasorChannel proto.InternalMessageInfo

func (m *PhasorChannel) GetChannelName() string {
	if m != nil {
		return m.ChannelName
	}
	return ""
}

func (m *PhasorChannel) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *PhasorChannel) GetData() []*Phasor {
	if m != nil {
		return m.Data
	}
	return nil
}

type Phasor struct {
	Time                 int64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	Angle                float64  `protobuf:"fixed64,2,opt,name=angle,proto3" json:"angle,omitempty"`
	Magnitude            float64  `protobuf:"fixed64,3,opt,name=magnitude,proto3" json:"magnitude,omitempty"`
	P                    float64  `protobuf:"fixed64,4,opt,name=P,json=p,proto3" json:"P,omitempty"`
	Q                    float64  `protobuf:"fixed64,5,opt,name=Q,json=q,proto3" json:"Q,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Phasor) Reset()         { *m = Phasor{} }
func (m *Phasor) String() string { return proto.CompactTextString(m) }
func (*Phasor) ProtoMessage()    {}
func (*Phasor) Descriptor() ([]byte, []int) {
	return fileDescriptor_eaf04cc48428f74f, []int{2}
}

func (m *Phasor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Phasor.Unmarshal(m, b)
}
func (m *Phasor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Phasor.Marshal(b, m, deterministic)
}
func (m *Phasor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Phasor.Merge(m, src)
}
func (m *Phasor) XXX_Size() int {
	return xxx_messageInfo_Phasor.Size(m)
}
func (m *Phasor) XXX_DiscardUnknown() {
	xxx_messageInfo_Phasor.DiscardUnknown(m)
}

var xxx_messageInfo_Phasor proto.InternalMessageInfo

func (m *Phasor) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *Phasor) GetAngle() float64 {
	if m != nil {
		return m.Angle
	}
	return 0
}

func (m *Phasor) GetMagnitude() float64 {
	if m != nil {
		return m.Magnitude
	}
	return 0
}

func (m *Phasor) GetP() float64 {
	if m != nil {
		return m.P
	}
	return 0
}

func (m *Phasor) GetQ() float64 {
	if m != nil {
		return m.Q
	}
	return 0
}

type ScalarChannel struct {
	ChannelName          string    `protobuf:"bytes,1,opt,name=channelName,proto3" json:"channelName,omitempty"`
	Unit                 string    `protobuf:"bytes,2,opt,name=unit,proto3" json:"unit,omitempty"`
	Data                 []*Scalar `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ScalarChannel) Reset()         { *m = ScalarChannel{} }
func (m *ScalarChannel) String() string { return proto.CompactTextString(m) }
func (*ScalarChannel) ProtoMessage()    {}
func (*ScalarChannel) Descriptor() ([]byte, []int) {
	return fileDescriptor_eaf04cc48428f74f, []int{3}
}

func (m *ScalarChannel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScalarChannel.Unmarshal(m, b)
}
func (m *ScalarChannel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScalarChannel.Marshal(b, m, deterministic)
}
func (m *ScalarChannel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScalarChannel.Merge(m, src)
}
func (m *ScalarChannel) XXX_Size() int {
	return xxx_messageInfo_ScalarChannel.Size(m)
}
func (m *ScalarChannel) XXX_DiscardUnknown() {
	xxx_messageInfo_ScalarChannel.DiscardUnknown(m)
}

var xxx_messageInfo_ScalarChannel proto.InternalMessageInfo

func (m *ScalarChannel) GetChannelName() string {
	if m != nil {
		return m.ChannelName
	}
	return ""
}

func (m *ScalarChannel) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

func (m *ScalarChannel) GetData() []*Scalar {
	if m != nil {
		return m.Data
	}
	return nil
}

type Scalar struct {
	Time                 int64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	Value                float64  `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Scalar) Reset()         { *m = Scalar{} }
func (m *Scalar) String() string { return proto.CompactTextString(m) }
func (*Scalar) ProtoMessage()    {}
func (*Scalar) Descriptor() ([]byte, []int) {
	return fileDescriptor_eaf04cc48428f74f, []int{4}
}

func (m *Scalar) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Scalar.Unmarshal(m, b)
}
func (m *Scalar) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Scalar.Marshal(b, m, deterministic)
}
func (m *Scalar) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Scalar.Merge(m, src)
}
func (m *Scalar) XXX_Size() int {
	return xxx_messageInfo_Scalar.Size(m)
}
func (m *Scalar) XXX_DiscardUnknown() {
	xxx_messageInfo_Scalar.DiscardUnknown(m)
}

var xxx_messageInfo_Scalar proto.InternalMessageInfo

func (m *Scalar) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *Scalar) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func init() {
	proto.RegisterType((*C37DataFrame)(nil), "xbospb.C37DataFrame")
	proto.RegisterType((*PhasorChannel)(nil), "xbospb.PhasorChannel")
	proto.RegisterType((*Phasor)(nil), "xbospb.Phasor")
	proto.RegisterType((*ScalarChannel)(nil), "xbospb.ScalarChannel")
	proto.RegisterType((*Scalar)(nil), "xbospb.Scalar")
}

func init() { proto.RegisterFile("c37.proto", fileDescriptor_eaf04cc48428f74f) }

var fileDescriptor_eaf04cc48428f74f = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0xbd, 0x6e, 0x83, 0x30,
	0x14, 0x85, 0xe5, 0x42, 0x90, 0xb8, 0xf9, 0x19, 0xac, 0xb6, 0xf2, 0xd0, 0x01, 0x31, 0x31, 0x31,
	0x84, 0x21, 0x53, 0x27, 0xaa, 0x8e, 0x55, 0xea, 0x3e, 0xc1, 0xe5, 0x47, 0x89, 0x25, 0xb0, 0x29,
	0x98, 0xaa, 0x4f, 0xd8, 0xe7, 0xaa, 0xb0, 0x51, 0x63, 0x22, 0x75, 0xcb, 0x76, 0xcf, 0x39, 0xbe,
	0x07, 0x7f, 0xc8, 0x10, 0x96, 0xd9, 0x21, 0xed, 0x7a, 0xa5, 0x15, 0x0d, 0xbe, 0x0b, 0x35, 0x74,
	0x45, 0xfc, 0x43, 0x60, 0x93, 0x67, 0x87, 0x17, 0xd4, 0xf8, 0xda, 0x63, 0x5b, 0xd3, 0x08, 0xd6,
	0x83, 0x46, 0x2d, 0x94, 0x7c, 0xc3, 0xb6, 0x66, 0x24, 0x22, 0x49, 0xc8, 0x5d, 0x8b, 0x3e, 0x42,
	0x20, 0xaa, 0x5c, 0x55, 0x35, 0xbb, 0x8b, 0x48, 0xb2, 0xe5, 0xb3, 0xa2, 0xcf, 0xb0, 0xeb, 0xce,
	0x38, 0xa8, 0x3e, 0x3f, 0xa3, 0x94, 0x75, 0x33, 0x30, 0x2f, 0xf2, 0x92, 0xf5, 0xfe, 0x21, 0xb5,
	0xdf, 0x4a, 0x8f, 0x6e, 0xca, 0xaf, 0x0e, 0x4f, 0xeb, 0x43, 0x89, 0x0d, 0x5e, 0xd6, 0xfd, 0xe5,
	0xfa, 0x87, 0x9b, 0xf2, 0xab, 0xc3, 0xb1, 0x80, 0xed, 0xa2, 0x7f, 0x02, 0x29, 0xed, 0xe8, 0x82,
	0x38, 0x16, 0xa5, 0xe0, 0x8f, 0x52, 0x68, 0x83, 0x11, 0x72, 0x33, 0xd3, 0x18, 0xfc, 0x0a, 0x35,
	0xce, 0x57, 0xdf, 0x2d, 0xaf, 0xce, 0x4d, 0x16, 0x37, 0x10, 0x58, 0x3d, 0x35, 0x68, 0x31, 0x97,
	0x7b, 0xdc, 0xcc, 0xf4, 0x1e, 0x56, 0x28, 0x4f, 0x8d, 0xfd, 0x3b, 0x84, 0x5b, 0x41, 0x9f, 0x20,
	0x6c, 0xf1, 0x24, 0x85, 0x1e, 0xab, 0x9a, 0x79, 0x26, 0xb9, 0x18, 0x74, 0x03, 0xe4, 0xc8, 0x7c,
	0xe3, 0x92, 0x6e, 0x52, 0xef, 0x6c, 0x65, 0xd5, 0xe7, 0x04, 0xb6, 0x20, 0xbf, 0x2d, 0x98, 0xad,
	0x9e, 0xc1, 0xf6, 0x10, 0x58, 0xfd, 0x1f, 0xd8, 0x17, 0x36, 0xe3, 0x1f, 0x98, 0x11, 0x45, 0x60,
	0xde, 0x53, 0xf6, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x64, 0x2d, 0xa8, 0x8d, 0x5c, 0x02, 0x00, 0x00,
}
