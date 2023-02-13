package metrics

import (
	"ECE461-Team1-Repository/api"
	"testing"
)

func TestBusFactor(t *testing.T) {
	url := "https://api.github.com/repos/cloudinary/cloudinary_npm/contributors"
	if getBusFactor(url) != 0.8113949 {
		t.Fatal("Bus Factor Failed")
	}
}

func TestResponsiveness(t *testing.T) {
	owner := "cloudinary"
	name := "cloudinary_npm"
	if getResponsivenessScore(owner, name) != 0.9563492 {
		t.Fatal("Responsiveness Failed")
	}
}

func TestGetLicenseScore(t *testing.T) {
	tst := api.Repo{FullName: "expressjs/express"}
	if getLicenseScore(tst) == 0 {
		t.Fatal("License Score Failed")
	}
}

func TestGetRampupAndCorrectnessScore(t *testing.T) {
	tst := api.Repo{CloneURL: "https://github.com/expressjs/express.git", Name: "express"}
	tst_ramp := getRampUpScore(tst)
	tst_correctness := getCorrectnessScore(tst)
	if tst_ramp != 0.48373964 || tst_correctness != 1.0 {
		t.Fatal("Cloning process Failed")
	}
}

func TestScaler(t *testing.T) {
	if RampUpScaler(0.0) != 0.0 {
		t.Fatal("Scaling process Failed")
	}
	if RampUpScaler(0.6) != 0.7822222 {
		t.Fatal("Scaling process Failed")
	}
}

func TestGetMetric(t *testing.T) {
	url := "https://www.npmjs.com/package/express"
	siteType := 0
	name := "express"
	netscore, _ := GetMetrics(url, siteType, name)
	if netscore != 0.7206798 {
		t.Fatal("GetMetric Failed")
	}
}
