package database

import (
	"database/sql"
	"time"
)

type MySQLConfig struct {
	// 数据源
	DataSource string
	// 最大连接数
	MaxOpenConns int
	// 最大空闲连接数
	MaxIdleConns int
	// 连接最大空闲时长
	ConnMaxLifetime time.Duration
}

func NewMySQL(cfg MySQLConfig) (*sql.DB, error) {
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 100
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 10
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = time.Hour
	}

	db, err := sql.Open("mysql", cfg.DataSource)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
