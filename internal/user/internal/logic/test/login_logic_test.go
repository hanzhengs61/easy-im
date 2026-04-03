package test

import (
	"context"
	"easy-im/internal/user/internal/config"
	"easy-im/internal/user/internal/logic"
	"easy-im/internal/user/internal/svc"
	"easy-im/internal/user/internal/types"
	"easy-im/internal/user/model"
	"easy-im/pkg/errorx"
	"easy-im/pkg/jwt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

// newTestSvcCtx 构造测试用 ServiceContext，注入 mock
func newLoginTestSvcCtx(t *testing.T, mockModel model.UsersModel) *svc.ServiceContext {
	t.Helper()
	return &svc.ServiceContext{
		Config:    config.Config{},
		UserModel: mockModel,
		JwtManager: jwt.NewManager(jwt.Config{
			Secret:          "test-secret-key",
			AccessTokenTTL:  2 * time.Hour,
			RefreshTokenTTL: 7 * 24 * time.Hour,
		}),
	}
}

// TestLoginLogic_Login 测试用户登录逻辑
func TestLoginLogic_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), 12)

	mockUsers := model.NewMockUsersModel(ctrl)
	ctx := context.Background()

	// Mokc: 查询用户信息成功
	mockUsers.EXPECT().FindOneByUsername(ctx, "testuser").
		Return(&model.Users{
			Id:       1,
			Username: "testuser",
			Password: string(hashed),
			Nickname: "测试用户",
			Status:   1,
		}, nil)

	svcCtx := newLoginTestSvcCtx(t, mockUsers)
	l := logic.NewLoginLogic(ctx, svcCtx)

	resp, err := l.Login(&types.LoginReq{
		Username: "testuser",
		Password: "correct_password",
	})

	assert.NoError(t, err)
	assert.Equal(t, "测试用户", resp.Nickname)
	assert.Equal(t, int64(1), resp.UserID)
}

// TestLogin_WrongPassword 测试用户密码错误
func TestLogin_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), 12)

	mockUsers := model.NewMockUsersModel(ctrl)
	ctx := context.Background()

	mockUsers.EXPECT().
		FindOneByUsername(ctx, "testuser").
		Return(&model.Users{
			Id:       1,
			Username: "testuser",
			Password: string(hashed),
			Status:   1,
		}, nil)

	svcCtx := newLoginTestSvcCtx(t, mockUsers)
	l := logic.NewLoginLogic(ctx, svcCtx)

	resp, err := l.Login(&types.LoginReq{
		Username: "testuser",
		Password: "wrong_password",
	})

	assert.Nil(t, resp)
	appErr, ok := errorx.IsAppError(err)
	assert.True(t, ok)
	assert.Equal(t, errorx.CodePasswordWrong, appErr.Code)
}

// TestLogin_UserNotFound 测试用户不存在
func TestLogin_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsers := model.NewMockUsersModel(ctrl)
	ctx := context.Background()

	mockUsers.EXPECT().
		FindOneByUsername(ctx, "ghost").
		Return(nil, model.ErrNotFound)

	svcCtx := newLoginTestSvcCtx(t, mockUsers)
	l := logic.NewLoginLogic(ctx, svcCtx)

	_, err := l.Login(&types.LoginReq{
		Username: "ghost",
		Password: "any",
	})

	appErr, ok := errorx.IsAppError(err)
	assert.True(t, ok)
	assert.Equal(t, errorx.CodeUserNotFound, appErr.Code)
}

// TestLogin_DisabledUser 测试用户禁用
func TestLogin_DisabledUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hashed, _ := bcrypt.GenerateFromPassword([]byte("password"), 12)

	mockUsers := model.NewMockUsersModel(ctrl)
	ctx := context.Background()

	mockUsers.EXPECT().
		FindOneByUsername(ctx, "banned").
		Return(&model.Users{
			Id:       2,
			Username: "banned",
			Password: string(hashed),
			Status:   2, // 禁用状态
		}, nil)

	svcCtx := newLoginTestSvcCtx(t, mockUsers)
	l := logic.NewLoginLogic(ctx, svcCtx)

	_, err := l.Login(&types.LoginReq{
		Username: "banned",
		Password: "password",
	})

	appErr, ok := errorx.IsAppError(err)
	assert.True(t, ok)
	assert.Equal(t, errorx.CodeUserDisabled, appErr.Code)
}
