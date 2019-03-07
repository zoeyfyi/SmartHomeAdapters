// Code generated by protoc-gen-go. DO NOT EDIT.
// source: thermostatserver/thermostatserver.proto

package thermostatserver

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

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

type SetThermostatStatus_Status int32

const (
	SetThermostatStatus_SETTING SetThermostatStatus_Status = 0
	SetThermostatStatus_WAITING SetThermostatStatus_Status = 1
	SetThermostatStatus_DONE    SetThermostatStatus_Status = 2
)

var SetThermostatStatus_Status_name = map[int32]string{
	0: "SETTING",
	1: "WAITING",
	2: "DONE",
}
var SetThermostatStatus_Status_value = map[string]int32{
	"SETTING": 0,
	"WAITING": 1,
	"DONE":    2,
}

func (x SetThermostatStatus_Status) String() string {
	return proto.EnumName(SetThermostatStatus_Status_name, int32(x))
}
func (SetThermostatStatus_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_thermostatserver_43af3b9f2460dee7, []int{3, 0}
}

type Thermostat struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Tempreture           int64    `protobuf:"varint,2,opt,name=tempreture,proto3" json:"tempreture,omitempty"`
	MinAngle             int64    `protobuf:"varint,3,opt,name=minAngle,proto3" json:"minAngle,omitempty"`
	MaxAngle             int64    `protobuf:"varint,4,opt,name=maxAngle,proto3" json:"maxAngle,omitempty"`
	MinTempreture        int64    `protobuf:"varint,5,opt,name=minTempreture,proto3" json:"minTempreture,omitempty"`
	MaxTempreture        int64    `protobuf:"varint,6,opt,name=maxTempreture,proto3" json:"maxTempreture,omitempty"`
	IsCalibrated         bool     `protobuf:"varint,7,opt,name=isCalibrated,proto3" json:"isCalibrated,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Thermostat) Reset()         { *m = Thermostat{} }
func (m *Thermostat) String() string { return proto.CompactTextString(m) }
func (*Thermostat) ProtoMessage()    {}
func (*Thermostat) Descriptor() ([]byte, []int) {
	return fileDescriptor_thermostatserver_43af3b9f2460dee7, []int{0}
}
func (m *Thermostat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Thermostat.Unmarshal(m, b)
}
func (m *Thermostat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Thermostat.Marshal(b, m, deterministic)
}
func (dst *Thermostat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Thermostat.Merge(dst, src)
}
func (m *Thermostat) XXX_Size() int {
	return xxx_messageInfo_Thermostat.Size(m)
}
func (m *Thermostat) XXX_DiscardUnknown() {
	xxx_messageInfo_Thermostat.DiscardUnknown(m)
}

var xxx_messageInfo_Thermostat proto.InternalMessageInfo

func (m *Thermostat) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Thermostat) GetTempreture() int64 {
	if m != nil {
		return m.Tempreture
	}
	return 0
}

func (m *Thermostat) GetMinAngle() int64 {
	if m != nil {
		return m.MinAngle
	}
	return 0
}

func (m *Thermostat) GetMaxAngle() int64 {
	if m != nil {
		return m.MaxAngle
	}
	return 0
}

func (m *Thermostat) GetMinTempreture() int64 {
	if m != nil {
		return m.MinTempreture
	}
	return 0
}

func (m *Thermostat) GetMaxTempreture() int64 {
	if m != nil {
		return m.MaxTempreture
	}
	return 0
}

func (m *Thermostat) GetIsCalibrated() bool {
	if m != nil {
		return m.IsCalibrated
	}
	return false
}

type ThermostatQuery struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ThermostatQuery) Reset()         { *m = ThermostatQuery{} }
func (m *ThermostatQuery) String() string { return proto.CompactTextString(m) }
func (*ThermostatQuery) ProtoMessage()    {}
func (*ThermostatQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_thermostatserver_43af3b9f2460dee7, []int{1}
}
func (m *ThermostatQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ThermostatQuery.Unmarshal(m, b)
}
func (m *ThermostatQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ThermostatQuery.Marshal(b, m, deterministic)
}
func (dst *ThermostatQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ThermostatQuery.Merge(dst, src)
}
func (m *ThermostatQuery) XXX_Size() int {
	return xxx_messageInfo_ThermostatQuery.Size(m)
}
func (m *ThermostatQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_ThermostatQuery.DiscardUnknown(m)
}

var xxx_messageInfo_ThermostatQuery proto.InternalMessageInfo

func (m *ThermostatQuery) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type SetThermostatRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Tempreture           int64    `protobuf:"varint,2,opt,name=tempreture,proto3" json:"tempreture,omitempty"`
	Unit                 string   `protobuf:"bytes,3,opt,name=unit,proto3" json:"unit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetThermostatRequest) Reset()         { *m = SetThermostatRequest{} }
