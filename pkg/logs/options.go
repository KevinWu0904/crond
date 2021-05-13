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
	if o == nil {
		return
	}

	fs.StringVar(&o.LogLevel, "log-level", o.LogLevel, "CronD service log level")
	fs.StringVar(&o.LogMode, "log-mode", o.LogMode, "Log mode decides debug or production log")
	fs.StringVar(&o.LogEncoder, "log-encoder", o.LogEncoder, "Log encoder decides json encoder or console encoder")
	fs.BoolVar(&o.EnableConsoleLog, "enable-console-log", o.EnableConsoleLog, "Determine whether to write console logs")
	fs.BoolVar(&o.EnableFileLog, "enable-file-log", o.EnableFileLog, "Determine whether to write file logs")
	fs.StringVar(&o.FileLogDir, "file-log-dir", o.FileLogDir, "If non-empty, write log files in this directory")
	fs.StringVar(&o.FileLogName, "file-log-name", o.FileLogName, "If non-empty, use this log file name")
	fs.IntVar(&o.FileLogNum, "file-log-num", o.FileLogNum, "Max file number")
	fs.IntVar(&o.FileLogSize, "file-log-size", o.FileLogSize, "File log max size, unit is MB")
	fs.IntVar(&o.FileLogAge, "file-log-age", o.FileLogAge, "File log max age, unit is day")
}
