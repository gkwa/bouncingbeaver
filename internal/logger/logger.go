package logger

import (
	"fmt"
	"os"
)

type Logger struct {
	level int
}

func New(verbosity int) *Logger {
	return &Logger{
		level: verbosity,
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(0, "ERROR", msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level >= 1 {
		l.log(1, "INFO", msg, args...)
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level >= 2 {
		l.log(2, "DEBUG", msg, args...)
	}
}

func (l *Logger) log(level int, levelName, msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[%s] %s", levelName, msg)
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fmt.Fprintf(os.Stderr, " %v=%v", args[i], args[i+1])
		}
	}
	fmt.Fprintln(os.Stderr)
}
