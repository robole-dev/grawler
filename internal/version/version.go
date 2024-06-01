package version

import (
	"fmt"
	"runtime"
)

// const GrawlerLogo = "\n\n                            .__                \n   ________________ __  _  _|  |   ___________ \n  / ___\\_  __ \\__  \\\\ \\/ \\/ /  | _/ __ \\_  __ \\\n / /_/  >  | \\// __ \\\\     /|  |_\\  ___/|  | \\/\n \\___  /|__|  (____  /\\/\\_/ |____/\\___  >__|   \n/_____/            \\/                 \\/       \n\n"

//                                __
//     ____ __________ __      __/ ___  _____
//    / __ `/ ___/ __ `| | /| / / / _ \/ ___/
//   / /_/ / /  / /_/ /| |/ |/ / /  __/ /
//   \__, /_/   \__,_/ |__/|__/_/\___/_/
//  /____/

//const GrawlerLogo = `
//                            .__
//   ________________ __  _  _|  |   ___________
//  / ___\_  __ \__  \\ \/ \/ /  | _/ __ \_  __ \
// / /_/  >  | \// __ \\     /|  |_\  ___/|  | \/
// \___  /|__|  (____  /\/\_/ |____/\___  >__|
///_____/            \/                 \/
//`

// https://patorjk.com/software/taag/#p=display&h=3&f=Slant&t=grawler

//const GrawlerLogo = "                              __         \n   ____ __________ __      __/ ___  _____\n  / __ `/ ___/ __ `| | /| / / / _ \\/ ___/\n / /_/ / /  / /_/ /| |/ |/ / /  __/ /    \n \\__, /_/   \\__,_/ |__/|__/_/\\___/_/     \n/____/                                   "

const GrawlerLogo = "                          _           \n  __ _ _ __ __ ___      _| | ___ _ __ \n / _` | '__/ _` \\ \\ /\\ / | |/ _ | '__|\n| (_| | | | (_| |\\ V  V /| |  __| |   \n \\__, |_|  \\__,_| \\_/\\_/ |_|\\___|_|   \n |___/                                "

type Info struct {
	Version string
	Commit  string
	Date    string
}

func NewVersion(version string, commit string, date string) *Info {
	return &Info{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}

func (v *Info) ToString() string {
	return ToString(v.Version, v.Commit, v.Date)
}

func ToString(version string, commit string, date string) string {
	return fmt.Sprintf(
		`%s

A simple url grawling application.
https://github.com/robole-dev/grawler

Version:  %s
Date:     %s 
Commit:   %s
Platform: %s/%s`,
		GrawlerLogo,
		version,
		date,
		commit,
		runtime.GOOS,
		runtime.GOARCH,
	)
}
