package validates

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/robole-dev/grawler/internal/grawl"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ValData struct {
	url         string
	statusCode  int
	contentType string
}
type Validator struct {
	config              Flags
	collector           *colly.Collector
	runningRequests     *grawl.RunningRequests
	headerAuth          string
	totalDuration       time.Duration
	visitMutex          sync.Mutex
	responseErrorRanges *grawl.ResponseCodeRanges
	requestCount        uint32
	responseCount       uint32
	errorCount          uint32
	redirections        uint32
	valDataByUrl        map[string]*ValData
}

func NewValidator(config Flags) *Validator {
	errorCodeRanges, err := grawl.NewResponseCodeRanges([]string{})
	if err != nil {
		panic(err)
	}

	return &Validator{
		config:              config,
		runningRequests:     grawl.NewRunningRequests(),
		requestCount:        0,
		responseCount:       0,
		errorCount:          0,
		redirections:        0,
		totalDuration:       time.Duration(0),
		responseErrorRanges: errorCodeRanges,
		valDataByUrl:        make(map[string]*ValData),
	}
}

func (v *Validator) ValidateCsv(csvPath string) {
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Error opening csv file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error parsing csv file: %v", err)
	}

	v.collector = colly.NewCollector()
	v.collector.MaxDepth = 10
	v.collector.Async = true
	//v.collector.SetRequestTimeout(time.Duration(v.config.FlagRequestTimeout * float32(time.Second)))
	v.collector.WithTransport(v)
	limitingRule := &colly.LimitRule{
		DomainGlob: "*",
		//Parallelism: v.config.FlagParallel,
	}
	//if v.config.FlagRandomDelay > 0 {
	//	limitingRule.RandomDelay = time.Duration(v.config.FlagRandomDelay) * time.Millisecond
	//} else {
	//	limitingRule.Delay = time.Duration(v.config.FlagDelay) * time.Millisecond
	//}

	v.collector.SetRedirectHandler(v.onRedirect)
	v.collector.OnRequest(v.onRequest)
	v.collector.OnResponse(v.onResponse)
	v.collector.OnError(v.onError)
	v.collector.OnResponseHeaders(v.onResponseHeaders)

	err = v.collector.Limit(limitingRule)
	if err != nil {
		fmt.Println("Error setting limits:", err)
		return
	}

	minColumnNumberNeeded := max(v.config.FlagColContentType, v.config.FlagColUrl, v.config.FlagColStatusCode) + 1
	for i, row := range records {
		if uint64(i) < v.config.FlagSkipRows {
			continue
		}
		if uint64(len(row)) < minColumnNumberNeeded {
			slog.Warn(fmt.Sprintf("Row %d does not have enough columns (have: %d, needed: %d). Skipping row.", i, len(row), minColumnNumberNeeded))
			continue
		}

		statusCode, err := strconv.Atoi(row[v.config.FlagColStatusCode])
		if err != nil {
			slog.Warn(fmt.Sprintf("Row %d does not have a valid status code (%s). Skipping row.", i, row[v.config.FlagColStatusCode]))
			continue
		}

		valData := &ValData{
			url:         row[v.config.FlagColUrl],
			statusCode:  statusCode,
			contentType: row[v.config.FlagColContentType],
		}

		v.visit(valData)
	}

	v.collector.Wait()
}

// RoundTrip implemnts the RoundTripper interface. Needed to measure roundtrip duration
func (v *Validator) RoundTrip(req *http.Request) (res *http.Response, err error) {
	reqResult, ok := v.runningRequests.LoadByUrl(req.URL.String())
	if !ok {
		fmt.Printf("No running request found for %s\n", req.URL)
		fmt.Printf("Urls: %v\n", v.runningRequests)
		os.Exit(1)
		return http.DefaultTransport.RoundTrip(req)
	}
	reqResult.UpdateOnRoundTripStart(time.Now())
	defer func() {
		reqResult.UpdateOnRoundTripEnd(time.Now())
	}()
	return http.DefaultTransport.RoundTrip(req)
}

func (v *Validator) onRedirect(req *http.Request, via []*http.Request) error {
	runningReq, ok := v.runningRequests.LoadByUrl(via[0].URL.String())
	atomic.AddUint32(&v.redirections, 1)
	if ok {
		fmt.Printf("Redirecting to %s from %s. ID: %d\n", req.URL, via[0].URL, runningReq.Id())
	} else {
		return fmt.Errorf("Could not find initial url of redirection to %s from %s\n", req.URL, via[0].URL)
	}
	return nil
}

func (v *Validator) onRequest(r *colly.Request) {
	requestUrl := r.URL.String()

	if v.headerAuth != "" {
		r.Headers.Set("Authorization", v.headerAuth)
	}

	foundOnUrl := v.runningRequests.GetFoundUrl(requestUrl)
	requestResult := grawl.NewResult(r.ID, requestUrl, foundOnUrl, v.responseErrorRanges)

	v.runningRequests.Store(r.ID, requestResult, requestUrl)
	v.requestCount++
}

