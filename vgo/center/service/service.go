package service

import (
	"log"

	"github.com/aiyun/openapm/proto"

	"github.com/kataras/iris"
	"google.golang.org/grpc"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

var gClient proto.AlarmClient

func (a *Service) Start() {
	initConfig()

	iris.Post("/add/group", addGroup)
	iris.Post("/add/users/", addUsers)
	iris.Post("/add/alerts", addAlerts)
	go iris.Listen(Conf.Center.Addr)

	//init grpc
	conn, err := grpc.Dial("localhost:50511", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	gClient = proto.NewAlarmClient(conn)
}
