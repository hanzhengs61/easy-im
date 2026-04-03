// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"easy-im/internal/user/internal/logic"
	"easy-im/pkg/response"
	"net/http"

	"easy-im/internal/user/internal/svc"
	"easy-im/internal/user/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// RegisterHandler 用户注册
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, err)
			return
		}

		l := logic.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			response.Fail(w, err)
		} else {
			response.OK(w, resp)
		}
	}
}
