package jwt

import (
	"easy-im/pkg/errorx"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义 JWT 载荷
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Config JWT 配置，从配置文件读取后传入
type Config struct {
	Secret          string
	AccessTokenTTL  time.Duration // 建议 2小时
	RefreshTokenTTL time.Duration // 建议 7天
}

// Manager JWT 管理器
type Manager struct {
	cfg Config
}

func NewManager(cfg Config) *Manager {
	return &Manager{cfg: cfg}
}

// GenerateAccessToken 签发 access token
func (m *Manager) GenerateAccessToken(userID int64, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.cfg.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "easy-im",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.cfg.Secret))
}

// GenerateRefreshToken 签发 refresh token（TTL 更长）
func (m *Manager) GenerateRefreshToken(userID int64, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.cfg.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "easy-im-refresh",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.cfg.Secret))
}

// ParseToken 解析并校验 token，返回 Claims
func (m *Manager) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		// 防止算法混淆攻击：强制要求 HMAC 系列算法
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errorx.New(errorx.CodeTokenInvalid)
		}
		return []byte(m.cfg.Secret), nil
	})
	if err != nil {
		// 区分过期和无效
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errorx.New(errorx.CodeTokenExpired)
		}
		return nil, errorx.New(errorx.CodeTokenInvalid)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errorx.New(errorx.CodeTokenInvalid)
	}
	return claims, nil
}
