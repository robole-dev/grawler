package grawl

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gocolly/colly/v2"
	"github.com/manifoldco/promptui"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DateFormat = "2006-01-02 15:04:05.000"
	ctxOrgUrl  = "initialRequestUrl"
)

type Grawler struct {
	flags               Flags
	headerAuth          string
	requestCount        uint32
	responseCount       int
	errorCount          uint32
	totalDuration       time.Duration
	runningRequests     *RunningRequests
	fileWriter          *FileWriter
	responseErrorRanges *responseCodeRanges
	collector           *colly.Collector
	redirections        uint32
	visitMutex          sync.Mutex
}

func NewGrawler(flags Flags) *Grawler {
	//errorCodeRanges, err := newResponseCodeRanges(flags.FlagResponseErrorCodes)
	errorCodeRanges, err := newResponseCodeRanges([]string{})
	if err != nil {
		panic(err)
	}

	return &Grawler{
		flags:               flags,
		requestCount:        0,
		responseCount:       0,
		errorCount:          0,
		totalDuration:       time.Duration(0),
		runningRequests:     NewRunningRequests(),
		responseErrorRanges: errorCodeRanges,
		redirections:        0,
	}
}

func (g *Grawler) Grawl(grawlUrl string) {

	fmt.Println("Grawling " + grawlUrl)

	parsedUrl, err := url.Parse(grawlUrl)
	if err != nil {
		fmt.Println("Error parsing the grawlUrl:", err)
		return
	}

	c := colly.NewCollector()
	g.collector = c
	c.MaxDepth = g.flags.FlagMaxDepth
	c.Async = true
	c.SetRequestTimeout(time.Duration(g.flags.FlagRequestTimeout * float32(time.Second)))

	if g.flags.FlagPath != "" {
		regexPatternPath := fmt.Sprintf(
			`^https?://%s%s.*$`,
			regexp.QuoteMeta(parsedUrl.Host),
			regexp.QuoteMeta(g.flags.FlagPath),
		)
		regexPath := regexp.MustCompile(regexPatternPath)
		c.URLFilters = append(c.URLFilters, regexPath)

		regexPatternUrl := fmt.Sprintf(
			`^%s$`,
			regexp.QuoteMeta(grawlUrl),
		)
		regexUrl := regexp.MustCompile(regexPatternUrl)
		c.URLFilters = append(c.URLFilters, regexUrl)
	}

	if g.flags.FlagUserAgent != "" {
		c.UserAgent = g.flags.FlagUserAgent
	}

	limitingRule := &colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: g.flags.FlagParallel,
	}

	if g.flags.FlagRandomDelay > 0 {
		limitingRule.RandomDelay = time.Duration(g.flags.FlagRandomDelay) * time.Millisecond
	} else {
		limitingRule.Delay = time.Duration(g.flags.FlagDelay) * time.Millisecond
	}

	err = c.Limit(limitingRule)
	if err != nil {
		fmt.Println("Error setting limits:", err)
		return
	}

	c.IgnoreRobotsTxt = !g.flags.FlagRespectRobotsTxt
	c.AllowURLRevisit = false
	c.AllowedDomains = slices.Concat(c.AllowedDomains, g.flags.FlagAllowedDomains)
	c.AllowedDomains = append(c.AllowedDomains, parsedUrl.Host)

	if len(g.flags.FlagURLFilters) > 0 {
		c.URLFilters = append(c.URLFilters, regexp.MustCompile("^"+grawlUrl+"$"))
		for _, filter := range g.flags.FlagURLFilters {
			c.URLFilters = append(c.URLFilters, regexp.MustCompile(filter))
		}
	}

	if len(g.flags.FlagDisallowedURLFilters) > 0 {
		for _, filter := range g.flags.FlagDisallowedURLFilters {
			c.DisallowedURLFilters = append(c.DisallowedURLFilters, regexp.MustCompile(filter))
		}
	}

	if g.flags.FlagUsername != "" {
		if g.flags.FlagPassword == "" {
			g.flags.FlagPassword, err = g.promptPassword()
			if err != nil {
				fmt.Println("error reading password:", err)
				return
			}
		}

		var auth = base64.StdEncoding.EncodeToString([]byte(g.flags.FlagUsername + ":" + g.flags.FlagPassword))
		g.headerAuth = fmt.Sprintf("Basic %s", auth)
	}

	c.OnRequest(g.onRequest)
	c.OnResponse(g.onResponse)
	c.OnError(g.onError)

	if g.flags.FlagSitemap {
		c.OnXML("//urlset/grawlUrl/loc", func(e *colly.XMLElement) {
			g.visit(c, e.Request, e.Request.AbsoluteURL(e.Text), e.Request.URL.String())
		})
		c.OnXML("//sitemapindex/sitemap/loc", func(e *colly.XMLElement) {
			g.visit(c, e.Request, e.Request.AbsoluteURL(e.Text), e.Request.URL.String())
		})
	} else {
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			if strings.HasPrefix(link, "mailto:") {
				return
			}
			if strings.HasPrefix(link, "tel:") {
				return
			}
			if g.flags.FlagRespectNofollow {
				if strings.ToLower(e.Attr("rel")) == "nofollow" {
					return
				}
			}

			g.visit(c, e.Request, e.Request.AbsoluteURL(link), e.Request.URL.String())
		})
	}

	if g.flags.FlagCheckAll {
		c.OnHTML("source[srcset]", func(e *colly.HTMLElement) {
			imgSrc := e.Attr("srcset")
			g.visit(c, e.Request, e.Request.AbsoluteURL(imgSrc), e.Request.URL.String())
		})

		c.OnHTML("img[src]", func(e *colly.HTMLElement) {
			imgSrc := e.Attr("src")
			g.visit(c, e.Request, e.Request.AbsoluteURL(imgSrc), e.Request.URL.String())
		})

		c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
			cssHref := e.Attr("href")
			g.visit(c, e.Request, e.Request.AbsoluteURL(cssHref), e.Request.URL.String())
		})

		c.OnHTML("script[src]", func(e *colly.HTMLElement) {
			scriptSrc := e.Attr("src")
			g.visit(c, e.Request, e.Request.AbsoluteURL(scriptSrc), e.Request.URL.String())
		})

		c.OnResponseHeaders(func(r *colly.Response) {
			if isHtmlResponse(r) || isXmlResponse(r) {
				return
			}

			//
			// Abort downloading all non-xml and non-html contents
			//
			r.Request.Abort()
			reqResult, ok := g.runningRequests.Load(r.Request.ID)
			if ok {
				duration := time.Since(reqResult.GetRequestAt())
				g.totalDuration += duration

				reqResult.UpdateOnResponse(r, g.responseCount, duration, nil)
				g.printResult(reqResult)
				g.checkStopOnError(reqResult)
			} else {
				fmt.Println("Request data not found", r.Request.URL)
			}
		})
	}

	c.SetRedirectHandler(g.onRedirect)

	if g.flags.FlagOutputFilename != "" {
		g.fileWriter = NewFileWriter(g.flags.FlagOutputFilename)
		g.fileWriter.InitFile()
	}

	err = c.Visit(grawlUrl)
	if err != nil {
		fmt.Printf("Could not visit grawlUrl: %v\n", err)
		return
	}
	c.Wait()

	g.printSummary()
}

