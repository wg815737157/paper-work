package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var (
	defaultLogger *zap.Logger
	sugarLogger   *zap.SugaredLogger

	customTimeFormat string
)

func Logger() *zap.Logger {
	return defaultLogger
}

func SugarLogger() *zap.SugaredLogger {
	return sugarLogger
}

func init() {
	defaultLogger = initDefaultLogger()
	sugarLogger = defaultLogger.Sugar()
}

// New initialize logger by input options
func initDefaultLogger() *zap.Logger {
	debug := -1
	timeFormat := "2006-01-02T15:04:05Z07:00"

	level := zapcore.Level(debug)
	// high priority
	hp := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapcore.ErrorLevel
	})
	// low priority
	lp := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= level && l < zapcore.ErrorLevel
	})
	cInfo := zapcore.Lock(os.Stdout)
	cError := zapcore.Lock(os.Stderr)

	useCustomTimeFormat := false
	cfg := zap.NewProductionEncoderConfig()
	useCustomTimeFormat = true
	customTimeFormat = timeFormat
	cfg.EncodeTime = customTimeEncoder
	cEncoder := zapcore.NewJSONEncoder(cfg)
	core := zapcore.NewTee(
		zapcore.NewCore(cEncoder, cError, hp),
		zapcore.NewCore(cEncoder, cInfo, lp),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel))
	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)
	if !useCustomTimeFormat {
		logger.Warn("time format for logger is not provided, use zap default")
	}
	defer logger.Sync()
	return logger
}

// customTimeEncoder encoder for time format
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(customTimeFormat))
}
