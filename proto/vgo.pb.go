// Code generated by protoc-gen-go.
// source: vgo.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	vgo.proto

It has these top-level messages:
	Group
	Hosts
	Users
	Alerts
	Reply
	Alert
	User
	Host
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type Group struct {
	Id     string            `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Alerts map[string]*Alert `protobuf:"bytes,2,rep,name=alerts" json:"alerts,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Users  map[string]*User  `protobuf:"bytes,3,rep,name=users" json:"users,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Hosts  map[string]*Host  `protobuf:"bytes,4,rep,name=hosts" json:"hosts,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Group) Reset()                    { *m = Group{} }
func (m *Group) String() string            { return proto1.CompactTextString(m) }
func (*Group) ProtoMessage()               {}
func (*Group) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Group) GetAlerts() map[string]*Alert {
	if m != nil {
		return m.Alerts
	}
	return nil
}

func (m *Group) GetUsers() map[string]*User {
	if m != nil {
		return m.Users
	}
	return nil
}

func (m *Group) GetHosts() map[string]*Host {
	if m != nil {
		return m.Hosts
	}
	return nil
}

type Hosts struct {
	GroupId string           `protobuf:"bytes,1,opt,name=group_id,json=groupId" json:"group_id,omitempty"`
	Hosts   map[string]*Host `protobuf:"bytes,2,rep,name=hosts" json:"hosts,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Hosts) Reset()                    { *m = Hosts{} }
func (m *Hosts) String() string            { return proto1.CompactTextString(m) }
func (*Hosts) ProtoMessage()               {}
func (*Hosts) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Hosts) GetHosts() map[string]*Host {
	if m != nil {
		return m.Hosts
	}
	return nil
}

type Users struct {
	GroupId string           `protobuf:"bytes,1,opt,name=group_id,json=groupId" json:"group_id,omitempty"`
	Users   map[string]*User `protobuf:"bytes,2,rep,name=users" json:"users,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Users) Reset()                    { *m = Users{} }
func (m *Users) String() string            { return proto1.CompactTextString(m) }
func (*Users) ProtoMessage()               {}
func (*Users) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Users) GetUsers() map[string]*User {
	if m != nil {
		return m.Users
	}
	return nil
}

type Alerts struct {
	GroupId string            `protobuf:"bytes,1,opt,name=group_id,json=groupId" json:"group_id,omitempty"`
	Alerts  map[string]*Alert `protobuf:"bytes,2,rep,name=alerts" json:"alerts,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Alerts) Reset()                    { *m = Alerts{} }
func (m *Alerts) String() string            { return proto1.CompactTextString(m) }
func (*Alerts) ProtoMessage()               {}
func (*Alerts) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Alerts) GetAlerts() map[string]*Alert {
	if m != nil {
		return m.Alerts
	}
	return nil
}

