package log

import (
	"log/syslog"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	logrusSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

type Config struct {
	Level  string
	Syslog Syslog
}

type Syslog struct {
	Enabled bool
	Network string
	Address string
	Tag     string
}

func InitLogrus(cfg Config) {
	logLevel, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logLevel = logrus.ErrorLevel
	}

	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)

	if cfg.Syslog.Enabled {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})

		priority := convertLevelToSyslogPriority(cfg.Level)

		hook, err := logrusSyslog.NewSyslogHook(
			cfg.Syslog.Network,
			cfg.Syslog.Address,
			priority,
			cfg.Syslog.Tag,
		)
		if err != nil {
			logrus.Errorf("Unable to connect to syslog daemon. error: %s", err.Error())
		} else {
			logrus.AddHook(hook)
		}
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:            true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
			TimestampFormat:        time.RFC3339,
			FullTimestamp:          true,
		})
	}
}

func convertLevelToSyslogPriority(level string) syslog.Priority {
	level = strings.ToLower(level)
	switch level {
	case "trace", "debug":
		return syslog.LOG_DEBUG
	case "info":
		return syslog.LOG_INFO
	case "notice":
		return syslog.LOG_NOTICE
	case "warn", "warning":
		return syslog.LOG_WARNING
	case "error":
		return syslog.LOG_ERR
	case "fatal":
		return syslog.LOG_CRIT
	case "alert":
		return syslog.LOG_ALERT
	case "panic":
		return syslog.LOG_EMERG
	default:
		return syslog.LOG_ERR
	}
}
