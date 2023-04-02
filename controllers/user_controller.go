package controllers

import (
	"ECE461-Team1-Repository/configs"
	"ECE461-Team1-Repository/metrics"
	"ECE461-Team1-Repository/models"
	"ECE461-Team1-Repository/responses"
	"context"
	"encoding/json"
	"fmt"

	//"io"
	"net/http"
	"strconv"
	"time"

	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson" //add this
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var repoCollection *mongo.Collection = configs.GetCollection(configs.DB, "repos")
var validate = validator.New()

type resp struct {
	Url string `json:"url"`
}

func CreateRepo() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var repo models.Repo

		defer cancel()

		var respbody resp
		c.BindJSON(&respbody)
		fmt.Println("in create repo")
		fmt.Println(respbody.Url)

		ndjson := metrics.GetMetrics("https://github.com", 1, respbody.Url)

		//if netscore is good...
		//upload zip file

		type NDJSON struct {
			Name                string  `json:"URL"`
			Rampup              float64 `json:"RAMP_UP_SCORE"`
			Correctness         float64 `json:"CORRECTNESS_SCORE"`
			Responsivemaintainer float64 `json:"RESPONSIVE_MAINTAINER_SCORE"`
			Busfactor           float64 `json:"BUS_FACTOR_SCORE"`
			Reviewcoverage      float64 `json:"REVIEW_COVERAGE_SCORE"`
			Dependancypinning   float64 `json:"DEPENDENCY_PINNING_RATE"`
			License             int     `json:"LICENSE_SCORE"`
			Net                 float64 `json:"NET_SCORE"`
		}
		
		
		var ndjsonData NDJSON
		
		err := json.Unmarshal([]byte(ndjson), &ndjsonData)
		if err != nil {
			return
		}

		// ndjsonData.Reviewcoverage = 0.32
		ndjsonData.Dependancypinning = 0.42
		
		//payload := strings.NewReader(fmt.Sprintf(`{ "name": "%s", "rampup": %f, "correctness": %d, "responsivemaintainer": %f, "busfactor": %f, "reviewcoverage": %f, "dependancypinning": %f, "license": %d, "net": %f }`, ndjsonData.Name, ndjsonData.Rampup, ndjsonData.Correctness, ndjsonData.Responsivemaintainer, ndjsonData.Busfactor, ndjsonData.Reviewcoverage, ndjsonData.Dependancypinning, ndjsonData.License, ndjsonData.Net))
		//payload := strings.NewReader(`{` + " " + ` "name": `+ ndjsonData.Name + `,` + " " + `"rampup": ` + strconv.FormatFloat(ndjsonData.Rampup, 'f', 2, 64) + `,` + "" + `"correctness": `+strconv.FormatFloat(ndjsonData.Correctness, 'f', 1, 64)+`,` + "" + `"responsivemaintainer": `+strconv.FormatFloat(ndjsonData.Responsivemaintainer, 'f', 2, 64)+`,` + "" + `
		//"busfactor": `+strconv.FormatFloat(ndjsonData.Busfactor, 'f', 2, 64)+`,` + "" + `"reviewcoverage": `+strconv.FormatFloat(ndjsonData.Reviewcoverage, 'f', 2, 64)+`,` + "" + `"dependancypinning": `+strconv.FormatFloat(ndjsonData.Dependancypinning, 'f', 2, 64)+`,` + "" + `"license": `+strconv.Itoa(ndjsonData.License)+`,` + "" + `"net": `+strconv.FormatFloat(ndjsonData.Net, 'f', 2, 64)+ "" + `}`)		
		payload := strings.NewReader(`{
			"name": "` + ndjsonData.Name + `",
			"rampup": ` + strconv.FormatFloat(ndjsonData.Rampup, 'f', 2, 64) + `,
			"correctness": ` + strconv.FormatFloat(ndjsonData.Correctness, 'f', 1, 64) + `,
			"responsivemaintainer": ` + strconv.FormatFloat(ndjsonData.Responsivemaintainer, 'f', 2, 64) + `,
			"busfactor": ` + strconv.FormatFloat(ndjsonData.Busfactor, 'f', 2, 64) + `,
			"reviewcoverage": ` + strconv.FormatFloat(ndjsonData.Reviewcoverage, 'f', 2, 64) + `,
			"dependancypinning": ` + strconv.FormatFloat(ndjsonData.Dependancypinning, 'f', 2, 64) + `,
			"license": ` + strconv.Itoa(ndjsonData.License) + `,
			"net": ` + strconv.FormatFloat(ndjsonData.Net, 'f', 2, 64) + `
		}`)
		fmt.Println(payload)
		//fmt.Println(payload2)


		// payload = io.reader(payload)
		// rc, ok := payload.(io.ReadCloser)
		// if !ok && payload != nil {
		//         rc = io.NopCloser(payload)
		// }
		c.Request.Body = ioutil.NopCloser(payload)
		if err := c.BindJSON(&repo); err != nil {
			// c.Request.Body = payload
			c.JSON(http.StatusBadRequest, responses.RepoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		//use the validator library to validate required fields
		if validationErr := validate.Struct(&repo); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.RepoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newRepo := models.Repo{
			//Id:       primitive.NewObjectID(),
			Name:                 repo.Name,
			RampUp:               repo.RampUp,
			Correctness:          repo.Correctness,
			ResponsiveMaintainer: repo.ResponsiveMaintainer,
			BusFactor:            repo.BusFactor,
			ReviewCoverage:       repo.ReviewCoverage,
			DependancyPinning:    repo.DependancyPinning,
			License:              repo.License,
			Net:                  repo.Net,
		}
		// newRepo := models.Repo{
		//     //Id:       primitive.NewObjectID(),
		//     Name:     "test",
		//     RampUp:   0.2,
		//     Correctness: 0.1,
		//     ResponsiveMaintainer: 0.3,
		//     BusFactor: 0.4,
		//     ReviewCoverage: 0.5,
		//     DependancyPinning: 0.6,
		//     License: 1,
		//     Net: 0.8,
		// }

		fmt.Println("inserting: " + newRepo.Name)
		result, err := repoCollection.InsertOne(ctx, newRepo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RepoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.RepoResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetARepo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		repoId := c.Param("repoId")
		var repo models.Repo
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(repoId)

		err := repoCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&repo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RepoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.RepoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": repo}})
	}
}

func EditARepo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		repoId := c.Param("repoId")
		var repo models.Repo
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(repoId)

		//validate the request body
		if err := c.BindJSON(&repo); err != nil {
			c.JSON(http.StatusBadRequest, responses.RepoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&repo); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.RepoResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": repo.Name, "rampup": repo.RampUp, "correctness": repo.Correctness, "responsivemaintainer": repo.ResponsiveMaintainer, "busfactor": repo.BusFactor, "reviewcoverage": repo.ReviewCoverage,
			"dependancypinning": repo.DependancyPinning, "license": repo.License, "net": repo.Net}
		result, err := repoCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RepoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated repo details
		var updatedRepo models.Repo
		if result.MatchedCount == 1 {
			err := repoCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedRepo)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.RepoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.RepoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedRepo}})
	}
}

func DeleteARepo() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		repoId := c.Param("repoId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(repoId)

		result, err := repoCollection.DeleteOne(ctx, bson.M{"_id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RepoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.RepoResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Repo with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.RepoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Repo successfully deleted!"}},
		)
	}
}

func GetAllRepos() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var repos []models.Repo
		defer cancel()

		results, err := repoCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.RepoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleRepo models.Repo
			if err = results.Decode(&singleRepo); err != nil {
				c.JSON(http.StatusInternalServerError, responses.RepoResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			repos = append(repos, singleRepo)
		}

		c.JSON(http.StatusOK,
			responses.RepoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": repos}},
		)
	}
}

func GetMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK,
			responses.RepoResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": c.Request.Body}},
		)
	}
}
