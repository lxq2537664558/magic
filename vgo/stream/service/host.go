package service

import (
	"log"
	"sync"

	"github.com/uber-go/zap"
)

type Host struct {
	Name string
	Addr string
}

func NewHost() *Host {
	return &Host{}
}

type HostsToGroup struct {
	sync.RWMutex
	hostsTogroups map[string]map[string]bool
}

func NewHostsToGroup() *HostsToGroup {
	hosts := &HostsToGroup{
		hostsTogroups: make(map[string]map[string]bool),
	}
	return hosts
}

func (hs *HostsToGroup) Show() {
	hs.RLock()
	defer hs.RUnlock()

	for hostname, groupIDs := range hs.hostsTogroups {
		VLogger.Info("HostsToGroup", zap.String("@hostname", hostname))
		VLogger.Info("HostsToGroup", zap.Object("@groupIDs", groupIDs))
	}

}

func (hs *HostsToGroup) Add(hostname string, gid string) {
	hs.Lock()
	if groups, ok := hs.hostsTogroups[hostname]; ok {
		groups[gid] = true
	} else {
		groups := make(map[string]bool)
		groups[gid] = true
		hs.hostsTogroups[hostname] = groups
	}
	hs.Unlock()
}

func (hs *HostsToGroup) Get(hostname string) map[string]bool {
	hs.RLock()
	if host, ok := hs.hostsTogroups[hostname]; ok {
		hs.RUnlock()
		return host
	}
	hs.RUnlock()
	return nil
}

func (hs *HostsToGroup) DeleGroupInHosts(hostname string, gid string) error {
	hs.Lock()
	if groups, ok := hs.hostsTogroups[hostname]; ok {
		delete(groups, gid)
	}
	hs.Unlock()
	return nil
}

func (hs *HostsToGroup) DelHost(hostname string) map[string]bool {
	hs.Lock()
	if _, ok := hs.hostsTogroups[hostname]; ok {
		delete(hs.hostsTogroups, hostname)
	}
	hs.Unlock()
	return nil
}

func HostTest() {
	hosts := NewHostsToGroup()
	hosts.Add("scc@Google", "zeus")
	hosts.Add("scc@Google", "room")
	hosts.Add("scc@Google", "cache")
	hosts.Add("scc@Google", "center")
	hosts.Add("scc@Google", "vgo")
	hosts.Add("scc@Google", "uuid")
	gs := hosts.Get("scc@Google")
	log.Println("Host get groups is ", gs)
}
