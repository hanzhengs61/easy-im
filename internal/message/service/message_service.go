package service

import (
	"context"
	"easy-im/internal/ws/msghandler"
	"easy-im/pkg/logger"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.uber.org/zap"
)

// MessageService 消息服务：消费 Kafka，持久化消息
type MessageService struct {
	db         sqlx.SqlConn
	mongoStore *MongoStore
}

func NewMessageService(db sqlx.SqlConn, mongoStore *MongoStore) *MessageService {
	return &MessageService{
		db:         db,
		mongoStore: mongoStore,
	}
}

// HandleKafkaMsg Kafka 消费回调，每条消息都走这里
func (s *MessageService) HandleKafkaMsg(ctx context.Context, _, value []byte) error {
	var event msghandler.KafkaMsgEvent
	if err := json.Unmarshal(value, &event); err != nil {
		logger.WithContext(ctx).Error("unmarshal kafka msg failed", zap.Error(err))
		// 格式错误直接丢弃，不重试
		return nil
	}

	return s.saveMessage(ctx, &event)
}

// saveMessage 保存消息
func (s *MessageService) saveMessage(ctx context.Context, event *msghandler.KafkaMsgEvent) error {
	now := time.Now().UnixMilli()

	// 1. 写 MySQL 消息索引（获取全局自增 ID）
	result, err := s.db.ExecCtx(ctx,
		`INSERT INTO messages (seq, chat_type, from_uid, to_id, msg_type, status, send_time, created_at)
		 VALUES (?, ?, ?, ?, ?, 1, ?, ?)`,
		event.Seq, event.ChatType, event.FromUID, event.ToID,
		event.MsgType, event.SendTime, now,
	)
	if err != nil {
		logger.WithContext(ctx).Error("insert message index failed", zap.Error(err))
		return err
	}

	msgID, _ := result.LastInsertId()

	// 2. 写 MongoDB 消息正文
	doc := &MsgDoc{
		MsgID:    msgID,
		Seq:      event.Seq,
		MsgType:  event.MsgType,
		ChatType: event.ChatType,
		FromUID:  event.FromUID,
		ToID:     event.ToID,
		Content:  string(event.Content),
		SendTime: event.SendTime,
	}
	if err = s.mongoStore.Insert(ctx, doc); err != nil {
		logger.WithContext(ctx).Error("insert message mongo failed",
			zap.Int64("msg_id", msgID),
			zap.Error(err),
		)
		// MongoDB 失败不回滚 MySQL（消息不丢，后续可补偿）
		return err
	}

	// 3. 更新会话表（upsert）
	if err = s.upsertConversation(ctx, event, msgID, now); err != nil {
		logger.WithContext(ctx).Warn("upsert conversation failed", zap.Error(err))
		// 会话更新失败不影响消息存储，记录告警即可
	}

	logger.WithContext(ctx).Info("message saved",
		zap.Int64("msg_id", msgID),
		zap.Int64("from_uid", event.FromUID),
		zap.Int64("to_id", event.ToID),
	)
	return nil
}

func (s *MessageService) upsertConversation(ctx context.Context, event *msghandler.KafkaMsgEvent, msgID, now int64) error {
	upsertSQL := `
		INSERT INTO conversations (owner_uid, target_id, chat_type, last_msg_id, last_msg_seq, unread_count, updated_at, created_at)
		VALUES (?, ?, ?, ?, ?, 1, ?, ?)
		ON DUPLICATE KEY UPDATE
			last_msg_id  = VALUES(last_msg_id),
			last_msg_seq = VALUES(last_msg_seq),
			unread_count = unread_count + 1,
			updated_at   = VALUES(updated_at)`

	// 更新接收方会话
	if _, err := s.db.ExecCtx(ctx, upsertSQL,
		event.ToID, event.FromUID, event.ChatType,
		msgID, event.Seq, now, now,
	); err != nil {
		return err
	}

	// 更新发送方会话（发送方不增加未读）
	_, err := s.db.ExecCtx(ctx, `
		INSERT INTO conversations (owner_uid, target_id, chat_type, last_msg_id, last_msg_seq, unread_count, updated_at, created_at)
		VALUES (?, ?, ?, ?, ?, 0, ?, ?)
		ON DUPLICATE KEY UPDATE
			last_msg_id  = VALUES(last_msg_id),
			last_msg_seq = VALUES(last_msg_seq),
			updated_at   = VALUES(updated_at)`,
		event.FromUID, event.ToID, event.ChatType,
		msgID, event.Seq, now, now,
	)
	return err
}
