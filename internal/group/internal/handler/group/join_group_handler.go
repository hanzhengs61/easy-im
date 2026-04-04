// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

import (
	"net/http"

	"easy-im/internal/group/internal/logic/group"
	"easy-im/internal/group/internal/svc"
	"easy-im/internal/group/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// JoinGroupHandler 加入群组
func JoinGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.JoinGroupReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewJoinGroupLogic(r.Context(), svcCtx)
		err := l.JoinGroup(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
