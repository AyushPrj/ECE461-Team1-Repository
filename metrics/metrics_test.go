package metrics

import (
	"ECE461-Team1-Repository/api"
	"testing"
)

func TestBusFactor(t *testing.T) {
	url := "https://api.github.com/repos/cloudinary/cloudinary_npm/contributors"
	if getBusFactor(url) > 1 {
		t.Fatal("Bus Factor Failed")
	}
}

func TestResponsivenessAndDepPinRate(t *testing.T) {
	owner := "cloudinary"
	name := "cloudinary_npm"
	if getResponsivenessScore(owner, name) > 1  || getDepPinRate(owner, name) > 1 {
		t.Fatal("Responsiveness Failed")
	}
}

func TestGetLicenseScore(t *testing.T) {
	tst := api.Repo{FullName: "expressjs/express"}
	if getLicenseScore(tst) == 0 {
		t.Fatal("License Score Failed")
	}
}

func TestGetRampupAndCorrectnessScoreAndReviewCoverage(t *testing.T) {
	tst := api.Repo{CloneURL: "https://github.com/expressjs/express.git", Name: "express"}
	tst_ramp, numLines := getRampUpScore(tst)
	tst_correctness := getCorrectnessScore(tst)
	tst_coverage := getReviewCoverage(tst, numLines)
	if tst_ramp >= 1 || tst_correctness != 1.0 || tst_coverage > 1{
		t.Fatal("Cloning process Failed")
	}
}

func TestScaler(t *testing.T) {
	if RampUpScaler(0.0) != 0.0 {
		t.Fatal("Scaling process Failed")
	}
	if RampUpScaler(0.6) >= 1 {
		t.Fatal("Scaling process Failed")
	}
}

func TestGetMetric(t *testing.T) {
	// url := "https://www.npmjs.com/package/express"
	// siteType := 0
	// name := "express"
	// netscore, _ := GetMetrics(url, siteType, name)
	// if netscore > 1 {
	// 	t.Fatal("GetMetric Failed")
	// }
}
