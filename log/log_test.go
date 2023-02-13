package log

import (
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

	LOG_LEVEL = "0"
}

func TestNoLog(t *testing.T) {
	Println(NONE, "nothing should print")
	Printf(NONE, "nothing should print")
}
