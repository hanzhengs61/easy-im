package test

import (
	"context"
	"database/sql"
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
)

// newTestSvcCtx 构造测试用 ServiceContext，注入 mock
func newRegisterTestSvcCtx(t *testing.T, mockModel model.UsersModel) *svc.ServiceContext {
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

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsers := model.NewMockUsersModel(ctrl)
	ctx := context.Background()

	// Mock：用户名不存在
	mockUsers.EXPECT().
		FindOneByUsername(ctx, "testuser").
		Return(nil, model.ErrNotFound)

	// Mock：插入成功
	mockUsers.EXPECT().
		Insert(ctx, gomock.Any()).
		DoAndReturn(func(_ context.Context, user *model.Users) (sql.Result, error) {
			return sqlResultMock{lastInsertID: 10001}, nil
		})

	svcCtx := newRegisterTestSvcCtx(t, mockUsers)
	l := logic.NewRegisterLogic(ctx, svcCtx)

	req := &types.RegisterReq{
		Username: "testuser",
		Password: "123456",
		Nickname: "测试用户",
	}
	resp, err := l.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(10001), resp.UserID)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsers := model.NewMockUsersModel(ctrl)
	ctx := context.Background()

	// mock：用户名已存在
	mockUsers.EXPECT().
		FindOneByUsername(ctx, "existuser").
		Return(&model.Users{Id: 1, Username: "existuser"}, nil)

	svcCtx := newRegisterTestSvcCtx(t, mockUsers)
	l := logic.NewRegisterLogic(ctx, svcCtx)

	resp, err := l.Register(&types.RegisterReq{
		Username: "existuser",
		Password: "123456",
		Nickname: "已存在用户",
	})

	assert.Nil(t, resp)
	assert.Error(t, err)

	appErr, ok := errorx.IsAppError(err)
	assert.True(t, ok)
	assert.Equal(t, errorx.CodeUserAlreadyExists, appErr.Code)
}

// sqlResultMock 实现 sql.Result 接口，用于 mock Insert 返回值
type sqlResultMock struct {
	lastInsertID int64
	rowsAffected int64
}

func (s sqlResultMock) LastInsertId() (int64, error) { return s.lastInsertID, nil }
func (s sqlResultMock) RowsAffected() (int64, error) { return s.rowsAffected, nil }
