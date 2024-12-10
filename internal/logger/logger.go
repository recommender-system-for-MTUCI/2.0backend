package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func jsonEncoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder // Формат времени: ISO 8601
	return zapcore.NewJSONEncoder(config)
}

func LoggerConfig() zapcore.Core {
	encoder := jsonEncoder()

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)
	return core
}
