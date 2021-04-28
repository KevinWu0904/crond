package main

import (
	"github.com/KevinWu0904/crond/cmd/crond-api/app"
	"github.com/KevinWu0904/crond/pkg/logs"
)

func main() {
	logs.Init()
	defer logs.Flush()

	app.NewAPIServerCommand().Execute()
}
