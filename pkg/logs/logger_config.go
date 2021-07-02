package logs

import "github.com/spf13/pflag"

// LoggerConfig stores all crond logger configurations.
type LoggerConfig struct {
	LogLevel         string
	LogEncoder       string
	EnableConsoleLog bool
	EnableFileLog    bool
	FileLogDir       string
	FileLogName      string
	FileLogNum       int
	FileLogSize      int
	FileLogAge       int
}

// DefaultLoggerConfig creates the LoggerConfig with sensible default settings.
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		LogLevel:         "info",
		LogEncoder:       "console",
		EnableConsoleLog: true,
		EnableFileLog:    false,
		FileLogDir:       "",
		FileLogName:      "",
		FileLogNum:       3,
		FileLogSize:      500,
		FileLogAge:       15,
	}
}

// BindLoggerFlags overwrites default logger configurations from CLI flags.
func BindLoggerFlags(c *LoggerConfig, fs *pflag.FlagSet) {
	fs.StringVar(&c.LogLevel, "log-level", c.LogLevel, "level")
	fs.StringVar(&c.LogEncoder, "log-encoder", c.LogEncoder, "encoder")
	fs.BoolVar(&c.EnableConsoleLog, "enable-console-log", c.EnableConsoleLog, "")
	fs.BoolVar(&c.EnableFileLog, "enable-file-log", c.EnableFileLog, "")
}
