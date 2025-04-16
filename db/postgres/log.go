package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Logger struct {
	slowThreshold         time.Duration
	skipErrRecordNotFound bool
}

func NewGormLogger(slowThreshold time.Duration) *Logger {
	return &Logger{
		slowThreshold:         slowThreshold,
		skipErrRecordNotFound: true,
	}
}

func (l *Logger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, s string, data ...any) {
	logrus.WithContext(ctx).WithField("type", "sql").Infof(s, data...)
}

func (l *Logger) Warn(ctx context.Context, s string, data ...any) {
	logrus.WithContext(ctx).WithField("type", "sql").Warnf(s, data...)
}

func (l *Logger) Error(ctx context.Context, s string, data ...any) {
	logrus.WithContext(ctx).WithField("type", "sql").Errorf(s, data...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := logrus.Fields{}
	fields["type"] = "sql"
	fields["caller"] = utils.FileWithLineNum()
	fields["elapsed"] = elapsed.Seconds()
	fields["rows"] = rows

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.skipErrRecordNotFound) {
		fields[logrus.ErrorKey] = err
		logrus.WithContext(ctx).WithFields(fields).Errorf("[RECORD NOT FOUND SQL] %s", sql)

		return
	}

	if l.slowThreshold != 0 && elapsed > l.slowThreshold {
		logrus.WithContext(ctx).WithFields(fields).Warnf("[SLOW SQL] %s", sql)

		return
	}

	logrus.WithContext(ctx).WithFields(fields).Debugf("[SQL] %s", sql)
}
