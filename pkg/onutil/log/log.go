package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Level() int {
	switch logLevel {
	case zapcore.DebugLevel:
		return 4
	case zapcore.InfoLevel:
		return 1
	default:
		return 0
	}
}

func IsDebug() bool {
	return logLevel <= zap.DebugLevel
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	log.Debugf(template, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(template string, args ...interface{}) {
	log.Infof(template, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	log.Warnf(template, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	log.Errorf(template, args...)
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	log.Panicf(template, args...)
}
