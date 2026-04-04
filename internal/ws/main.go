package main

import (
	"easy-im/internal/ws/config"
	"easy-im/internal/ws/connmgr"
	"easy-im/internal/ws/handler"
	"easy-im/internal/ws/msghandler"
	"easy-im/pkg/jwt"
	"easy-im/pkg/kafka"
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

	// Kafka 生产者
	producer := kafka.NewProducer(kafka.ProducerConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "easy-im-messages",
	})
	defer producer.Close()

	kafkaHandler := msghandler.NewKafkaHandler(producer)
	mgr := connmgr.NewManager(jwtMgr, kafkaHandler)
	wsHandler := handler.NewWSHandler(mgr)

	mux := http.NewServeMux()
	// WS 连接入口
	mux.Handle("/ws", wsHandler)
	// 其他普通 HTTP 接口（健康检查等）才走中间件
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"status":"ok","online":%d}`, mgr.OnlineCount())
	})
	chain := middleware.RecoveryHandler(middleware.LoggerHandler(mux))

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	logger.Info("ws server started", zap.String("addr", addr))
	if err := http.ListenAndServe(addr, chain); err != nil {
		logger.Fatal("ws server error", zap.Error(err))
	}
}
