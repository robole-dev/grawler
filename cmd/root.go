package cmd

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cache-warmer",
		Short: "A simple cache warming application.",
		Long: `This app scrapes the website of the given url and finds all relative links and visit these urls. So it 
forces the website to create the cached pages.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			warmItUp(url)
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cache-warmer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func warmItUp(url string) {

	fmt.Println("Warming " + url)

	requestCount := 0

	c := colly.NewCollector()

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		if strings.HasPrefix(link, "/") {
			//
		} else if strings.HasPrefix(link, url) {
			//
		} else {
			return
		}

		// Print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		_ = c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		requestCount++
		fmt.Printf("%d - Visiting: %s\n", requestCount, r.URL)
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Printf("Could not visit url: %v\n", err)
		return
	}
}
