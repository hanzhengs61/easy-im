// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"easy-im/internal/user/internal/config"
	"easy-im/internal/user/internal/handler"
	"easy-im/internal/user/internal/svc"
	"easy-im/pkg/middleware"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "internal/user/etc/user-api.yaml", "the config file")

func main() {
	stat.DisableLog()
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.Use(middleware.RecoveryMiddleware)
	server.Use(middleware.LoggerMiddleware)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("🀄🀄️🀄️：user server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
