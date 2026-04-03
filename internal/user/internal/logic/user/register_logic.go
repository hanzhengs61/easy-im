// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"easy-im/internal/user/model"
	"easy-im/pkg/errorx"
	"errors"
	"time"

	"easy-im/internal/user/internal/svc"
	"easy-im/internal/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
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

// Register 用户注册
func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
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
	now := time.Now().UnixMilli()
	user, err := l.svcCtx.UserModel.Insert(l.ctx, &model.Users{
		Username:  req.Username,
		Password:  string(hashed),
		Nickname:  req.Nickname,
		Status:    1,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		l.Logger.Error("insert user failed", zap.Error(err))
		return nil, errorx.Wrap(errorx.CodeServerError, err)
	}

	userID, _ := user.LastInsertId()
	l.Logger.Info("user registered", zap.Int64("user_id", userID), zap.String("username", req.Username))

	return &types.RegisterResp{
		UserID:   userID,
		Username: req.Username,
		Nickname: req.Nickname,
	}, nil
}
