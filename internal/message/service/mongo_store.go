package service

import (
	"context"
	"easy-im/pkg/logger"
	"easy-im/pkg/protocol"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// MsgDoc MongoDB 消息文档
type MsgDoc struct {
	MsgID    int64             `bson:"msg_id"`
	Seq      int64             `bson:"seq"`
	MsgType  protocol.MsgType  `bson:"msg_type"`
	ChatType protocol.ChatType `bson:"chat_type"`
	FromUID  int64             `bson:"from_uid"`
	ToID     int64             `bson:"to_id"`
	Content  string            `bson:"content"`
	SendTime int64             `bson:"send_time"`
	CreateAt int64             `bson:"created_at"`
}

// MongoStore MongoDB 消息存储
type MongoStore struct {
	col *mongo.Collection
}

func NewMongoStore(client *mongo.Client, dbName string) *MongoStore {
	col := client.Database(dbName).Collection("messages")

	// 创建索引
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "from_uid", Value: 1},
				{Key: "to_id", Value: 1},
				{Key: "send_time", Value: -1},
			},
		},
		{
			Keys:    bson.D{{Key: "msg_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	if _, err := col.Indexes().CreateMany(ctx, indexes); err != nil {
		logger.Error("create mongo indexes failed", zap.Error(err))
	}

	return &MongoStore{col: col}
}

// Insert 插入消息文档
func (s *MongoStore) Insert(ctx context.Context, doc *MsgDoc) error {
	doc.CreateAt = time.Now().UnixMilli()
	_, err := s.col.InsertOne(ctx, doc)
	return err
}

// GetHistory 获取会话历史消息（按时间倒序）
func (s *MongoStore) GetHistory(ctx context.Context, fromUID, toID int64, beforeTime int64, limit int) ([]*MsgDoc, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from_uid": fromUID, "to_id": toID},
			{"from_uid": toID, "to_id": fromUID},
		},
		"send_time": bson.M{"$lt": beforeTime},
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "send_time", Value: -1}}).
		SetLimit(int64(limit))

	cur, err := s.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var docs []*MsgDoc
	if err = cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}
