package validates

const (
	ViperValidatePrefix = "validate"
	FlagSkipRows        = "skip-rows"
	FlagColUrl          = "col-url"
	FlagColStatusCode   = "col-status-code"
	FlagColContentType  = "col-content-type"
)

type Flags struct {
	FlagSkipRows       uint64
	FlagColUrl         uint64
	FlagColStatusCode  uint64
	FlagColContentType uint64
	//FlagParallel
	//FlagRandomDelay
	//FlagDelay
	//FlagRequestTimeout
	//FlagStopOnError
	//FlagPauseOnError
}
