package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger() *zap.Logger {
	// Configure Zap logger
	config := zap.NewProductionConfig()
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	config.OutputPaths = []string{logFile.Name()}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Initialize logger
	l, err := config.Build()
	if err != nil {
		panic(err)
	}
	return l
}
