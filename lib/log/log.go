package log

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Init(level string, filePath string) error {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		return err
	}

	log.SetLevel(lvl)
	log.SetReportCaller(true)
	if filePath != "" {
		log.SetOutput(&lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    20, // megabytes
			MaxBackups: 10,
			MaxAge:     30,   //days
			Compress:   true, // disabled by default
		})

	}

	log.WithField("LogLevel", level).Warn("log inited")
	return nil
}
