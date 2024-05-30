package request_result

import (
	"strconv"
	"time"
)

type RequestResult struct {
	Index        int
	Url          string
	UrlHost      string
	UrlPath      string
	UrlParmeters string
	UrlFragment  string
	Duration     time.Duration
	RequestAt    time.Time
	ResponseAt   time.Time
	StatusCode   int
	//Success      bool
}

func GetCsvHeader() []string {
	return []string{
		//"#",
		"URL",
		"Host",
		"Path",
		"Parameters",
		"Fragment",
		"Success",
		"Status code",
		"Duration (ms)",
		//"Request at",
		"Response time",
	}
}

func (r *RequestResult) GetCsvRow() []string {
	success := "OK"

	if r.StatusCode == 404 {
		success = "Not found"
	} else if r.StatusCode >= 400 {
		success = "Not successful"
	}

	record := []string{
		//strconv.Itoa(r.index),
		r.Url,
		r.UrlHost,
		r.UrlPath,
		r.UrlParmeters,
		r.UrlFragment,
		success,
		strconv.Itoa(r.StatusCode),
		strconv.FormatInt(r.Duration.Milliseconds(), 10),
		//r.RequestAt.String(),
		r.ResponseAt.String(),
	}

	return record
}