type Reply struct {
	Msg string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

func (m *Reply) Reset()                    { *m = Reply{} }
func (m *Reply) String() string            { return proto1.CompactTextString(m) }
func (*Reply) ProtoMessage()               {}
func (*Reply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type Alert struct {
	// 1 : average 2 : gauge 3: alive status
	Type int32 `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	// 1:   >,   2:  = ,   3: <
	Operator   int32  `protobuf:"varint,2,opt,name=operator" json:"operator,omitempty"`
	WarnValue  int32  `protobuf:"varint,3,opt,name=warn_value,json=warnValue" json:"warn_value,omitempty"`
	CritValue  int32  `protobuf:"varint,4,opt,name=crit_value,json=critValue" json:"crit_value,omitempty"`
	WarnOutput string `protobuf:"bytes,5,opt,name=warn_output,json=warnOutput" json:"warn_output,omitempty"`
	CritOutput string `protobuf:"bytes,6,opt,name=crit_output,json=critOutput" json:"crit_output,omitempty"`
	Duration   int32  `protobuf:"varint,7,opt,name=duration" json:"duration,omitempty"`
	Template   string `protobuf:"bytes,8,opt,name=template" json:"template,omitempty"`
}

func (m *Alert) Reset()                    { *m = Alert{} }
func (m *Alert) String() string            { return proto1.CompactTextString(m) }
func (*Alert) ProtoMessage()               {}
func (*Alert) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type User struct {
	Name        string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Sms         string `protobuf:"bytes,2,opt,name=sms" json:"sms,omitempty"`
	Mail        string `protobuf:"bytes,3,opt,name=mail" json:"mail,omitempty"`
	MessagePush string `protobuf:"bytes,4,opt,name=message_push,json=messagePush" json:"message_push,omitempty"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto1.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type Host struct {
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Addr string `protobuf:"bytes,2,opt,name=addr" json:"addr,omitempty"`
}

func (m *Host) Reset()                    { *m = Host{} }
func (m *Host) String() string            { return proto1.CompactTextString(m) }
func (*Host) ProtoMessage()               {}
func (*Host) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func init() {
	proto1.RegisterType((*Group)(nil), "proto.Group")
	proto1.RegisterType((*Hosts)(nil), "proto.Hosts")
	proto1.RegisterType((*Users)(nil), "proto.Users")
	proto1.RegisterType((*Alerts)(nil), "proto.Alerts")
	proto1.RegisterType((*Reply)(nil), "proto.Reply")
	proto1.RegisterType((*Alert)(nil), "proto.Alert")
	proto1.RegisterType((*User)(nil), "proto.User")
	proto1.RegisterType((*Host)(nil), "proto.Host")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Alarm service

type AlarmClient interface {
	AddGroup(ctx context.Context, in *Group, opts ...grpc.CallOption) (*Reply, error)
	AddUsers(ctx context.Context, in *Users, opts ...grpc.CallOption) (*Reply, error)
	AddAlerts(ctx context.Context, in *Alerts, opts ...grpc.CallOption) (*Reply, error)
	AddHosts(ctx context.Context, in *Hosts, opts ...grpc.CallOption) (*Reply, error)
}

type alarmClient struct {
	cc *grpc.ClientConn
}

func NewAlarmClient(cc *grpc.ClientConn) AlarmClient {
	return &alarmClient{cc}
}

func (c *alarmClient) AddGroup(ctx context.Context, in *Group, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := grpc.Invoke(ctx, "/proto.Alarm/AddGroup", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alarmClient) AddUsers(ctx context.Context, in *Users, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := grpc.Invoke(ctx, "/proto.Alarm/AddUsers", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alarmClient) AddAlerts(ctx context.Context, in *Alerts, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := grpc.Invoke(ctx, "/proto.Alarm/AddAlerts", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *alarmClient) AddHosts(ctx context.Context, in *Hosts, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := grpc.Invoke(ctx, "/proto.Alarm/AddHosts", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Alarm service

type AlarmServer interface {
	AddGroup(context.Context, *Group) (*Reply, error)
	AddUsers(context.Context, *Users) (*Reply, error)
	AddAlerts(context.Context, *Alerts) (*Reply, error)
	AddHosts(context.Context, *Hosts) (*Reply, error)
}

func RegisterAlarmServer(s *grpc.Server, srv AlarmServer) {
	s.RegisterService(&_Alarm_serviceDesc, srv)
}

func _Alarm_AddGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Group)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).AddGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Alarm/AddGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).AddGroup(ctx, req.(*Group))
	}
	return interceptor(ctx, in, info, handler)
}

func _Alarm_AddUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Users)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).AddUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Alarm/AddUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).AddUsers(ctx, req.(*Users))
	}
	return interceptor(ctx, in, info, handler)
}

func _Alarm_AddAlerts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Alerts)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).AddAlerts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Alarm/AddAlerts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).AddAlerts(ctx, req.(*Alerts))
	}
	return interceptor(ctx, in, info, handler)
}

func _Alarm_AddHosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Hosts)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AlarmServer).AddHosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Alarm/AddHosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AlarmServer).AddHosts(ctx, req.(*Hosts))
	}
	return interceptor(ctx, in, info, handler)
}