func (g *Grawler) onRequest(r *colly.Request) {
	requestUrl := r.URL.String()

	if g.headerAuth != "" {
		r.Headers.Set("Authorization", g.headerAuth)
	}

	foundOnUrl := g.runningRequests.GetFoundUrl(requestUrl)
	requestResult := NewResult(r.ID, requestUrl, foundOnUrl, g.responseErrorRanges)

	g.runningRequests.Store(r.ID, requestResult, requestUrl)
	g.requestCount++
}

func (g *Grawler) onResponse(r *colly.Response) {
	g.responseCount++

	reqResult, ok := g.runningRequests.Load(r.Request.ID)
	if !ok {
		fmt.Printf("No start time found for %s\n", r.Request.URL)
	}

	duration := time.Since(reqResult.GetRequestAt())
	g.totalDuration += duration

	reqResult.UpdateOnResponse(r, g.responseCount, duration, nil)
	g.printResult(reqResult)
	g.checkStopOnError(reqResult)
}

func (g *Grawler) onRedirect(req *http.Request, via []*http.Request) error {
	runningReq, ok := g.runningRequests.LoadByUrl(via[0].URL.String())
	atomic.AddUint32(&g.redirections, 1)
	if ok {
		fmt.Printf("Redirecting to %s from %s. ID: %d\n", req.URL, via[0].URL, runningReq.id)
	} else {
		return fmt.Errorf("Could not find initial url of redirection to %s from %s\n", req.URL, via[0].URL)
	}
	return nil
}

func (g *Grawler) onError(r *colly.Response, err error) {
	g.errorCount++

	// Normal error on aborted binary files like images. Result is printed in OnResponseHeaders
	if err != nil && errors.Is(err, colly.ErrAbortedAfterHeaders) {
		return
	}

	//
	// Remove request if this url is filtered by colly
	//
	if r.StatusCode == 0 {
		fmt.Println("Error", r.Request.URL, err)
		g.runningRequests.Delete(r.Request.ID)
		return
	}

	reqResult, ok := g.runningRequests.Load(r.Request.ID)
	if ok {
		duration := time.Since(reqResult.GetRequestAt())
		reqResult.UpdateOnResponse(r, g.responseCount, duration, &err)
		g.totalDuration += duration
		g.printResult(reqResult)
		g.checkStopOnError(reqResult)
	} else {
		fmt.Println("Request data not found", r.Request.URL)
	}
}

