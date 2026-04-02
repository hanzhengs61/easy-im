// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"easy-im/internal/user/model"
	"easy-im/pkg/errorx"
	"errors"

	"easy-im/internal/user/internal/svc"
	"easy-im/internal/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	// 1. 检查用户名是否已存在
	_, err = l.svcCtx.UserModel.FindOneByUsername(l.ctx, req.Username)
	if err == nil {
		return nil, errorx.New(errorx.CodeUserAlreadyExists)
	}
	if !errors.Is(err, model.ErrNotFound) {
		return nil, err
	}

	// 2. bcrypt 加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// 3. 插入用户
	user := &model.User{
		Username: req.Username,
		Password: string(hashed),
		Nickname: req.Nickname,
		Avatar:   "", // 后面可加默认头像
		Gender:   0,
		Status:   1,
	}
	insertResult, err := l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		return nil, err
	}

	userID, _ := insertResult.LastInsertId()

	// 4. 生成 Token
	token, err := l.svcCtx.JwtManager.GenerateAccessToken(userID, req.Username)
	if err != nil {
		return nil, err
	}

	return &types.RegisterResponse{
		UserId: userID,
		Token:  token,
	}, nil
}
