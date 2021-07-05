package logs

import (
	"context"
	"errors"
	"os"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sugaredLogger *zap.SugaredLogger

var logLevelMapping = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
	"panic": zap.PanicLevel,
	"fatal": zap.FatalLevel,
}

// ErrInvalidLogLevel throws when input invalid log level.
var ErrInvalidLogLevel = errors.New("invalid log level")

// ErrInvalidLogEncoder throws when input invalid log encoder.
var ErrInvalidLogEncoder = errors.New("invalid log encoder")

// InitLogger initializes global zap logger instance.
func InitLogger(c *LoggerConfig) error {
	level, ok := logLevelMapping[c.LogLevel]
	if !ok {
		return ErrInvalidLogLevel
	}

	var encoder zapcore.Encoder

	switch c.LogEncoder {
	case "json":
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	case "console":
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	default:
		return ErrInvalidLogEncoder
	}

	cores := make([]zapcore.Core, 0)

	if c.EnableConsoleLog {
		normal := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= level && l < zapcore.ErrorLevel
		})
		critical := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= level && l >= zapcore.ErrorLevel
		})
		cores = append(cores,
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), normal),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), critical),
		)
	}
	if c.EnableFileLog {
		normal := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= level && l >= zapcore.DebugLevel
		})
		critical := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= level && l >= zapcore.ErrorLevel
		})
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(&lumberjack.Logger{
			Filename:   path.Join(c.FileLogDir, c.FileLogName+".log"),
			MaxSize:    c.FileLogSize,
			MaxAge:     c.FileLogAge,
			MaxBackups: c.FileLogNum,
			LocalTime:  true,
			Compress:   true,
		}), normal), zapcore.NewCore(encoder, zapcore.AddSync(&lumberjack.Logger{
			Filename:   path.Join(c.FileLogDir, c.FileLogName+".err.log"),
			MaxSize:    c.FileLogSize,
			MaxAge:     c.FileLogAge,
			MaxBackups: c.FileLogNum,
			LocalTime:  true,
			Compress:   true,
		}), critical),
		)
	}

	logger = zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel))
	sugaredLogger = logger.Sugar()

	return nil
}

// Flush suggests to be called with defer block, it flushes remaining buffered log records from memory into writer syncers.
func Flush() {
	sugaredLogger.Sync()
}

// GetLogger returns zap.Logger singleton.
func GetLogger() *zap.Logger {
	return logger
}

// GetSugaredLogger returns zap.SugaredLogger singleton.
func GetSugaredLogger() *zap.SugaredLogger {
	return sugaredLogger
}

// CtxDebug logs a debug level record with specific kvs.
func CtxDebug(ctx context.Context, template string, args ...interface{}) {
	sugaredLogger.With(GetAllKVs(ctx)...).Debugf(template, args...)
}

// CtxInfo logs a info level record with specific kvs.
func CtxInfo(ctx context.Context, template string, args ...interface{}) {
	sugaredLogger.With(GetAllKVs(ctx)...).Infof(template, args...)
}

// CtxWarn logs a warn level record with specific kvs.
func CtxWarn(ctx context.Context, template string, args ...interface{}) {
	sugaredLogger.With(GetAllKVs(ctx)...).Warnf(template, args...)
}

// CtxError logs an error level record with specific kvs.
func CtxError(ctx context.Context, template string, args ...interface{}) {
	sugaredLogger.With(GetAllKVs(ctx)...).Errorf(template, args...)
}

// CtxPanic logs a panic level record with specific kvs.
func CtxPanic(ctx context.Context, template string, args ...interface{}) {
	sugaredLogger.With(GetAllKVs(ctx)...).Panicf(template, args...)
}

// CtxFatal logs a fatal level record with specific kvs.
func CtxFatal(ctx context.Context, template string, args ...interface{}) {
	sugaredLogger.With(GetAllKVs(ctx)...).Fatalf(template, args...)
}

// Debug logs a debug level record.
func Debug(template string, args ...interface{}) {
	sugaredLogger.Debugf(template, args...)
}

// Info logs a info level record.
func Info(template string, args ...interface{}) {
	sugaredLogger.Infof(template, args...)
}

// Warn logs a warn level record.
func Warn(template string, args ...interface{}) {
	sugaredLogger.Warnf(template, args...)
}

// Error logs an error level record.
func Error(template string, args ...interface{}) {
	sugaredLogger.Errorf(template, args...)
}

// Panic logs a panic level record.
func Panic(template string, args ...interface{}) {
	sugaredLogger.Panicf(template, args...)
}

// Fatal logs a fatal level record.
func Fatal(template string, args ...interface{}) {
	sugaredLogger.Fatalf(template, args...)
}
