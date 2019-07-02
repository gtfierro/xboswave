// Code generated by protoc-gen-go. DO NOT EDIT.
// source: energise.proto

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

type EnergiseMessage struct {
	SPBC                 *SPBC        `protobuf:"bytes,1,opt,name=SPBC,json=sPBC,proto3" json:"SPBC,omitempty"`
	LPBCStatus           *LPBCStatus  `protobuf:"bytes,2,opt,name=LPBCStatus,json=lPBCStatus,proto3" json:"LPBCStatus,omitempty"`
	LPBCCommand          *LPBCCommand `protobuf:"bytes,3,opt,name=LPBCCommand,json=lPBCCommand,proto3" json:"LPBCCommand,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *EnergiseMessage) Reset()         { *m = EnergiseMessage{} }
func (m *EnergiseMessage) String() string { return proto.CompactTextString(m) }
func (*EnergiseMessage) ProtoMessage()    {}
func (*EnergiseMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a37b97c5e2bc6c22, []int{0}
}

func (m *EnergiseMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnergiseMessage.Unmarshal(m, b)
}
func (m *EnergiseMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnergiseMessage.Marshal(b, m, deterministic)
}
func (m *EnergiseMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnergiseMessage.Merge(m, src)
}
func (m *EnergiseMessage) XXX_Size() int {
	return xxx_messageInfo_EnergiseMessage.Size(m)
}
func (m *EnergiseMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_EnergiseMessage.DiscardUnknown(m)
}

var xxx_messageInfo_EnergiseMessage proto.InternalMessageInfo

func (m *EnergiseMessage) GetSPBC() *SPBC {
	if m != nil {
		return m.SPBC
	}
	return nil
}

func (m *EnergiseMessage) GetLPBCStatus() *LPBCStatus {
	if m != nil {
		return m.LPBCStatus
	}
	return nil
}

func (m *EnergiseMessage) GetLPBCCommand() *LPBCCommand {
	if m != nil {
		return m.LPBCCommand
	}
	return nil
}

type EnergiseError struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EnergiseError) Reset()         { *m = EnergiseError{} }
func (m *EnergiseError) String() string { return proto.CompactTextString(m) }
func (*EnergiseError) ProtoMessage()    {}
func (*EnergiseError) Descriptor() ([]byte, []int) {
	return fileDescriptor_a37b97c5e2bc6c22, []int{1}
}

func (m *EnergiseError) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnergiseError.Unmarshal(m, b)
}
func (m *EnergiseError) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnergiseError.Marshal(b, m, deterministic)
}
func (m *EnergiseError) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnergiseError.Merge(m, src)
}
func (m *EnergiseError) XXX_Size() int {
	return xxx_messageInfo_EnergiseError.Size(m)
}
func (m *EnergiseError) XXX_DiscardUnknown() {
	xxx_messageInfo_EnergiseError.DiscardUnknown(m)
}

var xxx_messageInfo_EnergiseError proto.InternalMessageInfo

func (m *EnergiseError) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

// EnergisePhasorTarget is a control target specified on a per-node basis
// upmu0 is the 'head of feeder'. This is reference for everything: SPBC, LPBCs
type EnergisePhasorTarget struct {
	NodeID               string   `protobuf:"bytes,1,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	Magnitude            float64  `protobuf:"fixed64,2,opt,name=magnitude,proto3" json:"magnitude,omitempty"`
	Angle                float64  `protobuf:"fixed64,3,opt,name=angle,proto3" json:"angle,omitempty"`
	Kvbase               *Double  `protobuf:"bytes,4,opt,name=kvbase,proto3" json:"kvbase,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EnergisePhasorTarget) Reset()         { *m = EnergisePhasorTarget{} }
func (m *EnergisePhasorTarget) String() string { return proto.CompactTextString(m) }
func (*EnergisePhasorTarget) ProtoMessage()    {}
func (*EnergisePhasorTarget) Descriptor() ([]byte, []int) {
	return fileDescriptor_a37b97c5e2bc6c22, []int{2}
}

func (m *EnergisePhasorTarget) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnergisePhasorTarget.Unmarshal(m, b)
}
func (m *EnergisePhasorTarget) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnergisePhasorTarget.Marshal(b, m, deterministic)
}
func (m *EnergisePhasorTarget) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnergisePhasorTarget.Merge(m, src)
}
func (m *EnergisePhasorTarget) XXX_Size() int {
	return xxx_messageInfo_EnergisePhasorTarget.Size(m)
}
func (m *EnergisePhasorTarget) XXX_DiscardUnknown() {
	xxx_messageInfo_EnergisePhasorTarget.DiscardUnknown(m)
}

var xxx_messageInfo_EnergisePhasorTarget proto.InternalMessageInfo

func (m *EnergisePhasorTarget) GetNodeID() string {
	if m != nil {
		return m.NodeID
	}
	return ""
}

func (m *EnergisePhasorTarget) GetMagnitude() float64 {
	if m != nil {
		return m.Magnitude
	}
	return 0
}

func (m *EnergisePhasorTarget) GetAngle() float64 {
	if m != nil {
		return m.Angle
	}
	return 0
}

func (m *EnergisePhasorTarget) GetKvbase() *Double {
	if m != nil {
		return m.Kvbase
	}
	return nil
}

// The SPBC message is sent by a supervisory controller (also called an SPBC)
// at regular intervals.  The expectation is the SPBC will send out a single
// message for each node that it is controlling, containing an
// EnergisePhasorTarget for that node. We restrict each message to a single
// phasor_target for now in order to maintain isolation between the nodes and
// bound what information they are allowed to see.
type SPBC struct {
	// current time of announcement in milliseconds
	Time int64 `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	// a phasor target for a specific node
	PhasorTarget *EnergisePhasorTarget `protobuf:"bytes,2,opt,name=phasor_target,json=phasorTarget,proto3" json:"phasor_target,omitempty"`
	// represents general errors in the SPBC
	Error                *EnergiseError `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *SPBC) Reset()         { *m = SPBC{} }
func (m *SPBC) String() string { return proto.CompactTextString(m) }
func (*SPBC) ProtoMessage()    {}
func (*SPBC) Descriptor() ([]byte, []int) {
	return fileDescriptor_a37b97c5e2bc6c22, []int{3}
}

func (m *SPBC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SPBC.Unmarshal(m, b)
}
func (m *SPBC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SPBC.Marshal(b, m, deterministic)
}
func (m *SPBC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SPBC.Merge(m, src)
}
func (m *SPBC) XXX_Size() int {
	return xxx_messageInfo_SPBC.Size(m)
}
func (m *SPBC) XXX_DiscardUnknown() {
	xxx_messageInfo_SPBC.DiscardUnknown(m)
}

var xxx_messageInfo_SPBC proto.InternalMessageInfo

func (m *SPBC) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *SPBC) GetPhasorTarget() *EnergisePhasorTarget {
	if m != nil {
		return m.PhasorTarget
	}
	return nil
}

func (m *SPBC) GetError() *EnergiseError {
	if m != nil {
		return m.Error
	}
	return nil
}

type LPBCStatus struct {
	// current time of announcement in milliseconds
	Time int64 `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	// contains potential errors
	Error *EnergiseError `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	// current P, Q errors of LPBC
	PhasorErrors *Phasor `protobuf:"bytes,3,opt,name=phasor_errors,json=phasorErrors,proto3" json:"phasor_errors,omitempty"`
	// true if LPBC P is saturated
	PSaturated bool `protobuf:"varint,4,opt,name=p_saturated,json=pSaturated,proto3" json:"p_saturated,omitempty"`
	// true if LPBC Q is saturated
	QSaturated bool `protobuf:"varint,5,opt,name=q_saturated,json=qSaturated,proto3" json:"q_saturated,omitempty"`
	// true if the LPBC is performing control
	DoControl bool `protobuf:"varint,6,opt,name=do_control,json=doControl,proto3" json:"do_control,omitempty"`
	// should be populated if p_saturated or q_saturated
	// gives the value at which p or q saturated
	PMax                 *Double  `protobuf:"bytes,7,opt,name=p_max,json=pMax,proto3" json:"p_max,omitempty"`
	QMax                 *Double  `protobuf:"bytes,8,opt,name=q_max,json=qMax,proto3" json:"q_max,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LPBCStatus) Reset()         { *m = LPBCStatus{} }
func (m *LPBCStatus) String() string { return proto.CompactTextString(m) }
func (*LPBCStatus) ProtoMessage()    {}
func (*LPBCStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_a37b97c5e2bc6c22, []int{4}
}

func (m *LPBCStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LPBCStatus.Unmarshal(m, b)
}
func (m *LPBCStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LPBCStatus.Marshal(b, m, deterministic)
}
func (m *LPBCStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LPBCStatus.Merge(m, src)
}
func (m *LPBCStatus) XXX_Size() int {
	return xxx_messageInfo_LPBCStatus.Size(m)
}
func (m *LPBCStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_LPBCStatus.DiscardUnknown(m)
}

var xxx_messageInfo_LPBCStatus proto.InternalMessageInfo

func (m *LPBCStatus) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *LPBCStatus) GetError() *EnergiseError {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *LPBCStatus) GetPhasorErrors() *Phasor {
	if m != nil {
		return m.PhasorErrors
	}
	return nil
}

