package controllers

import (
    "context"
    "ECE461-Team1-Repository/configs"
    "ECE461-Team1-Repository/models"
    "ECE461-Team1-Repository/responses"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson" //add this
)



var repoCollection *mongo.Collection = configs.GetCollection(configs.DB, "repos")
var validate = validator.New()

func CreateRepo() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        var repo models.Repo
        defer cancel()

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

        newRepo := models.Repo{
            //Id:       primitive.NewObjectID(),
            Name:     repo.Name,
            RampUp:   repo.RampUp,
            Correctness: repo.Correctness,
            ResponsiveMaintainer: repo.ResponsiveMaintainer,
            BusFactor: repo.BusFactor,
            ReviewCoverage: repo.ReviewCoverage,
            DependancyPinning: repo.DependancyPinning,
            License: repo.License,
            Net: repo.Net,
        }

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