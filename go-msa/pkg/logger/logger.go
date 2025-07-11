package logger

import (
	"log"
	"os"
	"strings"
)

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

type SimpleLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
	level       string
}

func New(service string, level string) Logger {
	prefix := "[" + strings.ToUpper(service) + "] "

	return &SimpleLogger{
		debugLogger: log.New(os.Stdout, prefix+"DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, prefix+"INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, prefix+"WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, prefix+"ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(os.Stdout, prefix+"FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
		level:       level,
	}
}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	l.infoLogger.Printf(msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	l.errorLogger.Printf(msg, args...)
}

func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if l.level == "debug" {
		l.debugLogger.Printf(msg, args...)
	}
}

func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	l.warnLogger.Printf(msg, args...)
}

func (l *SimpleLogger) Fatal(msg string, args ...interface{}) {
	l.fatalLogger.Printf(msg, args...)
	os.Exit(1)
}
