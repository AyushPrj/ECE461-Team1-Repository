package api

import (
	"testing"
)

func TestGetGraphQL(t *testing.T) {
	if getGraphQLData(`{"query" : "query{repository(owner: \"cloudinary\", name: \"cloudinary_npm\") {total: issues {totalCount} closed:issues(states: CLOSED) {totalCount}}}"}`)[0] != 123 {
		t.Fatal("Error getting data from npm url!")
	}
}

func TestGetRequest(t *testing.T) {
	if getRequest("https://api.github.com/repos/expressjs/express/readme")[0] != 123 {
		t.Fatal("Error getting data from request!")
	}
}

func TestGetRepo(t *testing.T) {
	tst := GetRepo("expressjs/express")
	if tst.ID != 237159 {
		t.Fatal("Error getting repository!")
	}
}

func TestGetTopContributor(t *testing.T) {
	tst := Contributor{Login: "tst", ID: 1, NodeID: "abc", Contributions: 10}
	tst_arr := []Contributor{tst}
	if getTopContributor(tst_arr) != tst {
		t.Fatal("Error getting top contributor!")
	}
}

func TestTotalContributions(t *testing.T) {
	tst := Contributor{Login: "tst", ID: 1, NodeID: "abc", Contributions: 10}
	tst_arr := []Contributor{tst}
	if getTotalNumContributions(tst_arr) != 10 {
		t.Fatal("Error getting total contributions!")
	}
}

func TestGetContrinutionRatio(t *testing.T) {
	tst := GetContributionRatio("https://api.github.com/repos/expressjs/express/contributors")
	if tst > 1 {
		t.Fatal("Error calculating contribution ratio!")
	}
}

func TestGetIssuesCount(t *testing.T) {
	tst1, tst2 := GetIssuesCount("expressjs", "express")
	if tst1 < 1 || tst2 <= 1 {
		t.Fatal("Error getting issue count!")
	}
}

func TestDepPinRate(t *testing.T) {
	tst := GetDepPinRate("expressjs", "express")
	if tst > 1 {
		t.Fatal("Error getting dependency pin rate!")
	}
}

func TestGetReadmeUrl(t *testing.T) {
	tst := Repo{FullName: "expressjs/express"}
	if getReadmeURL(tst) != "https://raw.githubusercontent.com/expressjs/express/master/Readme.md" {
		t.Fatal("Error getting readme url!")
	}
}

func TestGetRawReadme(t *testing.T) {
	tst := Repo{FullName: "expressjs/express"}
	if GetRawREADME(tst)[0] != '[' {
		t.Fatal("Error getting readme!")
	}
}

func TestGetLicense(t *testing.T) {
	tst := "IBM PowerPC Initialization and Boot Software"
	tst2 := "MIT License"
	tst3 := ""
	if GetLicenseFromREADME(tst) != "IBM-pibs" {
		t.Fatal("Error reading license!")
	}

	if GetLicenseFromREADME(tst2) != "MIT" {
		t.Fatal("Error reading license!")
	}

	if GetLicenseFromREADME(tst3) != "" {
		t.Fatal("Error reading license!")
	}
}

func TestCloning(t *testing.T) {
	tst := Repo{CloneURL: "https://github.com/expressjs/express.git", Name: "express"}
	if RunClocOnRepo(tst)[0] != ' ' {
		t.Fatal("Cloc: Error cloning repository!")
	}

	if CheckRepoForTest(tst) != 1.0 {
		t.Fatal("Checking: Error cloning repository!")
	}

	if CountReviewedLines(tst) < 0 {
		t.Fatal("Counting: Error cloning repository!")
	}
}
