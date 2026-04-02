package middleware

import (
	"context"
	"easy-im/pkg/errorx"
	"easy-im/pkg/jwt"
	"easy-im/pkg/response"
	"net/http"
	"strings"
)

// ctxUserKey 存入 context 的用户信息 key
type ctxUserKey struct{}

// UserInfo 从 token 中解析出的用户信息，注入 context 供 handler 使用
type UserInfo struct {
	UserID   int64
	Username string
}

// AuthMiddleware JWT 鉴权中间件
// 从 Authorization: Bearer <token> 中解析用户信息并注入 context
func AuthMiddleware(jwtMgr *jwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractToken(r)
			if tokenStr == "" {
				response.Fail(w, errorx.New(errorx.CodeTokenMissing))
				return
			}

			claims, err := jwtMgr.ParseToken(tokenStr)
			if err != nil {
				response.Fail(w, err)
				return
			}

			// 将用户信息注入 context，handler 层通过 GetUserFromCtx 获取
			ctx := context.WithValue(r.Context(), ctxUserKey{}, &UserInfo{
				UserID:   claims.UserID,
				Username: claims.Username,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromCtx 从 context 中取出用户信息，在 handler/logic 层调用
func GetUserFromCtx(ctx context.Context) (*UserInfo, bool) {
	user, ok := ctx.Value(ctxUserKey{}).(*UserInfo)
	return user, ok
}

// extractToken 从请求头提取 token
// 支持：Authorization: Bearer <token>
func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
