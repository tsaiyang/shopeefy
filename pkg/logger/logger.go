package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
)

var Logger *zap.Logger

func InitLogger() {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "/var/log/algoshop/app.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	writeSyncer := zapcore.AddSync(lumberjackLogger)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel)
	Logger = zap.New(core, zap.AddCaller())

	for i := range 5000 {
		Logger.Info("this is an info message", zap.String("key"+strconv.Itoa(i), "value"+strconv.Itoa(i)))
	}
}
