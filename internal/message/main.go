package main

import (
	"context"
	"easy-im/internal/message/config"
	"easy-im/internal/message/service"
	"easy-im/pkg/kafka"
	"easy-im/pkg/logger"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var configFile = flag.String("f", "internal/message/etc/message.yaml", "the config file")

func main() {
	var cfg config.Config
	conf.MustLoad(*configFile, &cfg)

	logger.Init(logger.Options{
		Level:       cfg.Log.Level,
		Format:      cfg.Log.Format,
		ServiceName: cfg.Name,
	})
	defer logger.Sync()

	// 初始化 MySQL 连接
	conn := sqlx.NewSqlConn("mysql", cfg.DB.DataSource)

	// MongoDB
	mongoCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		logger.Fatal("mongo connect failed", zap.Error(err))
	}
	defer mongoClient.Disconnect(context.Background())

	mongoStore := service.NewMongoStore(mongoClient, cfg.Mongo.Database)
	msgSvc := service.NewMessageService(conn, mongoStore)

	// Kafka 消费者
	consumer := kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		GroupID: cfg.Kafka.GroupID,
	})
	defer consumer.Close()

	ctx, stop := context.WithCancel(context.Background())

	// 启动消费
	go func() {
		logger.Info("message consumer starting...")
		consumer.StartConsume(ctx, msgSvc.HandleKafkaMsg)
	}()

	logger.Info("message service started")

	// 优雅退出：等待 SIGTERM / SIGINT
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down message service...")
	stop()
	time.Sleep(2 * time.Second) // 等待消费者处理完当前消息
}
