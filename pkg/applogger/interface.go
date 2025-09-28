package applogger

import "context"

type Logger interface {
	IsDebugLevel() bool
	IsInfoLevel() bool

	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)

	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})

	WithCtx(ctx context.Context) Logger
}
