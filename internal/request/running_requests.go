package request

import (
	"sync"
)

// @see https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c

type RunningRequests struct {
	sync.RWMutex
	results map[uint32]*Result
}

func NewRunningRequests() *RunningRequests {
	return &RunningRequests{
		results: make(map[uint32]*Result),
	}
}

func (rm *RunningRequests) GetValues() *[]*Result {
	rm.RLock()
	defer rm.RUnlock()
	values := make([]*Result, 0, len(rm.results))
	for _, value := range rm.results {
		values = append(values, value)
	}
	return &values
}

func (rm *RunningRequests) Load(key uint32) (value *Result, ok bool) {
	rm.RLock()
	result, ok := rm.results[key]
	rm.RUnlock()
	return result, ok
}

func (rm *RunningRequests) Delete(key uint32) {
	rm.Lock()
	delete(rm.results, key)
	rm.Unlock()
}

func (rm *RunningRequests) Store(key uint32, value *Result) *Result {
	rm.Lock()
	rm.results[key] = value
	rm.Unlock()
	return value
}
