package request

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"net/http"
	"strconv"
	"time"
)

const (
	DateFormat = "2006-01-02 15:04:05.000"
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
	status            string
	statusShort       string
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

	errorText := ""
	if r.Error != nil {
		errorText = fmt.Sprintf("%v", r.Error)
	}

	record := []string{
		//strconv.Itoa(r.index),
		r.Url,
		r.status,
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

func (r *Result) IsRedirected() bool {
	return r.UrlRedirectedFrom != ""
}

func (r *Result) UpdateOnResponse(response *colly.Response, index int, duration time.Duration, err *error) {

	binary, err2 := response.Ctx.MarshalBinary()
	if err2 != nil {
		return
	}

	fmt.Println("header", binary)
	if r.Url != response.Request.URL.String() {
		r.UrlRedirectedFrom = r.Url
		r.Url = response.Request.URL.String()
	}

	r.status = http.StatusText(r.StatusCode)
	r.statusShort = StatusAbbreviation(r.StatusCode)

	if r.Error != nil && r.StatusCode == 0 {
		r.status = "Skipped"
	} else if r.UrlRedirectedFrom != "" {
		r.status = "Redirect"
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

func (r *Result) GetPrintRow() string {
	row := ""
	row += "[" + r.ResponseAt.Format(DateFormat) + "]"
	row += " "
	row += strconv.Itoa(r.StatusCode)
	row += " "
	row += StatusAbbreviation(r.StatusCode)
	row += " - "
	row += r.Url

	if r.IsRedirected() {
		row += " - Redirected from: " + r.UrlRedirectedFrom
	}

	return row

}
