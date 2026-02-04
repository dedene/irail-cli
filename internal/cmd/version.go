package cmd

import (
	"fmt"
	"os"
	"runtime"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

type VersionCmd struct{}

func (c *VersionCmd) Run() error {
	fmt.Fprintf(os.Stdout, "irail %s\n", VersionString())
	fmt.Fprintf(os.Stdout, "commit: %s\n", commit)
	fmt.Fprintf(os.Stdout, "built:  %s\n", date)
	fmt.Fprintf(os.Stdout, "go:     %s\n", runtime.Version())

	return nil
}

func VersionString() string {
	return version
}
