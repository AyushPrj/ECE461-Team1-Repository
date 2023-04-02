package api

import (
	"ECE461-Team1-Repository/log"
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"strconv"
	"time"

	"fmt"
)

var GITHUB_TOKEN string

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

func init() {
	GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")
}

func getGraphQLData(query string) []byte {
	body := []byte(query)

	req, _ := http.NewRequest(http.MethodPost, "https://api.github.com/graphql", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+GITHUB_TOKEN)
	req.Header.Set("Accept", "application/vnd.github.hawkgirl-preview+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(log.DEBUG, err)
	}

	return responseData
}

func getRequest(url string) []byte {
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
		log.Println(log.DEBUG, err)
	}
	return responseData
}

func GetRepo(url string) Repo {
	responseData := getRequest("https://api.github.com/repos/" + url)

	var responseObject Repo
	json.Unmarshal(responseData, &responseObject)
	return responseObject
}

func getTopContributor(responseObject []Contributor) Contributor {
	// Return top contributor
	return responseObject[0]
}

/*
getTotalNumContributions is a simple function that iterates through a list of Contributors
and sums each Contributor's number of contributions.
*/

func getTotalNumContributions(responseObject []Contributor) int {

	totalNumContributions := 0

	for i := 0; i < len(responseObject); i++ {
		totalNumContributions += responseObject[i].Contributions
	}

	return totalNumContributions

}

// Number of contributions made by the top contributor by the total contributions
func GetContributionRatio(url string) float32 {
	respData := getRequest(url)

	var responseObject []Contributor
	json.Unmarshal(respData, &responseObject)

	top := getTopContributor(responseObject).Contributions
	total := getTotalNumContributions(responseObject)
	if float32(total) != 0.0 {
		return float32(top) / float32(total)
	} else {
		return 0
	}
}

// Takes in owner, name and TOKEN and outputs the (closed issues, total issues)
func GetIssuesCount(owner, name string) (int, int) {
	query := "{\"query\" : \"query{repository(owner: \\\"" + owner + "\\\", name: \\\"" + name + "\\\") {total: issues {totalCount} closed:issues(states: CLOSED) {totalCount}}}\"}"
	respData := (getGraphQLData(query))

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
	responseData := getRequest("https://api.github.com/repos/" + repo.FullName + "/readme")

	type ReadmeData struct {
		Name        string `json:"name"`
		DownloadURL string `json:"download_url"`
	}

	var responseObject ReadmeData
	json.Unmarshal(responseData, &responseObject)

	return responseObject.DownloadURL
	// return "https://raw.githubusercontent.com/" + repo.FullName + "/" + repo.DefaultBranch + "/README.md"
}

func GetRawREADME(repo Repo) string {
	url := getReadmeURL(repo)
	response, err := http.Get(url)
	if err != nil {
		log.Println(log.DEBUG, err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(log.DEBUG, err)
	}

	return string(responseData)
}

/*
GetLicenseFromREADME takes in the raw contents of a README.md file in a string variable. The
README is checked for one of many licenses compatible with the LGPLv2.1 license. If the function
finds a specific compatible license, it returns that license, otherwise it returns an empty string.
*/

func GetLicenseFromREADME(readmeText string) string {

	// parse readme for license, return specific license if found, return empty string if not found

	licenses := []string{
		"MIT", "Apache", "BSD 3-Clause",
		"BSD 2-Clause", "ISC", "BSD Zero Clause",
		"Boost Software", "UPL", "Universal Permissive",
		"JSON", "Simple Public", "Copyfree Open Innovation",
		"Xerox", "Sendmail", "All-Permissive", "Artistic",
		"Berkely Database", "Modified BSD", "CeCILL", "Cryptix General",
		"Zope Public", "XFree86", "X11", "WxWidgets Library", "WTFPL",
		"WebM", "Unlicense", "StandardMLofNJ", "Ruby", "SGI Free Software",
		"Python", "Ruby", "Perl", "OpenLDAP", "Netscape Javascript", "NCSA",
		"Mozilla Public", "Intel Open Source"}

	if strings.Contains(readmeText, "License") || strings.Contains(readmeText, "license") {

		for _, license := range licenses {

			if strings.Contains(readmeText, license) {
				return license
			}
		}
	}

	if strings.Contains(readmeText, "IBM PowerPC Initialization and Boot Software") || strings.Contains(readmeText, "IBM-pibs") {
		return "IBM-pibs"
	}

	return ""

}

/*
RunClocOnRepo makes use of the cloc bash command readily available on eceprog.
A repository is cloned into the directory, cloc is run on that directory, and
the cloc output is stored inside a string variable. The repository is cleaned up
in CheckRepoForTest, which also makes use of the cloned repository.
*/

func RunClocOnRepo(repo Repo) string {
	cloneString := repo.CloneURL

	// Get current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Println(log.DEBUG, "Error:", err)
	}

	// Navigate to the main folder
	if err := os.Chdir(dir); err != nil {
		log.Println("Error navigating to main folder:", err)
	}

	// Clone repo
	clone := exec.Command("git", "clone", cloneString)
	err = clone.Run()
	if err != nil {
		log.Println(log.DEBUG, err)
	}

	folderName := repo.Name + "/"
	cloc := exec.Command("cloc", folderName)
	out, err := cloc.CombinedOutput()

	if err != nil {
		log.Println(log.DEBUG, err)
	}

	stringOut := string(out)
	os.Chdir(dir)
	log.Println(log.DEBUG, stringOut)

	return stringOut
}

