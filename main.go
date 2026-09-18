package main

import (
	"context"
	"log"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/zeiss/builder/cmd"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	if err := fang.Execute(context.Background(), cmd.RootCmd); err != nil {
		os.Exit(1)
	}
}
