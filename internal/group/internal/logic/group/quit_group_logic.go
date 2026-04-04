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

type QuitGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuitGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuitGroupLogic {
	return &QuitGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// QuitGroup 退出群组
func (l *QuitGroupLogic) QuitGroup(req *types.QuitGroupReq) error {
	log := logger.WithContext(l.ctx)

	userInfo, ok := middleware.GetUserFromCtx(l.ctx)
	if !ok {
		return errorx.New(errorx.CodeUnauthorized)
	}

	// 1. 查成员记录
	member, err := l.svcCtx.MemberModel.FindOneByGroupIdUid(l.ctx, req.GroupID, userInfo.UserID)
	if errors.Is(err, model.ErrNotFound) {
		return errorx.New(errorx.CodeNotGroupMember)
	}
	if err != nil {
		return errorx.Wrap(errorx.CodeServerError, err)
	}

	// 2. 群主不能直接退出（必须先转让群主）
	if member.Role == 1 {
		return errorx.NewWithMsg(errorx.CodeForbidden, "群主不能退出群组，请先转让群主")
	}

	// 3. 删除成员记录
	if err = l.svcCtx.MemberModel.Delete(l.ctx, member.Id); err != nil {
		log.Error("quit group delete member failed", zap.Error(err))
		return errorx.Wrap(errorx.CodeServerError, err)
	}

	// 4. 更新群成员数
	group, err := l.svcCtx.GroupModel.FindOne(l.ctx, req.GroupID)
	if err == nil {
		group.MemberCount--
		group.UpdatedAt = time.Now().UnixMilli()
		_ = l.svcCtx.GroupModel.Update(l.ctx, group)
	}

	log.Info("user quit group",
		zap.Int64("uid", userInfo.UserID),
		zap.Int64("group_id", req.GroupID),
	)
	return nil
}