func (g *Grawler) visit(c *colly.Collector, r *colly.Request, url string, foundOnUrl string) {
	g.visitMutex.Lock()

	url = strings.Trim(url, " ")
	visited, err := c.HasVisited(url)
	if visited {
		//fmt.Println("Visited", url)
		g.visitMutex.Unlock()
		return
	}
	if err != nil {
		fmt.Println("Could not check if url has been visited: ", url)
	}

	hasFoundUrl := g.runningRequests.HasFoundUrl(url)
	if hasFoundUrl {
		g.visitMutex.Unlock()
		return
	}

	g.runningRequests.AddFoundUrl(url, foundOnUrl)
	//fmt.Println("Visit:", url)
	_ = r.Visit(url)
	g.visitMutex.Unlock()
}

func (g *Grawler) printSummary() {
	durationMin := time.Hour
	durationMax := time.Duration(0)
	returnCodes := map[int]int{}
	returnErrors := 0
	for _, result := range *g.runningRequests.GetValues() {
		durationMin = min(durationMin, result.duration)
		durationMax = max(durationMax, result.duration)
		if result.statusCode > 0 {
			returnCodes[result.statusCode]++
		} else {
			returnErrors++
		}
	}

	returnCodeKeys := make([]int, 0, len(returnCodes))
	for key := range returnCodes {
		returnCodeKeys = append(returnCodeKeys, key)
	}
	sort.Ints(returnCodeKeys)

	fmt.Println("")
	fmt.Println("Grawling finished at:", time.Now().Format(DateFormat))
	fmt.Println("Duration:            ", g.totalDuration.Round(time.Millisecond))
	fmt.Println("  - Min:             ", durationMin.Round(time.Millisecond))
	fmt.Println("  - Max:             ", durationMax.Round(time.Millisecond))
	fmt.Println("  - Avg:             ", time.Duration(int64(g.totalDuration)/int64(g.requestCount)).Round(time.Millisecond))
	fmt.Println("Requests:            ", g.requestCount)
	for _, code := range returnCodeKeys {
		fmt.Printf("  - Status code %d:  %d\n", code, returnCodes[code])
	}
	fmt.Printf("  - Other errors:     %d\n", returnErrors)
	fmt.Printf("  - Redirections:     %d\n", g.redirections)
}

func (g *Grawler) printResult(result *Result) {
	if result.IsRedirected() {
		color.Yellow(result.GetPrintRow())
	} else if result.HasError() {
		color.Red(result.GetPrintRow())
	} else {
		color.Green(result.GetPrintRow())
	}

	if g.fileWriter != nil {
		g.fileWriter.WriteResultLine(result)
	}

	g.updateStatusBar(result)
}

func (g *Grawler) updateStatusBar(result *Result) {
	//fmt.Printf("queue %d", g.collector.)
	//line := "cool: " + result.GetPrintRow()
	//
	//printFixedBottomLine(line)
}

func printFixedBottomLine(line string) {
	// ANSI Code um den Cursor auf die letzte Zeile zu setzen
	fmt.Printf("\033[s")              // Speichere die aktuelle Cursor-Position
	fmt.Printf("\033[H\033[2J")       // Lösche den Bildschirm
	fmt.Printf("\033[999;1H%s", line) // Setze Cursor an die letzte Zeile und drucke die Zeile
	fmt.Printf("\033[u")              // Stelle die gespeicherte Cursor-Position wieder her
}

func (g *Grawler) promptPassword() (string, error) {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Password",
		Validate: validate,
		Mask:     '*',
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}

	return result, nil
}

func (g *Grawler) checkStopOnError(result *Result) {
	if !result.HasError() {
		return
	}

	if g.flags.FlagStopOnError {
		fmt.Println("Stop grawling after error.")
		if result.error != nil {
			fmt.Printf("Grawling error: %s\n", result.error.Error())
		}
		os.Exit(1)
	}

	if g.flags.FlagPauseOnError {
		fmt.Println("Pause grawling after error.")
		prompt := g.promptResume()
		switch prompt {
		case "a":
			fmt.Println("Grawling aborted.")
			os.Exit(1)
		case "s":
			fmt.Println("Url skipped.")
			return
		case "r":
			fmt.Println("Retry url...")
			allowRevisit := g.collector.AllowURLRevisit
			g.collector.AllowURLRevisit = true
			_ = g.collector.Visit(result.url)
			g.collector.AllowURLRevisit = allowRevisit
		}
	}
}

func (g *Grawler) promptResume() interface{} {
	str := "Please choose: [a]bort the grawling or [r]etry the url or [s]kip the url."
	for {
		fmt.Println(str)
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			return nil
		}
		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "a", "s", "r":
			return input
		default:
			//fmt.Println(str)
		}
	}
}
