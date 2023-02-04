package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// useful data (name, full_name, default_branch, license, contributions_url)
type Repo struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"owner"`
	HTMLURL                  string      `json:"html_url"`
	Description              string      `json:"description"`
	Fork                     bool        `json:"fork"`
	URL                      string      `json:"url"`
	ForksURL                 string      `json:"forks_url"`
	KeysURL                  string      `json:"keys_url"`
	CollaboratorsURL         string      `json:"collaborators_url"`
	TeamsURL                 string      `json:"teams_url"`
	HooksURL                 string      `json:"hooks_url"`
	IssueEventsURL           string      `json:"issue_events_url"`
	EventsURL                string      `json:"events_url"`
	AssigneesURL             string      `json:"assignees_url"`
	BranchesURL              string      `json:"branches_url"`
	TagsURL                  string      `json:"tags_url"`
	BlobsURL                 string      `json:"blobs_url"`
	GitTagsURL               string      `json:"git_tags_url"`
	GitRefsURL               string      `json:"git_refs_url"`
	TreesURL                 string      `json:"trees_url"`
	StatusesURL              string      `json:"statuses_url"`
	LanguagesURL             string      `json:"languages_url"`
	StargazersURL            string      `json:"stargazers_url"`
	ContributorsURL          string      `json:"contributors_url"`
	SubscribersURL           string      `json:"subscribers_url"`
	SubscriptionURL          string      `json:"subscription_url"`
	CommitsURL               string      `json:"commits_url"`
	GitCommitsURL            string      `json:"git_commits_url"`
	CommentsURL              string      `json:"comments_url"`
	IssueCommentURL          string      `json:"issue_comment_url"`
	ContentsURL              string      `json:"contents_url"`
	CompareURL               string      `json:"compare_url"`
	MergesURL                string      `json:"merges_url"`
	ArchiveURL               string      `json:"archive_url"`
	DownloadsURL             string      `json:"downloads_url"`
	IssuesURL                string      `json:"issues_url"`
	PullsURL                 string      `json:"pulls_url"`
	MilestonesURL            string      `json:"milestones_url"`
	NotificationsURL         string      `json:"notifications_url"`
	LabelsURL                string      `json:"labels_url"`
	ReleasesURL              string      `json:"releases_url"`
	DeploymentsURL           string      `json:"deployments_url"`
	CreatedAt                time.Time   `json:"created_at"`
	UpdatedAt                time.Time   `json:"updated_at"`
	PushedAt                 time.Time   `json:"pushed_at"`
	GitURL                   string      `json:"git_url"`
	SSHURL                   string      `json:"ssh_url"`
	CloneURL                 string      `json:"clone_url"`
	SvnURL                   string      `json:"svn_url"`
	Homepage                 string      `json:"homepage"`
	Size                     int         `json:"size"`
	StargazersCount          int         `json:"stargazers_count"`
	WatchersCount            int         `json:"watchers_count"`
	Language                 string      `json:"language"`
	HasIssues                bool        `json:"has_issues"`
	HasProjects              bool        `json:"has_projects"`
	HasDownloads             bool        `json:"has_downloads"`
	HasWiki                  bool        `json:"has_wiki"`
	HasPages                 bool        `json:"has_pages"`
	HasDiscussions           bool        `json:"has_discussions"`
	ForksCount               int         `json:"forks_count"`
	MirrorURL                interface{} `json:"mirror_url"`
	Archived                 bool        `json:"archived"`
	Disabled                 bool        `json:"disabled"`
	OpenIssuesCount          int         `json:"open_issues_count"`
	License                  interface{} `json:"license"`
	AllowForking             bool        `json:"allow_forking"`
	IsTemplate               bool        `json:"is_template"`
	WebCommitSignoffRequired bool        `json:"web_commit_signoff_required"`
	Topics                   []string    `json:"topics"`
	Visibility               string      `json:"visibility"`
	Forks                    int         `json:"forks"`
	OpenIssues               int         `json:"open_issues"`
	Watchers                 int         `json:"watchers"`
	DefaultBranch            string      `json:"default_branch"`
	TempCloneToken           interface{} `json:"temp_clone_token"`
	Organization             struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"organization"`
	NetworkCount     int `json:"network_count"`
	SubscribersCount int `json:"subscribers_count"`
}

