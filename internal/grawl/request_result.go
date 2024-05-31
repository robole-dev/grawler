package grawl

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"net/http"
	"strconv"
	"time"
)

type Result struct {
	id                uint32
	Index             int
	orgUrl            string
	url               string
	urlHost           string
	urlPath           string
	urlParmeters      string
	urlFragment       string
	urlRedirectedFrom string
	duration          time.Duration
	requestAt         time.Time
	responseAt        time.Time
	statusCode        int
	error             error
	status            string
	statusShort       string
	foundOnUrl        string
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
		"Found on URL",
		"Redirected from",
		//"Request at",
		"Response time",
		"Info / error",
	}
}

func NewResult(id uint32, url string, foundOnUrl string) *Result {
	//fmt.Println("found on", foundOnUrl)
	return &Result{
		id:         id,
		orgUrl:     url,
		url:        url,
		foundOnUrl: foundOnUrl,
		requestAt:  time.Now(),
	}
}

func (r *Result) GetRequestAt() time.Time {
	return r.requestAt
}

func (r *Result) GetCsvRow() []string {

	errorText := ""
	if r.error != nil {
		errorText = fmt.Sprintf("%v", r.error)
	}

	record := []string{
		//strconv.Itoa(r.index),
		r.url,
		r.status,
		r.urlHost,
		r.urlPath,
		r.urlParmeters,
		r.urlFragment,
		strconv.FormatInt(r.duration.Milliseconds(), 10),
		strconv.Itoa(r.statusCode),
		r.foundOnUrl,
		r.urlRedirectedFrom,
		//r.requestAt.String(),
		r.responseAt.Format(DateFormat),
		errorText,
	}

	return record
}

func (r *Result) IsRedirected() bool {
	return r.urlRedirectedFrom != ""
}

func (r *Result) UpdateOnResponse(response *colly.Response, index int, duration time.Duration, err *error) {

	orgUrl := response.Request.Ctx.Get(ctxOrgUrl)

	if orgUrl != response.Request.URL.String() {
		r.urlRedirectedFrom = orgUrl
		r.url = response.Request.URL.String()
	}

	r.status = http.StatusText(r.statusCode)
	r.statusShort = StatusAbbreviation(r.statusCode)

	if r.error != nil && r.statusCode == 0 {
		r.status = "Skipped"
	} else if r.urlRedirectedFrom != "" {
		r.status = "Redirect"
	}

	r.duration = duration
	r.Index = index
	r.urlPath = response.Request.URL.Path
	r.urlHost = response.Request.URL.Host
	r.urlParmeters = response.Request.URL.RawQuery
	r.urlFragment = response.Request.URL.RawFragment
	r.statusCode = response.StatusCode
	r.responseAt = time.Now()

	if err != nil {
		r.error = *err
	}
}

func (r *Result) GetPrintRow() string {
	row := ""
	row += "[" + r.responseAt.Format(DateFormat) + "]"
	row += " "
	row += strconv.Itoa(r.statusCode)
	row += " "
	row += StatusAbbreviation(r.statusCode)
	row += " - "
	row += r.url

	if r.IsRedirected() {
		row += " - Redirected from: " + r.urlRedirectedFrom
	}

	return row

}

func (r *Result) HasError() bool {
	return r.error != nil
}
