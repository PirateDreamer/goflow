package middleware

import (
	"context"
	"goflow-api/internal/infrastructure/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const TraceHeader = "X-Trace-ID"

func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取或生成 Trace ID
		traceID := c.GetHeader(TraceHeader)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 2. 存入 Context
		ctx := context.WithValue(c.Request.Context(), logger.TraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		// 3. 存入响应 Header
		c.Header(TraceHeader, traceID)

		c.Next()
	}
}
