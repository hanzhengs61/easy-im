package msghandler

import (
	"context"
	"easy-im/pkg/kafka"
	"easy-im/pkg/logger"
	"easy-im/pkg/protocol"
	"fmt"

	"go.uber.org/zap"
)

// KafkaMsgEvent 发往 Kafka 的消息事件结构
type KafkaMsgEvent struct {
	Seq      int64             `json:"seq"`
	MsgType  protocol.MsgType  `json:"msg_type"`
	ChatType protocol.ChatType `json:"chat_type"`
	FromUID  int64             `json:"from_uid"`
	ToID     int64             `json:"to_id"`
	Content  []byte            `json:"content"`
	SendTime int64             `json:"send_time"`
}

// KafkaHandler 实现 connmgr.MsgHandler 接口
type KafkaHandler struct {
	producer *kafka.Producer
}

func NewKafkaHandler(producer *kafka.Producer) *KafkaHandler {
	return &KafkaHandler{producer: producer}
}

func (h *KafkaHandler) HandleChatMsg(msg *protocol.Message, fromUID int64) error {
	event := &KafkaMsgEvent{
		Seq:      msg.Seq,
		MsgType:  msg.Type,
		ChatType: msg.ChatType,
		FromUID:  fromUID,
		ToID:     msg.ToID,
		Content:  msg.Content,
		SendTime: msg.SendTime,
	}

	// Kafka key：保证同一会话的消息落到同一 partition，保序
	key := fmt.Sprintf("%d_%d_%d", msg.ChatType, fromUID, msg.ToID)

	ctx := context.Background()
	if err := h.producer.SendMessage(ctx, key, event); err != nil {
		logger.Error("send to kafka failed",
			zap.Int64("from_uid", fromUID),
			zap.Int64("to_id", msg.ToID),
			zap.Error(err),
		)
		return err
	}

	logger.Info("msg sent to kafka",
		zap.Int64("from_uid", fromUID),
		zap.Int64("to_id", msg.ToID),
		zap.Int64("seq", msg.Seq),
	)
	return nil
}
