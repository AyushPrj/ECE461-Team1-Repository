package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	//"github.com/joho/godotenv"
)

/*func init() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}*/

func build() {
	fmt.Println("Building...")
}

func install() {
	fmt.Println("Installing...")
}

func main() {
	args := os.Args[1:]
	str := strings.Join(args, "")
	//var file *os.File

	if str == "build" {
		build()
	} else if str == "install" {
		install()
	} else {
		file, err := os.Open(str)
		if err != nil {
			log.Fatalf("Failed to open file!")
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var text []string

		for scanner.Scan() {
			text = append(text, scanner.Text())
		}

		// The method os.File.Close() is called
		// on the os.File object to close the file
		file.Close()

		// and then a loop iterates through
		// and prints each of the slice values.
		for _, each_ln := range text {
			fmt.Println(each_ln)
		}
	}

	// githubToken := os.Getenv("GITHUB_TOKEN")
	// logLevel := os.Getenv("LOG_LEVEL")
	// logFile := os.Getenv("LOG_FILE")
	// fmt.Println(githubToken, logLevel, logFile)
}
