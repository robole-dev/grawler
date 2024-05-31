package request

import (
	"sort"
	"sync"
)

// @see https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c

type RunningRequests struct {
	sync.RWMutex
	results map[uint32]*Result
	idByUrl map[string]uint32
}

func NewRunningRequests() *RunningRequests {
	return &RunningRequests{
		results: make(map[uint32]*Result),
		idByUrl: make(map[string]uint32),
	}
}

func (rm *RunningRequests) GetValues() *[]*Result {
	rm.RLock()
	defer rm.RUnlock()
	values := make([]*Result, 0, len(rm.results))
	for _, value := range rm.results {
		values = append(values, value)
	}

	sort.Slice(values, func(i, j int) bool {
		return (values)[i].url < (values)[j].url
	})

	return &values
}

func (rm *RunningRequests) Load(key uint32) (value *Result, ok bool) {
	rm.RLock()
	result, ok := rm.results[key]
	rm.RUnlock()
	return result, ok
}

func (rm *RunningRequests) LoadByUrl(url string) (value *Result, ok bool) {
	rm.RLock()
	var result *Result = nil
	id, ok := rm.idByUrl[url]
	if ok {
		result, ok = rm.results[id]
	}
	rm.RUnlock()
	return result, ok
}

func (rm *RunningRequests) Delete(key uint32) {
	rm.Lock()
	value, ok := rm.Load(key)
	if ok && value.url != "" {
		delete(rm.idByUrl, value.url)
	}
	delete(rm.results, key)
	rm.Unlock()
}

func (rm *RunningRequests) Store(key uint32, value *Result, url string) *Result {
	rm.Lock()
	rm.results[key] = value
	rm.idByUrl[url] = key
	rm.Unlock()
	return value
}
