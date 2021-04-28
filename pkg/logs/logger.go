package logs

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger
var sugaredLog *zap.SugaredLogger

var (
	logLevel string

	logMode string

	enableConsole bool
	enableFile    bool

	logDir     string
	logFile    string
	logFileNum int
	logSize    int
	logAge     int
)

var logLevelMapping = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
	"fatal": zap.FatalLevel,
}

func init() {
	flag.StringVar(&logLevel, "log_level", "info", "CronD service log level")

	flag.StringVar(&logMode, "log_mode", "debug", "Log mode decides the log format, json or console encoder")

	flag.BoolVar(&enableConsole, "enable_console", true, "Determine whether to write console logs")
	flag.BoolVar(&enableFile, "enable_file", false, "Determine whether to write file logs")

	flag.StringVar(&logDir, "log_dir", "", "If non-empty, write log files in this directory")
	flag.StringVar(&logFile, "log_file", "", "If non-empty, use this log file")
	flag.IntVar(&logFileNum, "log_file_num", 3, "Max file number")
	flag.IntVar(&logSize, "log_size", 500, "File log max size, unit is MB")
	flag.IntVar(&logAge, "log_age", 15, "File log max age, unit is day")
}

func Init() {
	level, ok := logLevelMapping[logLevel]
	if !ok {
		panic(fmt.Sprintf("invalid log level: %s", logLevel))
	}

	syncers := make([]zapcore.WriteSyncer, 0)

	if enableConsole {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	if enableFile {
		syncers = append(syncers, zapcore.AddSync(&lumberjack.Logger{
			Filename:   path.Join(logDir, logFile),
			MaxSize:    logSize,
			MaxBackups: logFileNum,
			MaxAge:     logAge,
			Compress:   true,
		}))
	}

	var encoder zapcore.Encoder
	var multiWriteSyncer zapcore.WriteSyncer

	if logMode == "debug" {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	multiWriteSyncer = zapcore.NewMultiWriteSyncer(syncers...)

	log = zap.New(zapcore.NewCore(encoder, multiWriteSyncer, level), zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
	sugaredLog = log.Sugar()
}

func Flush() {
	sugaredLog.Sync()
}

func CtxDebug(ctx context.Context, template string, args ...interface{}) {
	sugaredLog.With(GetAllKVs(ctx)...).Debugf(template, args...)
}

func CtxInfo(ctx context.Context, template string, args ...interface{}) {
	sugaredLog.With(GetAllKVs(ctx)...).Infof(template, args...)
}

func CtxWarn(ctx context.Context, template string, args ...interface{}) {
	sugaredLog.With(GetAllKVs(ctx)...).Warnf(template, args...)
}

func CtxError(ctx context.Context, template string, args ...interface{}) {
	sugaredLog.With(GetAllKVs(ctx)...).Errorf(template, args...)
}

func CtxFatal(ctx context.Context, template string, args ...interface{}) {
	sugaredLog.With(GetAllKVs(ctx)...).Fatalf(template, args...)
}

func Debug(template string, args ...interface{}) {
	sugaredLog.Debugf(template, args...)
}

func Info(template string, args ...interface{}) {
	sugaredLog.Infof(template, args...)
}

func Warn(template string, args ...interface{}) {
	sugaredLog.Warnf(template, args...)
}

func Error(template string, args ...interface{}) {
	sugaredLog.Errorf(template, args...)
}

func Fatal(template string, args ...interface{}) {
	sugaredLog.Fatalf(template, args...)
}

func CtxAddKVs(ctx context.Context, kvs ...interface{}) context.Context {
	return ctxAddKVs(ctx, kvs...)
}
