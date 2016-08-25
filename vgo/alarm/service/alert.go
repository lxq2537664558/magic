package service

import (
	"sync"
)

type Groups struct {
	groups map[string]*Group
	*sync.RWMutex
}

var gs *Groups

type Group struct {
	ID     string
	Alerts map[string]*Alert
	Users  []*User
}

type Alert struct {
	Value       []float64 // index 0 : total value of warn , 1 : total value of critical
	Count       []int     //index 0: warn, 1 : critical
	NowCount    []int     //index 0: warn, 1 : critical
	AlarmOutput []string  // warn: mail, critical: mobile
}

type User struct {
	Name string
	Info map[string]string
}

// test, init some test datas
func init() {
	gs = &Groups{
		make(map[string]*Group),
		&sync.RWMutex{},
	}

	// group0
	gs.groups["group0"] = &Group{
		ID:     "group0",
		Alerts: make(map[string]*Alert),
	}

	gs.groups["group0"].Alerts["cpu.cpu_usage"] = &Alert{
		Count:       []int{5, 5},
		AlarmOutput: []string{"mail", "mobile"},
		NowCount:    []int{0, 0},
	}

	gs.groups["group0"].Alerts["mem.mem_usage"] = &Alert{
		Count:       []int{6, 6},
		AlarmOutput: []string{"mail", "mobile"},
		NowCount:    []int{0, 0},
	}

	gs.groups["group0"].Users = append(gs.groups["group0"].Users, &User{
		Name: "sunfei",
		Info: map[string]string{
			"mail":   "cto@188.com",
			"mobile": "15880261185",
		},
	})

	gs.groups["group0"].Users = append(gs.groups["group0"].Users, &User{
		Name: "scc",
		Info: map[string]string{
			"mail":   "kugou.happy@163.com",
			"mobile": "13067779969",
		},
	})

	// group1
	gs.groups["group1"] = &Group{
		ID:     "group1",
		Alerts: make(map[string]*Alert),
	}

	gs.groups["group1"].Alerts["cpu.cpu_usage"] = &Alert{
		Count:       []int{5, 5},
		AlarmOutput: []string{"mail", "mobile"},
		NowCount:    []int{0, 0},
	}

	gs.groups["group1"].Alerts["mem.mem_usage"] = &Alert{
		Count:       []int{6, 6},
		AlarmOutput: []string{"mail", "mobile"},
		NowCount:    []int{0, 0},
	}

	gs.groups["group1"].Users = append(gs.groups["group1"].Users, &User{
		Name: "sunface",
		Info: map[string]string{
			"mail":   "cto@188.com",
			"mobile": "15880261185",
		},
	})

	gs.groups["group1"].Users = append(gs.groups["group1"].Users, &User{
		Name: "congcong",
		Info: map[string]string{
			"mail":   "kugou2008.happy@163.com",
			"mobile": "13067779969",
		},
	})
}
