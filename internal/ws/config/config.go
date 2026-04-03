package config

type Config struct {
	Host string
	Port int

	JWT struct {
		Secret          string
		AccessTokenTTL  int64
		RefreshTokenTTL int64
	}

	Log struct {
		Level  string
		Format string
	}
}
