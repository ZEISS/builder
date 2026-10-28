package main

import (
	"context"
	"log"
	"os"

	"github.com/zeiss/builder/cmd"
	"github.com/zeiss/builder/internal/versions"

	"github.com/charmbracelet/fang"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	err := fang.Execute(
		context.Background(),
		cmd.RootCmd,
		fang.WithVersion(versions.VersionString),
	)
	if err != nil {
		os.Exit(1)
	}
}
