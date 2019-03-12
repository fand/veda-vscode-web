package logger

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func InitLogger(outputPath string) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{outputPath}
	l, err := config.Build()
	if err != nil {
		panic(err)
	}
	logger = l.Sugar()
}

func LogFatal(m interface{}) {
	if logger == nil {
		log.Println("WARN: logger is uninitialized")
	} else {
		logger.Fatal(m)
	}
}

func FlushLogger() {
	if logger == nil {
		log.Println("WARN: logger is uninitialized")
	} else {
		logger.Sync()
	}
}
