package api

import (
	"testing"
)

func TestGetNpmData(t *testing.T) {
	result := getNPMData("browserify")
	if result.Name != "browserify" {
		t.Fatal("Error getting data from npm url!")
	}
}

func TestGetGithubUrl(t *testing.T) {
	tstObj := getNPMData("browserify")
	if tstObj.Repository.URL != "git+ssh://git@github.com/browserify/browserify.git" {
		t.Fatal("Error getting url from npm Data Structure!")
	}
}
