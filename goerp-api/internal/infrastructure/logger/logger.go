package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var log zerolog.Logger

const TraceIDKey = "trace_id"

func Init() {
	zerolog.TimeFieldFormat = time.RFC3339
	log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// L 从 Context 中提取 Trace ID 并返回带上下文的日志对象
func L(ctx context.Context) *zerolog.Event {
	traceID, ok := ctx.Value(TraceIDKey).(string)
	if ok {
		return log.Info().Str(TraceIDKey, traceID)
	}
	return log.Info()
}

// ErrorL 记录错误带 Trace ID
func ErrorL(ctx context.Context, err error) *zerolog.Event {
	traceID, ok := ctx.Value(TraceIDKey).(string)
	if ok {
		return log.Error().Err(err).Str(TraceIDKey, traceID)
	}
	return log.Error().Err(err)
}

// Info 直接记录全局 Info 日志
func Info() *zerolog.Event {
	return log.Info()
}

// Error 直接记录全局 Error 日志
func Error(err error) *zerolog.Event {
	return log.Error().Err(err)
}
