package grawl

import (
	"github.com/gocolly/colly/v2"
	"net/http"
	"strconv"
	"time"
)

type Result struct {
	id                  uint32
	Index               int
	orgUrl              string
	url                 string
	urlHost             string
	urlPath             string
	urlParmeters        string
	urlFragment         string
	urlRedirectedFrom   string
	duration            time.Duration
	requestAt           time.Time
	responseAt          time.Time
	statusCode          int
	error               error
	status              string
	statusShort         string
	foundOnUrl          string
	contentType         string
	depth               int
	httpErrorCodeRanges *responseCodeRanges
}

func NewResult(id uint32, url string, foundOnUrl string, httpErrorRanges *responseCodeRanges) *Result {
	//fmt.Println("found on", foundOnUrl)
	return &Result{
		id:                  id,
		orgUrl:              url,
		url:                 url,
		foundOnUrl:          foundOnUrl,
		requestAt:           time.Now(),
		httpErrorCodeRanges: httpErrorRanges,
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
	r.depth = response.Request.Depth

	//fmt.Println("CT", r.contentType, " - ", response.Headers.Get("Content-Type"))

	r.contentType = response.Headers.Get("Content-Type")

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

func (r *Result) HasError() bool {
	return r.error != nil || r.httpErrorCodeRanges.IsError(r.statusCode)
}
