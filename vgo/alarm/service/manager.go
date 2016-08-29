package service

import (
	"log"
	"net"
	"sync"

	"github.com/corego/vgo/proto"

	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type updater struct{}

//Hello(context.Context, *Request) (*Reply, error)

func (u *updater) AddGroup(c context.Context, r *proto.Group) (*proto.Reply, error) {
	g := &Group{
		ID:     r.Id,
		Alerts: make(map[string]*Alert),
		Users:  make(map[string]*User),
	}

	for k, a := range r.Alerts {
		wd := time.Duration(a.WarnDuration)
		cd := time.Duration(a.CritDuration)
		alarm := &Alert{
			Value:       []float64{0, 0},
			Count:       []int32{a.WarnCount, a.CritCount},
			NowCount:    []int32{0, 0},
			AlarmOutput: []string{a.WarnAlarm, a.CritAlarm},
			Duration:    []time.Duration{wd * time.Second, cd * time.Second},
			LastTime:    []time.Time{time.Now().Add(-1 * wd * time.Second), time.Now().Add(-1 * cd * time.Second)},
		}
		g.Alerts[k] = alarm
	}

	for _, u := range r.Users {
		user := &User{
			Name: u.Name,
			Info: u.Info,
		}
		g.Users[u.Name] = user
	}

	gs.groups[r.Id] = g
	log.Println(gs)
	return &proto.Reply{Msg: "add group ok"}, nil
}

func (u *updater) AddUsers(c context.Context, r *proto.Users) (*proto.Reply, error) {
	g, ok := gs.groups[r.GroupId]
	if !ok {
		return &proto.Reply{Msg: "no group found"}, nil
	}

	for _, u := range r.Users {
		g.Users[u.Name] = &User{
			Name: u.Name,
			Info: u.Info,
		}
	}

	log.Println(g.Users)

	return &proto.Reply{Msg: "add users ok"}, nil
}

func (u *updater) AddAlerts(c context.Context, r *proto.Alerts) (*proto.Reply, error) {
	g, ok := gs.groups[r.GroupId]
	if !ok {
		return &proto.Reply{Msg: "no group found"}, nil
	}

	for k, a := range r.Alerts {
		alert, ok := g.Alerts[k]
		if ok {
			alert.Count = []int32{a.WarnCount, a.CritCount}
			alert.AlarmOutput = []string{a.WarnAlarm, a.CritAlarm}
		} else {
			wd := time.Duration(a.WarnDuration)
			cd := time.Duration(a.CritDuration)
			alarm := &Alert{
				Value:       []float64{0, 0},
				Count:       []int32{a.WarnCount, a.CritCount},
				NowCount:    []int32{0, 0},
				AlarmOutput: []string{a.WarnAlarm, a.CritAlarm},
				Duration:    []time.Duration{wd * time.Second, cd * time.Second},
				LastTime:    []time.Time{time.Now().Add(-1 * wd * time.Second), time.Now().Add(-1 * cd * time.Second)},
			}
			g.Alerts[k] = alarm
		}
	}

	log.Println(g.Alerts["cpu.usage"].Duration)
	return &proto.Reply{Msg: "add alerts ok"}, nil
}

func startManager() {
	initGroup()
	l, err := net.Listen("tcp", ":50511")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	proto.RegisterAlarmServer(s, &updater{})
	s.Serve(l)
}

func initGroup() {
	gs = &Groups{
		make(map[string]*Group),
		&sync.RWMutex{},
	}
}
