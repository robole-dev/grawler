package request

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"strconv"
	"time"
)

type Result struct {
	Index             int
	Url               string
	UrlHost           string
	UrlPath           string
	UrlParmeters      string
	UrlFragment       string
	UrlRedirectedFrom string
	Duration          time.Duration
	RequestAt         time.Time
	ResponseAt        time.Time
	StatusCode        int
	Error             error
}

func GetCsvHeader() []string {
	return []string{
		//"#",
		"URL",
		"Status",
		"Host",
		"Path",
		"Parameters",
		"Fragment",
		"Duration (ms)",
		"Status code",
		"Redirected from",
		//"Request at",
		"Response time",
		"Info / Error",
	}
}

func NewResult() *Result {
	return &Result{}
}

func (r *Result) GetCsvRow() []string {
	success := "OK"
	if r.StatusCode == 404 {
		success = "Not found"
	} else if r.StatusCode >= 400 {
		success = "Not successful"
	} else if r.Error != nil && r.StatusCode == 0 {
		success = "Skipped"
	}

	errorText := ""
	if r.Error != nil {
		errorText = fmt.Sprintf("%v", r.Error)
	}

	record := []string{
		//strconv.Itoa(r.index),
		r.Url,
		success,
		r.UrlHost,
		r.UrlPath,
		r.UrlParmeters,
		r.UrlFragment,
		strconv.FormatInt(r.Duration.Milliseconds(), 10),
		strconv.Itoa(r.StatusCode),
		r.UrlRedirectedFrom,
		//r.RequestAt.String(),
		r.ResponseAt.String(),
		errorText,
	}

	return record
}

func (r *Result) UpdateOnResponse(response *colly.Response, index int, duration time.Duration, err *error) {

	if r.Url != response.Request.URL.String() {
		r.UrlRedirectedFrom = r.Url
		r.Url = response.Request.URL.String()
	}

	r.Duration = duration
	r.Index = index
	r.UrlPath = response.Request.URL.Path
	r.UrlHost = response.Request.URL.Host
	r.UrlParmeters = response.Request.URL.RawQuery
	r.UrlFragment = response.Request.URL.RawFragment
	r.StatusCode = response.StatusCode
	r.ResponseAt = time.Now()

	if err != nil {
		r.Error = *err
	}
}
