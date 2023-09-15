package logger

import (
	"os"
	"time"

	"github.com/isfk/xiaohui/config"
	"github.com/isfk/xiaohui/internal/pkg/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLogger(server string, logConf config.LogConf) *zap.Logger {
	var cores []zapcore.Core
	fields := []zap.Field{{
		Key:    "version",
		Type:   zapcore.StringType,
		String: global.Version,
	}, {
		Key:    "server",
		Type:   zapcore.StringType,
		String: server,
	}}

	// default
	defaultConfig := zap.NewProductionEncoderConfig()
	defaultConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	coreDefault := zapcore.NewCore(zapcore.NewJSONEncoder(defaultConfig), zapcore.AddSync(&lumberjack.Logger{
		Filename:  logConf.DefaultPath,
		MaxSize:   logConf.MaxSize, // MB
		LocalTime: true,
		Compress:  true,
	}), zap.NewAtomicLevelAt(zap.DebugLevel)).With(fields)

	// error+
	errConfig := zap.NewProductionEncoderConfig()
	errConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	coreError := zapcore.NewCore(zapcore.NewJSONEncoder(errConfig), zapcore.AddSync(&lumberjack.Logger{
		Filename:  logConf.ErrorPath,
		MaxSize:   logConf.MaxSize, // MB
		LocalTime: true,
		Compress:  true,
	}), zap.NewAtomicLevelAt(zap.ErrorLevel)).With(fields)

	// console
	consoleConfig := zap.NewProductionEncoderConfig()
	consoleConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	coreConsole := zapcore.NewCore(zapcore.NewConsoleEncoder(consoleConfig), zapcore.AddSync(os.Stdout), zap.NewAtomicLevelAt(zap.DebugLevel)).With(fields)

	log := zap.New(zapcore.NewTee(append(cores, coreConsole, coreDefault, coreError)...), zap.AddCaller())
	defer func(log *zap.Logger) {
		_ = log.Sync()
	}(log)
	return log
}
