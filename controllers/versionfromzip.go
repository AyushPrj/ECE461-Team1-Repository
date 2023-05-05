package controllers

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func extractVersionUrlFromZip(encodedZip string) (string, string, bool) {
	// Decode the base64-encoded string
	decoded, err := base64.StdEncoding.DecodeString(encodedZip)
	if err != nil {
		return "", "", false
	}

	// Create a temporary file for the zip contents
	tempFile, err := ioutil.TempFile("", "tempzip-*.zip")
	if err != nil {
		return "", "", false
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write the decoded zip contents to the temporary file
	_, err = tempFile.Write(decoded)
	if err != nil {
		return "", "", false
	}

	// Open the zip file for reading
	reader, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return "", "", false
	}
	defer reader.Close()

	// Search for the package.json file in the zip archive
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "package.json") {
			// Open the file from the zip archive
			zippedFile, err := file.Open()
			if err != nil {
				return "", "", false
			}
			defer zippedFile.Close()

			// Read the contents of the file into memory
			packageJsonBytes, err := ioutil.ReadAll(zippedFile)
			if err != nil {
				return "", "", false
			}

			// Unmarshal the JSON into a struct
			var packageJson map[string]interface{}
			err = json.Unmarshal(packageJsonBytes, &packageJson)
			if err != nil {
				return "", "", false
			}

			// Get the value of the "version" field
			version, _ := packageJson["version"].(string)
			// url, ok := packageJson["repository"].(string)
			// Extract the url from the Json field repository: {url: "
			repository, _ := packageJson["repository"].(map[string]interface{})
			url, ok := repository["url"].(string)
			if ok {
				return url, version, true
			}
		}
	}

	// If the package.json file was not found, return boolean false
	return "", "", false
}

func downloadZip(url string) (string, error) {
	// First convert github url to zip url
	// strip .git from the end of the url
	url = strings.TrimSuffix(url, ".git")
	// replace github.com with codeload.github.com
	url = strings.Replace(url, "github.com", "codeload.github.com", 1)
	url = url + "/zip/master"

	// Download the zip file from the url
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// Read the contents of the zip file into memory
	zipBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Encode the zip file as a base64-encoded string
	encodedZip := base64.StdEncoding.EncodeToString(zipBytes)

	return encodedZip, nil

}