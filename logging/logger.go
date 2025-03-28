package logging

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger holds the global zap logger instance.
var (
	logger *zap.Logger
)

// Init initializes the global logger instance based on the environment and configuration.
// It sets the log level, chooses between console (development) and JSON (production) encoding,
// disables timestamps (optional, change if needed), and sets up global convenience functions (Info, Debugf, etc.).
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

	// Optionally disable timestamp - remove or adjust if timestamps are desired.
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

// ctxKey defines the type for the context key used to store the logger.
type ctxKey struct{}

// NewContext embeds the global logger instance into a new context.
// This allows passing the logger implicitly through function calls.
func NewContext(parent context.Context) context.Context {
	return context.WithValue(parent, ctxKey{}, logger)
}

// fromContext retrieves the logger instance from the context.
// If the logger is not found in the context, it returns the global logger instance.
// It uses zap.AddCallerSkip(1) to ensure the caller information points to the
// user's code (e.g., InfoCtx) rather than this function itself.
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
