package cmd

import (
	"github.com/robole-dev/grawler/internal/grawl"
	"github.com/spf13/cobra"
)

var (
	flagParallel       int
	flagDelay          int64
	flagMaxDepth       int
	flagOutputFilename string
	flagUsername       string
	flagPassword       string
	flagUserAgent      string
	flagSitemap        bool

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
	grawler := grawl.NewGrawler()
	grawler.FlagParallel = flagParallel
	grawler.FlagDelay = flagDelay
	grawler.FlagMaxDepth = flagMaxDepth
	grawler.FlagOutputFilename = flagOutputFilename
	grawler.FlagUsername = flagUsername
	grawler.FlagPassword = flagPassword
	grawler.FlagUserAgent = flagUserAgent
	grawler.FlagSitemap = flagSitemap

	grawler.Grawl(url)
}
