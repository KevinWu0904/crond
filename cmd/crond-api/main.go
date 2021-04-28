package main

import (
	"os"

	"github.com/KevinWu0904/crond/cmd/crond-api/app"
	"github.com/KevinWu0904/crond/pkg/logs"
)

func main() {
	logs.Init()
	defer logs.Flush()

	if err := app.NewAPIServerCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