func (m *SetThermostatRequest) String() string { return proto.CompactTextString(m) }
func (*SetThermostatRequest) ProtoMessage()    {}
func (*SetThermostatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_thermostatserver_43af3b9f2460dee7, []int{2}
}
func (m *SetThermostatRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetThermostatRequest.Unmarshal(m, b)
}
func (m *SetThermostatRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetThermostatRequest.Marshal(b, m, deterministic)
}
func (dst *SetThermostatRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetThermostatRequest.Merge(dst, src)
}
func (m *SetThermostatRequest) XXX_Size() int {
	return xxx_messageInfo_SetThermostatRequest.Size(m)
}
func (m *SetThermostatRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetThermostatRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetThermostatRequest proto.InternalMessageInfo

func (m *SetThermostatRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *SetThermostatRequest) GetTempreture() int64 {
	if m != nil {
		return m.Tempreture
	}
	return 0
}

func (m *SetThermostatRequest) GetUnit() string {
	if m != nil {
		return m.Unit
	}
	return ""
}

type SetThermostatStatus struct {
	Status               SetThermostatStatus_Status `protobuf:"varint,1,opt,name=status,proto3,enum=SetThermostatStatus_Status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *SetThermostatStatus) Reset()         { *m = SetThermostatStatus{} }
func (m *SetThermostatStatus) String() string { return proto.CompactTextString(m) }
func (*SetThermostatStatus) ProtoMessage()    {}
func (*SetThermostatStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_thermostatserver_43af3b9f2460dee7, []int{3}
}
func (m *SetThermostatStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetThermostatStatus.Unmarshal(m, b)
}
func (m *SetThermostatStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetThermostatStatus.Marshal(b, m, deterministic)
}
func (dst *SetThermostatStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetThermostatStatus.Merge(dst, src)
}
func (m *SetThermostatStatus) XXX_Size() int {
	return xxx_messageInfo_SetThermostatStatus.Size(m)
}
func (m *SetThermostatStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_SetThermostatStatus.DiscardUnknown(m)
}

var xxx_messageInfo_SetThermostatStatus proto.InternalMessageInfo

func (m *SetThermostatStatus) GetStatus() SetThermostatStatus_Status {
	if m != nil {
		return m.Status
	}
	return SetThermostatStatus_SETTING
}

func init() {
	proto.RegisterType((*Thermostat)(nil), "Thermostat")
	proto.RegisterType((*ThermostatQuery)(nil), "ThermostatQuery")
	proto.RegisterType((*SetThermostatRequest)(nil), "SetThermostatRequest")
	proto.RegisterType((*SetThermostatStatus)(nil), "SetThermostatStatus")
	proto.RegisterEnum("SetThermostatStatus_Status", SetThermostatStatus_Status_name, SetThermostatStatus_Status_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ThermostatServerClient is the client API for ThermostatServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ThermostatServerClient interface {
	GetThermostat(ctx context.Context, in *ThermostatQuery, opts ...grpc.CallOption) (*Thermostat, error)
	SetThermostat(ctx context.Context, in *SetThermostatRequest, opts ...grpc.CallOption) (ThermostatServer_SetThermostatClient, error)
}

type thermostatServerClient struct {
	cc *grpc.ClientConn
}

func NewThermostatServerClient(cc *grpc.ClientConn) ThermostatServerClient {
	return &thermostatServerClient{cc}
}

func (c *thermostatServerClient) GetThermostat(ctx context.Context, in *ThermostatQuery, opts ...grpc.CallOption) (*Thermostat, error) {
	out := new(Thermostat)
	err := c.cc.Invoke(ctx, "/ThermostatServer/GetThermostat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *thermostatServerClient) SetThermostat(ctx context.Context, in *SetThermostatRequest, opts ...grpc.CallOption) (ThermostatServer_SetThermostatClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ThermostatServer_serviceDesc.Streams[0], "/ThermostatServer/SetThermostat", opts...)
	if err != nil {
		return nil, err
	}
	x := &thermostatServerSetThermostatClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ThermostatServer_SetThermostatClient interface {
	Recv() (*SetThermostatStatus, error)
	grpc.ClientStream
}

type thermostatServerSetThermostatClient struct {
	grpc.ClientStream
}

func (x *thermostatServerSetThermostatClient) Recv() (*SetThermostatStatus, error) {
	m := new(SetThermostatStatus)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ThermostatServerServer is the server API for ThermostatServer service.
type ThermostatServerServer interface {
	GetThermostat(context.Context, *ThermostatQuery) (*Thermostat, error)
	SetThermostat(*SetThermostatRequest, ThermostatServer_SetThermostatServer) error
}

func RegisterThermostatServerServer(s *grpc.Server, srv ThermostatServerServer) {
	s.RegisterService(&_ThermostatServer_serviceDesc, srv)
}

func _ThermostatServer_GetThermostat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ThermostatQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThermostatServerServer).GetThermostat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ThermostatServer/GetThermostat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThermostatServerServer).GetThermostat(ctx, req.(*ThermostatQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _ThermostatServer_SetThermostat_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SetThermostatRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ThermostatServerServer).SetThermostat(m, &thermostatServerSetThermostatServer{stream})
}

type ThermostatServer_SetThermostatServer interface {
	Send(*SetThermostatStatus) error
	grpc.ServerStream
}

type thermostatServerSetThermostatServer struct {
	grpc.ServerStream
}

func (x *thermostatServerSetThermostatServer) Send(m *SetThermostatStatus) error {
	return x.ServerStream.SendMsg(m)
}

var _ThermostatServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ThermostatServer",
	HandlerType: (*ThermostatServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetThermostat",
			Handler:    _ThermostatServer_GetThermostat_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SetThermostat",
			Handler:       _ThermostatServer_SetThermostat_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "thermostatserver/thermostatserver.proto",
}

func init() {
	proto.RegisterFile("thermostatserver/thermostatserver.proto", fileDescriptor_thermostatserver_43af3b9f2460dee7)
}

var fileDescriptor_thermostatserver_43af3b9f2460dee7 = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0xcd, 0x6a, 0xf2, 0x40,
	0x14, 0x75, 0xa2, 0x5f, 0xd4, 0xeb, 0xa7, 0x0d, 0xb7, 0x16, 0x82, 0x85, 0x62, 0x87, 0x42, 0x5d,
	0x94, 0x54, 0xf4, 0x05, 0x2a, 0xad, 0x88, 0x1b, 0x4b, 0x63, 0xa0, 0xd0, 0x5d, 0xc4, 0x4b, 0x3b,
	0x60, 0xa2, 0x9d, 0x4c, 0x8a, 0x7d, 0x81, 0x3e, 0x69, 0x1f, 0xa4, 0x74, 0x0c, 0x3a, 0x91, 0x6c,
	0xba, 0xca, 0x3d, 0x3f, 0x39, 0xc9, 0x3d, 0x33, 0x70, 0xad, 0xde, 0x48, 0x46, 0xeb, 0x44, 0x85,
	0x2a, 0x21, 0xf9, 0x41, 0xf2, 0xf6, 0x98, 0xf0, 0x36, 0x72, 0xad, 0xd6, 0xfc, 0x9b, 0x01, 0x04,
	0x7b, 0x09, 0x5b, 0x60, 0x89, 0xa5, 0xcb, 0xba, 0xac, 0x57, 0xf7, 0x2d, 0xb1, 0xc4, 0x0b, 0x00,
	0x45, 0xd1, 0x46, 0x92, 0x4a, 0x25, 0xb9, 0x56, 0x97, 0xf5, 0xca, 0xbe, 0xc1, 0x60, 0x07, 0x6a,
	0x91, 0x88, 0x47, 0xf1, 0xeb, 0x8a, 0xdc, 0xb2, 0x56, 0xf7, 0x58, 0x6b, 0xe1, 0x76, 0xa7, 0x55,
	0x32, 0x2d, 0xc3, 0x78, 0x05, 0xcd, 0x48, 0xc4, 0xc1, 0x21, 0xfa, 0x9f, 0x36, 0xe4, 0x49, 0xed,
	0x0a, 0xb7, 0x86, 0xcb, 0xce, 0x5c, 0x26, 0x89, 0x1c, 0xfe, 0x8b, 0xe4, 0x3e, 0x5c, 0x89, 0x85,
	0x0c, 0x15, 0x2d, 0xdd, 0x6a, 0x97, 0xf5, 0x6a, 0x7e, 0x8e, 0xe3, 0x97, 0x70, 0x72, 0xd8, 0xf2,
	0x29, 0x25, 0xf9, 0x79, 0xbc, 0x2a, 0x7f, 0x81, 0xf6, 0x9c, 0xd4, 0xc1, 0xe5, 0xd3, 0x7b, 0x4a,
	0xc9, 0xdf, 0x2b, 0x41, 0xa8, 0xa4, 0xb1, 0x50, 0xba, 0x8e, 0xba, 0xaf, 0x67, 0xbe, 0x85, 0xd3,
	0x5c, 0xf6, 0x5c, 0x85, 0x2a, 0x4d, 0x70, 0x08, 0x76, 0xa2, 0x27, 0x1d, 0xdf, 0x1a, 0x9c, 0x7b,
	0x05, 0x2e, 0x6f, 0xf7, 0xf0, 0x33, 0x2b, 0xbf, 0x01, 0x3b, 0x7b, 0xbd, 0x01, 0xd5, 0xf9, 0x38,
	0x08, 0xa6, 0xb3, 0x89, 0x53, 0xfa, 0x05, 0xcf, 0xa3, 0xa9, 0x06, 0x0c, 0x6b, 0x50, 0x79, 0x78,
	0x9c, 0x8d, 0x1d, 0x6b, 0xf0, 0xc5, 0xc0, 0x31, 0x12, 0xf5, 0xd1, 0x63, 0x1f, 0x9a, 0x13, 0xf3,
	0x43, 0xe8, 0x78, 0x47, 0xed, 0x74, 0x1a, 0x06, 0xc3, 0x4b, 0x78, 0x07, 0xcd, 0xdc, 0xaf, 0xe1,
	0x99, 0x57, 0x54, 0x56, 0xa7, 0x5d, 0xb4, 0x01, 0x2f, 0xf5, 0xd9, 0xc2, 0xd6, 0xf7, 0x6d, 0xf8,
	0x13, 0x00, 0x00, 0xff, 0xff, 0x16, 0x90, 0x9f, 0x7e, 0x9a, 0x02, 0x00, 0x00,
}
