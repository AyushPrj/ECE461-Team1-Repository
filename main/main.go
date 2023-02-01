package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		fmt.Println(arg)
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	logLevel := os.Getenv("LOG_LEVEL")
	logFile := os.Getenv("LOG_FILE")
	// fmt.Println(githubToken, logLevel, logFile)
}
