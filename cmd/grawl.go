package cmd

import (
	"github.com/robole-dev/grawler/internal/configs"
	"github.com/robole-dev/grawler/internal/grawl"
	"github.com/spf13/cobra"
)

var (
	grawlFlags = grawl.Flags{}
	grawlCmd   = &cobra.Command{
		Use:     "grawl",
		Aliases: []string{"crawl"},
		Short:   "Crawls the given url",
		Long:    `This command scrapes and visits all urls from a page or uses an existing sitemap.xml.`,
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			warmItUp(url)
		},
		Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}
)

func init() {
	keyPrefix := "grawl"

	grawlCmd.Flags().Int64VarP(&grawlFlags.FlagDelay, "delay", "d", 0, "Delay between requests in milliseconds. (default 0)")
	configs.BindViperFlag(grawlCmd, keyPrefix, "delay")

	grawlCmd.Flags().Int64Var(&grawlFlags.FlagRandomDelay, "random-delay", 0, "Max random delay between requests in milliseconds. (default 0 for no random delay)")
	configs.BindViperFlag(grawlCmd, keyPrefix, "random-delay")

	grawlCmd.Flags().IntVarP(&grawlFlags.FlagMaxDepth, "max-depth", "m", 0, "Set it to 0 for infinite recursion. (default 0)")
	configs.BindViperFlag(grawlCmd, keyPrefix, "max-depth")

	grawlCmd.Flags().StringVarP(&grawlFlags.FlagOutputFilename, "output-filepath", "o", "", "Write statistic data of each request to this file.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "output-filepath")

	grawlCmd.Flags().IntVarP(&grawlFlags.FlagParallel, "parallel", "l", 1, "Number of parallel requests.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "parallel")

	grawlCmd.Flags().StringVarP(&grawlFlags.FlagUsername, "username", "u", "", "Use this for HTTP Basic Authentication. If you omit the password-flag a prompt will ask for the password.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "username")

	grawlCmd.Flags().StringVarP(&grawlFlags.FlagPassword, "password", "p", "", "Use this for HTTP Basic Authentication.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "password")

	grawlCmd.Flags().StringVar(&grawlFlags.FlagUserAgent, "user-agent", "grawler", "Sets the user agent.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "user-agent")

	grawlCmd.Flags().BoolVarP(&grawlFlags.FlagSitemap, "sitemap", "s", false, "Checks the sitemap. If this is flag is set the url parameter has to be the url to the sitemap.xml.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "sitemap")

	grawlCmd.Flags().StringSliceVarP(&grawlFlags.FlagAllowedDomains, "allowed-domains", "a", nil, "A comma separated list of allowed domains to be crawled. The domain of the given url is always allowed.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "allowed-domains")

	grawlCmd.Flags().BoolVar(&grawlFlags.FlagRespectRobotsTxt, "respect-robots-txt", false, "Respect the robots.txt file.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "respect-robots-txt")

	grawlCmd.Flags().StringVar(&grawlFlags.FlagPath, "path", "", "Restrict the crawlings on a certain url path.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "path")

	grawlCmd.Flags().BoolVarP(&grawlFlags.FlagCheckAll, "check-all", "", false, "In addtion to html and xml-urls, also check image, js and css-urls, among others.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "check-all")

	grawlCmd.Flags().Float32Var(&grawlFlags.FlagRequestTimeout, "request-timeout", 10, "Timeout in seconds to wait for a response.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "request-timeout")

	grawlCmd.Flags().StringSliceVar(&grawlFlags.FlagURLFilters, "url-filters", nil, "Only visit urls that match the regular expressions given here.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "url-filters")

	grawlCmd.Flags().StringSliceVar(&grawlFlags.FlagDisallowedURLFilters, "disallowed-url-filters", nil, "Do not visit urls that match the regular expressions given here.")
	configs.BindViperFlag(grawlCmd, keyPrefix, "disallowed-url-filters")
}

func warmItUp(url string) {
	grawler := grawl.NewGrawler(grawlFlags)
	grawler.Grawl(url)
}
