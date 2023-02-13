package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type NPMData struct {
	ID          string `json:"_id,omitempty"`
	Rev         string `json:"_rev,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	DistTags    struct {
		Latest string `json:"latest,omitempty"`
	} `json:"dist-tags,omitempty"`
	Readme      string `json:"readme,omitempty"`
	Maintainers []struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"maintainers,omitempty"`

	Author struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
	} `json:"author,omitempty"`
	Repository struct {
		Type string `json:"type,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"repository,omitempty"`
	Homepage string `json:"homepage,omitempty"`
	Bugs     struct {
		URL string `json:"url,omitempty"`
	} `json:"bugs,omitempty"`
	ReadmeFilename string `json:"readmeFilename,omitempty"`
	License        string `json:"license,omitempty"`
}

func getNPMData(pkgName string) NPMData {
	url := "https://registry.npmjs.org/" + pkgName
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var responseObject NPMData
	json.Unmarshal(responseData, &responseObject)
	return responseObject
}

func GetGithubURL(pkgName string) string {
	data := getNPMData(pkgName)
	return data.Repository.URL
}
