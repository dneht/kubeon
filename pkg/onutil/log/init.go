package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var log *zap.SugaredLogger
var logLevel = zapcore.InfoLevel

func Init(level string) {
	 zapLevel := &logLevel
	if "" != level {
		err :=  logLevel.Set(level)
		if nil != err {
			os.Exit(1)
		}
	}
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     customTimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   noneCallerEncoder,
		}),
		zapcore.Lock(os.Stdout), zap.NewAtomicLevelAt(*zapLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	log = logger.Sugar()
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05"))
}

func noneCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
}
