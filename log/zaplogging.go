package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type zapLogger struct {
	logger *zap.Logger
}

func (z *zapLogger) Info(args ...interface{}) {
	z.logger.Info(fmt.Sprint(args...))
}

func (z *zapLogger) InfoF(s string, args ...interface{}) {
	z.logger.Info(fmt.Sprintf(s, args...))
}

func (z *zapLogger) Error(args ...interface{}) {
	z.logger.Error(fmt.Sprint(args...))
}

func (z *zapLogger) ErrorF(s string, args ...interface{}) {
	z.logger.Error(fmt.Sprintf(s, args...))
}

func (z *zapLogger) Warn(args ...interface{}) {
	z.logger.Warn(fmt.Sprint(args...))
}

func (z *zapLogger) WarnF(s string, args ...interface{}) {
	z.logger.Warn(fmt.Sprintf(s, args...))
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func NewZapLogger(config Configuration) (LoggerInterface, error) {
	// TODO 记录日志到文件

	// 设置日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.LogLevel)); err != nil {
		return nil, err
	}

	cfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), zapcore.NewMultiWriteSyncer(os.Stderr), level)
	logger := zap.New(core)
	return &zapLogger{logger: logger}, nil
}
