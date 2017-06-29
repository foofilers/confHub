// Code generated by protoc-gen-go. DO NOT EDIT.
// source: users/users.proto

/*
Package users is a generated protocol buffer package.

It is generated from these files:
	users/users.proto

It has these top-level messages:
	DeleteRequest
	UserListRequest
	Permission
	User
*/
package users

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

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

type DeleteRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
}

func (m *DeleteRequest) Reset()                    { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string            { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()               {}
func (*DeleteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DeleteRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type UserListRequest struct {
	Query string `protobuf:"bytes,1,opt,name=query" json:"query,omitempty"`
	Order string `protobuf:"bytes,2,opt,name=order" json:"order,omitempty"`
	Page  int32  `protobuf:"varint,3,opt,name=page" json:"page,omitempty"`
	Count int32  `protobuf:"varint,4,opt,name=count" json:"count,omitempty"`
}

func (m *UserListRequest) Reset()                    { *m = UserListRequest{} }
func (m *UserListRequest) String() string            { return proto.CompactTextString(m) }
func (*UserListRequest) ProtoMessage()               {}
func (*UserListRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *UserListRequest) GetQuery() string {
	if m != nil {
		return m.Query
	}
	return ""
}

func (m *UserListRequest) GetOrder() string {
	if m != nil {
		return m.Order
	}
	return ""
}

func (m *UserListRequest) GetPage() int32 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *UserListRequest) GetCount() int32 {
	if m != nil {
		return m.Count
	}
	return 0
}

type Permission struct {
	Application string   `protobuf:"bytes,1,opt,name=application" json:"application,omitempty"`
	Perm        []string `protobuf:"bytes,2,rep,name=perm" json:"perm,omitempty"`
}

func (m *Permission) Reset()                    { *m = Permission{} }
func (m *Permission) String() string            { return proto.CompactTextString(m) }
func (*Permission) ProtoMessage()               {}
func (*Permission) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Permission) GetApplication() string {
	if m != nil {
		return m.Application
	}
	return ""
}

func (m *Permission) GetPerm() []string {
	if m != nil {
		return m.Perm
	}
	return nil
}

type User struct {
	Id          string        `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Username    string        `protobuf:"bytes,2,opt,name=username" json:"username,omitempty"`
	Email       string        `protobuf:"bytes,3,opt,name=email" json:"email,omitempty"`
	Permissions []*Permission `protobuf:"bytes,4,rep,name=permissions" json:"permissions,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *User) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *User) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetPermissions() []*Permission {
	if m != nil {
		return m.Permissions
	}
	return nil
}

func init() {
	proto.RegisterType((*DeleteRequest)(nil), "users.DeleteRequest")
	proto.RegisterType((*UserListRequest)(nil), "users.UserListRequest")
	proto.RegisterType((*Permission)(nil), "users.Permission")
	proto.RegisterType((*User)(nil), "users.User")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Users service

type UsersClient interface {
	List(ctx context.Context, in *UserListRequest, opts ...grpc.CallOption) (Users_ListClient, error)
	Add(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
}

type usersClient struct {
	cc *grpc.ClientConn
}

func NewUsersClient(cc *grpc.ClientConn) UsersClient {
	return &usersClient{cc}
}

func (c *usersClient) List(ctx context.Context, in *UserListRequest, opts ...grpc.CallOption) (Users_ListClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Users_serviceDesc.Streams[0], c.cc, "/users.Users/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &usersListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Users_ListClient interface {
	Recv() (*User, error)
	grpc.ClientStream
}

type usersListClient struct {
	grpc.ClientStream
}

func (x *usersListClient) Recv() (*User, error) {
	m := new(User)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *usersClient) Add(ctx context.Context, in *User, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := grpc.Invoke(ctx, "/users.Users/Add", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/users.Users/Delete", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Users service

type UsersServer interface {
	List(*UserListRequest, Users_ListServer) error
	Add(context.Context, *User) (*User, error)
	Delete(context.Context, *DeleteRequest) (*google_protobuf.Empty, error)
}

func RegisterUsersServer(s *grpc.Server, srv UsersServer) {
	s.RegisterService(&_Users_serviceDesc, srv)
}

func _Users_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(UserListRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UsersServer).List(m, &usersListServer{stream})
}

type Users_ListServer interface {
	Send(*User) error
	grpc.ServerStream
}

type usersListServer struct {
	grpc.ServerStream
}

func (x *usersListServer) Send(m *User) error {
	return x.ServerStream.SendMsg(m)
}

func _Users_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.Users/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).Add(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.Users/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Users_serviceDesc = grpc.ServiceDesc{
	ServiceName: "users.Users",
	HandlerType: (*UsersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _Users_Add_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Users_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "List",
			Handler:       _Users_List_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "users/users.proto",
}

func init() { proto.RegisterFile("users/users.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 352 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0x5f, 0x6b, 0xa3, 0x40,
	0x14, 0xc5, 0xf1, 0x5f, 0xd8, 0x5c, 0xd9, 0x2c, 0x19, 0x42, 0x10, 0xf7, 0x21, 0x22, 0x2c, 0xf8,
	0xb2, 0x66, 0x49, 0x60, 0xdf, 0x9b, 0xb6, 0x6f, 0x7d, 0x08, 0x96, 0x7e, 0x00, 0xa3, 0x57, 0x33,
	0xa0, 0x8e, 0x99, 0xd1, 0x42, 0x1e, 0xfa, 0x19, 0xfa, 0x95, 0xcb, 0xcc, 0x98, 0xd4, 0xe4, 0x45,
	0xe6, 0xfc, 0x38, 0x77, 0xe6, 0x78, 0x0f, 0xcc, 0x7b, 0x81, 0x5c, 0xac, 0xd5, 0x37, 0x6e, 0x39,
	0xeb, 0x18, 0x71, 0x94, 0xf0, 0x7f, 0x97, 0x8c, 0x95, 0x15, 0xae, 0x15, 0x3c, 0xf4, 0xc5, 0x1a,
	0xeb, 0xb6, 0x3b, 0x6b, 0x4f, 0xb8, 0x82, 0x9f, 0x4f, 0x58, 0x61, 0x87, 0x09, 0x9e, 0x7a, 0x14,
	0x1d, 0x99, 0x81, 0x49, 0x73, 0xcf, 0x08, 0x8c, 0x68, 0x9a, 0x98, 0x34, 0x0f, 0x4b, 0xf8, 0xf5,
	0x26, 0x90, 0xbf, 0x50, 0xd1, 0x5d, 0x2c, 0x0b, 0x70, 0x4e, 0x3d, 0xf2, 0xf3, 0xe0, 0xd2, 0x42,
	0x52, 0xc6, 0x73, 0xe4, 0x9e, 0xa9, 0xa9, 0x12, 0x84, 0x80, 0xdd, 0xa6, 0x25, 0x7a, 0x56, 0x60,
	0x44, 0x4e, 0xa2, 0xce, 0xd2, 0x99, 0xb1, 0xbe, 0xe9, 0x3c, 0x5b, 0x41, 0x2d, 0xc2, 0x1d, 0xc0,
	0x1e, 0x79, 0x4d, 0x85, 0xa0, 0xac, 0x21, 0x01, 0xb8, 0x69, 0xdb, 0x56, 0x34, 0x4b, 0x3b, 0xca,
	0x9a, 0xe1, 0xa5, 0x31, 0x52, 0x37, 0x23, 0xaf, 0x3d, 0x33, 0xb0, 0xa2, 0x69, 0xa2, 0xce, 0xe1,
	0x07, 0xd8, 0x32, 0xec, 0xfd, 0x4f, 0x10, 0x1f, 0x7e, 0xc8, 0x5d, 0x34, 0x69, 0x8d, 0x43, 0xbc,
	0xab, 0x96, 0x69, 0xb0, 0x4e, 0x69, 0xa5, 0x22, 0x4e, 0x13, 0x2d, 0xc8, 0x16, 0xdc, 0xf6, 0x9a,
	0x46, 0x78, 0x76, 0x60, 0x45, 0xee, 0x66, 0x1e, 0xeb, 0xf5, 0x7e, 0xe7, 0x4c, 0xc6, 0xae, 0xcd,
	0xa7, 0x01, 0x8e, 0x7c, 0x5f, 0x90, 0xbf, 0x60, 0xcb, 0x8d, 0x91, 0xe5, 0x30, 0x71, 0xb7, 0x42,
	0xdf, 0x1d, 0xf1, 0x7f, 0x06, 0x59, 0x81, 0xf5, 0x90, 0xe7, 0x64, 0x4c, 0x6f, 0x2c, 0xe4, 0x3f,
	0x4c, 0x74, 0x4d, 0x64, 0x31, 0xe0, 0x9b, 0xd6, 0xfc, 0x65, 0xac, 0x4b, 0x8e, 0x2f, 0x25, 0xc7,
	0xcf, 0xb2, 0xe4, 0xdd, 0x1f, 0x98, 0x67, 0xac, 0x8e, 0x0b, 0xc6, 0x0a, 0x5a, 0xc9, 0xb1, 0xac,
	0x38, 0xee, 0x66, 0x8f, 0xc5, 0x51, 0xde, 0xfa, 0x8a, 0xfc, 0x9d, 0x66, 0xb8, 0x37, 0x0e, 0x13,
	0x35, 0xb6, 0xfd, 0x0a, 0x00, 0x00, 0xff, 0xff, 0xab, 0x1b, 0xa9, 0x05, 0x45, 0x02, 0x00, 0x00,
}
