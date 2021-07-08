package config

import (
	"github.com/KevinWu0904/crond/internal/crond"
	"github.com/KevinWu0904/crond/internal/raft"
	"github.com/KevinWu0904/crond/pkg/logs"
)

// Config stores all crond configurations.
type Config struct {
	Logger *logs.LoggerConfig
	Server *crond.ServerConfig
	Raft   *raft.LayerConfig
}

// DefaultConfig creates the Config with sensible default settings.
func DefaultConfig() *Config {
	return &Config{
		Logger: logs.DefaultLoggerConfig(),
		Server: crond.DefaultServerConfig(),
		Raft:   raft.DefaultLayerConfig(),
	}
}
