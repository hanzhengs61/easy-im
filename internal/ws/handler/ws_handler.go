package handler

import (
	"easy-im/internal/ws/client"
	"easy-im/internal/ws/connmgr"
	"easy-im/pkg/logger"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	// 生产环境必须校验 Origin，防止跨站 WebSocket 攻击
	// 这里开发阶段先放开，后续 K8s 时收紧
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
	manager *connmgr.Manager
}

func NewWSHandler(mgr *connmgr.Manager) *WSHandler {
	return &WSHandler{manager: mgr}
}

// ServeHTTP 处理 WebSocket 握手，升级 HTTP 连接
func (h *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// HTTP → WebSocket 协议升级
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("ws upgrade failed", zap.Error(err))
		return
	}

	// 为每个连接分配唯一 ID
	clientID := uuid.NewString()
	c := client.NewClient(clientID, conn, h.manager)

	// 先注册再启动，避免竞态
	h.manager.Register(c)
	go c.Run()

	logger.Info("new ws connection",
		zap.String("client_id", clientID),
		zap.String("remote_addr", conn.RemoteAddr().String()),
	)
}
