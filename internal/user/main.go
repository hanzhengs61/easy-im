// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"easy-im/pkg/middleware"
	"flag"
	"fmt"
	"net/http"

	"easy-im/internal/user/internal/config"
	"easy-im/internal/user/internal/handler"
	"easy-im/internal/user/internal/svc"

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

	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		// 将 LoggerMiddleware 包装为 http.HandlerFunc 签名
		return func(w http.ResponseWriter, r *http.Request) {
			middleware.LoggerMiddleware(next).ServeHTTP(w, r)
		}
	})

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("🀄🀄️🀄️：user server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
