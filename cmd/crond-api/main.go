package main

import (
	"github.com/KevinWu0904/crond/cmd/crond-api/app"
	"github.com/KevinWu0904/crond/pkg/logs"
)

func main() {
	if err := app.NewAPIServerCommand().Execute(); err != nil {
		logs.Error("Execute crond-api server command failed: err=%v", err)
	}
}
