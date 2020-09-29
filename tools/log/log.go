package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
	"context"
)

var logger *zap.Logger

func init() {
	zcfg := zap.NewProductionConfig()
	zcfg.OutputPaths = []string{"stdout"}
	zcfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	zcfg.EncoderConfig.EncodeTime = TimeEncoder
	zcfg.Sampling = nil
	var err error
	if logger, err = zcfg.Build(); err != nil {
		panic(err)
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
}

func GetLogger() *zap.Logger {
	return logger
}

type LogUtil struct {
	ctx context.Context
}

func WithContext(ctx context.Context) *LogUtil {
	return &LogUtil{ctx: ctx}
}

func (l *LogUtil) getEventID() zap.Field {
	var s string
	t := l.ctx.Value("event_id")
	switch t.(type) {
	case string:
		s = t.(string)
	}

	return zap.String("event_id", s)
}
func (l *LogUtil) Debug(msg string, fields ...zap.Field) *LogUtil {
	fields = append(fields, l.getEventID())
	logger.Debug(msg, fields...)
	return l
}