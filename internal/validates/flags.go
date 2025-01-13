package validates

const (
	ViperValidatePrefix    = "validate"
	FlagNameSkipRows       = "skip-rows"
	FlagNameColUrl         = "col-url"
	FlagNameColStatusCode  = "col-status-code"
	FlagNameColContentType = "col-content-type"
	FlagNameParallel       = "parallel"
	FlagNameRandomDelay    = "random-delay"
	FlagNameDelay          = "delay"
	FlagNameRequestTimeout = "request-timeout"
	FlagNameStopOnError    = "stop-on-error"
	FlagNamePauseOnError   = "pause-on-error"
)

type Flags struct {
	FlagSkipRows       uint64
	FlagColUrl         uint64
	FlagColStatusCode  uint64
	FlagColContentType uint64
	FlagParallel       int
	FlagRandomDelay    int64
	FlagDelay          int64
	FlagRequestTimeout float32
	FlagStopOnError    bool
	FlagPauseOnError   bool
}
