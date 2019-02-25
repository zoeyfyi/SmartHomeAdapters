// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infoserver/infoserver.proto

package infoserver

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import empty "github.com/golang/protobuf/ptypes/empty"

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

type RobotQuery struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RobotQuery) Reset()         { *m = RobotQuery{} }
func (m *RobotQuery) String() string { return proto.CompactTextString(m) }
func (*RobotQuery) ProtoMessage()    {}
func (*RobotQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_infoserver_9103bbc35f51024b, []int{0}
}
func (m *RobotQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RobotQuery.Unmarshal(m, b)
}
func (m *RobotQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RobotQuery.Marshal(b, m, deterministic)
}
func (dst *RobotQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RobotQuery.Merge(dst, src)
}
func (m *RobotQuery) XXX_Size() int {
	return xxx_messageInfo_RobotQuery.Size(m)
}
func (m *RobotQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_RobotQuery.DiscardUnknown(m)
}

var xxx_messageInfo_RobotQuery proto.InternalMessageInfo

func (m *RobotQuery) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type ToggleRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Value                bool     `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	Force                bool     `protobuf:"varint,3,opt,name=force,proto3" json:"force,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ToggleRequest) Reset()         { *m = ToggleRequest{} }
func (m *ToggleRequest) String() string { return proto.CompactTextString(m) }
func (*ToggleRequest) ProtoMessage()    {}
func (*ToggleRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_infoserver_9103bbc35f51024b, []int{1}
}
func (m *ToggleRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ToggleRequest.Unmarshal(m, b)
}
func (m *ToggleRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ToggleRequest.Marshal(b, m, deterministic)
}
func (dst *ToggleRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ToggleRequest.Merge(dst, src)
}
func (m *ToggleRequest) XXX_Size() int {
	return xxx_messageInfo_ToggleRequest.Size(m)
}
func (m *ToggleRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ToggleRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ToggleRequest proto.InternalMessageInfo

func (m *ToggleRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ToggleRequest) GetValue() bool {
	if m != nil {
		return m.Value
	}
	return false
}

func (m *ToggleRequest) GetForce() bool {
	if m != nil {
		return m.Force
	}
	return false
}

type RangeRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Value                int64    `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RangeRequest) Reset()         { *m = RangeRequest{} }
func (m *RangeRequest) String() string { return proto.CompactTextString(m) }
func (*RangeRequest) ProtoMessage()    {}
func (*RangeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_infoserver_9103bbc35f51024b, []int{2}
}
func (m *RangeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RangeRequest.Unmarshal(m, b)
}
func (m *RangeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RangeRequest.Marshal(b, m, deterministic)
}
func (dst *RangeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RangeRequest.Merge(dst, src)
}
func (m *RangeRequest) XXX_Size() int {
	return xxx_messageInfo_RangeRequest.Size(m)
}
func (m *RangeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RangeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RangeRequest proto.InternalMessageInfo

func (m *RangeRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *RangeRequest) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type Robot struct {
	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Nickname      string `protobuf:"bytes,2,opt,name=nickname,proto3" json:"nickname,omitempty"`
	RobotType     string `protobuf:"bytes,3,opt,name=robotType,proto3" json:"robotType,omitempty"`
	InterfaceType string `protobuf:"bytes,4,opt,name=interfaceType,proto3" json:"interfaceType,omitempty"`
	// Types that are valid to be assigned to RobotStatus:
	//	*Robot_ToggleStatus
	//	*Robot_RangeStatus
	RobotStatus          isRobot_RobotStatus `protobuf_oneof:"robotStatus"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Robot) Reset()         { *m = Robot{} }
func (m *Robot) String() string { return proto.CompactTextString(m) }
func (*Robot) ProtoMessage()    {}
func (*Robot) Descriptor() ([]byte, []int) {
	return fileDescriptor_infoserver_9103bbc35f51024b, []int{3}
}
func (m *Robot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Robot.Unmarshal(m, b)
}
func (m *Robot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Robot.Marshal(b, m, deterministic)
}
func (dst *Robot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Robot.Merge(dst, src)
}
func (m *Robot) XXX_Size() int {
	return xxx_messageInfo_Robot.Size(m)
}
func (m *Robot) XXX_DiscardUnknown() {
	xxx_messageInfo_Robot.DiscardUnknown(m)
}

var xxx_messageInfo_Robot proto.InternalMessageInfo

func (m *Robot) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Robot) GetNickname() string {
	if m != nil {
		return m.Nickname
	}
	return ""
}

func (m *Robot) GetRobotType() string {
	if m != nil {
		return m.RobotType
	}
	return ""
}

func (m *Robot) GetInterfaceType() string {
	if m != nil {
		return m.InterfaceType
	}
	return ""
}

type isRobot_RobotStatus interface {
	isRobot_RobotStatus()
}

