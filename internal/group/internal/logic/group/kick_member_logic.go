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

	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
)

type KickMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewKickMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickMemberLogic {
	return &KickMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// KickMember 踢出成员（仅群主/管理员）
func (l *KickMemberLogic) KickMember(req *types.KickMemberReq) error {
	log := logger.WithContext(l.ctx)

	userInfo, ok := middleware.GetUserFromCtx(l.ctx)
	if !ok {
		return errorx.New(errorx.CodeUnauthorized)
	}

	// 1. 校验操作者权限（必须是群主或管理员）
	operatorMember, err := l.svcCtx.MemberModel.FindOneByGroupIdUid(l.ctx, req.GroupID, userInfo.UserID)
	if errors.Is(err, model.ErrNotFound) {
		return errorx.New(errorx.CodeNotGroupMember)
	}
	if err != nil {
		return errorx.Wrap(errorx.CodeServerError, err)
	}
	if operatorMember.Role > 2 {
		return errorx.New(errorx.CodeNotGroupOwner)
	}

	// 2. 不能踢自己
	if req.TargetUID == userInfo.UserID {
		return errorx.NewWithMsg(errorx.CodeForbidden, "不能踢出自己")
	}

	// 3. 查被踢成员
	targetMember, err := l.svcCtx.MemberModel.FindOneByGroupIdUid(l.ctx, req.GroupID, req.TargetUID)
	if errors.Is(err, model.ErrNotFound) {
		return errorx.New(errorx.CodeNotGroupMember)
	}
	if err != nil {
		return errorx.Wrap(errorx.CodeServerError, err)
	}

	// 4. 管理员不能踢群主，普通管理员不能踢其他管理员
	if targetMember.Role == 1 {
		return errorx.NewWithMsg(errorx.CodeForbidden, "不能踢出群主")
	}
	if operatorMember.Role == 2 && targetMember.Role == 2 {
		return errorx.NewWithMsg(errorx.CodeForbidden, "管理员不能踢出其他管理员")
	}

	// 5. 执行踢出
	if err = l.svcCtx.MemberModel.Delete(l.ctx, targetMember.Id); err != nil {
		log.Error("kick member failed", zap.Error(err))
		return errorx.Wrap(errorx.CodeServerError, err)
	}

	// 6. 更新群人数
	group, err := l.svcCtx.GroupModel.FindOne(l.ctx, req.GroupID)
	if err == nil {
		group.MemberCount--
		group.UpdatedAt = time.Now().UnixMilli()
		_ = l.svcCtx.GroupModel.Update(l.ctx, group)
	}

	log.Info("member kicked",
		zap.Int64("operator_uid", userInfo.UserID),
		zap.Int64("target_uid", req.TargetUID),
		zap.Int64("group_id", req.GroupID),
	)
	return nil
}
