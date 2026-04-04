// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"easy-im/internal/user/model"
	"easy-im/pkg/errorx"
	"errors"

	"easy-im/internal/user/internal/svc"
	"easy-im/internal/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 用户登录
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 1. 查用户
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if errors.Is(err, model.ErrNotFound) {
		return nil, errorx.New(errorx.CodeUserNotFound)
	}
	if err != nil {
		l.Logger.Error("find user failed", zap.Error(err))
		return nil, errorx.Wrap(errorx.CodeServerError, err)
	}

	// 2. 检查账户状态
	if user.Status != 1 {
		return nil, errorx.New(errorx.CodeUserDisabled)
	}

	// 3. 校验密码
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errorx.New(errorx.CodePasswordWrong)
	}

	// 4. 签发 token
	accessToken, err := l.svcCtx.JwtManager.GenerateAccessToken(user.Id, user.Username)
	if err != nil {
		l.Logger.Error("generate access token failed", zap.Error(err))
		return nil, errorx.Wrap(errorx.CodeServerError, err)
	}

	l.Logger.Info("user logged in", zap.Int64("user_id", user.Id))

	return &types.LoginResp{
		AccessToken: accessToken,
		ExpiresIn:   l.svcCtx.Config.Jwt.AccessTokenTTL,
		UserID:      user.Id,
		Nickname:    user.Nickname,
	}, nil
}
