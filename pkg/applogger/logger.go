package applogger

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"wn/pkg/constants"
	"strings"
)

type logger struct {
	log   *zap.Logger
	level string
}

func NewLogger(level string) (Logger, error) {

	lvl, err := zap.ParseAtomicLevel(strings.ToLower(level))
	if err != nil {
		return nil, err
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	cfg := zap.Config{
		Level:             lvl,
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	defer l.Sync()

	return &logger{log: l, level: strings.ToLower(level)}, nil
}

func (l *logger) IsDebugLevel() bool {
	return l.level == "debug"
}

func (l *logger) IsInfoLevel() bool {
	return l.level == "info"
}

func (l *logger) Debug(msg string) {
	l.log.Debug(msg)
}

func (l *logger) Info(msg string) {
	l.log.Info(msg)
}

func (l *logger) Warn(msg string) {
	l.log.Warn(msg)
}

func (l *logger) Error(msg string) {
	l.log.Error(msg)
}

func (l *logger) Debugf(msg string, args ...interface{}) {
	l.log.Debug(fmt.Sprintf(msg, args...))
}

func (l *logger) Infof(msg string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(msg, args...))
}

func (l *logger) Errorf(msg string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(msg, args...))
}

func (l *logger) Warnf(msg string, args ...interface{}) {
	l.log.Warn(fmt.Sprintf(msg, args...))
}

func (l *logger) WithCtx(ctx context.Context) Logger {
	clone := l.log.With(getLogParamFromCtx(ctx)...)
	return &logger{log: clone}
}

func getLogParamFromCtx(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0, 6)
	if ctx != nil {
		if ctxRqId, ok := ctx.Value(constants.RequestIdCtx).(string); ok {
			fields = append(fields, zap.String(constants.RequestIdCtx, ctxRqId))
		}
		if userIdCtx, ok := ctx.Value(constants.UserIdCtx).(string); ok {
			fields = append(fields, zap.String(constants.UserIdCtx, userIdCtx))
		}
		if roleCtx, ok := ctx.Value(constants.UserRoleCtx).(string); ok {
			fields = append(fields, zap.String(constants.UserRoleCtx, roleCtx))
		}
		if apiNameCtx, ok := ctx.Value(constants.ApiNameCtx).(string); ok {
			fields = append(fields, zap.String(constants.ApiNameCtx, apiNameCtx))
		}

		spanCtx := trace.SpanFromContext(ctx)
		if spanCtx.SpanContext().IsValid() {
			fields = append(fields, zap.String(constants.TraceIdCtx, spanCtx.SpanContext().TraceID().String()))
			fields = append(fields, zap.String(constants.SpanIdCtx, spanCtx.SpanContext().SpanID().String()))
		}
	}
	return fields
}