func (v *Validator) onResponse(r *colly.Response) {
	v.responseCount++

	reqResult, ok := v.runningRequests.Load(r.Request.ID)
	if !ok {
		fmt.Printf("No start time found for %s\n", r.Request.URL)
	}

	reqResult.UpdateOnResponse(r, v.responseCount, nil, v.requestCount)
	v.totalDuration += reqResult.GetDuration()

	// @todo: mutex
	valData := v.valDataByUrl[r.Request.URL.String()]
	if valData == nil {
		fmt.Printf("No validation data found for %s\n", r.Request.URL)
	} else {
		if valData.statusCode != r.StatusCode {
			fmt.Printf("Validation failed for %s: status code asserted: %d - Got: %d\n", r.Request.URL, valData.statusCode, r.StatusCode)
		}
		contentType := r.Headers.Get("Content-Type")
		if valData.contentType != contentType {
			fmt.Printf("Validation failed for %s: status code asserted: %s - Got: %s\n", r.Request.URL, valData.contentType, contentType)
		}
	}

	reqResult.PrintRowColored()
	v.checkStopOnError(reqResult)
}

func (v *Validator) onError(r *colly.Response, err error) {
	v.errorCount++
	v.responseCount++

	// Normal error on aborted binary files like images. Result is printed in OnResponseHeaders
	if err != nil && errors.Is(err, colly.ErrAbortedAfterHeaders) {
		return
	}

	//
	// Remove request if this url is filtered by colly
	//
	if r.StatusCode == 0 {
		fmt.Println("Error", r.Request.URL, err)
		v.runningRequests.Delete(r.Request.ID)
		return
	}

	reqResult, ok := v.runningRequests.Load(r.Request.ID)
	if ok {
		reqResult.UpdateOnResponse(r, v.responseCount, &err, v.requestCount)
		v.totalDuration += reqResult.GetDuration()
		v.printResult(reqResult)
		v.checkStopOnError(reqResult)
	} else {
		fmt.Println("Request data not found", r.Request.URL)
	}
}

func (v *Validator) onResponseHeaders(r *colly.Response) {
	if grawl.IsHtmlResponse(r) || grawl.IsXmlResponse(r) {
		return
	}

	//
	// Abort downloading all non-xml and non-html contents
	//
	r.Request.Abort()
	reqResult, ok := v.runningRequests.Load(r.Request.ID)
	if ok {
		reqResult.UpdateOnResponse(r, v.responseCount, nil, v.requestCount)
		v.totalDuration += reqResult.GetDuration()
		v.printResult(reqResult)
		v.checkStopOnError(reqResult)
	} else {
		fmt.Println("Request data not found", r.Request.URL)
	}
}

func (v *Validator) checkStopOnError(result *grawl.Result) {
	if !result.HasError() {
		return
	}

	//if v.config.FlagStopOnError {
	//	fmt.Println("Stop validating after error.")
	//	if result.error != nil {
	//		fmt.Printf("Validating error: %s\n", result.error.Error())
	//	}
	//	os.Exit(1)
	//}
	//
	//if v.config.FlagPauseOnError {
	//	fmt.Println("Pause validating after error.")
	//	prompt := grawl.PromptResume()
	//	switch prompt {
	//	case "a":
	//		fmt.Println("Validating aborted.")
	//		os.Exit(1)
	//	case "s":
	//		fmt.Println("Url skipped.")
	//		return
	//	case "r":
	//		fmt.Println("Retry url...")
	//		allowRevisit := g.collector.AllowURLRevisit
	//		v.collector.AllowURLRevisit = true
	//		_ = v.collector.Visit(result.url)
	//		v.collector.AllowURLRevisit = allowRevisit
	//	}
	//}
}

func (v *Validator) visit(valData *ValData) {
	v.visitMutex.Lock()
	defer v.visitMutex.Unlock()

	url := strings.Trim(valData.url, " ")

	//fmt.Println("Visiting", url)

	visited, err := v.collector.HasVisited(url)
	if visited {
		//fmt.Println("Visited", url)
		return
	}
	if err != nil {
		fmt.Println("Could not check if url has been visited: ", url)
	}
	hasFoundUrl := v.runningRequests.HasFoundUrl(url)
	if hasFoundUrl {
		return
	}

	v.runningRequests.AddFoundUrl(url, url)
	v.valDataByUrl[url] = valData
	//fmt.Println("Visit:", url)
	_ = v.collector.Visit(url)
}

func (v *Validator) printResult(result *grawl.Result) {
	result.PrintRowColored()
}