type Robot_ToggleStatus struct {
	ToggleStatus *ToggleStatus `protobuf:"bytes,5,opt,name=toggleStatus,proto3,oneof"`
}

type Robot_RangeStatus struct {
	RangeStatus *RangeStatus `protobuf:"bytes,6,opt,name=rangeStatus,proto3,oneof"`
}

func (*Robot_ToggleStatus) isRobot_RobotStatus() {}

func (*Robot_RangeStatus) isRobot_RobotStatus() {}

func (m *Robot) GetRobotStatus() isRobot_RobotStatus {
	if m != nil {
		return m.RobotStatus
	}
	return nil
}

func (m *Robot) GetToggleStatus() *ToggleStatus {
	if x, ok := m.GetRobotStatus().(*Robot_ToggleStatus); ok {
		return x.ToggleStatus
	}
	return nil
}

func (m *Robot) GetRangeStatus() *RangeStatus {
	if x, ok := m.GetRobotStatus().(*Robot_RangeStatus); ok {
		return x.RangeStatus
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Robot) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Robot_OneofMarshaler, _Robot_OneofUnmarshaler, _Robot_OneofSizer, []interface{}{
		(*Robot_ToggleStatus)(nil),
		(*Robot_RangeStatus)(nil),
	}
}

func _Robot_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Robot)
	// robotStatus
	switch x := m.RobotStatus.(type) {
	case *Robot_ToggleStatus:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ToggleStatus); err != nil {
			return err
		}
	case *Robot_RangeStatus:
		b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.RangeStatus); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Robot.RobotStatus has unexpected type %T", x)
	}
	return nil
}

