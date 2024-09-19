package cmd

import (
	"fmt"
	"github.com/robole-dev/grawler/internal/grawl"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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

const (
	viperGrawlPrefix             = "grawl"
	flagNameDelay                = "delay"
	flagNameRandomDelay          = "random-delay"
	flagNameMaxDepth             = "max-depth"
	flagNameOutputFilepath       = "output-filepath"
	flagNameParallel             = "parallel"
	flagNameUsername             = "username"
	flagNamePassword             = "password"
	flagNameUserAgent            = "user-agent"
	flagNameSitemap              = "sitemap"
	flagNameAllowedDomains       = "allowed-domains"
	flagNameRespectRobotsTxt     = "respect-robots-txt"
	flagNamePath                 = "path"
	flagNameCheckAll             = "check-all"
	flagNameRequestTimeout       = "request-timeout"
	flagNameUrlFilters           = "url-filters"
	flagNameDisallowedURLFilters = "disallowed-url-filters"
)

func init() {
	grawlCmd.Flags().Int64VarP(&grawlFlags.FlagDelay, flagNameDelay, "d", 0, "Delay between requests in milliseconds. (default 0)")
	bindViperFlag(flagNameDelay)

	grawlCmd.Flags().Int64Var(&grawlFlags.FlagRandomDelay, flagNameRandomDelay, 0, "Max random delay between requests in milliseconds. (default 0 for no random delay)")
	bindViperFlag(flagNameRandomDelay)

	grawlCmd.Flags().IntVarP(&grawlFlags.FlagMaxDepth, flagNameMaxDepth, "m", 0, "Set it to 0 for infinite recursion. (default 0)")
	bindViperFlag(flagNameMaxDepth)

	grawlCmd.Flags().StringVarP(&grawlFlags.FlagOutputFilename, flagNameOutputFilepath, "o", "", "Write statistic data of each request to this file.")
	bindViperFlag(flagNameOutputFilepath)

	grawlCmd.Flags().IntVarP(&grawlFlags.FlagParallel, flagNameParallel, "l", 1, "Number of parallel requests.")
	bindViperFlag(flagNameParallel)

	grawlCmd.Flags().StringVarP(&grawlFlags.FlagUsername, flagNameUsername, "u", "", "Use this for HTTP Basic Authentication. If you omit the password-flag a prompt will ask for the password.")
	bindViperFlag(flagNameUsername)

	grawlCmd.Flags().StringVarP(&grawlFlags.FlagPassword, flagNamePassword, "p", "", "Use this for HTTP Basic Authentication.")
	bindViperFlag(flagNamePassword)

	grawlCmd.Flags().StringVar(&grawlFlags.FlagUserAgent, flagNameUserAgent, "grawler", "Sets the user agent.")
	bindViperFlag(flagNameUserAgent)

	grawlCmd.Flags().BoolVarP(&grawlFlags.FlagSitemap, flagNameSitemap, "s", false, "Checks the sitemap. If this is flag is set the url parameter has to be the url to the sitemap.xml.")
	bindViperFlag(flagNameSitemap)

	grawlCmd.Flags().StringSliceVarP(&grawlFlags.FlagAllowedDomains, flagNameAllowedDomains, "a", nil, "A comma separated list of allowed domains to be crawled. The domain of the given url is always allowed.")
	bindViperFlag(flagNameAllowedDomains)

	grawlCmd.Flags().BoolVar(&grawlFlags.FlagRespectRobotsTxt, flagNameRespectRobotsTxt, false, "Respect the robots.txt file.")
	bindViperFlag(flagNameRespectRobotsTxt)

	grawlCmd.Flags().StringVar(&grawlFlags.FlagPath, flagNamePath, "", "Restrict the crawlings on a certain url path.")
	bindViperFlag(flagNamePath)

	grawlCmd.Flags().BoolVarP(&grawlFlags.FlagCheckAll, flagNameCheckAll, "", false, "In addtion to html and xml-urls, also check image, js and css-urls, among others.")
	bindViperFlag(flagNameCheckAll)

	grawlCmd.Flags().Float32Var(&grawlFlags.FlagRequestTimeout, flagNameRequestTimeout, 10, "Timeout in seconds to wait for a response.")
	bindViperFlag(flagNameRequestTimeout)

	grawlCmd.Flags().StringSliceVar(&grawlFlags.FlagURLFilters, flagNameUrlFilters, nil, "Only visit urls that match the regular expressions given here.")
	bindViperFlag(flagNameUrlFilters)

	grawlCmd.Flags().StringSliceVar(&grawlFlags.FlagDisallowedURLFilters, flagNameDisallowedURLFilters, nil, "Do not visit urls that match the regular expressions given here.")
	bindViperFlag(flagNameDisallowedURLFilters)
}

func warmItUp(url string) {

	// Get values from viper back to flag vars
	grawlFlags.FlagDelay = viper.GetInt64(viperGrawlPrefix + "." + flagNameDelay)
	grawlFlags.FlagRandomDelay = viper.GetInt64(viperGrawlPrefix + "." + flagNameRandomDelay)
	grawlFlags.FlagMaxDepth = viper.GetInt(viperGrawlPrefix + "." + flagNameMaxDepth)
	grawlFlags.FlagOutputFilename = viper.GetString(viperGrawlPrefix + "." + flagNameOutputFilepath)
	grawlFlags.FlagParallel = viper.GetInt(viperGrawlPrefix + "." + flagNameParallel)
	grawlFlags.FlagUsername = viper.GetString(viperGrawlPrefix + "." + flagNameUsername)
	grawlFlags.FlagPassword = viper.GetString(viperGrawlPrefix + "." + flagNamePassword)
	grawlFlags.FlagUserAgent = viper.GetString(viperGrawlPrefix + "." + flagNameUserAgent)
	grawlFlags.FlagSitemap = viper.GetBool(viperGrawlPrefix + "." + flagNameSitemap)
	grawlFlags.FlagAllowedDomains = viper.GetStringSlice(viperGrawlPrefix + "." + flagNameAllowedDomains)
	grawlFlags.FlagRespectRobotsTxt = viper.GetBool(viperGrawlPrefix + "." + flagNameRespectRobotsTxt)
	grawlFlags.FlagPath = viper.GetString(viperGrawlPrefix + "." + flagNamePath)
	grawlFlags.FlagCheckAll = viper.GetBool(viperGrawlPrefix + "." + flagNameCheckAll)
	grawlFlags.FlagRequestTimeout = cast.ToFloat32(viper.Get(viperGrawlPrefix + "." + flagNameRequestTimeout))
	grawlFlags.FlagURLFilters = viper.GetStringSlice(viperGrawlPrefix + "." + flagNameUrlFilters)
	grawlFlags.FlagDisallowedURLFilters = viper.GetStringSlice(viperGrawlPrefix + "." + flagNameDisallowedURLFilters)

	if flagConfigInfo {
		fmt.Println("")
		fmt.Println("Grawl configuration values")
		fmt.Println("==========================")
		fmt.Println("Url:", url)
		fmt.Println("Delay:", grawlFlags.FlagDelay)
		fmt.Println("RandomDelay:", grawlFlags.FlagRandomDelay)
		fmt.Println("MaxDepth:", grawlFlags.FlagMaxDepth)
		fmt.Println("OutputFilepath:", grawlFlags.FlagOutputFilename)
		fmt.Println("Parallel:", grawlFlags.FlagParallel)
		fmt.Println("Username:", grawlFlags.FlagUsername)
		fmt.Println("Password:", grawlFlags.FlagPassword)
		fmt.Println("UserAgent:", grawlFlags.FlagUserAgent)
		fmt.Println("Sitemap:", grawlFlags.FlagSitemap)
		fmt.Println("AllowedDomains:", grawlFlags.FlagAllowedDomains)
		fmt.Println("RespectRobotsTxt:", grawlFlags.FlagRespectRobotsTxt)
		fmt.Println("Path:", grawlFlags.FlagPath)
		fmt.Println("CheckAll:", grawlFlags.FlagCheckAll)
		fmt.Println("RequestTimeout:", grawlFlags.FlagRequestTimeout)
		fmt.Println("URLFilters:", grawlFlags.FlagURLFilters)
		fmt.Println("DisallowedURLFilters:", grawlFlags.FlagDisallowedURLFilters)
		fmt.Println("")
	}

	grawler := grawl.NewGrawler(grawlFlags)
	grawler.Grawl(url)
}

func bindViperFlag(flagLookup string) {
	key := viperGrawlPrefix + "." + flagLookup
	err := viper.BindPFlag(key, grawlCmd.Flags().Lookup(flagLookup))
	if err != nil {
		log.Fatalln(fmt.Errorf("error binding config option to flag: %v", err))
		return
	}
}
