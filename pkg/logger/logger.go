package logger

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// contextKey 避免语其他包的 context key 冲突
type contextKey string

const traceIDKey contextKey = "trace_id"

var (
	global *zap.Logger
	once   sync.Once
)

// Options 日志初始化配置
type Options struct {
	Level       string // debug | info | warn | error
	Format      string // json | console
	ServiceName string
}

// Init 初始化全局 logger，应在 main 函数最开始调用
func Init(opts Options) {
	once.Do(func() {
		global = newLogger(opts)
	})
}

func newLogger(opts Options) *zap.Logger {
	// 解析日志级别
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(opts.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	// 编码器配置
	encCfg := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 选择编码格式：生产用 json，开发用 console
	var enc zapcore.Encoder
	if opts.Format == "console" {
		enc = zapcore.NewConsoleEncoder(encCfg)
	} else {
		enc = zapcore.NewJSONEncoder(encCfg)
	}

	// 创建日志核心
	core := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), level)

	return zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		// 固定字段：服务名，所有日志都携带
		zap.Fields(zap.String("service", opts.ServiceName)),
	)
}

// WithContext 从 context 中提取 trace_id，返回携带该字段的 logger
// 用法：logger.WithContext(ctx).Info("处理请求")
func WithContext(ctx context.Context) *zap.Logger {
	l := getGlobal()
	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		return l.With(zap.String("trace_id", traceID))
	}
	return l
}

// NewContext 将 trace_id 注入 context，在请求入口处调用
func NewContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// 以下是对 zap 方法的直接封装，全局可用

func Debug(msg string, fields ...zap.Field) {
	getGlobal().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	getGlobal().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	getGlobal().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	getGlobal().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	getGlobal().Fatal(msg, fields...)
}

// Sync 程序退出前调用，刷新缓冲区
func Sync() {
	_ = getGlobal().Sync()
}

func getGlobal() *zap.Logger {
	if global == nil {
		// 防止未初始化时 panic，降级为默认 logger
		l, _ := zap.NewDevelopment()
		return l
	}
	return global
}