type Contributor struct {
	Login         string `json:"login"`
	ID            int    `json:"id"`
	NodeID        string `json:"node_id"`
	Contributions int    `json:"contributions"`
}

func getGraphQLData(query, GITHUB_TOKEN string) []byte {
	body := []byte(query)

	req, _ := http.NewRequest(http.MethodPost, "https://api.github.com/graphql", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+GITHUB_TOKEN)
	// req.Header.Add("Accept", "application/json")
	// req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	// fmt.Println(string(responseData))
	if err != nil {
		log.Fatal(err)
	}

	return responseData
}

func getRequest(url, GITHUB_TOKEN string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+GITHUB_TOKEN)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseData
}

func GetRepo(url, GITHUB_TOKEN string) Repo {

	responseData := getRequest("https://api.github.com/repos/"+url, GITHUB_TOKEN)

	var responseObject Repo
	json.Unmarshal(responseData, &responseObject)
	// fmt.Println(responseObject.License == nil)
	// topContributor := getTopContributor(responseObject.ContributorsURL, GITHUB_TOKEN)
	// fmt.Println(topContributor)
	return responseObject
}

func getTopContributor(responseObject []Contributor) Contributor {

	// Return top contributor
	return responseObject[0]
}

func getTotalNumContributions(responseObject []Contributor) int {

	totalNumContributions := 0

	for i := 0; i < len(responseObject); i++ {

		totalNumContributions += responseObject[i].Contributions
	}

	return totalNumContributions

}

// Number of contributions made by the top contributor by the total contributions
func GetContributionRatio(url, TOKEN string) float32 {
	respData := getRequest(url, TOKEN)

	var responseObject []Contributor
	json.Unmarshal(respData, &responseObject)

	top := getTopContributor(responseObject).Contributions
	total := getTotalNumContributions(responseObject)

	return float32(top) / float32(total)
}

// Takes in owner, name and TOKEN and outputs the (closed issues, total issues)
func GetIssuesCount(owner, name, GITHUB_TOKEN string) (int, int) {
	// query := "{\"query\" : \"query{repository(owner: \\\"cloudinary\\\", name: \\\"cloudinary_npm\\\") {total: issues {totalCount} closed:issues(states: CLOSED) {totalCount}}}\"}"
	query := "{\"query\" : \"query{repository(owner: \\\"" + owner + "\\\", name: \\\"" + name + "\\\") {total: issues {totalCount} closed:issues(states: CLOSED) {totalCount}}}\"}"

	respData := (getGraphQLData(query, GITHUB_TOKEN))
	// fmt.Println(string(respData))

	type Issue struct {
		Data struct {
			Repository struct {
				Total struct {
					TotalCount int `json:"totalCount"`
				} `json:"total"`
				Closed struct {
					TotalCount int `json:"totalCount"`
				} `json:"closed"`
			} `json:"repository"`
		} `json:"data"`
	}

	var respObj Issue
	json.Unmarshal(respData, &respObj)

	return respObj.Data.Repository.Closed.TotalCount, respObj.Data.Repository.Total.TotalCount
}

func getReadmeURL(repo Repo) string {
	return "https://raw.githubusercontent.com/" + repo.FullName + "/" + repo.DefaultBranch + "/README.md"
}

func GetRawREADME(repo Repo) string {
	url := getReadmeURL(repo)
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(responseData))
	return string(responseData)
}

func GetLicenseFromREADME(readmeText string) string {

	// parse readme for license, return specific license if found, return empty string if not found

	/* unsure about which specific licenses are compatible with LGPLv2.1 license, would like to go through a list of all
	the compatible licenses if they were known and return the specific license that was found
	*/
	if (strings.Contains(readmeText, "License")) && (strings.Contains(readmeText, "MIT")) {

		return "MIT"
	}

	return ""

}
