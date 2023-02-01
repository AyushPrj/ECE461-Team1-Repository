package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		fmt.Println(arg)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	github_token := os.Getenv("GITHUB_TOKEN")
	log_level := os.Getenv("LOG_LEVEL")
	log_file := os.Getenv("LOG_FILE")
	// fmt.Println(github_token, log_level, log_file)

}
