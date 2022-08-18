package server

import (
	"container/ring"
	"sort"
	"sync"
)

type ItemHistory map[string]*ring.Ring

type HostValueCache struct {
	sync.Mutex
	history ItemHistory
}

func NewHostValueCache(items []ActiveItem) (hvc *HostValueCache) {
	hvc = new(HostValueCache)
	hvc.history = make(ItemHistory)
	for _, item := range items {
		hvc.history[item.Key] = ring.New(10)
	}
	return
}

func (hvc *HostValueCache) SaveValue(value ActiveItemValue) {
	hvc.Lock()
	r := hvc.history[value.Key].Next()
	r.Value = value
	hvc.history[value.Key] = r
	hvc.Unlock()
}

func (hvc *HostValueCache) GetValues(item string) (flat_buf []ActiveItemValue) {
	i := 0

	hvc.Lock()
	buf := hvc.history[item]
	flat_buf = make([]ActiveItemValue, buf.Len())
	buf.Do(func(p interface{}) {
		if p != nil {
			flat_buf[i] = p.(ActiveItemValue)
		} else {
			flat_buf[i] = *new(ActiveItemValue)
		}
		i++
	})
	hvc.Unlock()

	sort.Slice(flat_buf, func(i, j int) bool {
		return flat_buf[i].Clock > flat_buf[j].Clock
	})
	return
}
