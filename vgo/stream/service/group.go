package service

// Group monitor Group
type Group struct {
	ID     string
	Alerts map[string]*Alert
	Parent *Group
	Child  []string
	Users  []*User
	Hosts  []string
}

func NewGroup() *Group {
	return &Group{}
}

func (g *Group) AddChild(child string) {
	g.Child = append(g.Child, child)
}
