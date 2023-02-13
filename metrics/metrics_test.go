package metrics

import (
	"fmt"
	"testing"
)

/*func TestNetScore1(t *testing.T) {
	url := "https://github.com/cloudinary/cloudinary_npm"
	siteType := 1
	name := "cloudinary/cloudinary_npm"

	netscore, _ := GetMetrics(url, siteType, name)

	if netscore == 0.8807582 {
		t.Logf("Net Score Passed")
	} else {
		t.Errorf("Net Score Failed")
	}
}

func TestNetScore2(t *testing.T) {
	url := "https://www.npmjs.com/package/express"
	siteType := 0
	name := "express"

	netscore, _ := GetMetrics(url, siteType, name)
	if netscore == 0.72071 {
		t.Logf("Net Score Passed")
	} else {
		t.Errorf("Net Score Failed")
	}
}*/

func TestBusFactor(t *testing.T) {
	url := "https://api.github.com/repos/cloudinary/cloudinary_npm/contributors"
	if getBusFactor(url) != 0.8113949 {
		t.Errorf("Bus Factor Failed")
	}
}

func TestResponsiveness(t *testing.T) {
	owner := "cloudinary"
	name := "cloudinary_npm"
	fmt.Println(getResponsivenessScore(owner, name))
	if getResponsivenessScore(owner, name) != 0.9563492 {
		t.Errorf("Responsiveness Failed")
	}
}

/*
		siteType := 1
		name := "cloudinary/cloudinary_npm"
		_, ndjson := GetMetrics(url, siteType, name)
		//extract BUS_FACTOR_SCORE from ndjson
		regMatch := regexp.MustCompile(`{"URL":".*", "NET_SCORE":(.*), "RAMP_UP_SCORE":(.*), "CORRECTNESS_SCORE":(.*), "BUS_FACTOR_SCORE":(.*), "RESPONSIVE_MAINTAINER_SCORE":(.*), "LICENSE_SCORE":(.*)}`)
		result := regMatch.FindStringSubmatch(ndjson)
		busFactorScore, _ := strconv.ParseFloat(result[4], 64)
		if busFactorScore == 0.81 {
			t.Logf("Bus Factor Score Passed")
		} else {
			t.Errorf("Bus Factor Score Failed")
		}
}

func TestRampUpScaler1(t *testing.T) {
	score := float32(0.1)
	result := RampUpScaler(score)
	if result == 0.1 {
		t.Logf("Ramp Up Scaler Passed")
	} else {
		t.Errorf("Ramp Up Scaler Failed")
	}
}

func TestRampUpScaler2(t *testing.T) {
	score := float32(0.25)
	result := RampUpScaler(score)
	if result == 1 {
		t.Logf("Ramp Up Scaler Passed")
	} else {
		t.Errorf("Ramp Up Scaler Failed")
	}
}

func TestRampUpScaler3(t *testing.T) {
	score := float32(0.5)
	result := RampUpScaler(score)
	fmt.Println(result)
	if result == 0.8888889 {
		t.Logf("Ramp Up Scaler Passed")
	} else {
		t.Errorf("Ramp Up Scaler Failed")
	}
}*/
