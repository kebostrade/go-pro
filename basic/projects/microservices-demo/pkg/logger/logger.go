package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the global logger
func Init(service string, debug bool) error {
	var config zap.Config

	if debug {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	var err error
	log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// Get returns the global logger
func Get() *zap.Logger {
	if log == nil {
		// Fallback to development logger
		log, _ = zap.NewDevelopment()
	}
	return log
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}

