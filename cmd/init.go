package cmd

import (
	"fmt"
	"github.com/robole-dev/grawler/internal/configs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var (
	initFlagUseHome = false
	initCmd         = &cobra.Command{
		Use:   "init",
		Short: "Writes a config file with default values",
		Long:  "Used without any flags it writes the config file \"grawler.yaml\" to the current directory.",
		Run: func(cmd *cobra.Command, args []string) {
			configFilePath := ""
			if len(args) > 0 {
				configFilePath = args[0]
			}
			writeConfigFile(configFilePath)
		},
		Args: cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	}
)

func init() {
	initCmd.Flags().BoolVar(&initFlagUseHome, "home", false, fmt.Sprintf("If enabled it writes the config file to the path \"%s\"", configs.DefaultConfFile()))
}

func writeConfigFile(configFilePathParam string) {
	curDir, _ := os.Getwd()
	localConfig := filepath.Join(curDir, configs.DefaultLocalConfFile())

	writePath := localConfig

	if configFilePathParam != "" {
		writePath = configFilePathParam
	}

	fmt.Println("Writing default config to", writePath)

	err := viper.SafeWriteConfigAs(writePath)
	if err != nil {
		log.Fatalln(fmt.Errorf("write config file error: %v", err))
		return
	}
}
