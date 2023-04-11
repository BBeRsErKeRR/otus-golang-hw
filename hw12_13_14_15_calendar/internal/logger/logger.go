package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level    string   `mapstructure:"level"`
	OutPaths []string `mapstructure:"out_paths"`
	ErrPaths []string `mapstructure:"err_paths"`
}

type Logger interface {
	Info(msg string, args ...zapcore.Field)
	Error(msg string, args ...zapcore.Field)
	Debug(msg string, args ...zapcore.Field)
	Warn(msg string, args ...zapcore.Field)
}

type logger struct {
	log *zap.Logger
}

func (l logger) Info(msg string, args ...zapcore.Field) {
	l.log.Info(msg, args...)
}

func (l logger) Error(msg string, args ...zapcore.Field) {
	l.log.Error(msg, args...)
}

func (l logger) Debug(msg string, args ...zapcore.Field) {
	l.log.Debug(msg, args...)
}

func (l logger) Warn(msg string, args ...zapcore.Field) {
	l.log.Warn(msg, args...)
}

func New(conf *Config) (Logger, error) {
	logLevel, err := zap.ParseAtomicLevel(conf.Level)
	if err != nil {
		return nil, err
	}
	cfg := zap.Config{
		Level:            logLevel,
		Encoding:         "json",
		Development:      false,
		OutputPaths:      conf.OutPaths,
		ErrorOutputPaths: conf.ErrPaths,
		// "initialFields": {"foo": "bar"},
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &logger{
		log: zapLogger,
	}, nil
}
