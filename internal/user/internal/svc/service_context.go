// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"easy-im/internal/user/internal/config"
	"easy-im/internal/user/model"
	"easy-im/pkg/jwt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	JwtManager *jwt.Manager
	UserModel  model.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 MySQL 连接
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,
		JwtManager: jwt.NewManager(jwt.Config{
			Secret:          c.Jwt.Secret,
			AccessTokenTTL:  time.Duration(c.Jwt.AccessTokenTTL) * time.Second,
			RefreshTokenTTL: time.Duration(c.Jwt.RefreshTokenTTL) * time.Second,
		}),
		UserModel: model.NewUsersModel(conn),
	}
}
