// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

import (
	"net/http"

	"easy-im/internal/group/internal/logic/group"
	"easy-im/internal/group/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取我的群组列表
func GetMyGroupsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := group.NewGetMyGroupsLogic(r.Context(), svcCtx)
		resp, err := l.GetMyGroups()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
