// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

import (
	"context"
	"easy-im/pkg/errorx"
	"easy-im/pkg/logger"
	"easy-im/pkg/middleware"

	"easy-im/internal/group/internal/svc"
	"easy-im/internal/group/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
)

type GetMyGroupsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyGroupsLogic {
	return &GetMyGroupsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetMyGroups 获取我的群组列表
func (l *GetMyGroupsLogic) GetMyGroups() (resp *types.GetMyGroupsResp, err error) {
	log := logger.WithContext(l.ctx)

	userInfo, ok := middleware.GetUserFromCtx(l.ctx)
	if !ok {
		return nil, errorx.New(errorx.CodeUnauthorized)
	}

	// 查用户所在的所有群
	memberships, err := l.svcCtx.MemberModel.FindByUid(l.ctx, userInfo.UserID)
	if err != nil {
		log.Error("find user groups failed", zap.Error(err))
		return nil, errorx.Wrap(errorx.CodeServerError, err)
	}

	groups := make([]types.GroupInfo, 0, len(memberships))
	for _, m := range memberships {
		g, err := l.svcCtx.GroupModel.FindOne(l.ctx, m.GroupId)
		if err != nil || g.Status != 1 {
			continue
		}
		groups = append(groups, types.GroupInfo{
			GroupID:     g.Id,
			Name:        g.Name,
			Avatar:      g.Avatar,
			Description: g.Description,
			OwnerUID:    g.OwnerUid,
			MemberCount: int(g.MemberCount),
		})
	}

	return &types.GetMyGroupsResp{Groups: groups}, nil
}
