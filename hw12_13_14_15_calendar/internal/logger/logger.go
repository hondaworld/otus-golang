package logger

import (
	"fmt"
	"os"
	"strings"
)

type Logger struct {
	level int
}

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func New(level string) *Logger {
	return &Logger{
		level: parseLogLevel(level),
	}
}

func parseLogLevel(level string) int {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		fmt.Fprintln(os.Stderr, "Unknown log level:", level, ", defaulting to InfoLevel")
		return InfoLevel
	}
}

func (l Logger) Info(msg string) {
	if l.level <= InfoLevel {
		fmt.Println("[INFO] ", msg)
	}
}

func (l Logger) Error(msg string) {
	if l.level <= ErrorLevel {
		fmt.Println("[ERROR]", msg)
	}
}

func (l Logger) Debug(msg string) {
	if l.level <= DebugLevel {
		fmt.Println("[DEBUG]", msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.level <= WarnLevel {
		fmt.Println("[WARN] ", msg)
	}
}
