// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

import (
	"easy-im/pkg/errorx"
	"easy-im/pkg/response"
	"net/http"
	"strconv"

	"easy-im/internal/group/internal/svc"

	"easy-im/internal/group/internal/logic/group"

	"github.com/gorilla/mux"
)

// GetGroupMembersHandler 获取群成员列表
func GetGroupMembersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupIDStr := mux.Vars(r)["groupId"]
		groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil || groupID <= 0 {
			response.Fail(w, errorx.New(errorx.CodeInvalidParam))
			return
		}

		l := group.NewGetGroupMembersLogic(r.Context(), svcCtx)
		resp, err := l.GetGroupMembers(groupID)
		if err != nil {
			response.Fail(w, err)
			return
		}
		response.OK(w, resp)
	}
}
