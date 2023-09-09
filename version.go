package main

import (
	"fmt"
)

var (
	// BuildDate date string of when build was performed filled in by -X compile flag
	BuildDate string

	// LatestCommit date string of when build was performed filled in by -X compile flag
	LatestCommit string

	// Version string of build filled in by -X compile flag
	Version string

	// GitRepo string of the git repo url when build was performed filled in by -X compile flag
	GitRepo string

	// GitBranch string of branch in the git repo filled in by -X compile flag
	GitBranch string
)

func getVersion() string {
	return fmt.Sprintf("tldr %s", Version)
}
