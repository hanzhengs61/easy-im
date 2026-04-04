// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

import (
	"context"
	"easy-im/internal/group/model"
	"easy-im/pkg/errorx"
	"easy-im/pkg/middleware"
	"errors"

	"easy-im/internal/group/internal/svc"
	"easy-im/internal/group/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMembersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMembersLogic {
	return &GetGroupMembersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetGroupMembers 获取群成员列表
func (l *GetGroupMembersLogic) GetGroupMembers(groupID int64) (*types.GetGroupMembersResp, error) {
	_, ok := middleware.GetUserFromCtx(l.ctx)
	if !ok {
		return nil, errorx.New(errorx.CodeUnauthorized)
	}

	group, err := l.svcCtx.GroupModel.FindOne(l.ctx, groupID)
	if errors.Is(err, model.ErrNotFound) {
		return nil, errorx.New(errorx.CodeGroupNotFound)
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeServerError, err)
	}

	members, err := l.svcCtx.MemberModel.FindByGroupId(l.ctx, groupID)
	if err != nil {
		return nil, errorx.Wrap(errorx.CodeServerError, err)
	}

	memberList := make([]types.GroupMember, 0, len(members))
	for _, m := range members {
		memberList = append(memberList, types.GroupMember{
			UID:      m.Uid,
			Role:     int(m.Role),
			JoinedAt: m.JoinedAt,
		})
	}

	return &types.GetGroupMembersResp{
		GroupID:     group.Id,
		Name:        group.Name,
		MemberCount: int(group.MemberCount),
		Members:     memberList,
	}, nil
}
