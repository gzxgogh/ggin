package logs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/logger"
)

type dbLog struct{}

func LogForDB() logger.Interface {
	newLog := logger.New(
		Log,
		logger.Config{},
	)
	return newLog
}

func (d *dbLog) LogMode(logger.LogLevel) logger.Interface {
	// config.Log.Printf("db logger level: %s", l)
	return d
}

func (d *dbLog) Info(ctx context.Context, s string, is ...interface{}) {
	t := fmt.Sprintf(s, is...)
	Log.Info(strings.TrimSpace(t))
}

func (d *dbLog) Warn(ctx context.Context, s string, is ...interface{}) {
	t := fmt.Sprintf(s, is...)
	Log.Warnf(strings.TrimSpace(t))
}

func (d *dbLog) Error(ctx context.Context, s string, is ...interface{}) {
	t := fmt.Sprintf(s, is...)
	Log.Error(strings.TrimSpace(t))
}

func (d *dbLog) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	r, i := fc()
	s := fmt.Sprintf("begin: %s, %s, %d, %s", begin, r, i, err)
	Log.Trace(s)
}
