// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string `json:",optional"`
		Table      string `json:",optional"`
	} `json:"Mysql"`
	Jwt struct {
		Secret          string `json:",optional"`
		AccessTokenTTL  int64  `json:",optional"` // 秒
		RefreshTokenTTL int64  `json:",optional"` // 秒
	} `json:"Jwt"`
}
