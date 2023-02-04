package metrics

import (
	"fmt"
	"proj/api"
)

func getBusFactor(url, TOKEN string) float32 {
	// TODO: might have to scale this someway
	return 1 - api.GetContributionRatio(url, TOKEN)
}

func getResponsivenessScore(owner, name, TOKEN string) float32 {
	closed, total := api.GetIssuesCount(owner, name, TOKEN)
	return float32(closed) / float32(total)
}

func GetMetrics(url, TOKEN string) {
	repo := api.GetRepo(url, TOKEN)

	rampUp := -1
	busFactor := getBusFactor(repo.ContributorsURL, TOKEN)
	correctness := -1
	responsiveness := getResponsivenessScore(repo.Owner.Login, repo.Name, TOKEN)
	license := -1

	fmt.Println("Ramp-up Time:", rampUp)
	fmt.Println("Bus Factor:", busFactor)
	fmt.Println("Correctness:", correctness)
	fmt.Println("Responsiveness:", responsiveness)
	fmt.Println("License:", license)
}
