package logs

import (
	"context"
	"errors"
	"io"
	"os"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var atomicLevel zap.AtomicLevel
var logger *zap.Logger
var sugaredLogger *zap.SugaredLogger

var ginWriteSyncer = io.Discard
var ginErrorWriteSyncer = io.Discard
var raftWriteSyncer = io.Discard

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
	atomicLevel = zap.NewAtomicLevelAt(level)

	var encoder zapcore.Encoder
	switch c.LogEncoder {
	case "json":
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	case "console":
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	default:
		return ErrInvalidLogEncoder
	}

	namedWS := func(c *LoggerConfig, name string) zapcore.WriteSyncer {
		var fullName string
		if name == "" {
			fullName = c.FileLogName + ".log"
		} else {
			fullName = c.FileLogName + "." + name + ".log"
		}

		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   path.Join(c.FileLogDir, fullName),
			MaxSize:    c.FileLogSize,
			MaxAge:     c.FileLogAge,
			MaxBackups: c.FileLogNum,
			LocalTime:  true,
			Compress:   true,
		})
	}
	crondWS := namedWS(c, "")
	crondErrWS := namedWS(c, "err")
	ginWS := namedWS(c, "gin")
	ginErrWS := namedWS(c, "gin.err")
	raftErrWS := namedWS(c, "raft.err")
	stdoutWS := zapcore.Lock(os.Stdout)
	stderrWS := zapcore.Lock(os.Stderr)

	cores := make([]zapcore.Core, 0)
	ginWSs := make([]zapcore.WriteSyncer, 0)
	ginErrWSs := make([]zapcore.WriteSyncer, 0)
	raftErrWSs := make([]zapcore.WriteSyncer, 0)
	if c.EnableConsoleLog {
		normal := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= atomicLevel.Level() && l < zapcore.ErrorLevel
		})
		critical := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= atomicLevel.Level() && l >= zapcore.ErrorLevel
		})

		cores = append(cores,
			zapcore.NewCore(encoder, stdoutWS, normal),
			zapcore.NewCore(encoder, stderrWS, critical),
		)

		ginWSs = append(ginWSs, stdoutWS)
		ginErrWSs = append(ginErrWSs, stderrWS)
		raftErrWSs = append(raftErrWSs, stderrWS)
	}
	if c.EnableFileLog {
		normal := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= atomicLevel.Level() && l >= zapcore.DebugLevel
		})
		critical := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
			return l >= atomicLevel.Level() && l >= zapcore.ErrorLevel
		})

		cores = append(cores,
			zapcore.NewCore(encoder, crondWS, normal),
			zapcore.NewCore(encoder, crondErrWS, critical),
		)

		ginWSs = append(ginWSs, ginWS)
		ginErrWSs = append(ginErrWSs, ginErrWS)
		raftErrWSs = append(raftErrWSs, raftErrWS)
	}

	logger = zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel))
	sugaredLogger = logger.Sugar()

	if len(ginWSs) > 0 {
		ginWriteSyncer = zapcore.NewMultiWriteSyncer(ginWSs...)
	}
	if len(ginErrWSs) > 0 {
		ginErrorWriteSyncer = zapcore.NewMultiWriteSyncer(ginErrWSs...)
	}
	if len(raftErrWSs) > 0 {
		raftWriteSyncer = zapcore.NewMultiWriteSyncer(raftErrWSs...)
	}

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

// GetGinWriter returns io.Writer for gin normal log.
func GetGinWriter() io.Writer {
	return ginWriteSyncer
}

// GetGinErrorWriter returns io.Writer for gin critical log.
func GetGinErrorWriter() io.Writer {
	return ginErrorWriteSyncer
}

// GetRaftErrorWriter returns io.Writer for raft log.
func GetRaftErrorWriter() io.Writer {
	return raftWriteSyncer
}

// Level returns zapcore.Level.
func Level() zapcore.Level {
	return atomicLevel.Level()
}

// SetLevel changes zapcore.Level.
func SetLevel(level zapcore.Level) {
	atomicLevel.SetLevel(level)
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
