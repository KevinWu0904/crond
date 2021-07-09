package raft

import (
	"os"

	"github.com/spf13/pflag"
)

// LayerConfig stores all crond raft layer configurations.
type LayerConfig struct {
	RaftProduction bool   `mapstructure:"raft-production"`
	RaftNode       string `mapstructure:"raft-node"`
	RaftBootstrap  bool   `mapstructure:"raft-bootstrap"`
	RaftDataDir    string `mapstructure:"raft-data-dir"`
}

// DefaultLayerConfig creates the LayerConfig with sensible default settings.
func DefaultLayerConfig() *LayerConfig {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return &LayerConfig{
		RaftProduction: false,
		RaftNode:       name,
		RaftBootstrap:  false,
		RaftDataDir:    "data",
	}
}

// BindLayerFlags overwrites default raft layer configurations from CLI flags.
func BindLayerFlags(c *LayerConfig, fs *pflag.FlagSet) {
	fs.BoolVar(&c.RaftProduction, "raft-production", c.RaftProduction, "if true, raft layer runs in production mode")
	fs.StringVar(&c.RaftNode, "raft-node", c.RaftNode, "raft layer node name")
	fs.BoolVar(&c.RaftBootstrap, "raft-bootstrap", c.RaftBootstrap, "if true, raft layer will bootstrap cluster")
	fs.StringVar(&c.RaftDataDir, "raft-data-dir", c.RaftDataDir, "raft layer persists data in this specific directory")
}
