package log

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func InitLog(logFile string, logLevel string) {
	SetLogFile(logFile)
	SetLogLevel(logLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    true,
	})
}

// logWay: such as file or console
func SetLogFile(logFile string) {
	if logFile == "" {
		log.SetOutput(os.Stdout)
	} else {
		f, e := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if nil != e {
			log.Fatalf("Error opening file %s for log: %s", logFile, e)
			os.Exit(0)
		} else {
			log.SetOutput(f)
		}
	}
}

// value: error, warning, info, debug
func SetLogLevel(logLevel string) {
	level := log.WarnLevel // warning

	switch logLevel {
	case "error":
		level = log.ErrorLevel
	case "warn":
		level = log.WarnLevel
	case "info":
		level = log.InfoLevel
	case "debug":
		level = log.DebugLevel
	default:
		level = log.DebugLevel
	}

	log.SetLevel(level)
}
