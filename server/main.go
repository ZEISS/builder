package main

import (
	"github.com/zeiss/builder/server/cmd"
)

func main() {
	err := cmd.Root.Execute()
	if err != nil {
		panic(err)
	}
}
