package metrics

import (
	"ECE461-Team1-Repository/api"
	"fmt"
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

func GetMetrics(siteType int, url string, TOKEN string) {
	var repo api.Repo

	if siteType == NPM {
		api.GetGithubURL(url)
		return
	} else if siteType == GITHUB {
		repo = api.GetRepo(url, TOKEN)

	}

	rampUp := -1
	busFactor := getBusFactor(repo.ContributorsURL, TOKEN)
	correctness := -1
	responsiveness := getResponsivenessScore(repo.Owner.Login, repo.Name, TOKEN)
	license := getLicenseScore(repo)

	fmt.Println("Ramp-up Time:", rampUp)
	fmt.Println("Bus Factor:", busFactor)
	fmt.Println("Correctness:", correctness)
	fmt.Println("Responsiveness:", responsiveness)
	fmt.Println("License:", license)
}
