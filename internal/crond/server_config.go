package crond

import (
	"github.com/spf13/pflag"
)

// ServerConfig stores all crond server configurations.
type ServerConfig struct {
	ServerPort int
}

// DefaultServerConfig creates the ServerConfig with sensible default settings.
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		ServerPort: 5281,
	}
}

// BindServerFlags overwrites default server configurations from CLI flags.
func BindServerFlags(c *ServerConfig, fs *pflag.FlagSet) {
	fs.IntVar(&c.ServerPort, "server-port", c.ServerPort, "crond server port")
}
