// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"easy-im/internal/user/internal/svc"
	"easy-im/internal/user/internal/types"
	"easy-im/internal/user/model"
	"easy-im/pkg/errorx"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *LoginLogic) Login(ctx context.Context, req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	user, err := l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if errors.Is(err, model.ErrNotFound) {
		return nil, errorx.New(errorx.CodeUserNotFound)
	}
	if err != nil {
		return nil, err
	}

	// 密码比对
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errorx.New(errorx.CodePasswordWrong)
	}

	id := int64(user.Id)

	// 生成 Token
	token, err := l.svcCtx.JwtManager.GenerateAccessToken(id, user.Username)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		UserId: id,
		Token:  token,
	}, nil
}
