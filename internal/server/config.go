package server

import (
	"os"

	"github.com/spf13/pflag"
)

// Config stores all crond server configurations.
type Config struct {
	ServerPort     int    `mapstructure:"server-port"`
	RaftProduction bool   `mapstructure:"raft-production"`
	RaftNode       string `mapstructure:"raft-node"`
	RaftBootstrap  bool   `mapstructure:"raft-bootstrap"`
	RaftDataDir    string `mapstructure:"raft-data-dir"`
}

// DefaultConfig creates the Config with sensible default settings.
func DefaultConfig() *Config {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return &Config{
		ServerPort:     5281,
		RaftProduction: false,
		RaftNode:       name,
		RaftBootstrap:  false,
		RaftDataDir:    "data",
	}
}

// BindFlags overwrites default crond server configurations from CLI flags.
func BindFlags(c *Config, fs *pflag.FlagSet) {
	fs.IntVar(&c.ServerPort, "server-port", c.ServerPort, "server server port")
	fs.BoolVar(&c.RaftProduction, "raft-production", c.RaftProduction, "if true, raft layer runs in production mode")
	fs.StringVar(&c.RaftNode, "raft-node", c.RaftNode, "raft layer node name")
	fs.BoolVar(&c.RaftBootstrap, "raft-bootstrap", c.RaftBootstrap, "if true, raft layer will bootstrap cluster")
	fs.StringVar(&c.RaftDataDir, "raft-data-dir", c.RaftDataDir, "raft layer persists data in this specific directory")
}
