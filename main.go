package main

import "github.com/jmull3n/issuemetrics/cmd"

// these values go from the go build, do not change them
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute(version, commit, date)
}
