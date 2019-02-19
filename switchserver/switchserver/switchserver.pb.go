// Code generated by protoc-gen-go. DO NOT EDIT.
// source: switchserver/switchserver.proto

package switchserver

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"
import wrappers "github.com/golang/protobuf/ptypes/wrappers"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SetSwitchStatus_Status int32

const (
	SetSwitchStatus_SETTING   SetSwitchStatus_Status = 0
	SetSwitchStatus_WAITING   SetSwitchStatus_Status = 1
	SetSwitchStatus_RETURNING SetSwitchStatus_Status = 2
	SetSwitchStatus_DONE      SetSwitchStatus_Status = 3
)

var SetSwitchStatus_Status_name = map[int32]string{
	0: "SETTING",
	1: "WAITING",
	2: "RETURNING",
	3: "DONE",
}
var SetSwitchStatus_Status_value = map[string]int32{
	"SETTING":   0,
	"WAITING":   1,
	"RETURNING": 2,
	"DONE":      3,
}

func (x SetSwitchStatus_Status) String() string {
	return proto.EnumName(SetSwitchStatus_Status_name, int32(x))
}
func (SetSwitchStatus_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_switchserver_391a31fcd7eed6bf, []int{5, 0}
}

type Switch struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	IsOn                 bool     `protobuf:"varint,2,opt,name=isOn,proto3" json:"isOn,omitempty"`
	OnAngle              int64    `protobuf:"varint,3,opt,name=onAngle,proto3" json:"onAngle,omitempty"`
	OffAngle             int64    `protobuf:"varint,4,opt,name=offAngle,proto3" json:"offAngle,omitempty"`
	RestAngle            int64    `protobuf:"varint,5,opt,name=restAngle,proto3" json:"restAngle,omitempty"`
	IsCalibrated         bool     `protobuf:"varint,6,opt,name=isCalibrated,proto3" json:"isCalibrated,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Switch) Reset()         { *m = Switch{} }
func (m *Switch) String() string { return proto.CompactTextString(m) }
func (*Switch) ProtoMessage()    {}
func (*Switch) Descriptor() ([]byte, []int) {
	return fileDescriptor_switchserver_391a31fcd7eed6bf, []int{0}
}
func (m *Switch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Switch.Unmarshal(m, b)
}
func (m *Switch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Switch.Marshal(b, m, deterministic)
}
func (dst *Switch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Switch.Merge(dst, src)
}
func (m *Switch) XXX_Size() int {
	return xxx_messageInfo_Switch.Size(m)
}
func (m *Switch) XXX_DiscardUnknown() {
	xxx_messageInfo_Switch.DiscardUnknown(m)
}

var xxx_messageInfo_Switch proto.InternalMessageInfo

func (m *Switch) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Switch) GetIsOn() bool {
	if m != nil {
		return m.IsOn
	}
	return false
}

func (m *Switch) GetOnAngle() int64 {
	if m != nil {
		return m.OnAngle
	}
	return 0
}

func (m *Switch) GetOffAngle() int64 {
	if m != nil {
		return m.OffAngle
	}
	return 0
}

func (m *Switch) GetRestAngle() int64 {
	if m != nil {
		return m.RestAngle
	}
	return 0
}

func (m *Switch) GetIsCalibrated() bool {
	if m != nil {
		return m.IsCalibrated
	}
	return false
}

type AddSwitchRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	IsOn                 bool     `protobuf:"varint,2,opt,name=isOn,proto3" json:"isOn,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddSwitchRequest) Reset()         { *m = AddSwitchRequest{} }
func (m *AddSwitchRequest) String() string { return proto.CompactTextString(m) }
func (*AddSwitchRequest) ProtoMessage()    {}
func (*AddSwitchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_switchserver_391a31fcd7eed6bf, []int{1}
}
func (m *AddSwitchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddSwitchRequest.Unmarshal(m, b)
}
func (m *AddSwitchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddSwitchRequest.Marshal(b, m, deterministic)
}
func (dst *AddSwitchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddSwitchRequest.Merge(dst, src)
}
func (m *AddSwitchRequest) XXX_Size() int {
	return xxx_messageInfo_AddSwitchRequest.Size(m)
}
func (m *AddSwitchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddSwitchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddSwitchRequest proto.InternalMessageInfo

func (m *AddSwitchRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *AddSwitchRequest) GetIsOn() bool {
	if m != nil {
		return m.IsOn
	}
	return false
}

type RemoveSwitchRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RemoveSwitchRequest) Reset()         { *m = RemoveSwitchRequest{} }
func (m *RemoveSwitchRequest) String() string { return proto.CompactTextString(m) }
func (*RemoveSwitchRequest) ProtoMessage()    {}
func (*RemoveSwitchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_switchserver_391a31fcd7eed6bf, []int{2}
}
func (m *RemoveSwitchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemoveSwitchRequest.Unmarshal(m, b)
}
func (m *RemoveSwitchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemoveSwitchRequest.Marshal(b, m, deterministic)
}
func (dst *RemoveSwitchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoveSwitchRequest.Merge(dst, src)
}
func (m *RemoveSwitchRequest) XXX_Size() int {
	return xxx_messageInfo_RemoveSwitchRequest.Size(m)
}
func (m *RemoveSwitchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoveSwitchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RemoveSwitchRequest proto.InternalMessageInfo

func (m *RemoveSwitchRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type SwitchQuery struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SwitchQuery) Reset()         { *m = SwitchQuery{} }
func (m *SwitchQuery) String() string { return proto.CompactTextString(m) }
func (*SwitchQuery) ProtoMessage()    {}
func (*SwitchQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_switchserver_391a31fcd7eed6bf, []int{3}
}
func (m *SwitchQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchQuery.Unmarshal(m, b)
}
func (m *SwitchQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchQuery.Marshal(b, m, deterministic)
}
func (dst *SwitchQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchQuery.Merge(dst, src)
}
func (m *SwitchQuery) XXX_Size() int {
	return xxx_messageInfo_SwitchQuery.Size(m)
}
func (m *SwitchQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchQuery.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchQuery proto.InternalMessageInfo

func (m *SwitchQuery) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type SetSwitchRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	On                   bool     `protobuf:"varint,2,opt,name=on,proto3" json:"on,omitempty"`
	Force                bool     `protobuf:"varint,3,opt,name=force,proto3" json:"force,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetSwitchRequest) Reset()         { *m = SetSwitchRequest{} }
func (m *SetSwitchRequest) String() string { return proto.CompactTextString(m) }
func (*SetSwitchRequest) ProtoMessage()    {}
func (*SetSwitchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_switchserver_391a31fcd7eed6bf, []int{4}
}
func (m *SetSwitchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetSwitchRequest.Unmarshal(m, b)
}
func (m *SetSwitchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetSwitchRequest.Marshal(b, m, deterministic)
}
func (dst *SetSwitchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetSwitchRequest.Merge(dst, src)
}
func (m *SetSwitchRequest) XXX_Size() int {
	return xxx_messageInfo_SetSwitchRequest.Size(m)
}
func (m *SetSwitchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetSwitchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetSwitchRequest proto.InternalMessageInfo

func (m *SetSwitchRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SetSwitchRequest) GetOn() bool {
	if m != nil {
		return m.On
	}
	return false
}

func (m *SetSwitchRequest) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
}

type SetSwitchStatus struct {
	Status               SetSwitchStatus_Status `protobuf:"varint,1,opt,name=status,proto3,enum=SetSwitchStatus_Status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *SetSwitchStatus) Reset()         { *m = SetSwitchStatus{} }
func (m *SetSwitchStatus) String() string { return proto.CompactTextString(m) }
func (*SetSwitchStatus) ProtoMessage()    {}
func (*SetSwitchStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_switchserver_391a31fcd7eed6bf, []int{5}
}
func (m *SetSwitchStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetSwitchStatus.Unmarshal(m, b)
}
func (m *SetSwitchStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetSwitchStatus.Marshal(b, m, deterministic)
}
func (dst *SetSwitchStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetSwitchStatus.Merge(dst, src)
}
func (m *SetSwitchStatus) XXX_Size() int {
	return xxx_messageInfo_SetSwitchStatus.Size(m)
}
func (m *SetSwitchStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_SetSwitchStatus.DiscardUnknown(m)
}

var xxx_messageInfo_SetSwitchStatus proto.InternalMessageInfo

func (m *SetSwitchStatus) GetStatus() SetSwitchStatus_Status {
	if m != nil {
		return m.Status
	}
	return SetSwitchStatus_SETTING
}

type SwitchCalibrationParameters struct {
	Id                   string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	OnAngle              *wrappers.Int64Value `protobuf:"bytes,2,opt,name=onAngle,proto3" json:"onAngle,omitempty"`
	OffAngle             *wrappers.Int64Value `protobuf:"bytes,3,opt,name=offAngle,proto3" json:"offAngle,omitempty"`
	RestAngle            *wrappers.Int64Value `protobuf:"bytes,4,opt,name=restAngle,proto3" json:"restAngle,omitempty"`
	IsCalibrated         *wrappers.BoolValue  `protobuf:"bytes,5,opt,name=isCalibrated,proto3" json:"isCalibrated,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *SwitchCalibrationParameters) Reset()         { *m = SwitchCalibrationParameters{} }
func (m *SwitchCalibrationParameters) String() string { return proto.CompactTextString(m) }
func (*SwitchCalibrationParameters) ProtoMessage()    {}
func (*SwitchCalibrationParameters) Descriptor() ([]byte, []int) {
	return fileDescriptor_switchserver_391a31fcd7eed6bf, []int{6}
}
func (m *SwitchCalibrationParameters) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchCalibrationParameters.Unmarshal(m, b)
}
func (m *SwitchCalibrationParameters) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchCalibrationParameters.Marshal(b, m, deterministic)
}
func (dst *SwitchCalibrationParameters) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchCalibrationParameters.Merge(dst, src)
}
func (m *SwitchCalibrationParameters) XXX_Size() int {
	return xxx_messageInfo_SwitchCalibrationParameters.Size(m)
}
func (m *SwitchCalibrationParameters) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchCalibrationParameters.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchCalibrationParameters proto.InternalMessageInfo

func (m *SwitchCalibrationParameters) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SwitchCalibrationParameters) GetOnAngle() *wrappers.Int64Value {
	if m != nil {
		return m.OnAngle
	}
	return nil
}

func (m *SwitchCalibrationParameters) GetOffAngle() *wrappers.Int64Value {
	if m != nil {
		return m.OffAngle
	}
	return nil
}

func (m *SwitchCalibrationParameters) GetRestAngle() *wrappers.Int64Value {
	if m != nil {
		return m.RestAngle
	}
	return nil
}

func (m *SwitchCalibrationParameters) GetIsCalibrated() *wrappers.BoolValue {
	if m != nil {
		return m.IsCalibrated
	}
	return nil
}

func init() {
	proto.RegisterType((*Switch)(nil), "Switch")
	proto.RegisterType((*AddSwitchRequest)(nil), "AddSwitchRequest")
	proto.RegisterType((*RemoveSwitchRequest)(nil), "RemoveSwitchRequest")
	proto.RegisterType((*SwitchQuery)(nil), "SwitchQuery")
	proto.RegisterType((*SetSwitchRequest)(nil), "SetSwitchRequest")
	proto.RegisterType((*SetSwitchStatus)(nil), "SetSwitchStatus")
	proto.RegisterType((*SwitchCalibrationParameters)(nil), "SwitchCalibrationParameters")
	proto.RegisterEnum("SetSwitchStatus_Status", SetSwitchStatus_Status_name, SetSwitchStatus_Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SwitchServerClient is the client API for SwitchServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SwitchServerClient interface {
	AddSwitch(ctx context.Context, in *AddSwitchRequest, opts ...grpc.CallOption) (*Switch, error)
	RemoveSwitch(ctx context.Context, in *RemoveSwitchRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	GetSwitch(ctx context.Context, in *SwitchQuery, opts ...grpc.CallOption) (*Switch, error)
	SetSwitch(ctx context.Context, in *SetSwitchRequest, opts ...grpc.CallOption) (SwitchServer_SetSwitchClient, error)
	CalibrateSwitch(ctx context.Context, in *SwitchCalibrationParameters, opts ...grpc.CallOption) (*empty.Empty, error)
}

type switchServerClient struct {
	cc *grpc.ClientConn
}

func NewSwitchServerClient(cc *grpc.ClientConn) SwitchServerClient {
	return &switchServerClient{cc}
}

func (c *switchServerClient) AddSwitch(ctx context.Context, in *AddSwitchRequest, opts ...grpc.CallOption) (*Switch, error) {
	out := new(Switch)
	err := c.cc.Invoke(ctx, "/SwitchServer/AddSwitch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *switchServerClient) RemoveSwitch(ctx context.Context, in *RemoveSwitchRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/SwitchServer/RemoveSwitch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *switchServerClient) GetSwitch(ctx context.Context, in *SwitchQuery, opts ...grpc.CallOption) (*Switch, error) {
	out := new(Switch)
	err := c.cc.Invoke(ctx, "/SwitchServer/GetSwitch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *switchServerClient) SetSwitch(ctx context.Context, in *SetSwitchRequest, opts ...grpc.CallOption) (SwitchServer_SetSwitchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_SwitchServer_serviceDesc.Streams[0], "/SwitchServer/SetSwitch", opts...)
	if err != nil {
		return nil, err
	}
	x := &switchServerSetSwitchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SwitchServer_SetSwitchClient interface {
	Recv() (*SetSwitchStatus, error)
	grpc.ClientStream
}

type switchServerSetSwitchClient struct {
	grpc.ClientStream
}

func (x *switchServerSetSwitchClient) Recv() (*SetSwitchStatus, error) {
	m := new(SetSwitchStatus)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *switchServerClient) CalibrateSwitch(ctx context.Context, in *SwitchCalibrationParameters, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/SwitchServer/CalibrateSwitch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SwitchServerServer is the server API for SwitchServer service.
type SwitchServerServer interface {
	AddSwitch(context.Context, *AddSwitchRequest) (*Switch, error)
	RemoveSwitch(context.Context, *RemoveSwitchRequest) (*empty.Empty, error)
	GetSwitch(context.Context, *SwitchQuery) (*Switch, error)
	SetSwitch(*SetSwitchRequest, SwitchServer_SetSwitchServer) error
	CalibrateSwitch(context.Context, *SwitchCalibrationParameters) (*empty.Empty, error)
}

func RegisterSwitchServerServer(s *grpc.Server, srv SwitchServerServer) {
	s.RegisterService(&_SwitchServer_serviceDesc, srv)
}

func _SwitchServer_AddSwitch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddSwitchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwitchServerServer).AddSwitch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SwitchServer/AddSwitch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwitchServerServer).AddSwitch(ctx, req.(*AddSwitchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwitchServer_RemoveSwitch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveSwitchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwitchServerServer).RemoveSwitch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SwitchServer/RemoveSwitch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwitchServerServer).RemoveSwitch(ctx, req.(*RemoveSwitchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwitchServer_GetSwitch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SwitchQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwitchServerServer).GetSwitch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SwitchServer/GetSwitch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwitchServerServer).GetSwitch(ctx, req.(*SwitchQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _SwitchServer_SetSwitch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SetSwitchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SwitchServerServer).SetSwitch(m, &switchServerSetSwitchServer{stream})
}

type SwitchServer_SetSwitchServer interface {
	Send(*SetSwitchStatus) error
	grpc.ServerStream
}

type switchServerSetSwitchServer struct {
	grpc.ServerStream
}

func (x *switchServerSetSwitchServer) Send(m *SetSwitchStatus) error {
	return x.ServerStream.SendMsg(m)
}

func _SwitchServer_CalibrateSwitch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SwitchCalibrationParameters)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwitchServerServer).CalibrateSwitch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SwitchServer/CalibrateSwitch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwitchServerServer).CalibrateSwitch(ctx, req.(*SwitchCalibrationParameters))
	}
	return interceptor(ctx, in, info, handler)
}

var _SwitchServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "SwitchServer",
	HandlerType: (*SwitchServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddSwitch",
			Handler:    _SwitchServer_AddSwitch_Handler,
		},
		{
			MethodName: "RemoveSwitch",
			Handler:    _SwitchServer_RemoveSwitch_Handler,
		},
		{
			MethodName: "GetSwitch",
			Handler:    _SwitchServer_GetSwitch_Handler,
		},
		{
			MethodName: "CalibrateSwitch",
			Handler:    _SwitchServer_CalibrateSwitch_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SetSwitch",
			Handler:       _SwitchServer_SetSwitch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "switchserver/switchserver.proto",
}

func init() {
	proto.RegisterFile("switchserver/switchserver.proto", fileDescriptor_switchserver_391a31fcd7eed6bf)
}

var fileDescriptor_switchserver_391a31fcd7eed6bf = []byte{
	// 513 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0x4f, 0x6f, 0xd3, 0x30,
	0x18, 0xc6, 0x9b, 0xf4, 0x6f, 0xde, 0x96, 0xae, 0x33, 0x13, 0x54, 0xe9, 0x80, 0xca, 0x02, 0xa9,
	0x5c, 0x5c, 0x54, 0xc6, 0x10, 0x42, 0x42, 0x2a, 0x50, 0x8d, 0x5e, 0x3a, 0x48, 0x0b, 0x9c, 0xd3,
	0xd5, 0x2d, 0x91, 0xd2, 0xb8, 0xd8, 0xce, 0xa6, 0x9d, 0xf8, 0x1a, 0x1c, 0x39, 0xf0, 0x41, 0x51,
	0xec, 0x24, 0xcd, 0x52, 0xe8, 0x38, 0x25, 0x8f, 0x1f, 0x3f, 0xef, 0xfb, 0xda, 0xf9, 0x05, 0x1e,
	0x89, 0x2b, 0x4f, 0x5e, 0x7c, 0x13, 0x94, 0x5f, 0x52, 0xde, 0xcf, 0x0a, 0xb2, 0xe1, 0x4c, 0x32,
	0xbb, 0xb3, 0x62, 0x6c, 0xe5, 0xd3, 0xbe, 0x52, 0xf3, 0x70, 0xd9, 0xa7, 0xeb, 0x8d, 0xbc, 0x8e,
	0xcd, 0x87, 0x79, 0xf3, 0x8a, 0xbb, 0x9b, 0x0d, 0xe5, 0x42, 0xfb, 0xf8, 0xb7, 0x01, 0x95, 0xa9,
	0xaa, 0x89, 0x9a, 0x60, 0x7a, 0x8b, 0xb6, 0xd1, 0x35, 0x7a, 0x96, 0x63, 0x7a, 0x0b, 0x84, 0xa0,
	0xe4, 0x89, 0xf3, 0xa0, 0x6d, 0x76, 0x8d, 0x5e, 0xcd, 0x51, 0xef, 0xa8, 0x0d, 0x55, 0x16, 0x0c,
	0x83, 0x95, 0x4f, 0xdb, 0xc5, 0xae, 0xd1, 0x2b, 0x3a, 0x89, 0x44, 0x36, 0xd4, 0xd8, 0x72, 0xa9,
	0xad, 0x92, 0xb2, 0x52, 0x8d, 0x8e, 0xc1, 0xe2, 0x54, 0x48, 0x6d, 0x96, 0x95, 0xb9, 0x5d, 0x40,
	0x18, 0x1a, 0x9e, 0x78, 0xe7, 0xfa, 0xde, 0x9c, 0xbb, 0x92, 0x2e, 0xda, 0x15, 0xd5, 0xef, 0xc6,
	0x1a, 0x3e, 0x85, 0xd6, 0x70, 0xb1, 0xd0, 0x83, 0x3a, 0xf4, 0x7b, 0x48, 0x85, 0xfc, 0x9f, 0x79,
	0xf1, 0x13, 0xb8, 0xeb, 0xd0, 0x35, 0xbb, 0xa4, 0x7b, 0xa3, 0xf8, 0x01, 0xd4, 0xf5, 0x86, 0x4f,
	0x21, 0xe5, 0xd7, 0x3b, 0xf6, 0x07, 0x68, 0x4d, 0xa9, 0xdc, 0xdf, 0xbd, 0x09, 0x26, 0x4b, 0x7a,
	0x9b, 0x2c, 0x40, 0x47, 0x50, 0x5e, 0x32, 0x7e, 0xa1, 0xef, 0xa9, 0xe6, 0x68, 0x81, 0x7f, 0xc0,
	0x41, 0x5a, 0x69, 0x2a, 0x5d, 0x19, 0x0a, 0xd4, 0x87, 0x8a, 0x50, 0x6f, 0xaa, 0x58, 0x73, 0x70,
	0x9f, 0xe4, 0x76, 0x10, 0xfd, 0x70, 0xe2, 0x6d, 0xf8, 0x35, 0x54, 0xe2, 0x68, 0x1d, 0xaa, 0xd3,
	0xd1, 0x6c, 0x36, 0x9e, 0x9c, 0xb5, 0x0a, 0x91, 0xf8, 0x3a, 0x1c, 0x2b, 0x61, 0xa0, 0x3b, 0x60,
	0x39, 0xa3, 0xd9, 0x67, 0x67, 0x12, 0x49, 0x13, 0xd5, 0xa0, 0xf4, 0xfe, 0x7c, 0x32, 0x6a, 0x15,
	0xf1, 0x4f, 0x13, 0x3a, 0xba, 0x78, 0x72, 0xbb, 0x1e, 0x0b, 0x3e, 0xba, 0xdc, 0x5d, 0x53, 0x49,
	0xb9, 0xd8, 0x39, 0xd6, 0x8b, 0xed, 0x07, 0x8f, 0xce, 0x56, 0x1f, 0x74, 0x88, 0x26, 0x8a, 0x24,
	0x44, 0x91, 0x71, 0x20, 0x4f, 0x4f, 0xbe, 0xb8, 0x7e, 0x48, 0xb7, 0x34, 0xbc, 0xcc, 0xd0, 0x50,
	0xbc, 0x3d, 0xb7, 0x45, 0xe5, 0x55, 0x16, 0x95, 0xd2, 0xed, 0xc9, 0x0c, 0x47, 0x6f, 0x72, 0x1c,
	0x95, 0x55, 0xda, 0xde, 0x49, 0xbf, 0x65, 0xcc, 0xd7, 0xe1, 0x1b, 0xfb, 0x07, 0xbf, 0x4c, 0x68,
	0xc4, 0xf7, 0xae, 0x7e, 0x2f, 0xf4, 0x14, 0xac, 0x14, 0x3a, 0x74, 0x48, 0xf2, 0x00, 0xda, 0x55,
	0xa2, 0x35, 0x2e, 0x44, 0xbd, 0xb3, 0x9c, 0xa1, 0x23, 0xf2, 0x17, 0xec, 0xec, 0x7b, 0x3b, 0xb3,
	0x8c, 0xa2, 0x5f, 0x15, 0x17, 0xd0, 0x63, 0xb0, 0xce, 0x92, 0xaf, 0x8e, 0x1a, 0x24, 0x03, 0x63,
	0xb6, 0xcb, 0x09, 0x58, 0x29, 0x1b, 0xe8, 0x90, 0xe4, 0x99, 0xb4, 0x5b, 0x79, 0x74, 0x70, 0xe1,
	0x99, 0x81, 0xc6, 0x70, 0x90, 0x9e, 0x32, 0xce, 0x1e, 0x93, 0x3d, 0x0c, 0xfc, 0x7b, 0xcc, 0x79,
	0x45, 0xad, 0x3c, 0xff, 0x13, 0x00, 0x00, 0xff, 0xff, 0x40, 0xd0, 0x29, 0x81, 0x94, 0x04, 0x00,
	0x00,
}
