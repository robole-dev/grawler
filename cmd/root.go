package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/cobra"
	url2 "net/url"
	"os"
	"robole-dev/cache-warmer/internal/request_result"
	"sort"
	"sync"
	"time"
)

var (
	// Used for flags.
	flagParallel       int
	flagDelay          int64
	flagMaxDepth       int
	flagOutputFilename string

	rootCmd = &cobra.Command{
		Use:   "cache-warmer <url>",
		Short: "A simple cache warming application.",
		Long: `This app scrapes the website of the given url and finds all relative links and visit these urls. So it 
forces the website to create the cached pages.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			warmItUp(url, flagParallel, flagDelay)
		},
		Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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

	fmt.Println("Warming " + url)

	var (
		requestCount   = 0
		responseCount  = 0
		requestResults []request_result.RequestResult
		startTimes     sync.Map
	)

	c := colly.NewCollector()
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: parallel,
		Delay:       time.Duration(delay) * time.Millisecond,
	})
	if err != nil {
		fmt.Println("Error setting limits:", err)
		//return
	}

	c.MaxDepth = flagMaxDepth
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		//if strings.HasPrefix(link, "/") {
		//	//
		//} else if strings.HasPrefix(link, url) {
		//	//
		//} else {
		//	return
		//}

		// Print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		_ = c.Visit(e.Request.AbsoluteURL(link))
	})

	parsedUrl, err := url2.Parse(url)
	if err != nil {
		fmt.Println("Error parsing the url:", err)
		return
	}

	c.AllowedDomains = []string{
		parsedUrl.Host,
	}

	c.OnRequest(func(r *colly.Request) {
		startTimes.Store(r.ID, time.Now())
		requestCount++
		fmt.Printf("%d - Visiting: %s\n", requestCount, r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		responseCount++

		reqResult := request_result.RequestResult{}
		reqResult.Index = responseCount
		reqResult.Url = r.Request.URL.String()
		reqResult.UrlPath = r.Request.URL.Path
		reqResult.UrlHost = r.Request.URL.Host
		reqResult.UrlParmeters = r.Request.URL.RawQuery
		reqResult.UrlFragment = r.Request.URL.RawFragment
		reqResult.StatusCode = r.StatusCode
		reqResult.ResponseAt = time.Now()

		if startTime, ok := startTimes.Load(r.Request.ID); ok {
			duration := time.Since(startTime.(time.Time))
			reqResult.Duration = duration
			startTimes.Delete(r.Request.ID)
			//fmt.Printf("Request to %s took %v\n", r.Request.URL, duration)
		} else {
			fmt.Printf("No start time found for %s\n", r.Request.URL)
		}

		requestResults = append(requestResults, reqResult)
	})

	err = c.Visit(url)
	if err != nil {
		fmt.Printf("Could not visit url: %v\n", err)
		return
	}

	if len(flagOutputFilename) > 0 {
		saveResult(&requestResults)
	}
}

func saveResult(results *[]request_result.RequestResult) {
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

	headers := request_result.GetCsvHeader()
	if err := writer.Write(headers); err != nil {
		panic(err)
	}

	for _, result := range *results {
		if err := writer.Write(result.GetCsvRow()); err != nil {
			panic(err)
		}
	}
}
