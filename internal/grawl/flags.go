package grawl

type Flags struct {
	FlagParallel         int
	FlagDelay            int64
	FlagMaxDepth         int
	FlagOutputFilename   string
	FlagUsername         string
	FlagPassword         string
	FlagUserAgent        string
	FlagSitemap          bool
	FlagAllowedDomains   []string
	FlagRespectRobotsTxt bool
	FlagPath             string
	FlagCheckAll         bool
}
