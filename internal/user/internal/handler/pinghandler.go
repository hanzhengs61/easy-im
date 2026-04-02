// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"easy-im/pkg/response"
	"net/http"

	"easy-im/internal/user/internal/logic"
	"easy-im/internal/user/internal/svc"
	"easy-im/internal/user/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PingRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, err)
			return
		}

		l := logic.NewPingLogic(r.Context(), svcCtx)
		resp, err := l.Ping(&req)
		if err != nil {
			response.Fail(w, err)
		} else {
			response.OK(w, resp)
		}
	}
}
