package kafka

import (
	"context"
	"easy-im/pkg/logger"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

// Consumer Kafka 消费者封装
type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg ConsumerConfig) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		GroupID:  cfg.GroupID,
		MinBytes: 1,
		MaxBytes: 10e6, // 10MB
	})
	return &Consumer{reader: r}
}

// ConsumeHandler 消息处理函数类型
type ConsumeHandler func(ctx context.Context, key, value []byte) error

// StartConsume 阻塞式消费，应在独立 goroutine 中调用
func (c *Consumer) StartConsume(ctx context.Context, handler ConsumeHandler) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				// context 取消，正常退出
				return
			}
			logger.WithContext(ctx).Error("kafka fetch failed", zap.Error(err))
			continue
		}

		if err = handler(ctx, msg.Key, msg.Value); err != nil {
			logger.WithContext(ctx).Error("kafka handler failed",
				zap.Error(err),
				zap.ByteString("key", msg.Key),
			)
			// 处理失败不提交 offset，会重试
			continue
		}

		// 手动提交 offset（至少一次语义）
		if err = c.reader.CommitMessages(ctx, msg); err != nil {
			logger.WithContext(ctx).Error("kafka commit failed", zap.Error(err))
		}
		logger.WithContext(ctx).Info("kafka consume success",
			zap.ByteString("key", msg.Key),
			zap.ByteString("value", msg.Value),
		)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
