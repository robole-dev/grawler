package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/robole-dev/grawler/internal/request"
	"github.com/spf13/cobra"
	url2 "net/url"
	"os"
	"sort"
	"time"
)

var (
	flagParallel       int
	flagDelay          int64
	flagMaxDepth       int
	flagOutputFilename string

	rootCmd = &cobra.Command{
		Use:   "grawler <url>",
		Short: "A simple url scraping application.",
		Long:  `This app scrapes the website of the given url and finds all relative links and visit these urls.`,
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			warmItUp(url, flagParallel, flagDelay)
		},
		Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&flagParallel, "parallel", "p", 1, "Number of parallel requests. (default 1)")
	rootCmd.Flags().Int64VarP(&flagDelay, "delay", "d", 0, "Delay between requests in milliseconds. (default 0)")
	rootCmd.Flags().IntVarP(&flagMaxDepth, "max-depth", "m", 0, "Set it to 0 for infinite recursion (default 0).")
	rootCmd.Flags().StringVarP(&flagOutputFilename, "output-filepath", "o", "", "The statistic data is written to this file. Leave empty for no ouput file (default).")
}

func warmItUp(url string, parallel int, delay int64) {

	fmt.Println("Grawling " + url)

	var (
		requestCount  = 0
		responseCount = 0
		errorCount    = 0
		totalDuration = time.Duration(0)
		//requestResults  []result.Result
		//startTimes      sync.Map
		runningRequests = request.NewRunningRequests()
	)

	c := colly.NewCollector()
	c.MaxDepth = flagMaxDepth

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: parallel,
		Delay:       time.Duration(delay) * time.Millisecond,
	})
	if err != nil {
		fmt.Println("Error setting limits:", err)
		return
	}

	parsedUrl, err := url2.Parse(url)
	if err != nil {
		fmt.Println("Error parsing the url:", err)
		return
	}

	c.AllowURLRevisit = false
	c.AllowedDomains = []string{
		parsedUrl.Host,
	}

	c.OnRequest(func(r *colly.Request) {

		if ref := r.Ctx.Get("_referer"); ref != "" {
			fmt.Println("Header:", ref)
			//r.Headers.Set("Referer", ref)
		}

		requestResult := request.NewResult()
		requestResult.RequestAt = time.Now()
		requestResult.Url = r.URL.String()

		runningRequests.Store(r.ID, requestResult)
		requestCount++
		//fmt.Printf("%d - Visiting: %s\n", requestCount, r.URL)
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		responseCount++

		if reqResult, ok := runningRequests.Load(r.Request.ID); ok {
			duration := time.Since(reqResult.RequestAt)
			totalDuration += duration
			reqResult.UpdateOnResponse(r, responseCount, duration, nil)
		} else {
			fmt.Printf("No start time found for %s\n", r.Request.URL)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		errorCount++

		if reqResult, ok := runningRequests.Load(r.Request.ID); ok {
			duration := time.Since(reqResult.RequestAt)
			totalDuration += duration
			reqResult.UpdateOnResponse(r, responseCount, duration, &err)
			fmt.Println("Error:", err)
		} else {
			fmt.Printf("Could not find request: %s\n", r.Request.URL)
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		_ = c.Visit(e.Request.AbsoluteURL(link))
	})

	err = c.Visit(url)
	if err != nil {
		fmt.Printf("Could not visit url: %v\n", err)
		return
	}

	if len(flagOutputFilename) > 0 {
		saveResult(runningRequests.GetValues())
	}

	runningRequests = nil

	fmt.Println("")
	fmt.Println("Grawling finished at:    ", time.Now())
	fmt.Println("Total request num:       ", requestCount)
	fmt.Println("Total request duration:  ", totalDuration.Round(time.Millisecond))
	fmt.Println("Average request duration:", time.Duration(int64(totalDuration)/int64(requestCount)).Round(time.Millisecond))
	fmt.Println("Total response num:      ", responseCount)
	fmt.Println("Total error/skipped num: ", errorCount)
}

func saveResult(results *[]*request.Result) {
	fmt.Printf("Saving file \"%s\".\n", flagOutputFilename)

	// Sort by URL
	sort.Slice(*results, func(i, j int) bool {
		return (*results)[i].Url < (*results)[j].Url
	})

	file, err := os.Create(flagOutputFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	headers := request.GetCsvHeader()
	if err := writer.Write(headers); err != nil {
		panic(err)
	}

	for _, result := range *results {
		if err := writer.Write(result.GetCsvRow()); err != nil {
			panic(err)
		}
	}
}
