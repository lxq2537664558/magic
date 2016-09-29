package service

import (
	"net"

	context "golang.org/x/net/context"

	"github.com/corego/vgo/proto"
	"github.com/uber-go/zap"
	"google.golang.org/grpc"
)

type updater struct{}

type Grpc struct {
	server *grpc.Server
	addr   string
}

func NewGrpc() *Grpc {
	grpc := &Grpc{}
	return grpc
}

func (gp *Grpc) Init(addr string) {
	gp.addr = addr
}

func (gp *Grpc) Start() {
	l, err := net.Listen("tcp", gp.addr)
	if err != nil {
		VLogger.Panic("Grpc", zap.String("@Listen", err.Error()))
	}
	s := grpc.NewServer()
	proto.RegisterAlarmServer(s, &updater{})
	s.Serve(l)
}

func (gp *Grpc) Close() error {
	if gp.server != nil {
		gp.server.Stop()
	}
	return nil
}

func (u *updater) AddGroup(c context.Context, group *proto.Group) (*proto.Reply, error) {
	VLogger.Info("AddGroup", zap.Object("Group", group))
	// add group to system groups
	if err := streamer.AddGroup(group); err != nil {
		return &proto.Reply{Msg: "add Group failed, err message is" + err.Error()}, nil
	}
	return &proto.Reply{Msg: "add Group ok"}, nil
}

func (u *updater) AddAlerts(c context.Context, alerts *proto.Alerts) (*proto.Reply, error) {
	VLogger.Info("AddAlerts", zap.Object("Alerts", alerts))

	if err := streamer.AddAlerts(alerts); err != nil {
		return &proto.Reply{Msg: "add Alerts failed, err message is" + err.Error()}, nil
	}

	return &proto.Reply{Msg: "add Alerts ok"}, nil
}

func (u *updater) AddUsers(c context.Context, users *proto.Users) (*proto.Reply, error) {
	VLogger.Info("AddUsers", zap.Object("Users", users))

	if err := streamer.AddUsers(users); err != nil {
		return &proto.Reply{Msg: "add Users failed, err message is" + err.Error()}, nil
	}
	return &proto.Reply{Msg: "add users ok"}, nil
}

func (u *updater) AddHosts(c context.Context, hosts *proto.Hosts) (*proto.Reply, error) {
	VLogger.Info("AddHosts", zap.Object("Hosts", hosts))

	if err := streamer.AddHosts(hosts); err != nil {
		return &proto.Reply{Msg: "add Hosts failed, err message is" + err.Error()}, nil
	}

	return &proto.Reply{Msg: "add Hosts ok"}, nil
}
