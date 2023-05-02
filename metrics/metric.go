package metrics

import (
	"ECE461-Team1-Repository/api"
	"ECE461-Team1-Repository/log"
	"fmt"
	"math"
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

/*
getLicenseScore checks to see if license is in the README, if it is not, then check
if there is a file containing a license
*/
func getLicenseScore(repo api.Repo) int {
	readme_string := api.GetRawREADME(repo)
	license_string := api.GetLicenseFromREADME(readme_string)

	if license_string == "" {
		return api.GetLicenseFromFile(repo.Owner.Login, repo.Name)
	}
	return 1
}

/*
getCorrectnessScore calls api to check repo for correctness
*/

func getCorrectnessScore(repo api.Repo) float64 {
	return api.CheckRepoForTest(repo)
}

/*
getRampUpScore uses the output of RunClocOnRepo, which is the output of cloc. The function
parses the string using a regex and pulls the last 2 numerical values from the string, the
number of commented lines of code and the total number of lines of code. The ratio of these
2 values is used as the score for ramp-up time.
*/
func getRampUpScore(repo api.Repo) (float32, int) {

	clocString := api.RunClocOnRepo(repo)
	regMatch := regexp.MustCompile(`.*SUM:\s*\d*\s*\d*\s*(\d*)\s*(\d*)`).FindStringSubmatch(clocString)
	if len(regMatch) < 3 {
		log.Println(log.DEBUG, "Regex could find no match")
		return 0, 0
	}

	commentLines := regMatch[1]
	codeLines := regMatch[2]

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

	return score, codeLinesVal
}

/*
RampUpScaler inputs the raw ramp up score (lines of comments / lines of code) and scales it
through a piecewise function based on common comment standards. The output remains between 0 and 1.
*/

func RampUpScaler(score float32) float32 {

	if score <= 0.1 {
		return score
	} else if score > 0.1 && score <= 0.25 {
		return 6*score - 0.5
	} else {
		var denomConst float32 = 0.5625
		score = (score - 0.25) * (score - 0.25)
		score *= -1
		score = score/denomConst + 1
		return score
	}
}

/*
getDepPinRate returns a ratio of pinned dependencies to total dependencies
*/

func getDepPinRate(owner, name string) float32 {
	return float32(api.GetDepPinRate(owner, name))
}

/*
getReviewCoverage returns lines added from pull requests divided by total lines. Output lies in [0.0, 1.0]
*/

func getReviewCoverage(repo api.Repo, numLines int) float32 {
	reviewLines := api.CountReviewedLines(repo)

	if (reviewLines < numLines) {
		return float32(reviewLines) / float32(numLines)
	}
	return float32(numLines) / float32(reviewLines)
}

/*
GetMetrics calculates rating for the input repo
*/

func GetMetrics(baseURL string, siteType int, name string) string {
	var repo api.Repo
	// fmt.Printf("net score \n")

	if siteType == api.NPM {
		giturl := api.GetGithubURL(name)
		gitLinkMatch := regexp.MustCompile(".*github.com/(.*).git")
		githubURL := gitLinkMatch.FindStringSubmatch(giturl)[1]
		repo = api.GetRepo(githubURL)
	} else if siteType == api.GITHUB {
		repo = api.GetRepo(name)
	}

	rampUp, numLines := getRampUpScore(repo)

	correctness := getCorrectnessScore(repo)

	busFactor := getBusFactor(repo.ContributorsURL)

	responsiveness := getResponsivenessScore(repo.Owner.Login, repo.Name)

	license := getLicenseScore(repo)
  
	depPinRate := getDepPinRate(repo.Owner.Login, repo.Name)

	reviewCoverage := getReviewCoverage(repo, numLines)

	// delete the cloned repo
	api.DeleteClonedRepo(repo)

	if(math.IsNaN(float64(rampUp))){
		rampUp = 0;
	}
	if(math.IsNaN(float64(correctness))){
		correctness = 0;
	}
	if(math.IsNaN(float64(busFactor))){
		busFactor = 0;
	}
	if(math.IsNaN(float64(responsiveness))){
		responsiveness = 0;
	}
	
	if(math.IsNaN(float64(depPinRate))){
		depPinRate = 0;
	}
	if(math.IsNaN(float64(reviewCoverage))){
		reviewCoverage = 0;
	}
	
	// OLD FORMULA: (.1 * rampUp + .1 * correctness + .3 * busFactor + .3 * responsiveness + .2 * license) * license
	//netScore := (0.1*float32(rampUp) + 0.1*float32(correctness) + 0.3*float32(busFactor) + 0.3*responsiveness + 0.2*float32(license)) * float32(license)
	// NEW FORMULA: (.1 * rampUp + .1 * correctness + .3 * busFactor + .2 * responsiveness + .1 * depPinRate + .2 * reviewCoverage) * license
	netScore := (0.1*float32(rampUp) + 0.1*float32(correctness) + 0.3*float32(busFactor) + 0.3*responsiveness + 0.1*depPinRate + 0.1*reviewCoverage) * float32(license)
	/*
	NEW FORMULA SUMMARY:
		- Ramp-up time: 10%
		- Correctness: 10%
		- Bus Factor: 30%
		- Responsiveness: 30%
		- Dependency Pinning Rate: 10%
		- Code Review Coverage: 10%

		All or nothing for license: (if license is 0, then the entire score is 0)
		- License: 100%
	*/


	// Log (info)
	log.Printf(log.INFO, "Name: %v", name)
	log.Printf(log.INFO, "Net Score: %v", netScore)
	log.Printf(log.INFO, "Ramp-up Time: %v", rampUp)
	log.Printf(log.INFO, "Bus Factor: %v", busFactor)
	log.Printf(log.INFO, "Correctness: %v", correctness)
	log.Printf(log.INFO, "Responsiveness: %#v", responsiveness)
	log.Printf(log.INFO, "License: %v", license)
	log.Printf(log.INFO, "Dependency Pinning Rate: %v", depPinRate)
	log.Printf(log.INFO, "Code Review Coverage: %v", reviewCoverage)

	//ndjson := `{"URL":"` + name + `", "NET_SCORE":` + fmt.Sprintf("%.2f", netScore) + `, "RAMP_UP_SCORE":` + fmt.Sprintf("%.2f", rampUp) +
	//	`, "CORRECTNESS_SCORE":` + fmt.Sprintf("%.1f", correctness) + `, "BUS_FACTOR_SCORE":` + fmt.Sprintf("%.2f", busFactor) + `, "RESPONSIVE_MAINTAINER_SCORE":` + fmt.Sprintf("%.2f", responsiveness) +
	//	`, "LICENSE_SCORE":` + fmt.Sprintf("%d", license) + `, "DEPENDENCY_PINNING_RATE":` + fmt.Sprintf("%.2f", depPinRate) + `, "REVIEW_COVERAGE_SCORE":` + fmt.Sprintf("%.2f", reviewCoverage) +  `}`

	ndjson := `{"NetScore":` + fmt.Sprintf("%.2f", netScore) + `, "RampUp":` + fmt.Sprintf("%.2f", rampUp) +
		`, "Correctness":` + fmt.Sprintf("%.1f", correctness) + `, "BusFactor":` + fmt.Sprintf("%.2f", busFactor) + `, "ResponsiveMaintainer":` + fmt.Sprintf("%.2f", responsiveness) +
		`, "LicenseScore":` + fmt.Sprintf("%d", license) + `, "GoodPinningPractice":` + fmt.Sprintf("%.2f", depPinRate) + `, "PullRequest":` + fmt.Sprintf("%.2f", reviewCoverage) + `}`

	log.Printf(log.DEBUG, ndjson)
	fmt.Println(netScore)

	return ndjson
}
