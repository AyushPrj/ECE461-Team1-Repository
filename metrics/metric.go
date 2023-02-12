package metrics

import (
	"ECE461-Team1-Repository/api"
	"fmt"
	"regexp"
	"strconv"
	"log"
)

const (
	NPM    = 0
	GITHUB = 1
)

func getBusFactor(url, TOKEN string) float32 {
	// TODO: might have to scale this someway
	return 1 - api.GetContributionRatio(url, TOKEN)
}

func getResponsivenessScore(owner, name, TOKEN string) float32 {
	closed, total := api.GetIssuesCount(owner, name, TOKEN)
	return float32(closed) / float32(total)
}

func getLicenseScore(repo api.Repo) int {
	readme_string := api.GetRawREADME(repo)
	license_string := api.GetLicenseFromREADME(readme_string)

	if license_string == "" {
		return 0
	}

	return 1

}

func getCorrectnessScore(repo api.Repo) float64 {

	return api.CheckRepoForTest(repo)
}

func getRampUpScore(repo api.Repo) float32 {

	clocString := api.RunClocOnRepo(repo)
	regMatch := regexp.MustCompile(`.*SUM:\s*\d*\s*\d*\s*(\d*)\s*(\d*)`)
	commentLines := regMatch.FindStringSubmatch(clocString)[1]
	codeLines := regMatch.FindStringSubmatch(clocString)[2]

	commentLinesVal, err := strconv.Atoi(commentLines)

	if err != nil {
		log.Fatal(err)
	}

	codeLinesVal, err := strconv.Atoi(codeLines)

	if err != nil {
		log.Fatal(err)
	}

	var score float32
	score = float32(commentLinesVal) / float32(codeLinesVal)
	// fmt.Printf("score: %f\n", score)
	// insert scaling factor here

	return score
}

func GetMetrics(baseURL string, siteType int, name string, TOKEN string) (float32, string) {
	var repo api.Repo

	if siteType == NPM {
		giturl := api.GetGithubURL(name)
		// parse the github url
		gitLinkMatch := regexp.MustCompile(".*github.com/(.*).git")
		githubURL := gitLinkMatch.FindStringSubmatch(giturl)[1]
		repo = api.GetRepo(githubURL, TOKEN)
		// fmt.Println(repo.FullName)
	} else if siteType == GITHUB {
		repo = api.GetRepo(name, TOKEN)
		// fmt.Println(repo.Name)
	}

	// fmt.Println(repo.CloneURL)
	rampUp := getRampUpScore(repo)
	//rampUp := -1.0
	correctness := getCorrectnessScore(repo)
	busFactor := getBusFactor(repo.ContributorsURL, TOKEN)
	responsiveness := getResponsivenessScore(repo.Owner.Login, repo.Name, TOKEN)
	license := getLicenseScore(repo)
	netScore := (0.1*float32(rampUp) + 0.1*float32(correctness) + 0.3*float32(busFactor) + 0.3*responsiveness + 0.2*float32(license)) * float32(license)
	// multiply by license score

	// TODO: Add to log (info)
	// fmt.Println("Ramp-up Time:", rampUp)
	// fmt.Println("Bus Factor:", busFactor)
	// fmt.Println("Correctness:", correctness)
	// fmt.Println("Responsiveness:", responsiveness)
	// fmt.Println("License:", license)

	ndjson := `{"URL":"` + baseURL + `", "NET_SCORE":` + fmt.Sprintf("%v", netScore) + `, "RAMP_UP_SCORE":` + fmt.Sprintf("%v", rampUp) +
		`, "CORRECTNESS_SCORE":` + fmt.Sprintf("%.1f", correctness) + `, "BUS_FACTOR_SCORE":` + fmt.Sprintf("%.2f", busFactor) + `, "RESPONSIVE_MAINTAINER_SCORE":` + fmt.Sprintf("%.2f", responsiveness) + `, "LICENSE_SCORE":` + fmt.Sprintf("%d", license) + `}`
	// fmt.Println(ndjson)

	return netScore, ndjson
}