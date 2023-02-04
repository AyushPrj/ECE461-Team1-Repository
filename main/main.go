package main

import (
	"fmt"
	"log"
	"os"
	"proj/metrics"

	"github.com/joho/godotenv"
)

var GITHUB_TOKEN string
var LOG_LEVEL string
var LOG_FILE string

func init() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")
	LOG_LEVEL = os.Getenv("LOG_LEVEL")
	LOG_FILE = os.Getenv("LOG_FILE")

}

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		fmt.Println(arg)
	}

	// fmt.Println(GITHUB_TOKEN, LOG_LEVEL, LOG_FILE)

	// GETS ALL THE METRICS IN THIS FUNCTION GIVEN THE URL (ONLY WORKS FOR GITHUB CURRENTLY)
	metrics.GetMetrics("cloudinary/cloudinary_npm", GITHUB_TOKEN)

	// api.GetRepo("cloudinary/cloudinary_npm", GITHUB_TOKEN)
	// fmt.Println(api.GetRawREADME(test))
	// fmt.Println(test.License != nil)
	// api.GetIssuesCount("", "", GITHUB_TOKEN)

	// test := api.GetNPMData("nodist").License
	// fmt.Println(test)

}
