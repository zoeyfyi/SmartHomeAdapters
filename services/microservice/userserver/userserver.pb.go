// Code generated by protoc-gen-go. DO NOT EDIT.
// source: userserver/userserver.proto

package userserver

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

type RegisterRequest struct {
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterRequest) Reset()         { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterRequest) ProtoMessage()    {}
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_userserver_0681c4be0fe3ba17, []int{0}
}
func (m *RegisterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRequest.Unmarshal(m, b)
}
func (m *RegisterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRequest.Marshal(b, m, deterministic)
}
func (dst *RegisterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRequest.Merge(dst, src)
}
func (m *RegisterRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterRequest.Size(m)
}
func (m *RegisterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRequest proto.InternalMessageInfo

func (m *RegisterRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RegisterRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *RegisterRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type Credentials struct {
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Credentials) Reset()         { *m = Credentials{} }
func (m *Credentials) String() string { return proto.CompactTextString(m) }
func (*Credentials) ProtoMessage()    {}
func (*Credentials) Descriptor() ([]byte, []int) {
	return fileDescriptor_userserver_0681c4be0fe3ba17, []int{1}
}
func (m *Credentials) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Credentials.Unmarshal(m, b)
}
func (m *Credentials) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Credentials.Marshal(b, m, deterministic)
}
func (dst *Credentials) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Credentials.Merge(dst, src)
}
func (m *Credentials) XXX_Size() int {
	return xxx_messageInfo_Credentials.Size(m)
}
func (m *Credentials) XXX_DiscardUnknown() {
	xxx_messageInfo_Credentials.DiscardUnknown(m)
}

var xxx_messageInfo_Credentials proto.InternalMessageInfo

func (m *Credentials) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *Credentials) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type Token struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Token) Reset()         { *m = Token{} }
func (m *Token) String() string { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()    {}
func (*Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_userserver_0681c4be0fe3ba17, []int{2}
}
func (m *Token) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Token.Unmarshal(m, b)
}
func (m *Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Token.Marshal(b, m, deterministic)
}
func (dst *Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Token.Merge(dst, src)
}
func (m *Token) XXX_Size() int {
	return xxx_messageInfo_Token.Size(m)
}
func (m *Token) XXX_DiscardUnknown() {
	xxx_messageInfo_Token.DiscardUnknown(m)
}

var xxx_messageInfo_Token proto.InternalMessageInfo

func (m *Token) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type User struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_userserver_0681c4be0fe3ba17, []int{3}
}
func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (dst *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(dst, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *User) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Email struct {
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Email) Reset()         { *m = Email{} }
func (m *Email) String() string { return proto.CompactTextString(m) }
func (*Email) ProtoMessage()    {}
func (*Email) Descriptor() ([]byte, []int) {
	return fileDescriptor_userserver_0681c4be0fe3ba17, []int{4}
}
func (m *Email) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Email.Unmarshal(m, b)
}
func (m *Email) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Email.Marshal(b, m, deterministic)
}
func (dst *Email) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Email.Merge(dst, src)
}
func (m *Email) XXX_Size() int {
	return xxx_messageInfo_Email.Size(m)
}
func (m *Email) XXX_DiscardUnknown() {
	xxx_messageInfo_Email.DiscardUnknown(m)
}

var xxx_messageInfo_Email proto.InternalMessageInfo

func (m *Email) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type UserId struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserId) Reset()         { *m = UserId{} }
func (m *UserId) String() string { return proto.CompactTextString(m) }
func (*UserId) ProtoMessage()    {}
func (*UserId) Descriptor() ([]byte, []int) {
	return fileDescriptor_userserver_0681c4be0fe3ba17, []int{5}
}
func (m *UserId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserId.Unmarshal(m, b)
}
func (m *UserId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserId.Marshal(b, m, deterministic)
}
func (dst *UserId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserId.Merge(dst, src)
}
func (m *UserId) XXX_Size() int {
	return xxx_messageInfo_UserId.Size(m)
}
func (m *UserId) XXX_DiscardUnknown() {
	xxx_messageInfo_UserId.DiscardUnknown(m)
}

var xxx_messageInfo_UserId proto.InternalMessageInfo

func (m *UserId) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*RegisterRequest)(nil), "RegisterRequest")
	proto.RegisterType((*Credentials)(nil), "Credentials")
	proto.RegisterType((*Token)(nil), "Token")
	proto.RegisterType((*User)(nil), "User")
	proto.RegisterType((*Email)(nil), "Email")
	proto.RegisterType((*UserId)(nil), "UserId")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// UserServerClient is the client API for UserServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UserServerClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	CheckCredentials(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*User, error)
	GetUserID(ctx context.Context, in *Email, opts ...grpc.CallOption) (*User, error)
	GetUserByID(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*User, error)
}

type userServerClient struct {
	cc *grpc.ClientConn
}

func NewUserServerClient(cc *grpc.ClientConn) UserServerClient {
	return &userServerClient{cc}
}

