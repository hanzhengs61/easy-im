// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"easy-im/internal/group/internal/config"
	"easy-im/internal/group/internal/handler"
	"easy-im/internal/group/internal/svc"
	"easy-im/pkg/logger"
	"easy-im/pkg/middleware"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "internal/group/etc/group-api.yaml", "the config file")

func main() {
	flag.Parse()

	var cfg config.Config
	conf.MustLoad(*configFile, &cfg)

	server := rest.MustNewServer(cfg.RestConf)
	defer server.Stop()

	server.Use(middleware.RecoveryMiddleware)
	server.Use(middleware.LoggerMiddleware)

	logger.Init(logger.Options{
		Level:       cfg.Log.Level,
		ServiceName: cfg.Name,
	})
	defer logger.Sync()

	ctx := svc.NewServiceContext(cfg)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", cfg.Host, cfg.Port)
	server.Start()
}