func (m *LPBCStatus) GetPSaturated() bool {
	if m != nil {
		return m.PSaturated
	}
	return false
}

func (m *LPBCStatus) GetQSaturated() bool {
	if m != nil {
		return m.QSaturated
	}
	return false
}

func (m *LPBCStatus) GetDoControl() bool {
	if m != nil {
		return m.DoControl
	}
	return false
}

func (m *LPBCStatus) GetPMax() *Double {
	if m != nil {
		return m.PMax
	}
	return nil
}

func (m *LPBCStatus) GetQMax() *Double {
	if m != nil {
		return m.QMax
	}
	return nil
}

type LPBCCommand struct {
	// current time of announcement in milliseconds
	Time int64 `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	// new phasor target (direct actuation). LPBC will also listen to phasor targets
	// from SPBC messages
	PhasorTarget *Phasor `protobuf:"bytes,2,opt,name=phasor_target,json=phasorTarget,proto3" json:"phasor_target,omitempty"`
	// set whether or not the LPBC is performing control
	DoControl            bool     `protobuf:"varint,3,opt,name=do_control,json=doControl,proto3" json:"do_control,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LPBCCommand) Reset()         { *m = LPBCCommand{} }
func (m *LPBCCommand) String() string { return proto.CompactTextString(m) }
func (*LPBCCommand) ProtoMessage()    {}
func (*LPBCCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_a37b97c5e2bc6c22, []int{5}
}

func (m *LPBCCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LPBCCommand.Unmarshal(m, b)
}
func (m *LPBCCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LPBCCommand.Marshal(b, m, deterministic)
}
func (m *LPBCCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LPBCCommand.Merge(m, src)
}
func (m *LPBCCommand) XXX_Size() int {
	return xxx_messageInfo_LPBCCommand.Size(m)
}
func (m *LPBCCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_LPBCCommand.DiscardUnknown(m)
}

var xxx_messageInfo_LPBCCommand proto.InternalMessageInfo

func (m *LPBCCommand) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *LPBCCommand) GetPhasorTarget() *Phasor {
	if m != nil {
		return m.PhasorTarget
	}
	return nil
}

