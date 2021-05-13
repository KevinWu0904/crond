package logs

import (
	"github.com/spf13/pflag"
)

type LoggerOptions struct {
	LogLevel         string
	LogMode          string
	LogEncoder       string
	EnableConsoleLog bool
	EnableFileLog    bool
	FileLogDir       string
	FileLogName      string
	FileLogNum       int
	FileLogSize      int
	FileLogAge       int
}

func NewLoggerOptions() *LoggerOptions {
	return &LoggerOptions{
		LogLevel:         "info",
		LogMode:          "debug",
		LogEncoder:       "console",
		EnableConsoleLog: true,
		EnableFileLog:    false,
		FileLogDir:       ".",
		FileLogName:      "server.log",
		FileLogNum:       3,
		FileLogSize:      500,
		FileLogAge:       15,
	}
}

func (o *LoggerOptions) BindFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.LogLevel, "log_level", "info", "CronD service log level")
	fs.StringVar(&o.LogMode, "log_mode", "debug", "Log mode decides debug or production log")
	fs.StringVar(&o.LogEncoder, "log_encoder", "console", "Log encoder decides json encoder or console encoder")
	fs.BoolVar(&o.EnableConsoleLog, "enable_console_log", true, "Determine whether to write console logs")
	fs.BoolVar(&o.EnableFileLog, "enable_file_log", false, "Determine whether to write file logs")
	fs.StringVar(&o.FileLogDir, "file_log_dir", ".", "If non-empty, write log files in this directory")
	fs.StringVar(&o.FileLogName, "file_log_name", "server.log", "If non-empty, use this log file name")
	fs.IntVar(&o.FileLogNum, "file_log_num", 3, "Max file number")
	fs.IntVar(&o.FileLogSize, "file_log_size", 500, "File log max size, unit is MB")
	fs.IntVar(&o.FileLogAge, "file_log_age", 15, "File log max age, unit is day")
}
