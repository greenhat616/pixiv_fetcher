package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// initLog 用于初始化日记模块
func initLog() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.999Z",
	})
	if Debug {
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetOutput(os.Stderr)
		log.SetLevel(log.ErrorLevel)
	}
}
