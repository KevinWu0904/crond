package crond

import (
	"github.com/KevinWu0904/crond/pkg/logs"
)

// Config stores all crond configurations.
type Config struct {
	Logger *logs.LoggerConfig
}

// DefaultConfig creates the Config with sensible default settings.
func DefaultConfig() *Config {
	return &Config{
		Logger: logs.DefaultLoggerConfig(),
	}
}
