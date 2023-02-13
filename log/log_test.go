package log

import (
	golog "log"
	"os"
	"testing"
)

func TestShouldLog(t *testing.T) {
	// setup
	LOG_LEVEL = "2"

	// tests
	if !shouldLog("2") {
		t.Error("expected shouldLog(\"2\") to return true")
	}
	if !shouldLog("1") {
		t.Error("expected shouldLog(\"1\") to return true")
	}
	if shouldLog("0") {
		t.Error("expected shouldLog(\"0\") to return false")
	}

}

func TestNoLog(t *testing.T) {
	LOG_LEVEL = "1"

	f, _ := os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	defer f.Close()

	golog.SetOutput(f)

	Println(INFO, "Running Test Case (println)")
	Printf(INFO, "Running Test Case (printf)")

}
