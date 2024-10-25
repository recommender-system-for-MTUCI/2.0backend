package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func customColorEncoder() zapcore.Encoder {
	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeLevel = zapcore.LevelEncoder(func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch level {
		case zapcore.DebugLevel:
			enc.AppendString(Cyan + level.CapitalString() + Reset)
		case zapcore.InfoLevel:
			enc.AppendString(Green + level.CapitalString() + Reset)
		case zapcore.WarnLevel:
			enc.AppendString(Yellow + level.CapitalString() + Reset)
		case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
			enc.AppendString(Red + level.CapitalString() + Reset)
		default:
			enc.AppendString(level.CapitalString())
		}
	})

	return zapcore.NewConsoleEncoder(config)
}
func LoggerConfig() zapcore.Core {
	core := zapcore.NewCore(
		customColorEncoder(),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)
	return core
}