var _Alarm_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Alarm",
	HandlerType: (*AlarmServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddGroup",
			Handler:    _Alarm_AddGroup_Handler,
		},
		{
			MethodName: "AddUsers",
			Handler:    _Alarm_AddUsers_Handler,
		},
		{
			MethodName: "AddAlerts",
			Handler:    _Alarm_AddAlerts_Handler,
		},
		{
			MethodName: "AddHosts",
			Handler:    _Alarm_AddHosts_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 545 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x54, 0xcd, 0x8e, 0xd3, 0x30,
	0x10, 0x26, 0x69, 0xd2, 0x6d, 0x26, 0x01, 0x21, 0x5f, 0xc8, 0x56, 0x42, 0x2c, 0x39, 0x55, 0x48,
	0x54, 0xb0, 0x5c, 0x10, 0xb7, 0x3d, 0xac, 0x16, 0x4e, 0xa0, 0x48, 0x70, 0xad, 0xcc, 0xc6, 0xea,
	0x56, 0x24, 0x4d, 0x64, 0x3b, 0x8b, 0xfa, 0x26, 0xdc, 0x78, 0x06, 0x9e, 0x8c, 0x37, 0x40, 0x78,
	0xc6, 0x6e, 0x9b, 0x74, 0xc3, 0xcf, 0xa1, 0xa7, 0xcc, 0xcf, 0x37, 0x9f, 0x3f, 0xcf, 0x4c, 0x0c,
	0xd1, 0xed, 0xb2, 0x9e, 0x37, 0xb2, 0xd6, 0x35, 0x0b, 0xe9, 0x93, 0xfd, 0xf2, 0x21, 0xbc, 0x92,
	0x75, 0xdb, 0xb0, 0x07, 0xe0, 0xaf, 0x8a, 0xd4, 0x3b, 0xf3, 0x66, 0x51, 0x6e, 0x2c, 0xf6, 0x02,
	0xc6, 0xbc, 0x14, 0x52, 0xab, 0xd4, 0x3f, 0x1b, 0xcd, 0xe2, 0xf3, 0xd4, 0x16, 0xce, 0x09, 0x3d,
	0xbf, 0xa0, 0xd4, 0xe5, 0x5a, 0xcb, 0x4d, 0xee, 0x70, 0xec, 0x39, 0x84, 0xad, 0x12, 0x52, 0xa5,
	0x23, 0x2a, 0x78, 0xd4, 0x2b, 0xf8, 0x88, 0x19, 0x8b, 0xb7, 0x28, 0x84, 0xdf, 0xd4, 0xca, 0xf0,
	0x07, 0x03, 0xf0, 0xb7, 0x98, 0x71, 0x70, 0x42, 0x4d, 0xaf, 0x20, 0xee, 0x1c, 0xca, 0x1e, 0xc2,
	0xe8, 0x8b, 0xd8, 0x38, 0xbd, 0x68, 0xb2, 0x0c, 0xc2, 0x5b, 0x5e, 0xb6, 0xc2, 0xe8, 0xf5, 0x0c,
	0x5f, 0xe2, 0xf8, 0xa8, 0x28, 0xb7, 0xa9, 0x37, 0xfe, 0x6b, 0x6f, 0x7a, 0x09, 0xb0, 0x17, 0x33,
	0xc0, 0xf3, 0xb4, 0xcf, 0x13, 0x3b, 0x1e, 0xac, 0x39, 0xa0, 0xd9, 0x8b, 0xfc, 0x7f, 0x1a, 0xac,
	0xe9, 0xd0, 0x64, 0xdf, 0x3c, 0x08, 0x89, 0x87, 0x9d, 0xc2, 0x64, 0x89, 0x77, 0x5f, 0xec, 0xc6,
	0x70, 0x42, 0xfe, 0xbb, 0x62, 0xdf, 0x2a, 0xbf, 0xd7, 0x2a, 0xaa, 0x1b, 0x68, 0xd5, 0x11, 0xa5,
	0x51, 0xa7, 0xfe, 0x21, 0xcd, 0x0e, 0xbd, 0x2f, 0x8d, 0xea, 0xee, 0x0e, 0xfd, 0x48, 0xcd, 0xcf,
	0xbe, 0x7b, 0x30, 0xb6, 0xdb, 0xf0, 0x37, 0x6d, 0x2f, 0x0f, 0x56, 0xf8, 0xb4, 0xbb, 0x12, 0x6a,
	0x68, 0x87, 0x8f, 0xb6, 0x65, 0xd9, 0x29, 0x84, 0xb9, 0x68, 0x4a, 0xa2, 0xa8, 0xd4, 0x72, 0x4b,
	0x61, 0xcc, 0xec, 0xa7, 0xe9, 0x2b, 0xe1, 0x19, 0x83, 0x40, 0x6f, 0x1a, 0x41, 0xc9, 0x30, 0x27,
	0x9b, 0x4d, 0x61, 0x52, 0x37, 0x42, 0x72, 0x5d, 0x4b, 0x3a, 0x23, 0xcc, 0x77, 0x3e, 0x7b, 0x0c,
	0xf0, 0x95, 0xcb, 0xf5, 0xc2, 0x2a, 0x18, 0x51, 0x36, 0xc2, 0xc8, 0x27, 0x0c, 0x60, 0xfa, 0x5a,
	0xae, 0xb4, 0x4b, 0x07, 0x36, 0x8d, 0x11, 0x9b, 0x7e, 0x02, 0x31, 0x55, 0xd7, 0xad, 0x6e, 0x5a,
	0x9d, 0x86, 0xa4, 0x88, 0x08, 0xdf, 0x53, 0x04, 0x01, 0x54, 0xef, 0x00, 0x63, 0x0b, 0xc0, 0x90,
	0x03, 0x18, 0x6d, 0x45, 0x6b, 0xa4, 0xac, 0xea, 0x75, 0x7a, 0x62, 0xb5, 0x6d, 0x7d, 0xcc, 0x69,
	0x51, 0x35, 0x25, 0xd7, 0x22, 0x9d, 0x50, 0xe5, 0xce, 0xcf, 0xae, 0x21, 0xc0, 0x09, 0xe2, 0x7d,
	0xd7, 0xbc, 0x12, 0xae, 0x19, 0x64, 0x63, 0x7f, 0x54, 0xa5, 0xe8, 0xaa, 0xa6, 0x3f, 0xc6, 0x44,
	0x54, 0xc5, 0x57, 0x25, 0xdd, 0xcf, 0xa0, 0xd0, 0x36, 0x7b, 0x91, 0x54, 0x42, 0x29, 0xbe, 0x14,
	0x8b, 0xa6, 0x55, 0x37, 0x74, 0xb9, 0x28, 0x8f, 0x5d, 0xec, 0x83, 0x09, 0x65, 0x73, 0x08, 0x70,
	0x83, 0x07, 0x0f, 0x31, 0x31, 0x5e, 0x14, 0xd2, 0x9d, 0x42, 0xf6, 0xf9, 0x0f, 0x1a, 0x03, 0x97,
	0x15, 0x9b, 0xc1, 0xe4, 0xa2, 0x28, 0xec, 0x33, 0x98, 0x74, 0x9f, 0xa1, 0xe9, 0xd6, 0xa3, 0x51,
	0x66, 0xf7, 0x1c, 0xd2, 0xfe, 0x14, 0x49, 0x77, 0xd5, 0xef, 0x20, 0x9f, 0x41, 0x64, 0x90, 0x6e,
	0x47, 0xef, 0xf7, 0x16, 0xef, 0x0f, 0xac, 0xf6, 0x15, 0x48, 0xba, 0xff, 0xf6, 0x21, 0xf2, 0xf3,
	0x98, 0xdc, 0x57, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x52, 0x97, 0x01, 0x33, 0xc9, 0x05, 0x00,
	0x00,
}
