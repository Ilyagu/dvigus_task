package logger

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Path                string
	Name                string
	MaxSize             int
	MaxBackups          int
	MaxAge              int
	RotateCheckInterval time.Duration
}

type Logger struct {
	*logrus.Logger
	conf               Config
	currentLogFileName string
}

func NewLogger(conf Config) *Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	currentDate := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("%s%s-%s.log", conf.Path, conf.Name, currentDate)

	logger.SetOutput(&lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
	})

	return &Logger{logger, conf, logFileName}
}

func (l *Logger) RotateLogByDate() {
	currentDate := time.Now().Format("2006-01-02")
	newLogFileName := fmt.Sprintf("%s%s-%s.log", l.conf.Path, l.conf.Name, currentDate)
	if newLogFileName != l.currentLogFileName {
		l.SetOutput(&lumberjack.Logger{
			Filename:   newLogFileName,
			MaxSize:    l.conf.MaxSize,
			MaxBackups: l.conf.MaxBackups,
			MaxAge:     l.conf.MaxAge,
		})
		l.currentLogFileName = newLogFileName
	}
}
