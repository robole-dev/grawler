package cmd

import (
	"fmt"
	"github.com/robole-dev/grawler/internal/version"
	"github.com/spf13/cobra"
	"os"
)

var (
	versionInfo *version.Info = nil
	flagVersion bool
	rootCmd     = &cobra.Command{
		Use:   "grawler",
		Short: "A simple url scraping application.",
		Long:  `This app scrapes the website of the given url and finds all relative links and visit these urls.`,
		Run: func(cmd *cobra.Command, args []string) {
			if flagVersion {
				printVersion()
			}
		},
		//Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}
)

func init() {
	rootCmd.AddCommand(grawlCmd)
	rootCmd.Flags().BoolVarP(&flagVersion, "version", "v", false, "Show version")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func SetVersionInfo(vers string, commit string, date string) {
	versionInfo = version.NewVersion(vers, commit, date)
}

func printVersion() {
	fmt.Println(versionInfo.ToString())
}
