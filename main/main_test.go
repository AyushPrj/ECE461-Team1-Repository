package main

import (
	"os"
	"testing"
)

// var coverProfile = flag.String("coverprofile", "", "write coverage profile to `file`")
var total int = 20
var pass int = 0

// Integration Test
func TestMain(t *testing.T) {
	os.Args = []string{"main", "test.txt"}
	main()
}
