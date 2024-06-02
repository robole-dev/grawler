package cmd

import (
	"github.com/robole-dev/grawler/internal/grawl"
	"github.com/spf13/cobra"
)

var (
	flags    = grawl.Flags{}
	grawlCmd = &cobra.Command{
		Use:     "grawl",
		Aliases: []string{"crawl"},
		Short:   "Crawls the given url",
		Long:    `The grawler searches for href-attributes and crawls these urls too. It also crawles sitmap.xmls.`,
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			warmItUp(url)
		},
		Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}
)

func init() {
	grawlCmd.Flags().Int64VarP(&flags.FlagDelay, "delay", "d", 0, "Delay between requests in milliseconds. (default 0)")
	grawlCmd.Flags().IntVarP(&flags.FlagMaxDepth, "max-depth", "m", 0, "Set it to 0 for infinite recursion. (default 0)")
	grawlCmd.Flags().StringVarP(&flags.FlagOutputFilename, "output-filepath", "o", "", "Write statistic data of each request to this file.")
	grawlCmd.Flags().IntVarP(&flags.FlagParallel, "parallel", "l", 1, "Number of parallel requests.")
	grawlCmd.Flags().StringVarP(&flags.FlagUsername, "username", "u", "", "Use this for HTTP Basic Authentication. If you omit the password-flag a prompt will ask for the password.")
	grawlCmd.Flags().StringVarP(&flags.FlagPassword, "password", "p", "", "Use this for HTTP Basic Authentication.")
	grawlCmd.Flags().StringVar(&flags.FlagUserAgent, "user-agent", "grawler", "Sets the user agent.")
	grawlCmd.Flags().BoolVarP(&flags.FlagSitemap, "sitemap", "s", false, "Checks the sitemap. If this is flag is set the url parameter has to be the url to the sitemap.xml.")
	grawlCmd.Flags().StringSliceVarP(&flags.FlagAllowedDomains, "allowed-domains", "a", nil, "A comma separated list of allowed domains to be crawled. The domain of the given url is always allowed.")
	grawlCmd.Flags().BoolVar(&flags.FlagRespectRobotsTxt, "respect-robots-txt", false, "Respect the robots.txt file.")
	grawlCmd.Flags().StringVar(&flags.FlagPath, "path", "", "Restrict the crawlings on a certain url path.")
	grawlCmd.Flags().BoolVarP(&flags.FlagCheckAll, "check-all", "", false, "In addtion to html and xml-urls, also check image, js and css-urls, among others.")
}

func warmItUp(url string) {
	grawler := grawl.NewGrawler(flags)
	grawler.Grawl(url)
}
