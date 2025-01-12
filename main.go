package main

import "github.com/robole-dev/grawler/cmd"

var (
	commit  string
	version string
	date    string
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}
