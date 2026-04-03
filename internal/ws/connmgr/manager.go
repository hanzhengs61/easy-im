package connmgr

import (
	"easy-im/internal/ws/client"
	"easy-im/pkg/errorx"
	"easy-im/pkg/jwt"
	"easy-im/pkg/logger"
	"easy-im/pkg/protocol"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Manager 连接管理器
// 核心职责：
//  1. 维护 clientID → *Client 映射（所有连接）
//  2. 维护 userID → []*Client 映射（一个用户可多端登录）
//  3. 消息路由：根据 toID 找到目标连接并投递
type Manager struct {
	// clientID → *Client（所有连接，含未鉴权的）
	clients sync.Map
	// userID → []*Client（已鉴权用户的连接列表）
	userClients sync.Map
	// Jwt 管理器，用于连接鉴权
	jwtManager *jwt.Manager
	// 消息下游处理器（发给 Kafka 持久化）
	msgHandler MsgHandler
}

// MsgHandler 消息持久化/转发处理器接口
type MsgHandler interface {
	HandleChatMsg(msg *protocol.Message, formUID int64) error
}

func NewManager(jwtManager *jwt.Manager, msgHandler MsgHandler) *Manager {
	manager := &Manager{
		jwtManager: jwtManager,
		msgHandler: msgHandler,
	}
	// 启动定时清理僵尸连接
	go manager.cleanDeadClients()
	return manager
}

// Register 注册新连接（ws 握手成功后调用）
func (m *Manager) Register(c *client.Client) {
	m.clients.Store(c.ID, c)
	logger.Info("client registered", zap.String("client_id", c.ID))
}

// Unregister 注销连接（连接断开时由 readPump defer 调用）
func (m *Manager) Unregister(c *client.Client) {
	m.clients.Delete(c.ID)
	if c.UserId != 0 {
		m.removeUserClient(c.UserId, c.ID)
	}
}

// HandleMessage 处理客户端上行消息，这是消息路由的核心入口
func (m *Manager) HandleMessage(c *client.Client, msg *protocol.Message) {
	switch msg.Type {
	case protocol.MsgTypeAuth:
		m.handleAuth(c, msg)
	case protocol.MsgTypeText, protocol.MsgTypeImage,
		protocol.MsgTypeAudio, protocol.MsgTypeVideo, protocol.MsgTypeFile:
		// 未鉴权的连接不能发消息
		if c.UserId == 0 {
			m.sendAuthRequired(c)
			return
		}
		m.handleChatMsg(c, msg)
	}

}

// handleChatMsg 处理聊天消息
func (m *Manager) handleChatMsg(c *client.Client, msg *protocol.Message) {
	msg.FromUID = c.UserId
	msg.SendTime = time.Now().UnixMilli()

	// 1. 先给发送方 ACK（告诉客户端消息已到服务端）
	m.sendAck(c, msg.Seq, msg.SendTime)

	// 2. 在线投递：找目标用户的连接直接推
	if msg.ChatType == protocol.ChatTypeSingle {
		m.deliverToUser(msg.ToID, msg)
	}

	// 3. 交给下游处理器（写 Kafka → Message 服务持久化）
	// todo: 预留
	if m.msgHandler != nil {
		if err := m.msgHandler.HandleChatMsg(msg, c.UserId); err != nil {
			logger.Error("msg handler failed",
				zap.Error(err),
				zap.Int64("from_uid", c.UserId),
			)
		}
	}
}

// DeliverToUser 投递消息给指定用户的所有在线连接（可从外部调用，如离线消息补发）
func (m *Manager) DeliverToUser(userID int64, msg *protocol.Message) {
	m.deliverToUser(userID, msg)
}

// deliverToUser 投递消息给指定用户的所有在线连接
func (m *Manager) deliverToUser(userID int64, msg *protocol.Message) {
	clientIDs, ok := m.userClients.Load(userID)
	if !ok {
		// 用户不在线，由 Message 服务存离线消息
		return
	}

	ids := clientIDs.([]string)
	for _, id := range ids {
		if c, ok := m.clients.Load(id); ok {
			c.(*client.Client).SendMsg(msg)
		}
	}
}

// OnlineCount 当前在线连接数
func (m *Manager) OnlineCount() int {
	count := 0
	m.clients.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

// IsUserOnline 判断用户是否在线
func (m *Manager) IsUserOnline(userID int64) bool {
	_, ok := m.userClients.Load(userID)
	return ok
}

// ───────────────────────── 内部辅助方法 ─────────────────────────

// cleanDeadClients 定时清理心跳超时的僵尸连接（双重保险）
func (m *Manager) cleanDeadClients() {
	ticker := time.NewTimer(60 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		deadline := time.Now().Add(-90 * time.Second)
		m.clients.Range(func(_, val any) bool {
			c := val.(*client.Client)
			if c.LastHeartbeat.Before(deadline) {
				logger.Info("clean dead client",
					zap.String("client_id", c.ID),
					zap.Int64("user_id", c.UserId),
				)
				c.Close()
			}
			return true
		})
	}
}

// removeUserClient 删除用户连接
func (m *Manager) removeUserClient(userID int64, clientID string) {
	actual, ok := m.userClients.Load(userID)
	if !ok {
		return
	}
	ids := actual.([]string)
	newIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != clientID {
			newIDs = append(newIDs, id)
		}
	}
	if len(newIDs) == 0 {
		m.userClients.Delete(userID)
	} else {
		m.userClients.Store(userID, newIDs)
	}
}

// handleAuth 处理连接鉴权
// 客户端连上后必须在 10s 内发送 Auth 消息，否则心跳超时会断开
func (m *Manager) handleAuth(c *client.Client, msg *protocol.Message) {
	var authContent protocol.AuthContent
	if err := json.Unmarshal(msg.Content, &authContent); err != nil {
		m.sendAuthAck(c, errorx.CodeInvalidParam, "消息格式错误", 0)
		c.Close()
		return
	}
	claims, err := m.jwtManager.ParseToken(authContent.Token)
	if err != nil {
		code := errorx.CodeTokenInvalid
		msg := "token 无效"
		if errors.Is(err, errorx.New(errorx.CodeTokenExpired)) { // 或者直接判断 err
			code = errorx.CodeTokenExpired
			msg = "token 已过期"
		}
		m.sendAuthAck(c, code, msg, 0)
		c.Close()
		return
	}
	// 设置用户 ID，注册到 user → clients 映射
	c.SetUserID(claims.UserID)
	m.addUserClient(claims.UserID, c)
}

// sendAuthAck 发送鉴权结果
func (m *Manager) sendAuthAck(c *client.Client, code int, msg string, userID int64) {
	content, _ := json.Marshal(&protocol.AuthAckContent{
		Code:   code,
		Msg:    msg,
		UserID: userID,
	})
	c.SendMsg(&protocol.Message{
		Type:    protocol.MsgTypeAuthAck,
		Content: content,
	})
}

func (m *Manager) addUserClient(userID int64, c *client.Client) {
	actual, _ := m.userClients.LoadOrStore(userID, []string{})
	ids := actual.([]string)
	ids = append(ids, c.ID)
	m.userClients.Store(userID, ids)
}

func (m *Manager) sendAuthRequired(c *client.Client) {
	m.sendAuthAck(c, errorx.CodeUnauthorized, "请先完成鉴权", 0)
}

func (m *Manager) sendAck(c *client.Client, seq int64, msgID int64) {
	content, _ := json.Marshal(&protocol.AckContent{
		AckSeq: seq,
		MsgID:  msgID,
	})
	c.SendMsg(&protocol.Message{
		Type:    protocol.MsgTypeAck,
		Content: content,
	})
}
