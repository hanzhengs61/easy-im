// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

import (
	"easy-im/pkg/response"
	"net/http"

	"easy-im/internal/group/internal/logic/group"
	"easy-im/internal/group/internal/svc"
	"easy-im/internal/group/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// CreateGroupHandler 创建群组
func CreateGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateGroupReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Fail(w, err)
			return
		}

		l := group.NewCreateGroupLogic(r.Context(), svcCtx)
		resp, err := l.CreateGroup(&req)
		if err != nil {
			response.Fail(w, err)
		} else {
			response.OK(w, resp)
		}
	}
}
