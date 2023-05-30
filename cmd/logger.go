package cmd

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

var log = logrus.New()

// 初始化日志
func initLogger() {
	// 输出Json
	useJSON := loggerConfig.UseJSON
	if useJSON {
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	// 日志级别
	level := loggerConfig.Level
	switch level {
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	}

	log.WithFields(logrus.Fields{
		"app": "chainstorage-cli",
	})

	// default
	log.SetReportCaller(true)

	// 输出日志文件
	isOutPutFile := loggerConfig.IsOutPutFile
	if isOutPutFile {
		maxAgeDay := time.Duration(24 * loggerConfig.MaxAgeDay)
		rotationTime := time.Duration(24 * loggerConfig.RotationTime)

		logPath := loggerConfig.LogPath
		loggerFile := loggerConfig.LoggerFile
		linkName := path.Join(logPath, loggerFile)
		logPattern := linkName + ".%Y%m%d"
		logFile, err := rotatelogs.New(
			logPattern,
			rotatelogs.WithLinkName(linkName),
			rotatelogs.WithMaxAge(maxAgeDay*time.Hour),
			rotatelogs.WithRotationTime(rotationTime*time.Hour),
		)
		if err != nil {
			log.Printf("failed to create rotatelogs: %s", err)
			return
		}

		log.SetOutput(logFile)
	}
}
