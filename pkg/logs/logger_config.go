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
		FileLogDir:       ".",
		FileLogName:      "crond",
		FileLogNum:       3,
		FileLogSize:      500,
		FileLogAge:       15,
	}
}

// BindLoggerFlags overwrites default logger configurations from CLI flags.
func BindLoggerFlags(c *LoggerConfig, fs *pflag.FlagSet) {
	fs.StringVar(&c.LogLevel, "log-level", c.LogLevel, "crond log level, it can be one of "+
		"(debug|info|warn|error|panic|fatal)")
	fs.StringVar(&c.LogEncoder, "log-encoder", c.LogEncoder, "crond log encoder, it can be one of "+
		"(console|json)")
	fs.BoolVar(&c.EnableConsoleLog, "enable-console-log", c.EnableConsoleLog, "if true, crond will log "+
		"append standard console additionally")
	fs.BoolVar(&c.EnableFileLog, "enable-file-log", c.EnableFileLog, "if true, crond will log append "+
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
