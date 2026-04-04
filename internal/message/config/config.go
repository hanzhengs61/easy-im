package config

import "easy-im/pkg/config"

type Config struct {
	Name  string
	DB    config.DBConfig
	Log   config.LogConfig
	Mongo struct {
		URI      string
		Database string
	}
	Kafka struct {
		Brokers []string
		Topic   string
		GroupID string
	}
}