func (m *LPBCCommand) GetDoControl() bool {
	if m != nil {
		return m.DoControl
	}
	return false
}

func init() {
	proto.RegisterType((*EnergiseMessage)(nil), "xbospb.EnergiseMessage")
	proto.RegisterType((*EnergiseError)(nil), "xbospb.EnergiseError")
	proto.RegisterType((*EnergisePhasorTarget)(nil), "xbospb.EnergisePhasorTarget")
	proto.RegisterType((*SPBC)(nil), "xbospb.SPBC")
	proto.RegisterType((*LPBCStatus)(nil), "xbospb.LPBCStatus")
	proto.RegisterType((*LPBCCommand)(nil), "xbospb.LPBCCommand")
}

func init() { proto.RegisterFile("energise.proto", fileDescriptor_a37b97c5e2bc6c22) }

var fileDescriptor_a37b97c5e2bc6c22 = []byte{
	// 460 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x93, 0x51, 0x8b, 0xd3, 0x40,
	0x10, 0xc7, 0x49, 0x9b, 0xc6, 0xeb, 0xf4, 0xee, 0x94, 0xbd, 0x53, 0xc2, 0x71, 0xe2, 0x19, 0x41,
	0x04, 0xa1, 0x0f, 0x57, 0xc4, 0x67, 0x9b, 0xbb, 0x07, 0xc1, 0x42, 0xd9, 0xfa, 0x5e, 0x36, 0x97,
	0x21, 0x16, 0x93, 0xec, 0x76, 0x77, 0x23, 0xf5, 0x2b, 0xf8, 0x21, 0xfc, 0x06, 0x7e, 0x47, 0xc9,
	0x64, 0xd3, 0xe6, 0x4a, 0xc4, 0xb7, 0x9d, 0x99, 0xff, 0x4c, 0x7f, 0xf3, 0x9f, 0x06, 0xce, 0xb1,
	0x44, 0x9d, 0x6d, 0x0c, 0x4e, 0x95, 0x96, 0x56, 0xb2, 0x60, 0x97, 0x48, 0xa3, 0x92, 0xab, 0xf1,
	0xc3, 0xec, 0x63, 0x93, 0xba, 0xba, 0x28, 0xab, 0x3c, 0x17, 0x49, 0x8e, 0xf6, 0xa7, 0x42, 0xd3,
	0x24, 0xa3, 0xdf, 0x1e, 0x3c, 0xbd, 0x77, 0xad, 0x0b, 0x34, 0x46, 0x64, 0xc8, 0x6e, 0xc0, 0x5f,
	0x2d, 0xe7, 0x71, 0xe8, 0xdd, 0x78, 0xef, 0x26, 0xb7, 0xa7, 0xd3, 0x66, 0xd4, 0xb4, 0xce, 0x71,
	0xdf, 0x2c, 0xe7, 0x31, 0xbb, 0x05, 0xf8, 0xb2, 0x9c, 0xc7, 0x2b, 0x2b, 0x6c, 0x65, 0xc2, 0x01,
	0xe9, 0x58, 0xab, 0x3b, 0x54, 0x38, 0xe4, 0xfb, 0x37, 0xfb, 0x00, 0x93, 0xba, 0x12, 0xcb, 0xa2,
	0x10, 0x65, 0x1a, 0x0e, 0xa9, 0xe9, 0xa2, 0xdb, 0xe4, 0x4a, 0x7c, 0x92, 0x1f, 0x82, 0xe8, 0x35,
	0x9c, 0xb5, 0x7c, 0xf7, 0x5a, 0x4b, 0xcd, 0x9e, 0xc1, 0xb0, 0x30, 0x19, 0xc1, 0x8d, 0x79, 0xfd,
	0x8c, 0x7e, 0x79, 0x70, 0xd9, 0x6a, 0x96, 0xdf, 0x84, 0x91, 0xfa, 0xab, 0xd0, 0x19, 0x5a, 0xf6,
	0x02, 0x82, 0x52, 0xa6, 0xf8, 0xf9, 0xce, 0xa9, 0x5d, 0xc4, 0xae, 0x61, 0x5c, 0x88, 0xac, 0xdc,
	0xd8, 0x2a, 0x45, 0xa2, 0xf7, 0xf8, 0x21, 0xc1, 0x2e, 0x61, 0x24, 0xca, 0x2c, 0x47, 0x42, 0xf4,
	0x78, 0x13, 0xb0, 0xb7, 0x10, 0x7c, 0xff, 0x91, 0x08, 0x83, 0xa1, 0x4f, 0xe4, 0xe7, 0x2d, 0xf9,
	0x9d, 0xac, 0x92, 0x1c, 0xb9, 0xab, 0xd6, 0x30, 0xe4, 0x1e, 0x63, 0xe0, 0xdb, 0x4d, 0x81, 0xf4,
	0xd3, 0x43, 0x4e, 0x6f, 0xf6, 0x09, 0xce, 0x14, 0x01, 0xae, 0x2d, 0x11, 0x3a, 0xeb, 0xae, 0xdb,
	0x59, 0x7d, 0x5b, 0xf0, 0x53, 0xd5, 0xdd, 0xe9, 0x3d, 0x8c, 0xb0, 0xf6, 0xc1, 0x19, 0xf8, 0xfc,
	0xb8, 0x95, 0x4c, 0xe2, 0x8d, 0x26, 0xfa, 0x33, 0xe8, 0x1e, 0xaa, 0x17, 0x69, 0x3f, 0x6f, 0xf0,
	0xff, 0x79, 0x6c, 0xb6, 0xe7, 0xa7, 0xd8, 0x38, 0x88, 0xbd, 0x17, 0x0d, 0x77, 0x4b, 0x4c, 0xad,
	0x86, 0xbd, 0x82, 0x89, 0x5a, 0x1b, 0x61, 0x2b, 0x2d, 0x2c, 0xa6, 0x64, 0xdf, 0x09, 0x07, 0xb5,
	0x6a, 0x33, 0xb5, 0x60, 0xdb, 0x11, 0x8c, 0x1a, 0xc1, 0xf6, 0x20, 0x78, 0x09, 0x90, 0xca, 0xf5,
	0x83, 0x2c, 0xad, 0x96, 0x79, 0x18, 0x50, 0x7d, 0x9c, 0xca, 0xb8, 0x49, 0xb0, 0x37, 0x30, 0x52,
	0xeb, 0x42, 0xec, 0xc2, 0x27, 0xbd, 0x97, 0xf1, 0xd5, 0x42, 0xec, 0x6a, 0xd1, 0x96, 0x44, 0x27,
	0xfd, 0xa2, 0xed, 0x42, 0xec, 0xa2, 0xea, 0xd1, 0x7f, 0xb4, 0xd7, 0xaf, 0x59, 0xff, 0x09, 0xff,
	0x61, 0x81, 0x3b, 0xda, 0xe3, 0x05, 0x86, 0x47, 0x0b, 0x24, 0x01, 0x7d, 0x8b, 0xb3, 0xbf, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x31, 0x5f, 0x98, 0x83, 0xc5, 0x03, 0x00, 0x00,
}
