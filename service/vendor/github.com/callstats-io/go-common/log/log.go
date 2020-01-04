package log

import (
	"context"
	"os"
	"time"

	"github.com/uber-go/zap"
)

type loggerCtxKey int

// Logging keys
const (
	LogKeyPackage = "package"
	LogKeyContext = "context"
)

// Environment variable names
const (
	EnvLogLevel = "LOG_LEVEL"
)

// Log levels
const (
	DebugLevel = "DEBUG"
	InfoLevel  = "INFO"
	WarnLevel  = "WARN"
	ErrorLevel = "ERROR"
	FatalLevel = "FATAL"
	PanicLevel = "PANIC"
)

// alias internally a context key to int. A value is fetched from context by type + value so this ensures our key will never conflict with another object.
type ctxKey int

// Internal context keys. The logger should only be accessed with the WithLogger, FromContext functions.
const (
	ctxKeyLogger ctxKey = iota
)

// DynamicLogger is zap.Logger with dynamic log level support
type DynamicLogger struct {
	zap.Logger
	// DynamicLevel must be created with zap.DynamicLevel()
	DynamicLevel *zap.AtomicLevel
}

// Logger extends zap.Logger with SetLevel
type Logger interface {
	// Check returns a CheckedMessage if logging a message at the specified level
	// is enabled. It's a completely optional optimization; in high-performance
	// applications, Check can help avoid allocating a slice to hold fields.
	//
	// See CheckedMessage for an example.
	Check(zap.Level, string) *zap.CheckedMessage

	// Log a message at the given level. Messages include any context that's
	// accumulated on the logger, as well as any fields added at the log site.
	//
	// Calling Panic should panic() and calling Fatal should terminate the
	// process, but calling Log(PanicLevel, ...) or Log(FatalLevel, ...) should
	// not. It may not be possible for compatibility wrappers to comply with
	// this last part (e.g. the bark wrapper).
	Log(zap.Level, string, ...zap.Field)
	Debug(string, ...zap.Field)
	Info(string, ...zap.Field)
	Warn(string, ...zap.Field)
	Error(string, ...zap.Field)
	DPanic(string, ...zap.Field)
	Panic(string, ...zap.Field)
	Fatal(string, ...zap.Field)
	// SetLevel changes the log level
	SetLevel(logLevel string)
	// Create a child logger, and optionally add some context to that logger.
	With(fields ...zap.Field) Logger
}

// Alias typed field loggers.
// This isn't the most elegant way and we should follow how golang evolves to make this easier.
var (
	// Currently time formatter logs time as Float64. There are some breaking changes coming in zap dev to fix this,
	// but for now we need to have a custom wrapper for time to make it RFC3339 compatible
	//Time      = zap.Time
	Skip      = zap.Skip
	Base64    = zap.Base64
	Bool      = zap.Bool
	Float64   = zap.Float64
	Int       = zap.Int
	Int64     = zap.Int64
	Uint      = zap.Uint
	Uint64    = zap.Uint64
	Uintptr   = zap.Uintptr
	String    = zap.String
	Stringer  = zap.Stringer
	Error     = zap.Error
	Stack     = zap.Stack
	Duration  = zap.Duration
	Marshaler = zap.Marshaler
	Object    = zap.Object
	Nest      = zap.Nest
)

// Time encodes a time.Time in RFC3339Nano format
func Time(key string, val time.Time) zap.Field {
	return zap.String(key, val.Format(time.RFC3339Nano))
}

// internal root logger
var (
	rootLogger Logger
)

func init() {
	rootLogger = FromEnv()
}

// FromEnv returns a new stdout logger using the log level from environment
func FromEnv() Logger {
	return NewLogger(os.Getenv(EnvLogLevel))
}

// SetRootLogger sets the provided logger as the new rootLogger
func SetRootLogger(logger Logger) {
	rootLogger = logger
}

// WithLogger returns a context with the specified logger
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}

// RootLogger returns the current root logger
func RootLogger() Logger {
	return rootLogger
}

// WithContext returns a context with the specified logger
func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}

// FromContext returns a logger from the context. If a logger has not been set, it returns the rootLogger
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(ctxKeyLogger).(Logger); ok {
		return logger
	}
	return rootLogger
}

// NewLogger returns a new logger with the specified level
func NewLogger(logLevel string) Logger {
	return NewStdoutLogger(logLevel)
}

// NewStdoutLogger returns a new zap.Logger that logs output as JSON to stdout.
// It includes hostname, timestamp and level as string by default.
// Expects another component outside of this to pick up the log messages from stdout and forward them to e.g. kibana
func NewStdoutLogger(logLevel string) Logger {
	dyn := zap.DynamicLevel()
	dyn.SetLevel(zapLogLevel(logLevel))
	options := []zap.Option{
		zap.Output(os.Stdout),
		dyn,
	}

	encoder := zap.NewJSONEncoder(
		zap.RFC3339NanoFormatter("timestamp"),
		zap.LevelString("level"),
	)

	hostname, _ := os.Hostname()

	zapLogger := zap.New(encoder, options...).
		With(zap.String("hostname", hostname))

	logger := &DynamicLogger{
		Logger:       zapLogger,
		DynamicLevel: &dyn,
	}

	return logger
}

// FromContextWithPackageName returns a new Logger from context and sets the package name
// Eg. FromContextWithPackageName(ctx, "go-common/log")
func FromContextWithPackageName(ctx context.Context, pkg string) Logger {
	return FromContext(ctx).With(String(LogKeyPackage, pkg))
}

// SetLevel dynamically changes the log level
func (logger *DynamicLogger) SetLevel(logLevel string) {
	logger.DynamicLevel.SetLevel(zapLogLevel(logLevel))
}

// With creates a child logger, and optionally add some context to that logger.
func (logger *DynamicLogger) With(fields ...zap.Field) Logger {
	return &DynamicLogger{
		Logger:       logger.Logger.With(fields...),
		DynamicLevel: logger.DynamicLevel,
	}
}

func zapLogLevel(level string) zap.Level {
	switch level {
	case DebugLevel:
		return zap.DebugLevel
	case InfoLevel:
		return zap.InfoLevel
	case WarnLevel:
		return zap.WarnLevel
	case ErrorLevel:
		return zap.ErrorLevel
	case FatalLevel:
		return zap.FatalLevel
	case PanicLevel:
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}
