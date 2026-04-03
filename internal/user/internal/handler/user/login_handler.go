// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"easy-im/pkg/response"
	"net/http"

	"easy-im/internal/user/internal/logic/user"
	"easy-im/internal/user/internal/svc"
	"easy-im/internal/user/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// LoginHandler 用户登录
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, err)
			return
		}

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			response.Fail(w, err)
		} else {
			response.OK(w, resp)
		}
	}
}
