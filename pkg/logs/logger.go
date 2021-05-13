package logs

import (
	"context"
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger
var sugaredLog *zap.SugaredLogger

var logLevelMapping = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
	"panic": zap.PanicLevel,
	"fatal": zap.FatalLevel,
}

func Init(options *LoggerOptions) {
	level, ok := logLevelMapping[options.LogLevel]
	if !ok {
		panic("invalid log level")
	}

	syncers := make([]zapcore.WriteSyncer, 0)

	if options.EnableConsoleLog {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	if options.EnableFileLog {
		syncers = append(syncers, zapcore.AddSync(&lumberjack.Logger{
			Filename:   path.Join(options.FileLogDir, options.FileLogName),
			MaxSize:    options.FileLogSize,
			MaxBackups: options.FileLogNum,
			MaxAge:     options.FileLogAge,
			Compress:   true,
		}))
	}

	var encoderConfig zapcore.EncoderConfig
	var encoder zapcore.Encoder
	var multiWriteSyncer zapcore.WriteSyncer

	if options.LogMode == "debug" {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	switch options.LogEncoder {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
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

func CtxPanic(ctx context.Context, template string, args ...interface{}) {
	sugaredLog.With(GetAllKVs(ctx)...).Panicf(template, args...)
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

func Panic(template string, args ...interface{}) {
	sugaredLog.Panicf(template, args...)
}

func Fatal(template string, args ...interface{}) {
	sugaredLog.Fatalf(template, args...)
}

func CtxAddKVs(ctx context.Context, kvs ...interface{}) context.Context {
	return ctxAddKVs(ctx, kvs...)
}
