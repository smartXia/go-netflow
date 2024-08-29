package log

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const dateFormat = "2006-01-02 15:04:05"

var (
	logger *zap.Logger
	once   sync.Once
)

func UnixTimestampEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(dateFormat))
}

func CreateLogger(logLevel, filepath string) *zap.Logger {
	var w zapcore.WriteSyncer
	if filepath == "" {
		w = zapcore.AddSync(os.Stdout)
	} else {
		w = zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath,
			MaxSize:    50, // megabytes
			MaxBackups: 30,
			MaxAge:     1, // days
			Compress:   true,
		})
	}

	var level zapcore.Level
	level.Set(logLevel)

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = UnixTimestampEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(cfg), w, level)
	return zap.New(core, zap.AddCaller())
}

// Setup global logger
func Setup(logLevel, filepath string) {
	once.Do(func() {
		logger = CreateLogger(logLevel, filepath)
	})
}

func GetLogger() *zap.Logger {
	return logger
}
