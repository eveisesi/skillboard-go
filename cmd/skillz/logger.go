package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func buildLogger() {

	now := time.Now().Format("2006-01-02-15")

	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Panicf(fmt.Sprintf("failed to configure log level: %s", err))
	}

	logger = logrus.New()

	logger.SetOutput(ioutil.Discard)

	logger.AddHook(&writerHook{
		Writer:    os.Stdout,
		LogLevels: logrus.AllLevels,
	})

	logger.AddHook(&writerHook{
		Writer: &lumberjack.Logger{
			Filename: fmt.Sprintf("logs/%s-error.log", now),
			MaxSize:  10,
			Compress: false,
		},
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	logger.AddHook(&writerHook{
		Writer: &lumberjack.Logger{
			Filename:   fmt.Sprintf("logs/%s-info.log", now),
			MaxBackups: 3,
			MaxSize:    10,
			Compress:   false,
		},
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
		},
	})

	logger.SetLevel(level)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableQuote: true,
	})

}

type writerHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

func (w *writerHook) Fire(entry *logrus.Entry) error {

	data := entry.Data

	message := entry.Message

	var service string
	var ok bool
	if service, ok = data["service"].(string); ok {
		message = fmt.Sprintf("[%s] %s", strings.ToLower(service), message)
		delete(data, "service")
		entry.Data = data
		entry.Message = message
	}

	line, err := entry.Bytes()
	if err != nil {
		return err
	}

	_, err = w.Writer.Write(line)
	return err

}

func (w *writerHook) Levels() []logrus.Level {
	return w.LogLevels
}
