// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"easy-im/internal/user/model"
	"easy-im/pkg/errorx"
	"easy-im/pkg/middleware"
	"errors"

	"easy-im/internal/user/internal/svc"
	"easy-im/internal/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetUserInfo 获取当前用户信息
func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	// 从 context 中获取当前登录用户
	userInfo, ok := middleware.GetUserFromCtx(l.ctx)
	if !ok {
		return nil, errorx.New(errorx.CodeUnauthorized)
	}

	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userInfo.UserID)
	if errors.Is(err, model.ErrNotFound) {
		return nil, errorx.New(errorx.CodeUserNotFound)
	}
	if err != nil {
		l.Logger.Error("find user failed", zap.Error(err), zap.Int64("user_id", userInfo.UserID))
		return nil, errorx.Wrap(errorx.CodeServerError, err)
	}

	return &types.UserInfoResp{
		UserID:    user.Id,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
	}, nil
}
