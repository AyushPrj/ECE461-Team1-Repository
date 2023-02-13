package log

import (
	"os"
	"testing"
)

func TestPrintln(t *testing.T) {
	//test if TestPrintln is called
	// setup
	os.Setenv("LOG_LEVEL", "2")
	os.Setenv("LOG_FILE", "./")

	// tests
	//print it to log file
	Println("2", "this is a debug message")
	Println("1", "this is an info message")
	Println("0", "this is a none message")

}

func TestPrintf(t *testing.T) {
	// setup
	os.Setenv("LOG_LEVEL", "2")
	os.Setenv("LOG_FILE", "./")

	// tests
	Printf("2", "this is a %s message", "debug")
	Printf("1", "this is an %s message", "info")
	Printf("0", "this is a %s message", "none")
}

func TestShouldLog(t *testing.T) {
	// setup
	os.Setenv("LOG_LEVEL", "2")

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
