package middleware

import (
	"easy-im/pkg/logger"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// LoggerMiddleware 请求日志中间件
// 为每个请求生成唯一 trace_id，并记录：方法、路径、状态码、耗时、客户端 IP
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 生成 trace_id：优先读取上游传入的（微服务链路追踪），否则自行生成
		traceID := r.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.NewString()
		}

		// 将 trace_id 写入响应头，方便前端排查
		w.Header().Set("X-Trace-Id", traceID)

		// 注入 context，后续 handler 通过 logger.WithContext(ctx) 获取带 trace_id 的 logger
		ctx := logger.NewContext(r.Context(), traceID)
		r = r.WithContext(ctx)

		// 包装 ResponseWriter 以捕获状态码
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		// 请求完成后记录日志
		duration := time.Since(start)
		logger.WithContext(ctx).Info("http request",
			zap.String("trace_id", traceID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.Int("status", rw.status),
			zap.Duration("duration", duration),
			zap.String("ip", clientIP(r)),
			zap.String("user_agent", r.UserAgent()),
		)
	})
}

// responseWriter 包装 http.ResponseWriter，捕获状态码
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// clientIP 获取真实客户端 IP（兼容反向代理）
func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
