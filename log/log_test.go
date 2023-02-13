package log

import (
	"os"
	"testing"
)


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
