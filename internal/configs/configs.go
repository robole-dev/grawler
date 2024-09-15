package configs

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CrawlerConfDir() string {
	return filepath.Join(HomeDir(), ".config", "grawler")
}

func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error trying to find users home directory")
	}

	return home
}

func DefaultConfFile() string {
	return filepath.Join(CrawlerConfDir(), DefaultConfFileName()+"."+DefaultConfFileType())
}

func DefaultConfFileName() string {
	return "conf"
}

func DefaultConfFileType() string {
	return "yaml"
}

func BindViperFlag(cobraCommand *cobra.Command, keyPrefix string, flagLookup string) {
	//viper.BindPFlag(keyPrefix+"."+toSnakeCase(flagLookup), cobraCommand.Flags().Lookup(flagLookup))
	viper.BindPFlag(keyPrefix+"."+flagLookup, cobraCommand.Flags().Lookup(flagLookup))
}

func toCamelCase(key string) string {
	words := strings.Split(key, "-")
	caser := cases.Title(language.AmericanEnglish)
	for i := 1; i < len(words); i++ {
		words[i] = caser.String(words[i])
	}

	return strings.Join(words, "")
}

func toSnakeCase(key string) string {
	return strings.Replace(key, "-", "_", -1)
}

func DefaultLocalConfFile() string {
	return "grawler.yaml"
}
