package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "grawler",
	Short: "A simple url scraping application.",
	Long:  `This app scrapes the website of the given url and finds all relative links and visit these urls.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//},
	//Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
}

func init() {
	rootCmd.AddCommand(grawlCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
