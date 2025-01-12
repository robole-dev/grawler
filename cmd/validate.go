package cmd

import (
	"fmt"
	"github.com/robole-dev/grawler/internal/validates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

var (
	validateFlags = validates.Flags{}
	validateCmd   = &cobra.Command{
		Use:   "validate",
		Short: "Reads a given csv file and checks for the same return codes.",
		//Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			csvFilePath := args[0]

			validator := validates.NewValidator(validateFlags)
			validator.ValidateCsv(csvFilePath)
		},
		Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}
)

func init() {
	validateCmd.Flags().Uint64Var(&validateFlags.FlagSkipRows, validates.FlagSkipRows, 1, "The number of csv rows to skip when parsing.")
	bindValidateViperFlag(validates.FlagSkipRows)

	validateCmd.Flags().Uint64Var(&validateFlags.FlagColUrl, validates.FlagColUrl, 3, "The csv column containing the urls. The column number is zero based.")
	bindValidateViperFlag(validates.FlagColUrl)

	validateCmd.Flags().Uint64Var(&validateFlags.FlagColStatusCode, validates.FlagColStatusCode, 1, "The csv column containing the html status codes. The column number is zero based.")
	bindValidateViperFlag(validates.FlagColStatusCode)

	validateCmd.Flags().Uint64Var(&validateFlags.FlagColContentType, validates.FlagColContentType, 5, "The csv column containing the content type returned from the server. The column number is zero based.")
	bindValidateViperFlag(validates.FlagColContentType)
}

func bindValidateViperFlag(flagLookup string) {
	key := validates.ViperValidatePrefix + "." + flagLookup
	err := viper.BindPFlag(key, validateCmd.Flags().Lookup(flagLookup))
	if err != nil {
		log.Fatalln(fmt.Errorf("error binding config option to flag: %v", err))
		return
	}
}
