package main

import (
	"fmt"
	"log"
	"os"
	"proj/api"

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

	fmt.Println(GITHUB_TOKEN, LOG_LEVEL, LOG_FILE)
	test := api.GetRepo("cloudinary/cloudinary_npm")
	fmt.Println(test.License != nil)

}
