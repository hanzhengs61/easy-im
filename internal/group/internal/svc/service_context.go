// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"easy-im/internal/group/internal/config"
	"easy-im/internal/group/model"
	"easy-im/pkg/jwt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config      config.Config
	GroupModel  model.GroupsModel
	MemberModel model.GroupMembersModel
	JWTManager  *jwt.Manager
}

func NewServiceContext(c config.Config) *ServiceContext {

	conn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config:      c,
		GroupModel:  model.NewGroupsModel(conn),
		MemberModel: model.NewGroupMembersModel(conn),
		JWTManager: jwt.NewManager(jwt.Config{
			Secret:          c.Jwt.Secret,
			AccessTokenTTL:  time.Duration(c.Jwt.AccessTokenTTL) * time.Second,
			RefreshTokenTTL: time.Duration(c.Jwt.RefreshTokenTTL) * time.Second,
		}),
	}
}
