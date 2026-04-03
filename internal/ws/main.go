package main

import (
	"easy-im/internal/ws/config"
	"easy-im/internal/ws/connmgr"
	"easy-im/internal/ws/handler"
	"easy-im/pkg/jwt"
	"easy-im/pkg/logger"
	"easy-im/pkg/middleware"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"go.uber.org/zap"
)

var configFile = flag.String("f", "internal/ws/etc/ws.yaml", "the config file")

func main() {
	var cfg config.Config
	conf.MustLoad(*configFile, &cfg)

	logger.Init(logger.Options{
		Level:       cfg.Log.Level,
		Format:      cfg.Log.Format,
		ServiceName: "ws-server",
	})
	defer logger.Sync()

	jwtMgr := jwt.NewManager(jwt.Config{
		Secret:          cfg.JWT.Secret,
		AccessTokenTTL:  time.Duration(cfg.JWT.AccessTokenTTL) * time.Second,
		RefreshTokenTTL: time.Duration(cfg.JWT.RefreshTokenTTL) * time.Second,
	})

	// todo:后续 接入 Kafka 时传入真实 msgHandler，暂时传 nil
	mgr := connmgr.NewManager(jwtMgr, nil)
	wsHandler := handler.NewWSHandler(mgr)

	mux := http.NewServeMux()
	// WS 连接入口
	mux.Handle("/ws", wsHandler)
	// 其他普通 HTTP 接口（健康检查等）才走中间件
	healthHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"status":"ok","online":%d}`, mgr.OnlineCount())
	})
	// 为普通路由套中间件
	protectedMux := middleware.RecoveryMiddleware(
		middleware.LoggerMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/health":
					healthHandler.ServeHTTP(w, r)
				default:
					http.NotFound(w, r)
				}
			}),
		),
	)
	// 使用自定义的 mux，根据路径选择是否走中间件
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws" {
			mux.ServeHTTP(w, r) // WS 直连，不走中间件
		} else {
			protectedMux.ServeHTTP(w, r) // 其他路由走中间件
		}
	})

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info("ws server starting", zap.String("addr", addr))

	if err := http.ListenAndServe(addr, finalHandler); err != nil {
		logger.Fatal("ws server failed", zap.Error(err))
	}
}
