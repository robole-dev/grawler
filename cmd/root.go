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
		Short: "A simple web crawling application.",
		Long:  `The grawler scrapes and visit urls from websites and sitemaps and provides some statistics.`,
		Run: func(cmd *cobra.Command, args []string) {
			if flagVersion {
				printVersion()
			} else {
				_ = cmd.Help()
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
