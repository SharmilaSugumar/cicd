package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if os.Getenv("APP_ENV") == "development" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	Log = logger
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

func WithContext(ctx context.Context) *zap.Logger {
	l := Log
	if reqID, ok := ctx.Value("RequestID").(string); ok {
		l = l.With(zap.String("request_id", reqID))
	}
	if userID, ok := ctx.Value("userID").(string); ok {
		l = l.With(zap.String("user_id", userID))
	}
	// Note: Gin Context passes values differently (c.Get vs c.Value).
	// We map Gin keys to standard context inside middleware to use this universally.
	return l
}
