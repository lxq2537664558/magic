package service

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aiyun/openapm/proto"

	"github.com/uber-go/zap"
)

// Group monitor Group
type Group struct {
	ID     string
	Parent *Group
	Alerts map[string]*Alert
	Users  map[string]*User
	Hosts  map[string]*Host
	Child  map[string]bool
	// 第一个key为 metric.Name+"."+field ， 第二个key为子的路径
	// Alerts map[string]map[string]*Alert
}

func NewGroup() *Group {
	return &Group{
		Alerts: make(map[string]*Alert),
		Users:  make(map[string]*User),
		Hosts:  make(map[string]*Host),
		Child:  make(map[string]bool),
	}
}

func (g *Group) Show() {
	VLogger.Info("Group", zap.String("@ID", g.ID))
	VLogger.Info("Group", zap.Object("@Alerts", g.Alerts))
	VLogger.Info("Group", zap.Object("@Parent", g.Parent))
	VLogger.Info("Group", zap.Object("@Child", g.Child))
	VLogger.Info("Group", zap.Object("@Users", g.Users))
	VLogger.Info("Group", zap.Object("@Hosts", g.Hosts))
}

// ComputAlarm 计算是否需要报警,如果当前组节点找不到alert，返回false
func (g *Group) ComputAlarm(am *Alarmer, metric *MetricData, Interval int, originalGroup *Group) bool {
	find := false
	for field, _ := range metric.Fields {
		if alert, ok := g.Alerts[metric.Name+"."+field]; ok {
			if value, ok := metric.Fields[field].(float64); ok {
				am.compute(alert, metric, value, Interval, originalGroup)
				find = true
			} else {
				continue
			}
		}
	}
	return find
}

func (g *Group) AddChild(child string) {
	g.Child[child] = true
}

type Groups struct {
	sync.RWMutex
	groups map[string]*Group
}

func NewGroups() *Groups {
	groups := &Groups{
		groups: make(map[string]*Group),
	}
	return groups
}

func (gs *Groups) AddGroup(proGroup *proto.Group) error {
	gs.Lock()
	defer gs.Unlock()
	if _, ok := gs.groups[proGroup.Id]; ok {
		return fmt.Errorf("group %s existing", proGroup.Id)
	} else {
		VLogger.Info("AddGroup", zap.String("@Gid", proGroup.Id))
		group, err := proGroupToGroup(proGroup)
		if err != nil {
			VLogger.Error("AddGroup", zap.String("@Gid", proGroup.Id), zap.Error(err))
			return err
		}
		// 如果父节点不为空那么更新父节点到数据库
		if group.Parent != nil {
			parGRaw := groupToGroupRaw(group.Parent)
			streamer.db.Insert(parGRaw)
		}
		// 将新节点更新到数据库
		gs.groups[proGroup.Id] = group
		groupRaw := groupToGroupRaw(group)
		// add hosts
		for hname, _ := range proGroup.Hosts {
			streamer.hostsTogroups.Add(hname, proGroup.Id)
		}
		return streamer.db.Insert(groupRaw)
	}
}

func (gs *Groups) AddAlerts(alerts *proto.Alerts) error {
	gs.Lock()
	defer gs.Unlock()
	if group, ok := gs.groups[alerts.GroupId]; ok {
		log.Println(group)
		for rule, proAlert := range alerts.GetAlerts() {
			alertStatic := NewAlertStatic()
			alertStatic.Type = proAlert.Type
			alertStatic.Operator = proAlert.Operator
			alertStatic.WarnValue = proAlert.WarnValue
			alertStatic.CritValue = proAlert.CritValue
			alertStatic.WarnOutput = proAlert.WarnOutput
			alertStatic.CritOutput = proAlert.CritOutput
			alertStatic.Duration = proAlert.Duration
			alertStatic.Template = proAlert.Template
			alert := NewAlert()
			alert.AlertSt = alertStatic
			group.Alerts[rule] = alert

		}
		// updata db data
		groupRaw := groupToGroupRaw(group)
		return streamer.db.Insert(groupRaw)
	} else {
		return fmt.Errorf("group %s is not existing", alerts.GroupId)
	}
}

func (gs *Groups) AddUsers(users *proto.Users) error {
	gs.Lock()
	defer gs.Unlock()
	if group, ok := gs.groups[users.GroupId]; ok {
		log.Println(group)
		for uname, proUser := range users.Users {
			group.Users[uname] = &User{
				Name:        proUser.Name,
				Sms:         proUser.Sms,
				Mail:        proUser.Mail,
				MessagePush: proUser.MessagePush,
			}
		}
		// updata db data
		groupRaw := groupToGroupRaw(group)
		return streamer.db.Insert(groupRaw)
	} else {
		return fmt.Errorf("group %s is not existing", users.GroupId)
	}
}

