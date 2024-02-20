package cmd

import (
	"fmt"
	"os"
)

var (
	commit    = "unknown"
	version   = "unreleased"
	buildDate = "unknown"
)

func printVersion() {
	fmt.Fprintf(os.Stdout, "* Commit: %s\n", commit)
	fmt.Fprintf(os.Stdout, "* Version: %s\n", version)
	fmt.Fprintf(os.Stdout, "* Build Date: %s\n", buildDate)
}
