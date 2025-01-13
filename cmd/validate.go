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
	validateCmd.Flags().Uint64Var(&validateFlags.FlagSkipRows, validates.FlagNameSkipRows, 1, "The number of csv rows to skip when parsing.")
	bindValidateViperFlag(validates.FlagNameSkipRows)

	validateCmd.Flags().Uint64Var(&validateFlags.FlagColUrl, validates.FlagNameColUrl, 3, "The csv column containing the urls. The column number is zero based.")
	bindValidateViperFlag(validates.FlagNameColUrl)

	validateCmd.Flags().Uint64Var(&validateFlags.FlagColStatusCode, validates.FlagNameColStatusCode, 1, "The csv column containing the html status codes. The column number is zero based.")
	bindValidateViperFlag(validates.FlagNameColStatusCode)

	validateCmd.Flags().Uint64Var(&validateFlags.FlagColContentType, validates.FlagNameColContentType, 5, "The csv column containing the content type returned from the server. The column number is zero based.")
	bindValidateViperFlag(validates.FlagNameColContentType)

	validateCmd.Flags().IntVarP(&validateFlags.FlagParallel, validates.FlagNameParallel, "l", 1, "Number of parallel requests.")
	bindValidateViperFlag(validates.FlagNameParallel)

	validateCmd.Flags().Int64Var(&validateFlags.FlagRandomDelay, validates.FlagNameRandomDelay, 0, "Max random delay between requests in milliseconds. (default 0 for no random delay)")
	bindValidateViperFlag(validates.FlagNameRandomDelay)

	validateCmd.Flags().Int64VarP(&validateFlags.FlagDelay, validates.FlagNameDelay, "d", 0, "Delay between requests in milliseconds. (default 0)")
	bindValidateViperFlag(validates.FlagNameDelay)

	validateCmd.Flags().Float32Var(&validateFlags.FlagRequestTimeout, validates.FlagNameRequestTimeout, 10, "Timeout in seconds to wait for a response.")
	bindValidateViperFlag(validates.FlagNameRequestTimeout)

	validateCmd.Flags().BoolVar(&validateFlags.FlagStopOnError, validates.FlagNameStopOnError, false, "The validation stops on errors.")
	bindValidateViperFlag(validates.FlagNameStopOnError)

	validateCmd.Flags().BoolVar(&validateFlags.FlagPauseOnError, validates.FlagNamePauseOnError, false, "The validation pauses on errors and you have the option to cancel, skip or try again.")
	bindValidateViperFlag(validates.FlagNamePauseOnError)
}

func bindValidateViperFlag(flagLookup string) {
	key := validates.ViperValidatePrefix + "." + flagLookup
	err := viper.BindPFlag(key, validateCmd.Flags().Lookup(flagLookup))
	if err != nil {
		log.Fatalln(fmt.Errorf("error binding config option to flag: %v", err))
		return
	}
}
