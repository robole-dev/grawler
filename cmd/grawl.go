package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/manifoldco/promptui"
	"github.com/robole-dev/grawler/internal/request"
	"github.com/spf13/cobra"
	url2 "net/url"
	"os"
	"sort"
	"time"
)

type promptContent struct {
	errorMsg string
	label    string
}

var (
	flagParallel       int
	flagDelay          int64
	flagMaxDepth       int
	flagOutputFilename string
	flagUsername       string
	flagPassword       string
	headerAuth         string

	emptyFlagValue string = string(rune(0))

	grawlCmd = &cobra.Command{
		Use:     "grawl",
		Aliases: []string{"crawl"},
		Short:   "Crawls the given url",
		Long:    `The grawler searches for href-attributes and crawls these urls too.`,
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			warmItUp(url)
		},
		Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}
)

func init() {
	grawlCmd.Flags().Int64VarP(&flagDelay, "delay", "d", 0, "Delay between requests in milliseconds. (default 0)")
	grawlCmd.Flags().IntVarP(&flagMaxDepth, "max-depth", "m", 0, "Set it to 0 for infinite recursion. (default 0)")
	grawlCmd.Flags().StringVarP(&flagOutputFilename, "output-filepath", "o", "", "Write statistic data of each request to this file.")
	grawlCmd.Flags().IntVarP(&flagParallel, "parallel", "l", 1, "Number of parallel requests.")
	grawlCmd.Flags().StringVarP(&flagUsername, "username", "u", "", "Use this for HTTP Basic Authentication.")

	grawlCmd.Flags().StringVarP(&flagPassword, "password", "p", "", "Use this for HTTP Basic Authentication. Leave empty for password prompt.")
	grawlCmd.Flags().Lookup("password").NoOptDefVal = emptyFlagValue
	//grawlCmd.Flags().Lookup("password").DefValue = ""

	//BindPFlags(grawlCmd.Flags())
	//BindPFlags("port", grawlCmd.Flags().Lookup("port"))
}

func warmItUp(url string) {

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
		Parallelism: flagParallel,
		Delay:       time.Duration(flagDelay) * time.Millisecond,
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

	//if flagUsername != "" {
	//	fmt.Printf("pass '%v' \n", flagPassword)
	//	if flagPassword == emptyFlagValue {
	//		flagPassword, err = promptPassword()
	//		if err != nil {
	//			fmt.Println("Error reading password:", err)
	//			return
	//		}
	//
	//		var auth = base64.StdEncoding.EncodeToString([]byte(flagUsername + ":" + flagPassword))
	//		headerAuth = fmt.Sprintf("Basic %s", auth)
	//	}
	//}

	c.OnRequest(func(r *colly.Request) {
		if headerAuth != "" {
			r.Headers.Set("Authorization", headerAuth)
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

func promptPassword() (string, error) {
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
