package controllers

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

func extractVersionFromZip(encodedZip string) (string, bool) {
	// Decode the base64-encoded string
	decoded, err := base64.StdEncoding.DecodeString(encodedZip)
	if err != nil {
		return "", false
	}

	// Create a temporary file for the zip contents
	tempFile, err := ioutil.TempFile("", "tempzip-*.zip")
	if err != nil {
		return "", false
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write the decoded zip contents to the temporary file
	_, err = tempFile.Write(decoded)
	if err != nil {
		return "", false
	}

	// Open the zip file for reading
	reader, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return "", false
	}
	defer reader.Close()

	// Search for the package.json file in the zip archive
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, "package.json") {
			// Open the file from the zip archive
			zippedFile, err := file.Open()
			if err != nil {
				return "", false
			}
			defer zippedFile.Close()

			// Read the contents of the file into memory
			packageJsonBytes, err := ioutil.ReadAll(zippedFile)
			if err != nil {
				return "", false
			}

			// Unmarshal the JSON into a struct
			var packageJson map[string]interface{}
			err = json.Unmarshal(packageJsonBytes, &packageJson)
			if err != nil {
				return "", false
			}

			// Get the value of the "version" field
			version, ok := packageJson["version"].(string)
			if ok {
				return version, true
			}
		}
	}

	// If the package.json file was not found, return boolean false
	return "", false
}
