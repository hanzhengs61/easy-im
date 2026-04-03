package protocol

import "encoding/json"

// MsgType 消息类型
type MsgType int32

const (
	MsgTypeHeartbeat    MsgType = 101 // 心跳 ping
	MsgTypeHeartbeatAck MsgType = 102 // 心跳 pong
	MsgTypeAuth         MsgType = 103 // 连接后鉴权
	MsgTypeAuthAck      MsgType = 104 // 鉴权响应
	MsgTypeKickOut      MsgType = 105 // 被踢下线

	MsgTypeText  MsgType = 201 // 文本消息
	MsgTypeImage MsgType = 202 // 图片消息
	MsgTypeAudio MsgType = 203 // 语音消息
	MsgTypeVideo MsgType = 204 // 视频消息
	MsgTypeFile  MsgType = 205 // 文件消息

	MsgTypeAck          MsgType = 301 // 消息已送达 ACK
	MsgTypeRead         MsgType = 302 // 消息已读
	MsgTypeTyping       MsgType = 303 // 正在输入
	MsgTypeOnlineStatus MsgType = 304 // 在线状态变更
)

// ChatType 会话类型
type ChatType int32

const (
	ChatTypeSingle ChatType = 1 // 单聊
	ChatTypeGroup  ChatType = 2 // 群聊
)

// Message WebSocket 消息帧，所有消息统一格式
// 客户端和服务端都使用此结构收发消息
type Message struct {
	Seq      int64           `json:"seq"`                 // 消息序号，客户端自增，用于 ACK 对应（防重）
	Type     MsgType         `json:"type"`                // 消息类型
	ChatType ChatType        `json:"chat_type,omitempty"` // 会话相关（聊天消息必填）
	FromUID  int64           `json:"from_uid,omitempty"`  // 发送方（服务端下发时填充，客户端上行不需要）
	ToID     int64           `json:"to_id,omitempty"`     // 接收方：单聊填用户ID，群聊填群组ID
	Content  json.RawMessage `json:"content,omitempty"`   // 消息内容（JSON 格式，根据 Type 反序列化为不同结构）
	SendTime int64           `json:"send_time,omitempty"` // 服务端时间戳（毫秒），客户端以此为准显示时间
}

// AuthContent 鉴权消息内容
type AuthContent struct {
	Token string `json:"token"`
}

// AuthAckContent 鉴权响应
type AuthAckContent struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	UserID int64  `json:"user_id,omitempty"`
}

// TextContent 文本消息内容
type TextContent struct {
	Text string `json:"text"`
}

// AckContent ACK 消息内容
type AckContent struct {
	// 对应客户端发送的 seq
	AckSeq int64 `json:"ack_seq"`
	// 服务端为该消息分配的全局唯一 ID
	MsgID int64 `json:"msg_id"`
}

// KickContent 踢下线原因
type KickContent struct {
	Reason string `json:"reason"`
}

// Encode 序列化为 JSON 字节
func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

// Decode 从 JSON 字节反序列化
func Decode(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
