package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// var coverProfile = flag.String("coverprofile", "", "write coverage profile to `file`")
var total int = 20
var pass int = 0

func TestMain2(t *testing.T) {
	os.Args = []string{"main", "test.txt"}
	main()
}

func checkFormat(ndjson string) bool {
	inp := []byte(ndjson)
	regMatch := regexp.MustCompile(`{"URL":".*", "NET_SCORE":(.*), "RAMP_UP_SCORE":(.*), "CORRECTNESS_SCORE":(.*), "BUS_FACTOR_SCORE":(.*), "RESPONSIVE_MAINTAINER_SCORE":(.*), "LICENSE_SCORE":(.*)}`)
	result := regMatch.Find(inp)

	if result == nil {
		return false
	} else {
		return true
	}
}

func checkNetScore(netScore int) bool {
	//check if netScore is between 0 and 1
	if netScore >= 0 && netScore <= 1 {
		return true
	} else {
		return false
	}
}

func TestFormat(t *testing.T) {
	testurl := `{"URL":"https://www.npmjs.com/package/browserify", "NET_SCORE":-1, "RAMP_UP_SCORE":0.11, "CORRECTNESS_SCORE":0.02, "BUS_FACTOR_SCORE":0.75, "RESPONSIVE_MAINTAINER_SCORE":0.76, "LICENSE_SCORE":0}`
	result := checkFormat(testurl)
	if result == false {
		t.Errorf("ndjson format Failed")
	} else {
		t.Logf("ndjson format Passed")
		pass = pass + 1
	}
}

func TestNetScore(t *testing.T) {
	testurl := `{"URL":"https://www.npmjs.com/package/browserify", "NET_SCORE":0.5, "RAMP_UP_SCORE":0.11, "CORRECTNESS_SCORE":0.02, "BUS_FACTOR_SCORE":0.75, "RESPONSIVE_MAINTAINER_SCORE":0.76, "LICENSE_SCORE":0}`
	regMatch := regexp.MustCompile(`{"URL":".*", "NET_SCORE":(.*), "RAMP_UP_SCORE":(.*), "CORRECTNESS_SCORE":(.*), "BUS_FACTOR_SCORE":(.*), "RESPONSIVE_MAINTAINER_SCORE":(.*), "LICENSE_SCORE":(.*)}`)
	result := regMatch.FindStringSubmatch(testurl)
	netScore, _ := strconv.Atoi(result[1])
	result2 := checkNetScore(netScore)
	if result2 == false {
		t.Errorf("Net Score Failed. Should be between 0 and 1")
	} else {
		t.Logf("Net Score Passed")
		//fmt.Println("Net Score Passed")
		pass = pass + 1
	}
}

func TestRampUpScore(t *testing.T) {
	//check if ramp_up_score is either -1 or between 0 and 1
	testurl := `{"URL":"https://www.npmjs.com/package/browserify", "NET_SCORE":0.5, "RAMP_UP_SCORE":0.11, "CORRECTNESS_SCORE":0.02, "BUS_FACTOR_SCORE":0.75, "RESPONSIVE_MAINTAINER_SCORE":0.76, "LICENSE_SCORE":0}`
	regMatch := regexp.MustCompile(`{"URL":".*", "NET_SCORE":(.*), "RAMP_UP_SCORE":(.*), "CORRECTNESS_SCORE":(.*), "BUS_FACTOR_SCORE":(.*), "RESPONSIVE_MAINTAINER_SCORE":(.*), "LICENSE_SCORE":(.*)}`)
	result := regMatch.FindStringSubmatch(testurl)
	rampUpScore, _ := strconv.ParseFloat(result[2], 64)
	if rampUpScore == -1 || (rampUpScore >= 0 && rampUpScore <= 1) {
		t.Logf("Ramp Up Score Passed")
		//fmt.Println("Ramp Up Score Passed")
		pass = pass + 1
	} else {
		t.Errorf("Ramp Up Score Failed. Should be either -1 or between 0 and 1")
	}
}

func TestCorrectnessScore(t *testing.T) {
	//check if correctness_score is either -1 or between 0 and 1
	testurl := `{"URL":"https://www.npmjs.com/package/browserify", "NET_SCORE":0.5, "RAMP_UP_SCORE":0.11, "CORRECTNESS_SCORE":0.02, "BUS_FACTOR_SCORE":0.75, "RESPONSIVE_MAINTAINER_SCORE":0.76, "LICENSE_SCORE":0}`
	regMatch := regexp.MustCompile(`{"URL":".*", "NET_SCORE":(.*), "RAMP_UP_SCORE":(.*), "CORRECTNESS_SCORE":(.*), "BUS_FACTOR_SCORE":(.*), "RESPONSIVE_MAINTAINER_SCORE":(.*), "LICENSE_SCORE":(.*)}`)
	result := regMatch.FindStringSubmatch(testurl)
	correctnessScore, _ := strconv.ParseFloat(result[3], 64)
	if correctnessScore == -1 || (correctnessScore >= 0 && correctnessScore <= 1) {
		t.Logf("Correctness Score Passed")
		//fmt.Println("Correctness Score Passed")
		pass = pass + 1
	} else {
		t.Errorf("Correctness Score Failed. Should be either -1 or between 0 and 1")
	}
}

