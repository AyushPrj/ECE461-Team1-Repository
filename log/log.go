package log

import (
	"log"
	"os"
)

const (
	NONE  = "0"
	INFO  = "1"
	DEBUG = "2"
)

var LOG_LEVEL string
var LOG_FILE string

func init() {
	LOG_LEVEL = os.Getenv("LOG_LEVEL")
	LOG_FILE = os.Getenv("LOG_FILE") + "/logfile.log"
}

func shouldLog(mode string) bool {
	return LOG_LEVEL == mode || (LOG_LEVEL == "2" && mode == "1")
}

func Println(mode string, v ...any) {
	if shouldLog(mode) {
		log.Println(v)
	}
}

func Printf(mode string, format string, v ...any) {
	if shouldLog(mode) {
		log.Printf(format, v...)
	}
}
