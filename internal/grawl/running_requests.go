package grawl

import (
	"sort"
	"sync"
)

// @see https://medium.com/@deckarep/the-new-kid-in-town-gos-sync-map-de24a6bf7c2c

type RunningRequests struct {
	sync.RWMutex
	results map[uint32]*Result
	//requests      map[uint32]*RunningRequest
	idByUrl       map[string]uint32
	foundUrlOnUrl map[string]string
}

//type RunningRequest struct {
//	result *Result
//	initialRequestUrl string
//}

func NewRunningRequests() *RunningRequests {
	return &RunningRequests{
		results:       make(map[uint32]*Result),
		idByUrl:       make(map[string]uint32),
		foundUrlOnUrl: make(map[string]string),
	}
}

func (rr *RunningRequests) GetValues() *[]*Result {
	rr.RLock()
	values := make([]*Result, 0, len(rr.results))
	for _, value := range rr.results {
		values = append(values, value)
	}

	// Sorting
	sort.Slice(values, func(i, j int) bool {
		return (values)[i].url < (values)[j].url
	})

	rr.RUnlock()
	return &values
}

func (rr *RunningRequests) Load(requestId uint32) (value *Result, ok bool) {
	rr.RLock()
	result, ok := rr.results[requestId]
	rr.RUnlock()
	return result, ok
}

func (rr *RunningRequests) LoadByUrl(url string) (value *Result, ok bool) {
	rr.RLock()
	var result *Result = nil
	id, ok := rr.idByUrl[url]
	if ok {
		result, ok = rr.results[id]
	}
	rr.RUnlock()
	return result, ok
}

func (rr *RunningRequests) Delete(key uint32) {
	rr.Lock()
	value, ok := rr.results[key]
	if ok && value.url != "" {
		delete(rr.idByUrl, value.url)
	}
	delete(rr.results, key)
	rr.Unlock()
}

func (rr *RunningRequests) Store(requestId uint32, value *Result, url string) *Result {
	rr.Lock()
	rr.results[requestId] = value
	rr.idByUrl[url] = requestId
	rr.Unlock()
	return value
}

func (rr *RunningRequests) AddFoundUrl(url string, foundOnUrl string) {
	rr.Lock()
	_, ok := rr.foundUrlOnUrl[url]
	if !ok {
		rr.foundUrlOnUrl[url] = foundOnUrl
	}
	rr.Unlock()
}

func (rr *RunningRequests) GetFoundUrl(url string) string {
	rr.Lock()
	result, ok := rr.foundUrlOnUrl[url]
	rr.Unlock()
	if ok {
		return result
	}
	return ""
}

func (rr *RunningRequests) HasFoundUrl(url string) bool {
	rr.Lock()
	_, ok := rr.foundUrlOnUrl[url]
	rr.Unlock()
	return ok
}