func TestBusFactorScore(t *testing.T) {
	//check if bus_factor_score is either -1 or between 0 and 1
	testurl := `{"URL":"https://www.npmjs.com/package/browserify", "NET_SCORE":0.5, "RAMP_UP_SCORE":0.11, "CORRECTNESS_SCORE":0.02, "BUS_FACTOR_SCORE":0.75, "RESPONSIVE_MAINTAINER_SCORE":0.76, "LICENSE_SCORE":0}`
	regMatch := regexp.MustCompile(`{"URL":".*", "NET_SCORE":(.*), "RAMP_UP_SCORE":(.*), "CORRECTNESS_SCORE":(.*), "BUS_FACTOR_SCORE":(.*), "RESPONSIVE_MAINTAINER_SCORE":(.*), "LICENSE_SCORE":(.*)}`)
	result := regMatch.FindStringSubmatch(testurl)
	busFactorScore, _ := strconv.ParseFloat(result[4], 64)
	if busFactorScore == -1 || (busFactorScore >= 0 && busFactorScore <= 1) {
		t.Logf("Bus Factor Score Passed")
		//fmt.Println("Bus Factor Score Passed")
		pass = pass + 1
	} else {
		t.Errorf("Bus Factor Score Failed. Should be either -1 or between 0 and 1")
	}
}

func TestResponsiveMaintainerScore(t *testing.T) {
	//check if responsive_maintainer_score is either -1 or between 0 and 1
	testurl := `{"URL":"https://www.npmjs.com/package/browserify", "NET_SCORE":0.5, "RAMP_UP_SCORE":0.11, "CORRECTNESS_SCORE":0.02, "BUS_FACTOR_SCORE":0.75, "RESPONSIVE_MAINTAINER_SCORE":0.76, "LICENSE_SCORE":0}`
	regMatch := regexp.MustCompile(`{"URL":".*", "NET_SCORE":(.*), "RAMP_UP_SCORE":(.*), "CORRECTNESS_SCORE":(.*), "BUS_FACTOR_SCORE":(.*), "RESPONSIVE_MAINTAINER_SCORE":(.*), "LICENSE_SCORE":(.*)}`)
	result := regMatch.FindStringSubmatch(testurl)
	responsiveMaintainerScore, _ := strconv.ParseFloat(result[5], 64)
	if responsiveMaintainerScore == -1 || (responsiveMaintainerScore >= 0 && responsiveMaintainerScore <= 1) {
		t.Logf("Responsive Maintainer Score Passed")
		//fmt.Println("Responsive Maintainer Score Passed")
		pass = pass + 1
	} else {
		t.Errorf("Responsive Maintainer Score Failed. Should be either -1 or between 0 and 1")
	}
}

func TestLicenseScore(t *testing.T) {
	//check if license_score is either 0 or 1
	testurl := `{"URL":"https://www.npmjs.com/package/browserify", "NET_SCORE":0.5, "RAMP_UP_SCORE":0.11, "CORRECTNESS_SCORE":0.02, "BUS_FACTOR_SCORE":0.75, "RESPONSIVE_MAINTAINER_SCORE":0.76, "LICENSE_SCORE":0}`
	regMatch := regexp.MustCompile(`{"URL":".*", "NET_SCORE":(.*), "RAMP_UP_SCORE":(.*), "CORRECTNESS_SCORE":(.*), "BUS_FACTOR_SCORE":(.*), "RESPONSIVE_MAINTAINER_SCORE":(.*), "LICENSE_SCORE":(.*)}`)
	result := regMatch.FindStringSubmatch(testurl)
	licenseScore, _ := strconv.ParseFloat(result[6], 64)
	if licenseScore == 1 || licenseScore == 0 {
		t.Logf("License Score Passed")
		//fmt.Println("License Score Passed")
		pass = pass + 1
	} else {
		t.Errorf("License Score Failed. Should be either 0 or 1")
	}
}

func TestURL(t *testing.T) {
	test := `{"URL":"https://www.npmjs.com/package/browserify", "NET_SCORE":0.5, "RAMP_UP_SCORE":0.11, "CORRECTNESS_SCORE":0.02, "BUS_FACTOR_SCORE":0.75, "RESPONSIVE_MAINTAINER_SCORE":0.76, "LICENSE_SCORE":0}`
	//check if URL contains github
	regMatch := regexp.MustCompile(`{"URL":"(.*)", "NET_SCORE":(.*), "RAMP_UP_SCORE":(.*), "CORRECTNESS_SCORE":(.*), "BUS_FACTOR_SCORE":(.*), "RESPONSIVE_MAINTAINER_SCORE":(.*), "LICENSE_SCORE":(.*)}`)
	result := regMatch.FindStringSubmatch(test)
	if strings.Contains(result[1], "github.com") || strings.Contains(result[1], "npmjs") {
		t.Logf("URL Passed")
		//fmt.Println("URL Passed")
		pass = pass + 1
	} else {
		t.Errorf("URL Failed")
	}
}

func TestSummary(t *testing.T) {
	fmt.Println("Total: ", total)
	fmt.Println("Passed: ", pass)
	fmt.Println(pass, "/", total, "test cases passed")
}

func TestMain(m *testing.M) {
	m.Run()
}
