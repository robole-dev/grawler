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

func (r *Result) IsRedirected() bool {
	return r.urlRedirectedFrom != ""
}

func (r *Result) UpdateOnResponse(response *colly.Response, index int, duration time.Duration, err *error) {

	orgUrl := response.Request.Ctx.Get(ctxOrgUrl)

	if orgUrl != response.Request.URL.String() {
		r.urlRedirectedFrom = orgUrl
		r.url = response.Request.URL.String()
	}

	r.status = http.StatusText(response.StatusCode)
	r.statusShort = StatusAbbreviation(response.StatusCode)

	if err != nil && response.StatusCode == 0 {
		r.status = "Skipped"
	} else if r.urlRedirectedFrom != "" {
		r.status += " (Redirected)"
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

	if r.HasError() {
		row += " - Found on: " + r.foundOnUrl
	}

	return row

}

func GetCsvHeader() []string {
	return []string{
		"Response time",
		"Status code",
		"Status",
		"URL",

		"Found on URL",
		"Duration (ms)",
		"Redirected from",

		"Host",
		"Path",
		"Parameters",
		"Fragment",

		"Info / error",
	}
}

func (r *Result) GetCsvRow() []string {

	errorText := ""
	if r.error != nil {
		errorText = fmt.Sprintf("%v", r.error)
	}

	record := []string{
		r.responseAt.Format(DateFormat),
		strconv.Itoa(r.statusCode),
		r.status,
		r.url,

		r.foundOnUrl,
		strconv.FormatInt(r.duration.Milliseconds(), 10),
		r.urlRedirectedFrom,

		r.urlHost,
		r.urlPath,
		r.urlParmeters,
		r.urlFragment,

		errorText,
	}

	return record
}

func (r *Result) HasError() bool {
	return r.error != nil || r.statusCode >= 400
}
