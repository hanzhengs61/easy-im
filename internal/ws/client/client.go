package client

import (
	"context"
	"easy-im/pkg/logger"
	"easy-im/pkg/protocol"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// 写超时：超时此时间写不出去就断开
	writeTimeout = 10 * time.Second
	// 心跳超时：超过此时间没有收到心跳就断开
	heartbeatTimeout = 90 * time.Second
	// 发送缓冲区大小：堆积超过此数量就断开（背压控制）
	sendBufSize = 256
)

// Client 代表一个在线的 ws 连接
type Client struct {
	// 连接唯一标识
	ID string
	// 鉴权后填充的用户信息
	UserId int64
	// 底层 ws 连接
	conn *websocket.Conn
	// 发送缓冲通道：其他 goroutine 向此 channer 投递消息
	sendCh chan []byte
	// 连接关闭信息
	closeCh chan struct{}
	// 保证 Close 只执行一次
	once sync.Once
	// 最后一次收到心跳的时间
	LastHeartbeat time.Time
	mu            sync.RWMutex
	// 连接所属的管理器（用于注销自己）
	manager Manager
	ctx     context.Context
	cancel  context.CancelFunc
}

// Manager 定义连接管理器接口，避免循环依赖
type Manager interface {
	// Unregister 注销一个连接
	Unregister(client *Client)
	// HandleMessage 处理一个消息
	HandleMessage(client *Client, msg *protocol.Message)
}

// NewClient 创建一个 Client
func NewClient(id string, conn *websocket.Conn, mgr Manager) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ID:            id,
		conn:          conn,
		sendCh:        make(chan []byte, sendBufSize),
		closeCh:       make(chan struct{}),
		LastHeartbeat: time.Now(),
		manager:       mgr,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Run 启动读写两个 goroutine，这是连接的主循环
// 必须在 goroutine 中调用：go client.Run()
func (c *Client) Run() {
	go c.readPump()
	go c.writePump()
}

// readPump 读循环：持续从 ws 读消息
// 每个连接独占一个 goroutine，阻塞式读
func (c *Client) readPump() {
	// 用于定时发送 Ping 帧（针对不发心跳的客户端做服务端主动探活）
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	// 设置读超时：每次读必须在心跳超时内完成
	c.conn.SetReadDeadline(time.Now().Add(heartbeatTimeout))
	// 收到 Pone 帧时重置读超时（浏览器心跳）
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(heartbeatTimeout))
		c.updateHeartbeat()
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.WithContext(c.ctx).Warn("ws read error",
					zap.String("client_id", c.ID),
					zap.Error(err))
			}
			return
		}
		// 重置读超时
		c.conn.SetReadDeadline(time.Now().Add(heartbeatTimeout))

		// 解析消息
		msg, err := protocol.Decode(data)
		if err != nil {
			logger.WithContext(c.ctx).Warn("invalid message format",
				zap.String("client_id", c.ID),
				zap.Error(err),
			)
			continue
		}

		// 心跳直接在这里处理，不用走消息路由
		if msg.Type == protocol.MsgTypeHeartbeat {
			c.updateHeartbeat()
			c.SendMsg(&protocol.Message{Type: protocol.MsgTypeHeartbeatAck})
			continue
		}

		// 其他消息交给 Manager 路由
		c.manager.HandleMessage(c, msg)
	}
}

// writePump 写循环：从 sendCh 取消息写入 ws
// 所有写操作都在这一个 goroutine 里，避免并发写冲突
func (c *Client) writePump() {
	// 用于定时发送 Ping 帧（针对不发心跳的客户端做服务端主动探活）
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case data, ok := <-c.sendCh:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if !ok {
				// sendCh 已关闭，发送关闭帧
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				logger.WithContext(c.ctx).Warn("ws write error",
					zap.String("client_id", c.ID),
					zap.Error(err),
				)
				return
			}
		case <-ticker.C:
			// 主动发送 Ping
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.closeCh:
			return
		}
	}
}

// updateHeartbeat 更新最后一次收到心跳的时间
func (c *Client) updateHeartbeat() {
	c.mu.Lock()
	c.LastHeartbeat = time.Now()
	c.mu.Unlock()
}

// SendMsg 向客户端发送消息（非阻塞，满了就断开）
func (c *Client) SendMsg(msg *protocol.Message) {
	data, err := msg.Encode()
	if err != nil {
		return
	}
	select {
	case c.sendCh <- data:
	default:
		// 发送缓冲区满，说明客户端消费太慢，主动断开
		logger.WithContext(c.ctx).Warn("send buffer full, closing client",
			zap.String("client_id", c.ID),
			zap.Int64("user_id", c.UserId),
		)
		c.Close()
	}
}

// Close 关闭连接，保证幂等（多次调用只执行一次）
func (c *Client) Close() {
	c.once.Do(func() {
		c.cancel()
		close(c.closeCh)
		close(c.sendCh)
	})
}

// SetUserID 鉴权成功后设置用户 Id
func (c *Client) SetUserID(uid int64) {
	c.mu.Lock()
	c.UserId = uid
	c.mu.Unlock()
}
