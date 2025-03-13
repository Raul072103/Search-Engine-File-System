package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger(logPath string) *zap.Logger {
	// Configure Zap logger
	config := zap.NewProductionConfig()
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Configure console encoder
	consoleEncoderConfig := zap.NewDevelopmentEncoderConfig() // Use Development config for stacktraces and color
	consoleEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(config.EncoderConfig), zapcore.AddSync(logFile), config.Level),
		zapcore.NewCore(zapcore.NewConsoleEncoder(consoleEncoderConfig), zapcore.AddSync(os.Stdout), config.Level),
	)

	// Initialize logger
	l := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)) // Add caller and stacktrace
	return l
}
