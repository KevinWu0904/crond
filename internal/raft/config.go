package raft

import (
	"os"

	"github.com/spf13/pflag"
)

// LayerConfig stores all crond raft layer configurations.
type LayerConfig struct {
	EnableDebug bool
	NodeName    string
	Bootstrap   bool
	DataDir     string
}

// DefaultLayerConfig creates the LayerConfig with sensible default settings.
func DefaultLayerConfig() *LayerConfig {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return &LayerConfig{
		EnableDebug: false,
		NodeName:    name,
		Bootstrap:   false,
		DataDir:     "data",
	}
}

// BindLayerFlags overwrites default raft layer configurations from CLI flags.
func BindLayerFlags(c *LayerConfig, fs *pflag.FlagSet) {
	fs.BoolVar(&c.EnableDebug, "enable-debug", c.EnableDebug, "raft layer enables debug")
	fs.StringVar(&c.NodeName, "node-name", c.NodeName, "raft layer os hostname")
	fs.BoolVar(&c.Bootstrap, "bootstrap", c.Bootstrap, "raft layer will bootstrap raft cluster")
	fs.StringVar(&c.DataDir, "data-dir", c.DataDir, "raft layer persists data in this specific directory")
}
