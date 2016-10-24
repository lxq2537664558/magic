package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aiyun/openapm/proto"

	"github.com/kataras/iris"
	"github.com/uber-go/zap"
)

type RawAlerts struct {
	Alerts []*RawAlert `json:"alerts"`
}

type RawAlert struct {
	Name       string `json:"name"`
	Type       int32  `json:"type"`
	Operator   int32  `json:"operator"`
	WarnValue  int32  `json:"warn_value"`
	CritValue  int32  `json:"crit_value"`
	WarnOutput string `json:"warn_output"`
	CritOutput string `json:"crit_output"`
	Duration   int32  `json:"duration"`
	Template   string `json:"template"`
}

type RawUsers struct {
	Users []*RawUser `json:"users"`
}

type RawUser struct {
	Name        string `json:"name"`
	MessagePush string `json:"message_push"`
	Mail        string `json:"mail"`
	Sms         string `json:"sms"`
}

func addGroup(c *iris.Context) {
	gid := c.PostValue("groupid")

	ralerts := &RawAlerts{}
	err := json.Unmarshal(c.FormValue("alerts"), &ralerts)
	if err != nil {
		vLogger.Warn("unmarshal alerts error", zap.Error(err))
	}

	rusers := &RawUsers{}
	err = json.Unmarshal(c.FormValue("users"), &rusers)
	if err != nil {
		vLogger.Warn("unmarshal users error", zap.Error(err))
	}

	alerts := make(map[string]*proto.Alert)
	for _, a := range ralerts.Alerts {
		alert := &proto.Alert{
			Type:       a.Type,
			Operator:   a.Operator,
			WarnValue:  a.WarnValue,
			CritValue:  a.CritValue,
			WarnOutput: a.WarnOutput,
			CritOutput: a.CritOutput,
			Duration:   a.Duration,
			Template:   a.Template,
		}

		alerts[a.Name] = alert
	}

	users := make(map[string]*proto.User)
	for _, u := range rusers.Users {
		user := &proto.User{
			Name:        u.Name,
			Sms:         u.Sms,
			Mail:        u.Mail,
			MessagePush: u.MessagePush,
		}
		users[u.Name] = user
	}

	g := &proto.Group{
		Id:     gid,
		Alerts: alerts,
		Users:  users,
	}

	r, err := gClient.AddGroup(context.Background(), g)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Println(r)
}

func addUsers(c *iris.Context) {
	gid := c.PostValue("groupid")

	rusers := &RawUsers{}
	err := json.Unmarshal(c.FormValue("users"), &rusers)
	if err != nil {
		vLogger.Warn("unmarshal users error", zap.Error(err))
	}

	users := make(map[string]*proto.User)
	for _, u := range rusers.Users {
		user := &proto.User{
			Name:        u.Name,
			Sms:         u.Sms,
			Mail:        u.Mail,
			MessagePush: u.MessagePush,
		}
		users[u.Name] = user
	}

	u := &proto.Users{
		GroupId: gid,
		Users:   users,
	}

	r, err := gClient.AddUsers(context.Background(), u)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Println(r)
}

func addAlerts(c *iris.Context) {
	gid := c.PostValue("groupid")

	ralerts := &RawAlerts{}
	err := json.Unmarshal(c.FormValue("alerts"), &ralerts)
	if err != nil {
		vLogger.Warn("unmarshal alerts error", zap.Error(err))
	}

	alerts := make(map[string]*proto.Alert)
	for _, a := range ralerts.Alerts {
		alert := &proto.Alert{
			Type:       a.Type,
			Operator:   a.Operator,
			WarnValue:  a.WarnValue,
			CritValue:  a.CritValue,
			WarnOutput: a.WarnOutput,
			CritOutput: a.CritOutput,
			Duration:   a.Duration,
			Template:   a.Template,
		}

		alerts[a.Name] = alert
	}

	a := &proto.Alerts{
		GroupId: gid,
		Alerts:  alerts,
	}

	r, err := gClient.AddAlerts(context.Background(), a)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Println(r)
}
