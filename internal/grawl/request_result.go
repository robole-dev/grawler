package grawl

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gocolly/colly/v2"
	"net/http"
	"strconv"
	"time"
)

type Result struct {
	id                uint32
	Index             uint32
	initialRequestUrl string
	url               string
	urlHost           string
	urlPath           string
	urlParmeters      string
	urlFragment       string
	urlRedirectedFrom string
	//duration            time.Duration
	requestAt           time.Time
	responseAt          time.Time
	statusCode          int
	error               error
	status              string
	statusShort         string
	foundOnUrl          string
	contentType         string
	depth               int
	httpErrorCodeRanges *ResponseCodeRanges
	requestCount        uint32
	updatedAtResponse   bool
}

func NewResult(id uint32, url string, foundOnUrl string, httpErrorRanges *ResponseCodeRanges) *Result {
	//fmt.Println("found on", foundOnUrl)
	return &Result{
		id:                  id,
		initialRequestUrl:   url,
		url:                 url,
		foundOnUrl:          foundOnUrl,
		httpErrorCodeRanges: httpErrorRanges,
		requestCount:        0,
		updatedAtResponse:   false,
	}
}

func (r *Result) Id() uint32 {
	return r.id
}

func (r *Result) GetRequestAt() time.Time {
	return r.requestAt
}

func (r *Result) IsRedirected() bool {
	return r.urlRedirectedFrom != ""
}

func (r *Result) UpdateOnRoundTripStart(requestTime time.Time) {
	r.requestAt = requestTime
}

func (r *Result) UpdateOnRoundTripEnd(responseTime time.Time) {
	r.responseAt = responseTime
}

func (r *Result) UpdateOnResponse(
	response *colly.Response,
	index uint32,
	err *error,
	requestCount uint32,
) {

	//requestId := response.Request.ID
	//initialRequestUrl := response.Request.Ctx.Get(ctxOrgUrl)
	initialRequestUrl := r.initialRequestUrl
	requestUrl := response.Request.URL.String()

	if initialRequestUrl != requestUrl {
		r.urlRedirectedFrom = initialRequestUrl
		r.url = requestUrl
	}

	r.status = http.StatusText(response.StatusCode)
	r.statusShort = StatusAbbreviation(response.StatusCode)

	if err != nil && response.StatusCode == 0 {
		r.status = "Skipped"
	} else if r.IsRedirected() {
		r.status += " (Redirected)"
	}

	r.requestCount = requestCount
	//r.duration = duration
	r.Index = index
	r.urlPath = response.Request.URL.Path
	r.urlHost = response.Request.URL.Host
	r.urlParmeters = response.Request.URL.RawQuery
	r.urlFragment = response.Request.URL.RawFragment
	r.statusCode = response.StatusCode
	r.depth = response.Request.Depth

	//fmt.Println("CT", r.contentType, " - ", response.Headers.Get("Content-Type"))

	r.contentType = response.Headers.Get("Content-Type")
	r.updatedAtResponse = true

	if err != nil {
		r.error = *err
	}
}

func (r *Result) GetPrintRow() string {
	row := ""
	row += "[" + r.responseAt.Format(DateFormat) + "]"
	row += " "
	row += fmt.Sprintf("%d/%d", r.Index, r.requestCount)
	row += " "
	row += strconv.Itoa(r.statusCode)
	row += " "
	row += StatusAbbreviation(r.statusCode)
	row += " "
	row += fmt.Sprintf("%dms", r.GetDuration().Milliseconds())
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

func (r *Result) GetDuration() time.Duration {
	if r.responseAt.IsZero() {
		return 0
	}
	return r.responseAt.Sub(r.requestAt)
}

func (r *Result) PrintRowColored() {
	if r.IsRedirected() {
		color.Yellow(r.GetPrintRow())
	} else if r.HasError() {
		color.Red(r.GetPrintRow())
	} else {
		color.Green(r.GetPrintRow())
	}
}
