package server

import (
	"sync"
)

///gopkg.in/yaml.v3

type CheckList map[string][]ActiveItem

type MonitoringConfig struct {
	sync.RWMutex
	config CheckList
}

func (mc *MonitoringConfig) GetConfig(hostname string) (items []ActiveItem) {
	mc.Lock()
	items = mc.config[hostname]
	mc.Unlock()
	return
}

func (mc *MonitoringConfig) LoadConfig(hostname string, checks []ActiveItem) {
	mc.Lock()
	mc.config[hostname] = checks
	mc.Unlock()
}

func (mc *MonitoringConfig) UpdateOrAddItem(hostname string, check ActiveItem) {
	mc.Lock()
	found := false
	for i, item := range mc.config[hostname] {
		if item.Key == check.Key {
			mc.config[hostname][i] = check
			found = true
			break
		}
	}
	if !found {
		mc.config[hostname] = append(mc.config[hostname], check)
	}
	mc.Unlock()

}

func (mc *MonitoringConfig) GetHosts() (hosts *CheckList) {
	hosts = &mc.config
	return
}

func (mc *MonitoringConfig) GetHostnames() (hostnames []string) {
	hostnames = make([]string, 0, len(mc.config))

	for k, _ := range mc.config {
		hostnames = append(hostnames, k)
	}
	return
}

func (mc *MonitoringConfig) GetKeys(hostname string) (keys []string) {
	for _, v := range mc.GetConfig(hostname) {
		keys = append(keys, v.Key)
	}
	return
}
