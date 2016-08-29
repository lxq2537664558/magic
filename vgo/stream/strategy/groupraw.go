package strategy

import "fmt"

type GroupRaw struct {
	ID     string
	Alerts map[string]*Alert
	Parent string
	Users  []*User
	Hosts  []string
}

func NewGroupRaw() *GroupRaw {
	gr := &GroupRaw{}
	return gr
}

// GetNewGroup 生成新的Group，并赋值父节点
func (gr *GroupRaw) GetNewGroup(data *GroupRaw) (*Group, error) {
	g := NewGroup()
	g.ID = gr.ID
	g.Alerts = gr.Alerts
	g.Hosts = gr.Hosts
	g.Users = gr.Users

	// Set Group Parent
	if gr.Parent == "" {
		return g, nil
	}
	strategyes.RLock()
	if parent, ok := strategyes.AllGroup[data.Parent]; ok {
		strategyes.RUnlock()
		g.Parent = parent
		parent.AddChild(g.ID)

	} else {
		strategyes.RUnlock()
		return nil, fmt.Errorf("unfind %s parsent", data.Parent)
	}
	return g, nil
}

// func (gr *GroupRaw) Add(data *GroupRaw) error {
// 	// strategyes.AllGroup

// }

// func (gr *GroupRaw) Init(dbfile string) {
// 	// Loading GroupRaw from boltdb
// }

// func (gr *GroupRaw) Add() {

// }
