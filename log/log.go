package log

import (
	"fmt"
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
	LOG_FILE = os.Getenv("LOG_FILE")
	if(LOG_FILE == "") {
		LOG_FILE = "logfile.log"
	}
	if(LOG_LEVEL == "") {
		LOG_LEVEL = "2"
	}
	file, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

    log.SetOutput(file)
    log.SetOutput(os.Stdout)
}

func shouldLog(mode string) bool {
	return LOG_LEVEL == mode || (LOG_LEVEL == "2" && mode == "1")
}

func Println(mode string, v ...any) {
	if shouldLog(mode) {
		if(mode == "2") {
			fmt.Println(v...)
		}
		// fmt.Printf("should write to %s with mode %s and level %s: ", LOG_FILE, mode, LOG_LEVEL)
		log.Println(v...)
	}
}

func Printf(mode string, format string, v ...any) {
	if shouldLog(mode) {
		if(mode == "2") {
			fmt.Printf(format, v...)
		}
		log.Printf(format, v...)
	}
}