/*
CheckRepoForTest works off of the repository cloned in RunClocOnRepo.
os.ReadDir returns a list of DirEntry objects of a folder, which lists
a few attributes of every file/folder in the given folder name, including
its name. We can use the names to check if a test suite/test folder exists
in the repository, assigning a score of 1 or 0 based on if it does or does
not exist. The cloned repository is also cleaned up in this function.
*/

func CheckRepoForTest(repo Repo) float64 {
	testFound := 0.0

	// Get current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Println(log.DEBUG, "Error:", err)
	}

	// Navigate to the main folder
	if err := os.Chdir(dir); err != nil {
		log.Println("Error navigating to main folder:", err)
	}

	// Go to repo folder
	temp, err := os.ReadDir(repo.Name)

	if err != nil {
		log.Println(log.DEBUG, err)
	}

	for _, val := range temp {

		currentFile := val.Name()

		if currentFile == "test" {
			testFound = 1.0
			break
		}
	}

	return testFound
}

/*
GetDepPinRate queries a GraphQL API to retrieve information about a repository's dependencies.
It then calculates the percentage of pinned dependencies (dependencies with version numbers explicitly specified)
out of the total number of dependencies for the repository.
*/

func GetDepPinRate(owner, name string) float32 {
	query := "{\"query\":\"{repository(owner:\\\"" + owner + "\\\", name:\\\"" + name + "\\\") { dependencyGraphManifests { totalCount, edges{ node { dependencies { totalCount , nodes { packageName, requirements, hasDependencies}}}}}}}\"}"
	respData := (getGraphQLData(query))

	type DependencyGraph struct {
		Data struct {
			Repository struct {
				DependencyGraphManifests struct {
					TotalCount int `json:"totalCount"`
					Edges []struct {
						Node struct {
							Dependencies struct {
								TotalCount int `json:"totalCount"`
								Nodes []struct {
									PackageName      string `json:"packageName"`
									Requirements     string `json:"requirements"`
									HasDependencies  bool   `json:"hasDependencies"`
								} `json:"nodes"`
							} `json:"dependencies"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"dependencyGraphManifests"`
			} `json:"repository"`
		} `json:"data"`
	}

	var respObj DependencyGraph
	fmt.Println(GetPackageRequirements(owner, name))

	if err := json.Unmarshal(respData, &respObj); err != nil {
		log.Println(log.DEBUG, err)
		return GetPackageRequirements(owner, name)
	}

	dgm := respObj.Data.Repository.DependencyGraphManifests
	if dgm.TotalCount == 0 {
		return GetPackageRequirements(owner, name)
	}

	pinnedReq := 0
	totDep := 0
	versionRegex := regexp.MustCompile(`\d+\.\d+`)
	for _, edge := range dgm.Edges {
		for _, dep := range edge.Node.Dependencies.Nodes {
			totDep++;
			if versionRegex.MatchString(dep.Requirements) {
				pinnedReq += 1
			}
		}
	}

	return float32(pinnedReq) / float32(totDep)
}

/*
GetPackageRequirements is a backup function for DepPinRate. If the funciton fails to find
a Dependency-Graph for a repo. This function will be called. It searches for a requirements.txt
and/or a package.json to find dependencies and determine if they are pinned.
*/

func GetPackageRequirements(owner, name string) float32 {

	fileName := ""
	numDependencies := 0
	numPinned := 0

	// Get current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Println(log.DEBUG, "Error:", err)
	}

	// Navigate to the main folder
	if err := os.Chdir(dir); err != nil {
		log.Println(log.DEBUG, "Error navigating to main folder:", err)
	}

	// Go to repo folder
	temp, _ := os.ReadDir(name)
	for _, val := range temp {

		currentFile := strings.ToLower(val.Name())
		if (currentFile == "requirements.txt" || currentFile == "package.json") {// Add more if more are found
			fileName = val.Name() 
		}
	}

	if fileName == "" { return 0 }

	pattern := regexp.MustCompile(`[=><]\d+\.\d+`)

	file, err := os.Open(name + "/" + fileName)
    if err != nil {
        log.Println(log.DEBUG, "Error opening file:", err)
        return 0
    }
    defer file.Close()

	scanner := bufio.NewScanner(file)
	if strings.ToLower(fileName) == "requirements.txt" {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line, "#") {
				numDependencies += 1
				if pattern.MatchString(line) {
					numPinned += 1
				}
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println(log.DEBUG, "Error scanning file:", err, numDependencies)
		}
	}
	if strings.ToLower(fileName) == "package.json" {
		// TODO: NEED TO IMPLEMENT
	}

	return float32(numPinned) / float32(numDependencies)
}

/*
CountReviewedLines counts the amount of lines that were merged into the repo's main via
pull request. For each identified pull request, it uses git diff to get the number of lines
added and deleted in the pull request. The function then adds up the lines added and
deleted to get the total number of reviewed lines.
*/

func CountReviewedLines(repo Repo) int {
	// Get current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Println(log.DEBUG, "Error:", err)
	}

	// Navigate to the main folder
	if err := os.Chdir(dir); err != nil {
		log.Println("Error navigating to main folder:", err)
	}

	// Go to repo folder
	err = os.Chdir(repo.Name)
	if err != nil {
		log.Println(log.DEBUG, err)
	}

	cmd := exec.Command("git", "log", "--merges", "--pretty=format:'%h %s'")
	out, err := cmd.Output()
	if err != nil {
		log.Println(log.DEBUG, err)
	}

	var totLinesReviewed int = 0
	commits := strings.Split(string(out), "\n")

	for _, commit := range commits {
		if (len(commit) > 1) {
			parts := strings.Split(commit, " ")
			hash := parts[0][1 : len(parts[0])-1]

			if "Merge" == parts[1] {
				cmd := exec.Command("git", "diff", hash+"^", hash, "--numstat")
				out, err := cmd.Output()

				if err != nil {
					log.Println(log.DEBUG, err)
				}

				for _, line := range strings.Split(string(out), "\n") {
					parts := strings.Fields(line)

					// 3 fields :: lines added, lines deleted, filename
					if len(parts) == 3 && parts[0] != "-" && parts[1] != "-" {
						added, err := strconv.Atoi(parts[0])
						if err != nil {
							log.Println(log.DEBUG, err)
						}

						deleted, err := strconv.Atoi(parts[1])
						if err != nil {
							log.Println(log.DEBUG, err)
						}

						totLinesReviewed += added + deleted // IDK if needs to be + or -
					}
				}
			}
		}
		// ELSE NOTHING - COMMIT IS NOT A PULL REQUEST
	}

	os.Chdir(dir)
	rem := exec.Command("rm", "-r", repo.Name)
	err = rem.Run()

	if err != nil {
		log.Println(log.DEBUG, err)
	}
	os.RemoveAll(repo.Name)

	return totLinesReviewed
}