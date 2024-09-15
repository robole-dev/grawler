package cmd

import (
	"errors"
	"fmt"
	"github.com/robole-dev/grawler/internal/configs"
	"github.com/robole-dev/grawler/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var (
	versionInfo    *version.Info = nil
	flagVersion    bool
	flagConfigPath string
	rootCmd        = &cobra.Command{
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
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(grawlCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.Flags().BoolVarP(&flagVersion, "version", "v", false, "Show version")
	rootCmd.PersistentFlags().StringVar(&flagConfigPath, "config", "", "Manually set the path to your config file.")
}

func initConfig() {
	curDir, _ := os.Getwd()
	localConfig := filepath.Join(curDir, configs.DefaultLocalConfFile())

	configFilePath := flagConfigPath

	if flagConfigPath != "" {
		viper.SetConfigFile(flagConfigPath)
	} else if _, err := os.Stat(localConfig); !errors.Is(err, os.ErrNotExist) {
		configFilePath = localConfig
		viper.SetConfigFile(localConfig)
	} else {
		configFilePath = configs.DefaultConfFile()
		viper.SetConfigType(configs.DefaultConfFileType())
		viper.AddConfigPath(configs.CrawlerConfDir())
		viper.SetConfigName(configs.DefaultConfFileName())
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			//fmt.Println("Config file not found.", flagConfigPath)
		} else {
			log.Fatalln(fmt.Sprintf("Something unexpected happened reading configuration file: %s, err: %s", configFilePath, err))
		}
	}
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
