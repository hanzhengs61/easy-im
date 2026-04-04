package config

// DBConfig 所有服务通用的数据库配置
type DBConfig struct {
	DataSource string
}

// JWTConfig 所有服务通用的 JWT 配置
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  int64
	RefreshTokenTTL int64
}

type LogConfig struct {
	Level  string // debug | info | warn | error
	Format string // json | console
}
