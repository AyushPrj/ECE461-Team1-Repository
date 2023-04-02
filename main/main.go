package main

import (
	"ECE461-Team1-Repository/api"
	"ECE461-Team1-Repository/log"
	//"ECE461-Team1-Repository/metrics"
	//"bufio"
	"fmt"
	//golog "log"
	//"os"
	//"regexp"
	//"sort"
	//"strings"

	//rest api
	"encoding/json"
	"net/http"

	"ECE461-Team1-Repository/configs"
	"ECE461-Team1-Repository/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Link struct {
	site     int
	name     string
	netScore float32
	ndjson   string
}

var GITHUB_TOKEN string
var LOG_LEVEL string
var LOG_FILE string

func init() {
	GITHUB_TOKEN = api.GITHUB_TOKEN
	LOG_LEVEL = log.LOG_LEVEL
	LOG_FILE = log.LOG_FILE
}

// func cli(toRateURL string) Link {
// 	// Initialize the log file
// 	f, err := os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
// 	if err != nil && LOG_LEVEL != "0" {
// 		golog.Fatalf("error opening log file: %v", err)
// 	}
// 	// Close log file after program is complete
// 	defer f.Close()
// 	defer log.Printf(log.INFO, "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")

// 	golog.SetOutput(f)

// 	log.Printf(log.INFO, "LOG LEVEL: %v", LOG_LEVEL)

// 	// args := os.Args[1:]
// 	//str := strings.Join(args, "")
// 	//file, err := os.Open(str)
// 	// if err != nil {
// 	// 	log.Println(log.INFO, "Failed to open input file")
// 	// }
// 	// scanner := bufio.NewScanner(file)
// 	// scanner.Split(bufio.ScanLines)
// 	// var text []string
// 	// for scanner.Scan() {
// 	// 	text = append(text, scanner.Text())
// 	// }

// 	// The method os.File.Close() is called on the os.File object to close the file
// 	//file.Close()

// 	var links Link

// 	// A loop iterates through and prints each of the slice values.
// 	//for _, each_ln := range text {
// 		var tmpSite int
// 		var tmpName string
// 		gitMatch := strings.Contains(toRateURL, "github")
// 		if gitMatch {
// 			gitLinkMatch := regexp.MustCompile(".*github.com/(.*)")
// 			tmpName = gitLinkMatch.FindStringSubmatch(toRateURL)[1]
// 			tmpSite = api.GITHUB
// 		} else {
// 			npmLinkMatch := regexp.MustCompile(".*package/(.*)")
// 			tmpName = npmLinkMatch.FindStringSubmatch(toRateURL)[1]
// 			tmpSite = api.NPM
// 		}

// 		// get the metrics in ndjson format for each link and add to list
// 		// fmt.Printf("%s\n", tmpName)
// 		netscore, ndjson := metrics.GetMetrics(toRateURL, tmpSite, tmpName)
// 		newLink := Link{site: tmpSite, name: tmpName, netScore: netscore, ndjson: ndjson}
// 		return newLink
// 	//}

// 	// Sort array of links by net score (descending)
// 	// sort.Slice(links, func(i, j int) bool {
// 	// 	return links[i].netScore > links[j].netScore
// 	// })

// 	// for _, tst_print := range links {
// 	// fmt.Printf("%+v\n", tst_print)
// 	// metrics.GetMetrics(tst_print.name, GITHUB_TOKEN)
// 	// if tst_print.site == GITHUB {
// 	// metrics.GetMetrics(tst_print.site, tst_print.name, GITHUB_TOKEN)
// 	// }
// 	// }

// 	// GETS ALL THE METRICS IN THIS FUNCTION GIVEN THE URL (ONLY WORKS FOR GITHUB CURRENTLY)
// 	// metrics.GetMetrics("cloudinary/cloudinary_npm", GITHUB_TOKEN)
// 	// metrics.GetMetrics("lodash/lodash", GITHUB_TOKEN)
// 	// metrics.GetMetrics("nullivex/nodist", GITHUB_TOKEN)

// 	// Display the metrics in the stored ndjson format
// 	// printOutput(links)
// 	return links
// }

func printOutput(links []Link) {
	//for _, link := range links {
	fmt.Println(links)
	//
}

var links []Link
type reposJson map[string]interface{}
type arr_repos []map[string]interface{}

func jsonOutput(c *gin.Context) {
	allrepos := make(arr_repos, 0) //necessary so that when you call the API again, it doesnt append the same stuff to the list 

	for _, link := range links {
		newjson := make(reposJson)
		json.Unmarshal([]byte(link.ndjson), &newjson)
		allrepos = append(allrepos, newjson)		
	}

	c.IndentedJSON(http.StatusOK, allrepos)
}

func main() {
	// args := os.Args[1:]
	// str := strings.Join(args, "")
	// file, err := os.Open(str)
	// if err != nil {
	// 	log.Println(log.INFO, "Failed to open input file")
	// }
	// scanner := bufio.NewScanner(file)
	// scanner.Split(bufio.ScanLines)
	// var text []string
	// for scanner.Scan() {
	// 	text = append(text, scanner.Text())
	// }
	
	// //The method os.File.Close() is called on the os.File object to close the file
	// file.Close()

	// for _, each_repo := range text {
	// 	links = append(links, cli(each_repo))
	// }

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))

	//run database
	 configs.ConnectDB()

	// router.Static("/assets", "./assets")
	// router.LoadHTMLGlob("views/*")

	routes.RepoRoute(router)

	router.Run("localhost:5500")
}
