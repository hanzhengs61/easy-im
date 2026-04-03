package kafka

import (
	"context"
	"easy-im/pkg/logger"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// Producer Kafka 生产者封装
type Producer struct {
	writer *kafka.Writer
}

// ProducerConfig 生产者配置
type ProducerConfig struct {
	Brokers []string
	Topic   string
}

func NewProducer(cfg ProducerConfig) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
		// 批量发送配置：提升吞吐
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		// 异步写失败重试
		MaxAttempts: 3,
		// 生产环境建议 RequiredAcks: kafka.RequireOne
		RequiredAcks: kafka.RequireOne,
	}
	return &Producer{writer: w}
}

// SendMessage 发送消息，key 用于保证同一会话消息顺序
func (p *Producer) SendMessage(ctx context.Context, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: data,
	})
	if err != nil {
		logger.WithContext(ctx).Error("kafka send failed",
			zap.String("topic", p.writer.Topic),
			zap.Error(err),
		)
		return err
	}
	logger.WithContext(ctx).Info("kafka send success",
		zap.String("topic", p.writer.Topic),
		zap.String("key", key),
	)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
