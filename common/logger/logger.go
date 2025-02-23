package logger

import (
	"log"
	"os"

	"github.com/rs/zerolog"
)

//Log ...
var Log zerolog.Logger

//InitLogger ...
func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel := zerolog.InfoLevel // Default log level
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		level, err := zerolog.ParseLevel(levelStr)
		if err != nil {
			log.Println("Invalid LOG_LEVEL environment variable, using default")
		} else {
			logLevel = level
		}
	}
	w := zerolog.ConsoleWriter{Out: os.Stdout}

	Log = zerolog.New(w).Level(logLevel).With().Timestamp().Logger()
}
