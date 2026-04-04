package config

import "easy-im/pkg/config"

type Config struct {
	Host string
	Port int
	JWT  config.JWTConfig
	Log  config.LogConfig
}