func (c *userServerClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/UserServer/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServerClient) CheckCredentials(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/UserServer/CheckCredentials", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServerClient) GetUserID(ctx context.Context, in *Email, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/UserServer/GetUserID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServerClient) GetUserByID(ctx context.Context, in *UserId, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/UserServer/GetUserByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServerServer is the server API for UserServer service.
type UserServerServer interface {
	Register(context.Context, *RegisterRequest) (*empty.Empty, error)
	CheckCredentials(context.Context, *Credentials) (*User, error)
	GetUserID(context.Context, *Email) (*User, error)
	GetUserByID(context.Context, *UserId) (*User, error)
}

func RegisterUserServerServer(s *grpc.Server, srv UserServerServer) {
	s.RegisterService(&_UserServer_serviceDesc, srv)
}

func _UserServer_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServerServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/UserServer/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServerServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserServer_CheckCredentials_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Credentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServerServer).CheckCredentials(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/UserServer/CheckCredentials",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServerServer).CheckCredentials(ctx, req.(*Credentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserServer_GetUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Email)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServerServer).GetUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/UserServer/GetUserID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServerServer).GetUserID(ctx, req.(*Email))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserServer_GetUserByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServerServer).GetUserByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/UserServer/GetUserByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServerServer).GetUserByID(ctx, req.(*UserId))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserServer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "UserServer",
	HandlerType: (*UserServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _UserServer_Register_Handler,
		},
		{
			MethodName: "CheckCredentials",
			Handler:    _UserServer_CheckCredentials_Handler,
		},
		{
			MethodName: "GetUserID",
			Handler:    _UserServer_GetUserID_Handler,
		},
		{
			MethodName: "GetUserByID",
			Handler:    _UserServer_GetUserByID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "userserver/userserver.proto",
}

func init() {
	proto.RegisterFile("userserver/userserver.proto", fileDescriptor_userserver_0681c4be0fe3ba17)
}

var fileDescriptor_userserver_0681c4be0fe3ba17 = []byte{
	// 296 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x90, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x9b, 0xd8, 0xc4, 0x76, 0x2a, 0x5a, 0x16, 0x91, 0x90, 0x2a, 0xca, 0x9e, 0xd4, 0xc3,
	0x06, 0x14, 0xbc, 0x0a, 0xb6, 0x45, 0x7a, 0x8d, 0x8a, 0xe7, 0xd4, 0x8c, 0x31, 0x34, 0xc9, 0xc6,
	0xdd, 0x8d, 0xd2, 0x3f, 0xe5, 0x6f, 0x94, 0xdd, 0x34, 0x49, 0x29, 0x9e, 0xbc, 0xcd, 0xdb, 0xb7,
	0xf3, 0xe6, 0xf1, 0xc1, 0xa4, 0x92, 0x28, 0x24, 0x8a, 0x2f, 0x14, 0x41, 0x37, 0xb2, 0x52, 0x70,
	0xc5, 0xfd, 0x49, 0xc2, 0x79, 0x92, 0x61, 0x60, 0xd4, 0xb2, 0x7a, 0x0f, 0x30, 0x2f, 0xd5, 0xba,
	0x36, 0xe9, 0x2b, 0x1c, 0x85, 0x98, 0xa4, 0x52, 0xa1, 0x08, 0xf1, 0xb3, 0x42, 0xa9, 0x08, 0x81,
	0x7e, 0x11, 0xe5, 0xe8, 0xed, 0x5d, 0x58, 0x97, 0xc3, 0xd0, 0xcc, 0xe4, 0x18, 0x1c, 0xcc, 0xa3,
	0x34, 0xf3, 0x2c, 0xf3, 0x58, 0x0b, 0xe2, 0xc3, 0xa0, 0x8c, 0xa4, 0xfc, 0xe6, 0x22, 0xf6, 0x6c,
	0x63, 0xb4, 0x9a, 0xde, 0xc3, 0x68, 0x2a, 0x30, 0xc6, 0x42, 0xa5, 0x51, 0x26, 0xff, 0x11, 0x70,
	0x06, 0xce, 0x33, 0x5f, 0x61, 0xa1, 0x57, 0x95, 0x1e, 0x9a, 0x55, 0x23, 0xe8, 0x35, 0xf4, 0x5f,
	0x24, 0x0a, 0x72, 0x08, 0x76, 0x1a, 0x6f, 0x2c, 0x3b, 0x8d, 0xdb, 0xf6, 0x76, 0xd7, 0x5e, 0x47,
	0xcd, 0xcd, 0xbd, 0x3f, 0x5b, 0x50, 0x0f, 0x5c, 0x1d, 0xb5, 0x88, 0x77, 0xc3, 0x6e, 0x7e, 0x2c,
	0x00, 0x6d, 0x3d, 0x19, 0x9e, 0xe4, 0x0e, 0x06, 0x0d, 0x2c, 0x32, 0x66, 0x3b, 0xdc, 0xfc, 0x13,
	0x56, 0x83, 0x66, 0x0d, 0x68, 0x36, 0xd7, 0xa0, 0x69, 0x8f, 0x5c, 0xc1, 0x78, 0xfa, 0x81, 0x6f,
	0xab, 0x6d, 0x20, 0x07, 0x6c, 0x4b, 0xf9, 0x0e, 0xd3, 0x67, 0x68, 0x8f, 0x9c, 0xc2, 0xf0, 0x11,
	0x95, 0xa9, 0x33, 0x23, 0x2e, 0x33, 0xb5, 0x3b, 0xf7, 0x1c, 0x46, 0x1b, 0xf7, 0x61, 0xbd, 0x98,
	0x91, 0x7d, 0x56, 0xf7, 0x6e, 0x3f, 0x2c, 0x5d, 0x73, 0xfb, 0xf6, 0x37, 0x00, 0x00, 0xff, 0xff,
	0x2f, 0x62, 0x49, 0xae, 0x11, 0x02, 0x00, 0x00,
}
