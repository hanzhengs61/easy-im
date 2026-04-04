// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

import (
	"context"
	"easy-im/internal/group/model"
	"easy-im/pkg/errorx"
	"easy-im/pkg/logger"
	"easy-im/pkg/middleware"
	"time"

	"easy-im/internal/group/internal/svc"
	"easy-im/internal/group/internal/types"

	"go.uber.org/zap"
)

type CreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateGroup 创建群组
func (l *CreateGroupLogic) CreateGroup(req *types.CreateGroupReq) (resp *types.CreateGroupResp, err error) {
	log := logger.WithContext(l.ctx)

	userInfo, ok := middleware.GetUserFromCtx(l.ctx)
	if !ok {
		return nil, errorx.New(errorx.CodeUnauthorized)
	}

	// 成员数量限制（群主 + 初始成员）
	totalCount := 1 + len(req.MemberUIDs)
	if totalCount > 500 {
		return nil, errorx.NewWithMsg(errorx.CodeGroupFull, "初始成员不能超过499人")
	}

	now := time.Now().UnixMilli()

	// 1. 创建群组
	result, err := l.svcCtx.GroupModel.Insert(l.ctx, &model.Groups{
		Name:        req.Name,
		Avatar:      req.Avatar,
		Description: req.Description,
		OwnerUid:    userInfo.UserID,
		MemberCount: int64(totalCount),
		MaxMember:   500,
		Status:      1,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		log.Error("create group failed", zap.Error(err))
		return nil, errorx.Wrap(errorx.CodeServerError, err)
	}

	groupID, _ := result.LastInsertId()

	// 2. 插入群主为成员（role=1）
	members := []*model.GroupMembers{
		{
			GroupId:   groupID,
			Uid:       userInfo.UserID,
			Role:      1,
			JoinedAt:  now,
			CreatedAt: now,
		},
	}

	// 3. 插入初始成员（role=3）
	for _, uid := range req.MemberUIDs {
		if uid == userInfo.UserID {
			continue // 跳过群主自己
		}
		members = append(members, &model.GroupMembers{
			GroupId:   groupID,
			Uid:       uid,
			Role:      3,
			JoinedAt:  now,
			CreatedAt: now,
		})
	}

	// 批量插入成员
	if err = l.batchInsertMembers(members); err != nil {
		log.Error("insert members failed", zap.Error(err))
		return nil, errorx.Wrap(errorx.CodeServerError, err)
	}

	log.Info("group created",
		zap.Int64("group_id", groupID),
		zap.Int64("owner_uid", userInfo.UserID),
		zap.Int("member_count", totalCount),
	)

	return &types.CreateGroupResp{
		GroupID:     groupID,
		Name:        req.Name,
		MemberCount: totalCount,
	}, nil
}

func (l *CreateGroupLogic) batchInsertMembers(members []*model.GroupMembers) error {
	for _, m := range members {
		if _, err := l.svcCtx.MemberModel.Insert(l.ctx, m); err != nil {
			return err
		}
	}
	return nil
}
