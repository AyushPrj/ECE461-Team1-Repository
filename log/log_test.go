package log

import (
	"testing"
	"os"
)

func TestShouldPrint(t *testing.T) {
	lvl := os.Getenv("LOG_LEVEL")
	if(lvl == 0 && shouldLog())
}