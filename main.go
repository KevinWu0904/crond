package main

import (
	"os"

	"github.com/KevinWu0904/crond/cmd"
)

func main() {
	if err := cmd.RootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
