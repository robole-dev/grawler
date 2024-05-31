package grawl

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/fatih/color"
	"github.com/gocolly/colly/v2"
	"github.com/manifoldco/promptui"
	url2 "net/url"
	"os"
	"strings"
	"time"
)

const (
	DateFormat = "2006-01-02 15:04:05.000"
	ctxOrgUrl  = "orgUrl"
)

type Grawler struct {
	FlagParallel       int
	FlagDelay          int64
	FlagMaxDepth       int
	FlagOutputFilename string
	FlagUsername       string
	FlagPassword       string
	FlagUserAgent      string
	FlagSitemap        bool

	headerAuth string

	requestCount    uint32
	responseCount   int
	errorCount      uint32
	totalDuration   time.Duration
	runningRequests *RunningRequests
}

func NewGrawler() *Grawler {
	return &Grawler{

		//FlagParallel       :int,
		//FlagDelay          :int64,
		FlagMaxDepth: 0,
		//FlagOutputFilename :string,
		//FlagUsername       :string,
		//FlagPassword       :string,
		//FlagUserAgent      :string,
		//FlagSitemap        :bool,

		requestCount:    0,
		responseCount:   0,
		errorCount:      0,
		totalDuration:   time.Duration(0),
		runningRequests: NewRunningRequests(),
	}
}

func (g *Grawler) Grawl(url string) {

	fmt.Println("Grawling " + url)

	c := colly.NewCollector()
	c.MaxDepth = g.FlagMaxDepth

	if g.FlagUserAgent != "" {
		c.UserAgent = g.FlagUserAgent
	}

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: g.FlagParallel,
		Delay:       time.Duration(g.FlagDelay) * time.Millisecond,
	})
	if err != nil {
		fmt.Println("error setting limits:", err)
		return
	}

	parsedUrl, err := url2.Parse(url)
	if err != nil {
		fmt.Println("error parsing the url:", err)
		return
	}

	c.AllowURLRevisit = false
	c.AllowedDomains = []string{
		parsedUrl.Host,
	}

	if g.FlagUsername != "" {
		if g.FlagPassword == "" {
			g.FlagPassword, err = g.promptPassword()
			if err != nil {
				fmt.Println("error reading password:", err)
				return
			}
		}

		var auth = base64.StdEncoding.EncodeToString([]byte(g.FlagUsername + ":" + g.FlagPassword))
		g.headerAuth = fmt.Sprintf("Basic %s", auth)
	}

	//c.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
	//	if len(via) > 0 {
	//		fmt.Println("Redirecting to:", req.URL)
	//	}
	//	return http.ErrUseLastResponse
	//})

	c.OnRequest(func(r *colly.Request) {
		url = r.URL.String()
		r.Ctx.Put(ctxOrgUrl, url)

		if g.headerAuth != "" {
			r.Headers.Set("Authorization", g.headerAuth)
		}

		foundOnUrl := g.runningRequests.GetFoundUrl(url)
		requestResult := NewResult(r.ID, url, foundOnUrl)

		g.runningRequests.Store(r.ID, requestResult, url)
		g.requestCount++
	})

	c.OnResponse(func(r *colly.Response) {
		g.responseCount++

		if reqResult, ok := g.runningRequests.Load(r.Request.ID); ok {
			duration := time.Since(reqResult.GetRequestAt())
			g.totalDuration += duration

			reqResult.UpdateOnResponse(r, g.responseCount, duration, nil)
			g.printResult(reqResult)
		} else {
			fmt.Printf("No start time found for %s\n", r.Request.URL)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		g.errorCount++

		//
		// Remove request if this url is filtered by colly
		//
		if r.StatusCode == 0 {
			g.runningRequests.Delete(r.Request.ID)
			return
		}

		if reqResult, ok := g.runningRequests.Load(r.Request.ID); ok {
			duration := time.Since(reqResult.GetRequestAt())
			g.totalDuration += duration
			reqResult.UpdateOnResponse(r, g.responseCount, duration, &err)
			g.printResult(reqResult)
			//fmt.Println("error:", err)
		} else {
			fmt.Printf("Could not find request: %s\n", r.Request.URL)
		}
	})

	if g.FlagSitemap {
		fmt.Println("Look for sitemap")
		c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
			g.visit(c, e.Request.AbsoluteURL(e.Text), e.Request.URL.String())
		})
		c.OnXML("//sitemapindex/sitemap/loc", func(e *colly.XMLElement) {
			g.visit(c, e.Request.AbsoluteURL(e.Text), e.Request.URL.String())
		})
	} else {
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			strings.Trim(link, " ")

			if strings.HasPrefix(link, "mailto:") {
				return
			}

			g.visit(c, e.Request.AbsoluteURL(link), e.Request.URL.String())
		})
	}

	err = c.Visit(url)
	if err != nil {
		fmt.Printf("Could not visit url: %v\n", err)
		return
	}

	if len(g.FlagOutputFilename) > 0 {
		g.saveResultFile(g.runningRequests)
	}

	g.printSummary()
}

func (g *Grawler) visit(c *colly.Collector, url string, foundOnUrl string) {
	visited, err := c.HasVisited(url)
	if err != nil || !visited {
		g.runningRequests.AddFoundUrl(url, foundOnUrl)
	}

	_ = c.Visit(url)
}

func (g *Grawler) printSummary() {
	fmt.Println("")
	fmt.Println("Grawling finished at:", time.Now().Format(DateFormat))
	fmt.Println("Duration:            ", g.totalDuration.Round(time.Millisecond))
	fmt.Println("Avg request duration:", time.Duration(int64(g.totalDuration)/int64(g.requestCount)).Round(time.Millisecond))
	fmt.Println("Responses:           ", g.responseCount)
	fmt.Println("Requests:            ", g.requestCount)
	fmt.Println("Errors/Skipped:      ", g.errorCount)
}

func (g *Grawler) saveResultFile(runningRequests *RunningRequests) {
	fmt.Printf("Saving file \"%s\".\n", g.FlagOutputFilename)

	results := runningRequests.GetValues()

	file, err := os.Create(g.FlagOutputFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	headers := GetCsvHeader()
	if err := writer.Write(headers); err != nil {
		panic(err)
	}

	for _, result := range *results {
		if err := writer.Write(result.GetCsvRow()); err != nil {
			panic(err)
		}
	}
}

func (g *Grawler) promptPassword() (string, error) {
	validate := func(input string) error {
		return nil
		//_, err := strconv.Parse(input, 64)
		//if err != nil {
		//	return errors.New("Invalid number")
		//}
		//return nil
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

	//fmt.Printf("You choose %q\n", result)
	return result, nil
}

func (g *Grawler) printResult(result *Result) {
	if result.IsRedirected() {
		color.Yellow(result.GetPrintRow())
	} else if result.HasError() {
		color.Red(result.GetPrintRow())
	} else {
		color.Green(result.GetPrintRow())
	}
}
