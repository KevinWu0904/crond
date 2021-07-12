package config

import (
	"github.com/KevinWu0904/crond/internal/crond"
	"github.com/KevinWu0904/crond/internal/crond/raft"
	"github.com/KevinWu0904/crond/pkg/logs"
)

// Config stores all crond configurations.
type Config struct {
	Logger *logs.LoggerConfig  `mapstructure:"logger"`
	Server *crond.ServerConfig `mapstructure:"server"`
	Raft   *raft.LayerConfig   `mapstructure:"raft"`
}

// DefaultConfig creates the Config with sensible default settings.
func DefaultConfig() *Config {
	return &Config{
		Logger: logs.DefaultLoggerConfig(),
		Server: crond.DefaultServerConfig(),
		Raft:   raft.DefaultLayerConfig(),
	}
}
