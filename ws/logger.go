// Package ws provides logging for WebSocket events.
//
// This file implements the Logger interface for WebSocket event logging.
package ws

import (
	"log"
	"time"
)

type LogLevel int

const (
	LogLevelTrace LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// DefaultLogger implements Logger interface using standard log package.
type DefaultLogger struct {
	level LogLevel
}

func NewDefaultLogger(level LogLevel) *DefaultLogger {
	log.SetFlags(0)
	return &DefaultLogger{level: level}
}

func (l *DefaultLogger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *DefaultLogger) GetLevel() LogLevel {
	return l.level
}

func (l *DefaultLogger) Infof(format string, args ...interface{}) {
	if l.level <= LogLevelInfo {
		log.Printf("%s [INFO] "+format, append([]interface{}{time.Now().Format("2006/01/02 15:04:05.000000000")}, args...)...)
	}
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
	if l.level <= LogLevelError {
		log.Printf("%s [ERROR] "+format, append([]interface{}{time.Now().Format("2006/01/02 15:04:05.000000000")}, args...)...)
	}
}

func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		log.Printf("%s [DEBUG] "+format, append([]interface{}{time.Now().Format("2006/01/02 15:04:05.000000000")}, args...)...)
	}
}

func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
	if l.level <= LogLevelWarn {
		log.Printf("%s [WARN] "+format, append([]interface{}{time.Now().Format("2006/01/02 15:04:05.000000000")}, args...)...)
	}
}

func (l *DefaultLogger) Tracef(format string, args ...interface{}) {
	if l.level <= LogLevelTrace {
		log.Printf("%s [TRACE] "+format, append([]interface{}{time.Now().Format("2006/01/02 15:04:05.000000000")}, args...)...)
	}
}
