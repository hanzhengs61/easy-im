package middleware

import (
	"easy-im/pkg/errorx"
	"easy-im/pkg/logger"
	"easy-im/pkg/response"
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

// RecoveryMiddleware panic 恢复中间件
// 捕获所有未处理的 panic，返回 500，并打印堆栈，防止服务崩溃
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				logger.WithContext(r.Context()).Error("panic recovered",
					zap.Any("error", err),
					zap.ByteString("stack", stack),
				)
				response.Fail(w, errorx.New(errorx.CodeServerError))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
