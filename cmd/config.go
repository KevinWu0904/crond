package cmd

import (
	"github.com/KevinWu0904/crond/internal/server"
	"github.com/KevinWu0904/crond/pkg/logs"
)

// Config stores all crond configurations.
type Config struct {
	*RootConfig   `mapstructure:",squash"`
	*ServerConfig `mapstructure:",squash"`
}

// DefaultConfig creates the Config with sensible default settings.
func DefaultConfig() *Config {
	return &Config{
		RootConfig: &RootConfig{
			Logger: logs.DefaultConfig(),
		},
		ServerConfig: &ServerConfig{
			Server: server.DefaultConfig(),
		},
	}
}

// RootConfig stores crond root command configurations.
type RootConfig struct {
	Logger *logs.Config `mapstructure:"logger"`
}

// ServerConfig stores crond server command configurations.
type ServerConfig struct {
	Server *server.Config `mapstructure:"server"`
}
