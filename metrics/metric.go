package metrics

import (
	"ECE461-Team1-Repository/api"
	"ECE461-Team1-Repository/log"
	"fmt"
	"regexp"
	"strconv"
)

func getBusFactor(url string) float32 {
	return 1 - api.GetContributionRatio(url)
}

func getResponsivenessScore(owner, name string) float32 {
	closed, total := api.GetIssuesCount(owner, name)
	if total != 0 {
		return float32(closed) / float32(total)
	} else {
		return 0
	}
}

// getLicenseScore checks if the output of GetLicenseFromREADME is blank, and assigns 0 or 1 accordingly
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

/*
getRampUpScore uses the output of RunClocOnRepo, which is the output of cloc. The function
parses the string using a regex and pulls the last 2 numerical values from the string, the
number of commented lines of code and the total number of lines of code. The ratio of these
2 values is used as the score for ramp-up time.
*/
func getRampUpScore(repo api.Repo) float32 {

	clocString := api.RunClocOnRepo(repo)
	regMatch := regexp.MustCompile(`.*SUM:\s*\d*\s*\d*\s*(\d*)\s*(\d*)`)
	commentLines := regMatch.FindStringSubmatch(clocString)[1]
	codeLines := regMatch.FindStringSubmatch(clocString)[2]

	commentLinesVal, err := strconv.Atoi(commentLines)

	if err != nil {
		log.Println(log.DEBUG, err)
	}

	codeLinesVal, err := strconv.Atoi(codeLines)

	if err != nil {
		log.Println(log.DEBUG, err)
	}

	var score float32
	if codeLinesVal != 0 {
		score = float32(commentLinesVal) / float32(codeLinesVal)
	} else {
		score = 0
	}
	// insert scaling factor here
	score = RampUpScaler(score)

	return score
}

/*
RampUpScaler inputs the raw ramp up score (lines of comments / lines of code) and scales it
through a piecewise function based on common comment standards. The output remains between 0 and 1.
*/

func RampUpScaler(score float32) float32 {

	if score <= 0.1 {
		return score
	} else if score > 0.1 && score <= 0.25 {
		return 6 * score - 0.5
	} else {
		var denomConst float32 = 0.5625
		score = (score - 0.25) * (score - 0.25)
		score *= -1
		score = score/denomConst + 1
		return score
	}

}

func GetMetrics(baseURL string, siteType int, name string) (float32, string) {
	var repo api.Repo

	if siteType == api.NPM {
		giturl := api.GetGithubURL(name)
		// parse the github url
		gitLinkMatch := regexp.MustCompile(".*github.com/(.*).git")
		githubURL := gitLinkMatch.FindStringSubmatch(giturl)[1]
		repo = api.GetRepo(githubURL)
	} else if siteType == api.GITHUB {
		repo = api.GetRepo(name)
	}

	rampUp := getRampUpScore(repo)
	correctness := getCorrectnessScore(repo)
	busFactor := getBusFactor(repo.ContributorsURL)
	responsiveness := getResponsivenessScore(repo.Owner.Login, repo.Name)
	license := getLicenseScore(repo)
	netScore := (0.1*float32(rampUp) + 0.1*float32(correctness) + 0.3*float32(busFactor) + 0.3*responsiveness + 0.2*float32(license)) * float32(license)
	// multiply by license score

	// Log (info)
	log.Printf(log.INFO, "Name: %v", name)
	log.Printf(log.INFO, "Net Score: %v", netScore)
	log.Printf(log.INFO, "Ramp-up Time: %v", rampUp)
	log.Printf(log.INFO, "Bus Factor: %v", busFactor)
	log.Printf(log.INFO, "Correctness: %v", correctness)
	log.Printf(log.INFO, "Responsiveness: %#v", responsiveness)
	log.Printf(log.INFO, "License: %v", license)

	ndjson := `{"URL":"` + baseURL + `", "NET_SCORE":` + fmt.Sprintf("%.2f", netScore) + `, "RAMP_UP_SCORE":` + fmt.Sprintf("%.2f", rampUp) +
		`, "CORRECTNESS_SCORE":` + fmt.Sprintf("%.1f", correctness) + `, "BUS_FACTOR_SCORE":` + fmt.Sprintf("%.2f", busFactor) + `, "RESPONSIVE_MAINTAINER_SCORE":` + fmt.Sprintf("%.2f", responsiveness) + `, "LICENSE_SCORE":` + fmt.Sprintf("%d", license) + `}`

	log.Printf(log.DEBUG, ndjson)
	// fmt.Println(netScore)

	return netScore, ndjson
}
