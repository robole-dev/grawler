package cmd

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/fatih/color"
	"github.com/gocolly/colly/v2"
	"github.com/manifoldco/promptui"
	"github.com/robole-dev/grawler/internal/request"
	"github.com/spf13/cobra"
	url2 "net/url"
	"os"
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
	flagUserAgent      string
	flagSitemap        bool

	headerAuth string

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
	grawlCmd.Flags().StringVarP(&flagUsername, "username", "u", "", "Use this for HTTP Basic Authentication. If you omit the password-flag a prompt will ask for the password.")
	grawlCmd.Flags().StringVarP(&flagPassword, "password", "p", "", "Use this for HTTP Basic Authentication.")
	grawlCmd.Flags().StringVar(&flagUserAgent, "user-agent", "", "Sets the user agent.")
	grawlCmd.Flags().BoolVarP(&flagSitemap, "sitemap", "s", false, "Checks the sitemap. If this is flag is set the url parameter has to be the url to the sitemap.xml.")
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

	if flagUserAgent != "" {
		c.UserAgent = flagUserAgent
	}

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: flagParallel,
		Delay:       time.Duration(flagDelay) * time.Millisecond,
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

	if flagUsername != "" {
		if flagPassword == "" {
			flagPassword, err = promptPassword()
			if err != nil {
				fmt.Println("error reading password:", err)
				return
			}
		}

		var auth = base64.StdEncoding.EncodeToString([]byte(flagUsername + ":" + flagPassword))
		headerAuth = fmt.Sprintf("Basic %s", auth)
	}

	c.OnRequest(func(r *colly.Request) {
		if headerAuth != "" {
			r.Headers.Set("Authorization", headerAuth)
		}

		requestResult := request.NewResult(r.ID, r.URL.String())

		runningRequests.Store(r.ID, requestResult, r.URL.String())
		requestCount++

		r.Ctx.Put("orgUrl", r.URL.String())

		//fmt.Printf("%d - Visiting: %s\n", requestCount, r.URL)
		//fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		responseCount++

		if reqResult, ok := runningRequests.Load(r.Request.ID); ok {
			duration := time.Since(reqResult.GetRequestAt())
			totalDuration += duration
			reqResult.UpdateOnResponse(r, responseCount, duration, nil)
			printResult(reqResult)
		} else {
			fmt.Printf("No start time found for %s\n", r.Request.URL)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		errorCount++

		if reqResult, ok := runningRequests.Load(r.Request.ID); ok {
			duration := time.Since(reqResult.GetRequestAt())
			totalDuration += duration
			reqResult.UpdateOnResponse(r, responseCount, duration, &err)
			//fmt.Println("error:", err)
		} else {
			fmt.Printf("Could not find request: %s\n", r.Request.URL)
		}
	})

	if flagSitemap {
		fmt.Println("Look for sitemap")
		c.OnXML("//urlset/url/loc", func(e *colly.XMLElement) {
			_ = c.Visit(e.Request.AbsoluteURL(e.Text))
		})
		c.OnXML("//sitemapindex/sitemap/loc", func(e *colly.XMLElement) {
			_ = c.Visit(e.Request.AbsoluteURL(e.Text))
		})
	} else {
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			_ = c.Visit(e.Request.AbsoluteURL(link))
		})
	}

	err = c.Visit(url)
	if err != nil {
		fmt.Printf("Could not visit url: %v\n", err)
		return
	}

	if len(flagOutputFilename) > 0 {
		saveResult(runningRequests)
	}

	fmt.Println("")
	fmt.Println("Grawling finished at:    ", time.Now().Format(request.DateFormat))
	fmt.Println("Requests:                ", requestCount)
	fmt.Println("duration:                ", totalDuration.Round(time.Millisecond))
	fmt.Println("Average request duration:", time.Duration(int64(totalDuration)/int64(requestCount)).Round(time.Millisecond))
	fmt.Println("Total response num:      ", responseCount)
	fmt.Println("Total error/skipped num: ", errorCount)
}

func saveResult(runningRequests *request.RunningRequests) {
	fmt.Printf("Saving file \"%s\".\n", flagOutputFilename)

	results := runningRequests.GetValues()

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

func printResult(result *request.Result) {
	if result.IsRedirected() {
		color.Yellow(result.GetPrintRow())
	} else if result.HasError() {
		color.Red(result.GetPrintRow())
	} else {
		color.Green(result.GetPrintRow())
	}
}
