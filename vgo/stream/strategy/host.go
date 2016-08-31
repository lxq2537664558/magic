package strategy

import (
	"log"
	"sync"
)

type Hosts struct {
	sync.RWMutex
	hosts map[string]map[string]bool
}

func NewHosts() *Hosts {
	hosts := &Hosts{
		hosts: make(map[string]map[string]bool),
	}
	return hosts
}

func (hs *Hosts) Add(hostname string, gid string) {
	hs.Lock()
	if groups, ok := hs.hosts[hostname]; ok {
		groups[gid] = true
	} else {
		groups := make(map[string]bool)
		groups[gid] = true
		hs.hosts[hostname] = groups
	}
	hs.Unlock()
}

func (hs *Hosts) Get(hostname string) map[string]bool {
	hs.RLock()
	if host, ok := hs.hosts[hostname]; ok {
		hs.RUnlock()
		return host
	}
	hs.RUnlock()
	return nil
}

func (hs *Hosts) DeleGroupInHosts(hostname string, gid string) error {
	hs.Lock()
	if groups, ok := hs.hosts[hostname]; ok {
		delete(groups, gid)
	}
	hs.Unlock()
	return nil
}

func (hs *Hosts) DelHost(hostname string) map[string]bool {
	hs.Lock()
	if _, ok := hs.hosts[hostname]; ok {
		// for k, _ := range groups {
		// 	delete(groups, k)
		// }
		delete(hs.hosts, hostname)
	}
	hs.Unlock()
	return nil
}

func HostTest() {
	hosts := NewHosts()
	hosts.Add("scc@Google", "zeus")
	hosts.Add("scc@Google", "room")
	hosts.Add("scc@Google", "cache")
	hosts.Add("scc@Google", "center")
	hosts.Add("scc@Google", "vgo")
	hosts.Add("scc@Google", "uuid")
	gs := hosts.Get("scc@Google")
	log.Println("Host get groups is ", gs)
}
