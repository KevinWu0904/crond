package main

import (
	"github.com/KevinWu0904/crond/cmd/crond-api/app"
	"github.com/KevinWu0904/crond/pkg/logs"
)

func main() {
	options := app.NewAPIServerOptions()
	command := app.NewAPIServerCommand(options)

	logs.Init(options.LoggerOptions)
	defer logs.Flush()

	if err := command.Execute(); err != nil {
		logs.Error("Execute crond-api server command failed: err=%v", err)
	}
}