func _Robot_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Robot)
	switch tag {
	case 5: // robotStatus.toggleStatus
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ToggleStatus)
		err := b.DecodeMessage(msg)
		m.RobotStatus = &Robot_ToggleStatus{msg}
		return true, err
	case 6: // robotStatus.rangeStatus
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RangeStatus)
		err := b.DecodeMessage(msg)
		m.RobotStatus = &Robot_RangeStatus{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Robot_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Robot)
	// robotStatus
	switch x := m.RobotStatus.(type) {
	case *Robot_ToggleStatus:
		s := proto.Size(x.ToggleStatus)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Robot_RangeStatus:
		s := proto.Size(x.RangeStatus)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type ToggleStatus struct {
	Value                bool     `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ToggleStatus) Reset()         { *m = ToggleStatus{} }
func (m *ToggleStatus) String() string { return proto.CompactTextString(m) }
func (*ToggleStatus) ProtoMessage()    {}
func (*ToggleStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_infoserver_9103bbc35f51024b, []int{4}
}
func (m *ToggleStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ToggleStatus.Unmarshal(m, b)
}
func (m *ToggleStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ToggleStatus.Marshal(b, m, deterministic)
}
func (dst *ToggleStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ToggleStatus.Merge(dst, src)
}
func (m *ToggleStatus) XXX_Size() int {
	return xxx_messageInfo_ToggleStatus.Size(m)
}
func (m *ToggleStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ToggleStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ToggleStatus proto.InternalMessageInfo

func (m *ToggleStatus) GetValue() bool {
	if m != nil {
		return m.Value
	}
	return false
}

type RangeStatus struct {
	Max                  int64    `protobuf:"varint,1,opt,name=max,proto3" json:"max,omitempty"`
	Min                  int64    `protobuf:"varint,2,opt,name=min,proto3" json:"min,omitempty"`
	Current              int64    `protobuf:"varint,3,opt,name=current,proto3" json:"current,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RangeStatus) Reset()         { *m = RangeStatus{} }
func (m *RangeStatus) String() string { return proto.CompactTextString(m) }
func (*RangeStatus) ProtoMessage()    {}
func (*RangeStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_infoserver_9103bbc35f51024b, []int{5}
}
func (m *RangeStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RangeStatus.Unmarshal(m, b)
}
func (m *RangeStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RangeStatus.Marshal(b, m, deterministic)
}
func (dst *RangeStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RangeStatus.Merge(dst, src)
}
func (m *RangeStatus) XXX_Size() int {
	return xxx_messageInfo_RangeStatus.Size(m)
}
func (m *RangeStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_RangeStatus.DiscardUnknown(m)
}

var xxx_messageInfo_RangeStatus proto.InternalMessageInfo

func (m *RangeStatus) GetMax() int64 {
	if m != nil {
		return m.Max
	}
	return 0
}

func (m *RangeStatus) GetMin() int64 {
	if m != nil {
		return m.Min
	}
	return 0
}

func (m *RangeStatus) GetCurrent() int64 {
	if m != nil {
		return m.Current
	}
	return 0
}

func init() {
	proto.RegisterType((*RobotQuery)(nil), "RobotQuery")
	proto.RegisterType((*ToggleRequest)(nil), "ToggleRequest")
	proto.RegisterType((*RangeRequest)(nil), "RangeRequest")
	proto.RegisterType((*Robot)(nil), "Robot")
	proto.RegisterType((*ToggleStatus)(nil), "ToggleStatus")
	proto.RegisterType((*RangeStatus)(nil), "RangeStatus")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// InfoServerClient is the client API for InfoServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type InfoServerClient interface {
	GetRobot(ctx context.Context, in *RobotQuery, opts ...grpc.CallOption) (*Robot, error)
	GetRobots(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (InfoServer_GetRobotsClient, error)
	ToggleRobot(ctx context.Context, in *ToggleRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	RangeRobot(ctx context.Context, in *RangeRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type infoServerClient struct {
	cc *grpc.ClientConn
}

func NewInfoServerClient(cc *grpc.ClientConn) InfoServerClient {
	return &infoServerClient{cc}
}

func (c *infoServerClient) GetRobot(ctx context.Context, in *RobotQuery, opts ...grpc.CallOption) (*Robot, error) {
	out := new(Robot)
	err := c.cc.Invoke(ctx, "/InfoServer/GetRobot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *infoServerClient) GetRobots(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (InfoServer_GetRobotsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_InfoServer_serviceDesc.Streams[0], "/InfoServer/GetRobots", opts...)
	if err != nil {
		return nil, err
	}
	x := &infoServerGetRobotsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type InfoServer_GetRobotsClient interface {
	Recv() (*Robot, error)
	grpc.ClientStream
}

type infoServerGetRobotsClient struct {
	grpc.ClientStream
}

func (x *infoServerGetRobotsClient) Recv() (*Robot, error) {
	m := new(Robot)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *infoServerClient) ToggleRobot(ctx context.Context, in *ToggleRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/InfoServer/ToggleRobot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *infoServerClient) RangeRobot(ctx context.Context, in *RangeRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/InfoServer/RangeRobot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InfoServerServer is the server API for InfoServer service.
type InfoServerServer interface {
	GetRobot(context.Context, *RobotQuery) (*Robot, error)
	GetRobots(*empty.Empty, InfoServer_GetRobotsServer) error
	ToggleRobot(context.Context, *ToggleRequest) (*empty.Empty, error)
	RangeRobot(context.Context, *RangeRequest) (*empty.Empty, error)
}

func RegisterInfoServerServer(s *grpc.Server, srv InfoServerServer) {
	s.RegisterService(&_InfoServer_serviceDesc, srv)
}

func _InfoServer_GetRobot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RobotQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InfoServerServer).GetRobot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/InfoServer/GetRobot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InfoServerServer).GetRobot(ctx, req.(*RobotQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _InfoServer_GetRobots_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(empty.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(InfoServerServer).GetRobots(m, &infoServerGetRobotsServer{stream})
}

type InfoServer_GetRobotsServer interface {
	Send(*Robot) error
	grpc.ServerStream
}

type infoServerGetRobotsServer struct {
	grpc.ServerStream
}

func (x *infoServerGetRobotsServer) Send(m *Robot) error {
	return x.ServerStream.SendMsg(m)
}

func _InfoServer_ToggleRobot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ToggleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InfoServerServer).ToggleRobot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/InfoServer/ToggleRobot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InfoServerServer).ToggleRobot(ctx, req.(*ToggleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _InfoServer_RangeRobot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InfoServerServer).RangeRobot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/InfoServer/RangeRobot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InfoServerServer).RangeRobot(ctx, req.(*RangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _InfoServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "InfoServer",
	HandlerType: (*InfoServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRobot",
			Handler:    _InfoServer_GetRobot_Handler,
		},
		{
			MethodName: "ToggleRobot",
			Handler:    _InfoServer_ToggleRobot_Handler,
		},
		{
			MethodName: "RangeRobot",
			Handler:    _InfoServer_RangeRobot_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetRobots",
			Handler:       _InfoServer_GetRobots_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "infoserver/infoserver.proto",
}

func init() {
	proto.RegisterFile("infoserver/infoserver.proto", fileDescriptor_infoserver_9103bbc35f51024b)
}

var fileDescriptor_infoserver_9103bbc35f51024b = []byte{
	// 398 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xd1, 0xce, 0xd2, 0x30,
	0x18, 0x5d, 0x99, 0x20, 0xfb, 0xb6, 0x19, 0xd3, 0x18, 0xb3, 0xec, 0xe7, 0x02, 0x1b, 0x2e, 0xb8,
	0xea, 0x08, 0x68, 0xbc, 0x37, 0x31, 0x62, 0xb8, 0xb2, 0xf0, 0x02, 0x63, 0x74, 0xcb, 0x22, 0xb4,
	0xd8, 0x75, 0x44, 0xde, 0xd2, 0x17, 0xf1, 0x1d, 0xcc, 0xda, 0xe1, 0x86, 0x86, 0xe4, 0xbf, 0xfb,
	0xce, 0xe1, 0x9c, 0x8f, 0xaf, 0x67, 0x07, 0x9e, 0x4a, 0x91, 0xcb, 0x8a, 0xab, 0x0b, 0x57, 0x49,
	0x37, 0xd2, 0xb3, 0x92, 0x5a, 0xc6, 0x4f, 0x85, 0x94, 0xc5, 0x91, 0x27, 0x06, 0xed, 0xeb, 0x3c,
	0xe1, 0xa7, 0xb3, 0xbe, 0xda, 0x1f, 0xc9, 0x04, 0x80, 0xc9, 0xbd, 0xd4, 0xdf, 0x6a, 0xae, 0xae,
	0xf8, 0x15, 0x0c, 0xca, 0x43, 0x84, 0xa6, 0x68, 0xee, 0xb1, 0x41, 0x79, 0x20, 0x1b, 0x08, 0x77,
	0xb2, 0x28, 0x8e, 0x9c, 0xf1, 0x1f, 0x35, 0xaf, 0xf4, 0xbf, 0x02, 0xfc, 0x06, 0x86, 0x97, 0xf4,
	0x58, 0xf3, 0x68, 0x30, 0x45, 0xf3, 0x31, 0xb3, 0xa0, 0x61, 0x73, 0xa9, 0x32, 0x1e, 0xb9, 0x96,
	0x35, 0x80, 0xbc, 0x87, 0x80, 0xa5, 0xa2, 0x78, 0xde, 0x2e, 0xb7, 0xdd, 0x45, 0x7e, 0x23, 0x18,
	0x9a, 0x0b, 0xff, 0xd3, 0xc7, 0x30, 0x16, 0x65, 0xf6, 0x5d, 0xa4, 0x27, 0x6b, 0xf1, 0xd8, 0x5f,
	0x8c, 0x27, 0xe0, 0xa9, 0xc6, 0xb4, 0xbb, 0x9e, 0xed, 0x15, 0x1e, 0xeb, 0x08, 0x3c, 0x83, 0xb0,
	0x14, 0x9a, 0xab, 0x3c, 0xcd, 0xb8, 0x51, 0xbc, 0x30, 0x8a, 0x7b, 0x12, 0xaf, 0x20, 0xd0, 0xe6,
	0xf1, 0x5b, 0x9d, 0xea, 0xba, 0x8a, 0x86, 0x53, 0x34, 0xf7, 0x97, 0x21, 0xdd, 0xf5, 0xc8, 0xb5,
	0xc3, 0xee, 0x44, 0x78, 0x01, 0xbe, 0x6a, 0x1e, 0xd9, 0x7a, 0x46, 0xc6, 0x13, 0x50, 0xd6, 0x71,
	0x6b, 0x87, 0xf5, 0x25, 0x9f, 0x42, 0xf0, 0xcd, 0x65, 0x16, 0x92, 0x19, 0x04, 0xfd, 0x3f, 0xe8,
	0x52, 0x41, 0xbd, 0x84, 0xc9, 0x06, 0xfc, 0xde, 0x4a, 0xfc, 0x1a, 0xdc, 0x53, 0xfa, 0xd3, 0x48,
	0x5c, 0xd6, 0x8c, 0x86, 0x29, 0x45, 0x1b, 0x65, 0x33, 0xe2, 0x08, 0x5e, 0x66, 0xb5, 0x52, 0x5c,
	0x68, 0x13, 0x88, 0xcb, 0x6e, 0x70, 0xf9, 0x0b, 0x01, 0x7c, 0x15, 0xb9, 0xdc, 0x9a, 0xd6, 0xe0,
	0x77, 0x30, 0xfe, 0xc2, 0xb5, 0xcd, 0xdc, 0xa7, 0x5d, 0x3b, 0xe2, 0x91, 0x05, 0xc4, 0xc1, 0x09,
	0x78, 0x37, 0x49, 0x85, 0xdf, 0x52, 0x5b, 0x30, 0x7a, 0x2b, 0x18, 0xfd, 0xdc, 0x14, 0xac, 0x93,
	0x2f, 0x10, 0xfe, 0x08, 0x7e, 0x5b, 0x24, 0xfb, 0x29, 0xe9, 0x5d, 0xad, 0xe2, 0x07, 0x2b, 0x88,
	0x83, 0x3f, 0x00, 0xd8, 0xd2, 0x18, 0x5f, 0x48, 0xfb, 0x0d, 0x7a, 0x6c, 0xdb, 0x8f, 0x0c, 0xb3,
	0xfa, 0x13, 0x00, 0x00, 0xff, 0xff, 0xfe, 0x89, 0x76, 0xd2, 0x19, 0x03, 0x00, 0x00,
}
