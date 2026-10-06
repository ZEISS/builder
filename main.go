package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/zeiss/builder/cmd"
)

const (
	versionFmt = "%s (%s %s)"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	v := fmt.Sprintf(versionFmt, version, commit, date)

	if err := fang.Execute(
		context.Background(),
		cmd.RootCmd,
		fang.WithVersion(v),
	); err != nil {
		os.Exit(1)
	}
}
