package zynra

import (
	"log"
)

type Logger interface {
	LogInfo(format string, v ...any)
	LogWarn(format string, v ...any)
	LogError(format string, v ...any)
}

type stdLogger struct{}

func (l *stdLogger) LogInfo(format string, v ...any) {
	log.Printf("[INFO] "+format, v...)
}

func (l *stdLogger) LogWarn(format string, v ...any) {
	log.Printf("[WARN] "+format, v...)
}

func (l *stdLogger) LogError(format string, v ...any) {
	log.Printf("[ERROR] "+format, v...)
}
