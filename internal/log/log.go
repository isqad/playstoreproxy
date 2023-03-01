package log

import (
	"fmt"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	// Logger default app logger
	Logger zerolog.Logger
	mu     sync.RWMutex
	// Ident is distinctive feature for log message: "from that app this log message"
	Ident string
)

const (
	timeFormat = "2006-01-02 15:04:05.000"
)

// Config defines parameters for the logger
type Config struct {
	Level string `mapstructure:"level"`
}

// Init initializes the package logger.
// Supported levels are: ["debug", "info", "warn", "error"]
func Init(level string, ident string) {
	l := zerolog.GlobalLevel()
	switch level {
	case "trace":
		l = zerolog.TraceLevel
	case "debug":
		l = zerolog.DebugLevel
	case "info":
		l = zerolog.InfoLevel
	case "warn":
		l = zerolog.WarnLevel
	case "error":
		l = zerolog.ErrorLevel
	}
	zerolog.TimeFieldFormat = timeFormat
	output := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false, TimeFormat: timeFormat}

	mu.Lock()
	Logger = zerolog.New(output).Level(l).With().Timestamp().Logger()
	Ident = ident
	mu.Unlock()
}

// Info logs a formatted info level log to the console
func Info(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Info().Str("ident", Ident).Str("msg", msg).Msg("")
}

// Trace logs a formatted trace level log to the console
func Trace(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Trace().Str("ident", Ident).Str("msg", msg).Msg("")
}

// Debug logs a formatted debug level log to the console
func Debug(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Debug().Str("ident", Ident).Str("msg", msg).Msg("")
}

// Error logs a formatted error level log to the console
func Error(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Error().Str("ident", Ident).Str("msg", msg).Msg("")
}

// Panic logs a formatted panic level log to the console.
// The panic() function is called, which stops the ordinary flow of a goroutine.
func Panic(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Panic().Str("ident", Ident).Str("msg", msg).Msg("")
}

// Warn logs a formatted warning level log to the console
func Warn(msg string) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Warn().Str("ident", Ident).Str("msg", msg).Msg("")
}

// Infof logs a formatted info level log to the console
func Infof(format string, v ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Info().Str("ident", Ident).Str("msg", fmt.Sprintf(format, v...)).Msg("")
}

// Tracef logs a formatted debug level log to the console
func Tracef(format string, v ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Trace().Str("ident", Ident).Str("msg", fmt.Sprintf(format, v...)).Msg("")
}

// Debugf logs a formatted debug level log to the console
func Debugf(format string, v ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Debug().Str("ident", Ident).Str("msg", fmt.Sprintf(format, v...)).Msg("")
}

// Warnf logs a formatted warn level log to the console
func Warnf(format string, v ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Warn().Str("ident", Ident).Str("msg", fmt.Sprintf(format, v...)).Msg("")
}

// Errorf logs a formatted error level log to the console
func Errorf(format string, v ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Error().Str("ident", Ident).Str("msg", fmt.Sprintf(format, v...)).Msg("")
}

// Panicf logs a formatted panic level log to the console.
// The panic() function is called, which stops the ordinary flow of a goroutine.
func Panicf(format string, v ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()
	Logger.Panic().Str("ident", Ident).Str("msg", fmt.Sprintf(format, v...)).Msg("")
}
