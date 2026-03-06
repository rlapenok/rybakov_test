package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// NewLogger creates a new logger with the specified level and format
// Console format includes colored output for better readability
func NewLogger(level string, format string) *zerolog.Logger {
	var logger zerolog.Logger

	writer := parseFormat(format)

	parsedLevel := parseLevel(level)
	zerolog.SetGlobalLevel(parsedLevel)
	zerolog.TimeFieldFormat = time.DateTime

	logger = zerolog.New(writer).
		Level(parsedLevel).
		With().
		Timestamp().
		Logger()

	return &logger
}

// parseLevel parses the level string to a zerolog.Level
func parseLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// parseFormat parses the format string to a io.Writer
func parseFormat(format string) io.Writer {
	switch format {
	case "json":
		return os.Stdout
	default:
		consoleWriter := zerolog.NewConsoleWriter()
		consoleWriter.NoColor = false
		consoleWriter.TimeFormat = time.DateTime
		return consoleWriter
	}
}
