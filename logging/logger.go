package logging

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

// Init initializes the logger.
func Init(isProd bool) {
	level := zapcore.DebugLevel
	if defaultConfig.logLevel != "" {
		if err := level.Set(defaultConfig.logLevel); err != nil {
			level = zapcore.DebugLevel // default level
		}
	}

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	getEncoder := zapcore.NewConsoleEncoder

	if isProd {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.MessageKey = "message"
		encoderConfig.LevelKey = "severity"
		getEncoder = zapcore.NewJSONEncoder
	}

	// disable timestamp
	encoderConfig.TimeKey = ""
	encoderConfig.EncodeTime = nil

	encoder := getEncoder(encoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
	)

	logger = zap.New(core, zap.AddCaller())

	sugar := logger.Sugar()

	Info = sugar.Info
	Infof = sugar.Infof
	Debug = sugar.Debug
	Debugf = sugar.Debugf
	Error = sugar.Error
	Errorf = sugar.Errorf
	Fatal = sugar.Fatal
	Fatalf = sugar.Fatalf
	Panic = sugar.Panic
	Panicf = sugar.Panicf
	Warn = sugar.Warn
	Warnf = sugar.Warnf
}

// ctxKey is the type of value for the context key.
type ctxKey struct{}

// NewContext returns a new context with the logger instance.
func NewContext(parent context.Context) context.Context {
	return context.WithValue(parent, ctxKey{}, logger)
}

func fromContext(ctx context.Context) *zap.Logger {
	c, ok := ctx.Value(ctxKey{}).(*zap.Logger)
	if !ok {
		return logger
	}

	instance := *c

	return instance.WithOptions(zap.AddCallerSkip(1))
}

type printWrapper func(args ...any)
type printfWrapper func(template string, args ...any)

var (
	Info   printWrapper
	Infof  printfWrapper
	Debug  printWrapper
	Debugf printfWrapper
	Error  printWrapper
	Errorf printfWrapper
	Fatal  printWrapper
	Fatalf printfWrapper
	Panic  printWrapper
	Panicf printfWrapper
	Warn   printWrapper
	Warnf  printfWrapper
)

// DebugCtx logs a message at level Debug on the logger associated with the context.
func DebugCtx(ctx context.Context, message any) {
	fromContext(ctx).Sugar().Debug(message)
}

// DebugfCtx logs a message at level Debug on the logger associated with the context.
func DebugfCtx(ctx context.Context, template string, args ...any) {
	fromContext(ctx).Sugar().Debugf(template, args...)
}

// DebugWithFieldCtx logs a message at level Debug on the logger associated with the context.
func DebugWithFieldCtx(ctx context.Context, message string, fields ...zap.Field) {
	fromContext(ctx).Debug(message, fields...)
}

// DebugWithDataCtx logs a message at level Debug on the logger associated with the context.
func DebugWithDataCtx(ctx context.Context, message any, data any) {
	fromContext(ctx).With(zap.Any("data", data)).Sugar().Debug(message)
}

// InfoCtx logs a message at level Info on the logger associated with the context.
func InfoCtx(ctx context.Context, message any) {
	fromContext(ctx).Sugar().Info(message)
}

// InfofCtx logs a message at level Info on the logger associated with the context.
func InfofCtx(ctx context.Context, template string, args ...any) {
	fromContext(ctx).Sugar().Infof(template, args...)
}

// InfoWithFieldCtx logs a message at level Info on the logger associated with the context.
func InfoWithFieldCtx(ctx context.Context, message string, fields ...zap.Field) {
	fromContext(ctx).Info(message, fields...)
}

// InfoWithDataCtx logs a message at level Info on the logger associated with the context.
func InfoWithDataCtx(ctx context.Context, message any, data any) {
	fromContext(ctx).With(zap.Any("data", data)).Sugar().Info(message)
}

// WarnCtx logs a message at level Warn on the logger associated with the context.
func WarnCtx(ctx context.Context, message any) {
	fromContext(ctx).Sugar().Warn(message)
}

// WarnfCtx logs a message at level Warn on the logger associated with the context.
func WarnfCtx(ctx context.Context, template string, args ...any) {
	fromContext(ctx).Sugar().Warnf(template, args...)
}

// WarnWithFieldCtx logs a message at level Warn on the logger associated with the context.
func WarnWithFieldCtx(ctx context.Context, message string, fields ...zap.Field) {
	fromContext(ctx).Warn(message, fields...)
}

// WarnWithDataCtx logs a message at level Warn on the logger associated with the context.
func WarnWithDataCtx(ctx context.Context, message any, data any) {
	fromContext(ctx).With(zap.Any("data", data)).Sugar().Warn(message)
}

// ErrorCtx logs a message at level Error on the logger associated with the context.
func ErrorCtx(ctx context.Context, message any) {
	fromContext(ctx).Sugar().Error(message)
}

// ErrorfCtx logs a message at level Error on the logger associated with the context.
func ErrorfCtx(ctx context.Context, template string, args ...any) {
	fromContext(ctx).Sugar().Errorf(template, args...)
}

// ErrorWithFieldCtx logs a message at level Error on the logger associated with the context.
func ErrorWithFieldCtx(ctx context.Context, message string, fields ...zap.Field) {
	fromContext(ctx).Error(message, fields...)
}

// ErrorWithDataCtx logs a message at level Error on the logger associated with the context.
func ErrorWithDataCtx(ctx context.Context, message any, data any) {
	fromContext(ctx).With(zap.Any("data", data)).Sugar().Error(message)
}

// FatalCtx logs a message at level Fatal on the logger associated with the context.
func FatalCtx(ctx context.Context, message any) {
	fromContext(ctx).Sugar().Fatal(message)
}

// FatalfCtx logs a message at level Fatal on the logger associated with the context.
func FatalfCtx(ctx context.Context, template string, args ...any) {
	fromContext(ctx).Sugar().Fatalf(template, args...)
}

// FatalWithFieldCtx logs a message at level Fatal on the logger associated with the context.
func FatalWithFieldCtx(ctx context.Context, message string, fields ...zap.Field) {
	fromContext(ctx).Fatal(message, fields...)
}

// FatalWithDataCtx logs a message at level Fatal on the logger associated with the context.
func FatalWithDataCtx(ctx context.Context, message any, data any) {
	fromContext(ctx).With(zap.Any("data", data)).Sugar().Fatal(message)
}

// PanicCtx logs a message at level Panic on the logger associated with the context.
func PanicCtx(ctx context.Context, message any) {
	fromContext(ctx).Sugar().Panic(message)
}

// PanicfCtx logs a message at level Panic on the logger associated with the context.
func PanicfCtx(ctx context.Context, template string, args ...any) {
	fromContext(ctx).Sugar().Panicf(template, args...)
}

// PanicWithFieldCtx logs a message at level Panic on the logger associated with the context.
func PanicWithFieldCtx(ctx context.Context, message string, fields ...zap.Field) {
	fromContext(ctx).Panic(message, fields...)
}

// PanicWithDataCtx logs a message at level Panic on the logger associated with the context.
func PanicWithDataCtx(ctx context.Context, message any, data any) {
	fromContext(ctx).With(zap.Any("data", data)).Sugar().Panic(message)
}
