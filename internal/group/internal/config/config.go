// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"easy-im/pkg/config"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	DB  config.DBConfig
	Jwt config.JWTConfig
}
