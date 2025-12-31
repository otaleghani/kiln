package log

import (
	"github.com/charmbracelet/log"
	"os"
)

// SetDefault replaces the global logger instance.
func SetDefault(l *Logger) {
	Default = l
}

// New creates a new Logger with your specific configuration.
func New(level Level) *Logger {
	l := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		Level:           level,
		// You can add ReportCaller: true here if you want line numbers
	})

	return &Logger{charm: l}
}

func Debug(msg string, keyvals ...any) {
	Default.Debug(msg, keyvals...)
}

func Info(msg string, keyvals ...any) {
	Default.Info(msg, keyvals...)
}

func Warn(msg string, keyvals ...any) {
	Default.Warn(msg, keyvals...)
}

func Error(msg string, keyvals ...any) {
	Default.Error(msg, keyvals...)
}

func Fatal(msg string, keyvals ...any) {
	Default.Fatal(msg, keyvals...)
}

func SetPrefix(prefix string) {
	Default.SetPrefix(prefix)
}

func SetLevel(level Level) {
	Default.SetLevel(level)
}

// With returns a new logger instance with the added fields (chainable)
func With(keyvals ...any) *Logger {
	return Default.With(keyvals...)
}

func (l *Logger) Debug(msg string, keyvals ...any) {
	l.charm.Debug(msg, keyvals...)
}

func (l *Logger) Info(msg string, keyvals ...any) {
	l.charm.Info(msg, keyvals...)
}

func (l *Logger) Warn(msg string, keyvals ...any) {
	l.charm.Warn(msg, keyvals...)
}

func (l *Logger) Error(msg string, keyvals ...any) {
	l.charm.Error(msg, keyvals...)
}

func (l *Logger) Fatal(msg string, keyvals ...any) {
	l.charm.Fatal(msg, keyvals...)
}

func (l *Logger) SetPrefix(prefix string) {
	l.charm.SetPrefix(prefix)
}

func (l *Logger) SetLevel(level Level) {
	l.charm.SetLevel(level)
}

// With returns a new Logger With the fields added to the context.
func (l *Logger) With(keyvals ...any) *Logger {
	newCharm := l.charm.With(keyvals...)
	return &Logger{charm: newCharm}
}

// File adds a file field to the logger
func (l *Logger) WithFile(path string) *Logger {
	return l.With(FieldFile, path)
}

// Error adds an error field to the logger
func (l *Logger) WithError(err error) *Logger {
	return l.With(FieldError, err)
}

func (l *Logger) WithReason(reason string) *Logger {
	return l.With(FieldReason, reason)
}

// Level is an alias of log.Level
type Level = log.Level

const (
	LevelDebug = log.DebugLevel
	LevelInfo  = log.InfoLevel
	LevelWarn  = log.WarnLevel
	LevelError = log.ErrorLevel
	LevelFatal = log.FatalLevel
)

// Field is used for the With function
type Field string

const (
	FieldFile   Field = "file"
	FieldError  Field = "error"
	FieldReason Field = "reason"
	FieldName   Field = "name"
	FieldType   Field = "type"
	FieldPath   Field = "path"
)

// Logger wraps the charm logger.
// 'charm' is private, so it won't pollute your LSP
type Logger struct {
	charm *log.Logger
}

// Default is the global logger instance used by the package-level functions.
// We initialize it with a sensible default (Info level)
var Default = New(LevelInfo)