func (gs *Groups) AddHosts(hosts *proto.Hosts) error {
	gs.Lock()
	defer gs.Unlock()
	if group, ok := gs.groups[hosts.GroupId]; ok {
		log.Println(group)
		for hostname, proHost := range hosts.Hosts {
			group.Hosts[hostname] = &Host{
				Name: proHost.Name,
				Addr: proHost.Addr,
			}
			streamer.hostsTogroups.Add(hostname, hosts.GroupId)
		}
		// updata db data
		groupRaw := groupToGroupRaw(group)
		return streamer.db.Insert(groupRaw)
	} else {
		return fmt.Errorf("group %s is not existing", hosts.GroupId)
	}
}

func (gs *Groups) GetGroup(gid string) *Group {
	gs.RLock()
	defer gs.RUnlock()
	if group, ok := gs.groups[gid]; ok {
		return group
	}
	return nil
}

func proGroupToGroup(proGroup *proto.Group) (*Group, error) {
	group := NewGroup()

	// set alertStatic
	for rule, proAlert := range proGroup.GetAlerts() {
		alertStatic := NewAlertStatic()
		alertStatic.Type = proAlert.Type
		alertStatic.Operator = proAlert.Operator
		alertStatic.WarnValue = proAlert.WarnValue
		alertStatic.CritValue = proAlert.CritValue
		alertStatic.WarnOutput = proAlert.WarnOutput
		alertStatic.CritOutput = proAlert.CritOutput
		alertStatic.Duration = proAlert.Duration
		alertStatic.Template = proAlert.Template
		alert := NewAlert()
		alert.AlertSt = alertStatic
		group.Alerts[rule] = alert
	}

	// set gid
	group.ID = proGroup.Id

	// set Parent
	index := strings.LastIndex(proGroup.Id, ".")
	if index != -1 {
		Parent := proGroup.Id[0:index]
		VLogger.Info("Parent", zap.String("@Parent", Parent))
		if parentGroup, ok := streamer.groups.groups[Parent]; ok {
			// Parent add child
			parentGroup.AddChild(proGroup.Id)
			group.Parent = parentGroup
		} else {
			return nil, fmt.Errorf("Parent  %s is not existing", parentGroup)
		}
	}

	// set child

	// set users
	for uname, proUser := range proGroup.Users {
		group.Users[uname] = &User{
			Name:        proUser.Name,
			Sms:         proUser.Sms,
			Mail:        proUser.Mail,
			MessagePush: proUser.MessagePush,
		}
	}

	// set hosts
	for hname, proHost := range proGroup.Hosts {
		group.Hosts[hname] = &Host{
			Name: proHost.Name,
			Addr: proHost.Addr,
		}
	}
	return group, nil
}

func groupToGroupRaw(group *Group) *GroupRaw {
	groupRaw := NewGroupRaw()

	// set ID
	groupRaw.ID = group.ID

	// set alertStatic
	for rule, groupAlert := range group.Alerts {
		alertStatic := NewAlertStatic()
		alertStatic.Type = groupAlert.AlertSt.Type
		alertStatic.Operator = groupAlert.AlertSt.Operator
		alertStatic.WarnValue = groupAlert.AlertSt.WarnValue
		alertStatic.CritValue = groupAlert.AlertSt.CritValue
		alertStatic.WarnOutput = groupAlert.AlertSt.WarnOutput
		alertStatic.CritOutput = groupAlert.AlertSt.CritOutput
		alertStatic.Duration = groupAlert.AlertSt.Duration
		alertStatic.Template = groupAlert.AlertSt.Template
		groupRaw.AlertStatics[rule] = alertStatic
	}

	// set Parent
	if group.Parent != nil {
		groupRaw.Parent = group.Parent.ID
	}

	// set Users
	for uname, groupUser := range group.Users {
		groupRaw.Users[uname] = &User{
			Name:        groupUser.Name,
			Sms:         groupUser.Sms,
			Mail:        groupUser.Mail,
			MessagePush: groupUser.MessagePush,
		}
	}

	// set hosts
	for hname, groupHost := range group.Hosts {
		groupRaw.Hosts[hname] = &Host{
			Name: groupHost.Name,
			Addr: groupHost.Addr,
		}
	}

	return groupRaw
}

// GroupRaw access to data
type GroupRaw struct {
	ID           string
	AlertStatics map[string]*AlertStatic
	Parent       string
	Users        map[string]*User
	Hosts        map[string]*Host
}

func NewGroupRaw() *GroupRaw {
	gr := &GroupRaw{
		AlertStatics: make(map[string]*AlertStatic),
		Users:        make(map[string]*User),
		Hosts:        make(map[string]*Host),
	}
	return gr
}
