package logs

import "github.com/spf13/pflag"

// Config stores logger configurations.
type Config struct {
	LogLevel         string `mapstructure:"log-level"`
	LogEncoder       string `mapstructure:"log-encoder"`
	EnableConsoleLog bool   `mapstructure:"enable-console-log"`
	EnableFileLog    bool   `mapstructure:"enable-file-log"`
	FileLogDir       string `mapstructure:"file-log-dir"`
	FileLogName      string `mapstructure:"file-log-name"`
	FileLogNum       int    `mapstructure:"file-log-num"`
	FileLogSize      int    `mapstructure:"file-log-size"`
	FileLogAge       int    `mapstructure:"file-log-age"`
}

// DefaultConfig creates the Config with sensible default settings.
func DefaultConfig() *Config {
	return &Config{
		LogLevel:         "info",
		LogEncoder:       "console",
		EnableConsoleLog: true,
		EnableFileLog:    false,
		FileLogDir:       "log",
		FileLogName:      "server",
		FileLogNum:       3,
		FileLogSize:      500,
		FileLogAge:       15,
	}
}

// BindFlags overwrites default logger configurations from CLI flags.
func BindFlags(c *Config, fs *pflag.FlagSet) {
	fs.StringVar(&c.LogLevel, "log-level", c.LogLevel, "server log level, it can be one of "+
		"(debug|info|warn|error|panic|fatal)")
	fs.StringVar(&c.LogEncoder, "log-encoder", c.LogEncoder, "server log encoder, it can be one of "+
		"(console|json)")
	fs.BoolVar(&c.EnableConsoleLog, "enable-console-log", c.EnableConsoleLog, "if true, server will log "+
		"append standard console additionally")
	fs.BoolVar(&c.EnableFileLog, "enable-file-log", c.EnableFileLog, "if true, server will log append "+
		"specific file additionally")
	fs.StringVar(&c.FileLogDir, "file-log-dir", c.FileLogDir, "when enable-file-log is true, this param "+
		"indicates log file path")
	fs.StringVar(&c.FileLogName, "file-log-name", c.FileLogName, "when enable-file-log is true, this param "+
		"indicates log file name")
	fs.IntVar(&c.FileLogNum, "file-log-num", c.FileLogNum, "when enable-file-log is true, we can reserve "+
		"at most file-log-num for log rotation files")
	fs.IntVar(&c.FileLogSize, "file-log-size", c.FileLogSize, "when enable-file-log is true, we should "+
		"rotate log once actual log size larger than file-log-size (unit is MB)")
	fs.IntVar(&c.FileLogAge, "file-log-age", c.FileLogAge, "when enable-file-log is true, we can reserve "+
		"at most file-log-age days for log rotation files (unit is Day)")
}
