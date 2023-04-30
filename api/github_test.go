package api

import (
	"os"
	"testing"
	"os/exec"
)

func TestGetGraphQL(t *testing.T) {
	if getGraphQLData(`{"query" : "query{repository(owner: \"cloudinary\", name: \"cloudinary_npm\") {total: issues {totalCount} closed:issues(states: CLOSED) {totalCount}}}"}`)[0] != 123 {
		t.Fatal("Error getting data from npm url (Probably invalid github token)!")
	}
}

func TestRESTler(t *testing.T) {
	// Set RESTler command and arguments
	cmd := exec.Command("restler", "-spec", "openapi.json", "-output", "test-results.json")

	// Run RESTler command
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Error running RESTler: %s\n%s", err, out)
	}

	// Read RESTler test results
	f, err := os.Open("test-results.json")
	if err != nil {
		t.Fatalf("Error reading test results: %s", err)
	}
	defer f.Close()

	// TODO: parse RESTler test results and check for errors
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

func TestGetPackageRequirements(t *testing.T) {
	// Testing package.json
	tst := GetPackageRequirements("expressjs", "express")
	if tst > 1 {
		t.Fatal("Error getting pin rate for package.json!")
	}
	// Testing requirements.txt
	tst = GetPackageRequirements("binder-examples", "requirements")
	if tst > 1 {
		t.Fatal("Error getting pin rate for requirements.txt!")
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

	if GetLicenseFromFile("expressjs", "express") != 1 {
		t.Fatal("License: Error cloning repository!")
	}

	if CountReviewedLines(tst) < 0 {
		t.Fatal("Counting: Error cloning repository!")
	}

	DeleteClonedRepo(tst)
	// Check if repo is deleted
	_, err := os.Stat(tst.Name)
	if err == nil {
		t.Fatal("Error deleting repository!")
	}
}
