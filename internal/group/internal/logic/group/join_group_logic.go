// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package group

import (
	"context"
	"easy-im/internal/group/model"
	"easy-im/pkg/errorx"
	"easy-im/pkg/logger"
	"easy-im/pkg/middleware"
	"errors"
	"time"

	"easy-im/internal/group/internal/svc"
	"easy-im/internal/group/internal/types"

	"go.uber.org/zap"
)

type JoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// JoinGroup 加入群组
func (l *JoinGroupLogic) JoinGroup(req *types.JoinGroupReq) error {
	log := logger.WithContext(l.ctx)

	userInfo, ok := middleware.GetUserFromCtx(l.ctx)
	if !ok {
		return errorx.New(errorx.CodeUnauthorized)
	}

	// 1. 群组是否存在
	group, err := l.svcCtx.GroupModel.FindOne(l.ctx, req.GroupID)
	if errors.Is(err, model.ErrNotFound) {
		return errorx.New(errorx.CodeGroupNotFound)
	}
	if err != nil {
		return errorx.Wrap(errorx.CodeServerError, err)
	}
	if group.Status != 1 {
		return errorx.NewWithMsg(errorx.CodeGroupNotFound, "群组已解散")
	}

	// 2. 是否已是成员
	_, err = l.svcCtx.MemberModel.FindOneByGroupIdUid(l.ctx, req.GroupID, userInfo.UserID)
	if err == nil {
		return errorx.New(errorx.CodeAlreadyInGroup)
	}
	if !errors.Is(err, model.ErrNotFound) {
		return errorx.Wrap(errorx.CodeServerError, err)
	}

	// 3. 人数是否已满
	if group.MemberCount >= group.MaxMember {
		return errorx.New(errorx.CodeGroupFull)
	}

	now := time.Now().UnixMilli()

	// 4. 加入群组（事务：插入成员 + 更新人数）
	if _, err = l.svcCtx.MemberModel.Insert(l.ctx, &model.GroupMembers{
		GroupId:   req.GroupID,
		Uid:       userInfo.UserID,
		Role:      3,
		JoinedAt:  now,
		CreatedAt: now,
	}); err != nil {
		log.Error("join group insert member failed", zap.Error(err))
		return errorx.Wrap(errorx.CodeServerError, err)
	}

	// 5. 更新群成员数
	group.MemberCount++
	group.UpdatedAt = now
	if err = l.svcCtx.GroupModel.Update(l.ctx, group); err != nil {
		log.Warn("update group member_count failed", zap.Error(err))
	}

	log.Info("user joined group",
		zap.Int64("uid", userInfo.UserID),
		zap.Int64("group_id", req.GroupID),
	)
	return nil
}
