package controllers

import (
	"ECE461-Team1-Repository/configs"
	"ECE461-Team1-Repository/metrics"
	models "ECE461-Team1-Repository/models"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	//"fmt"
	// "regexp"
	"strings"

	//"io"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson" //add this
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var repoCollection *mongo.Collection = configs.GetCollection(configs.DB, "repos")
var contentCollection *mongo.Collection = configs.GetCollection(configs.DB, "largeStrings")
var historyCollection *mongo.Collection = configs.GetCollection(configs.DB, "history")
var fschunksCollection *mongo.Collection = configs.GetCollection(configs.DB, "fs.chunks")
var fsfilesCollection *mongo.Collection = configs.GetCollection(configs.DB, "fs.files")

func CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var requestBody models.AuthenticationRequest
	var defaultUser models.User
	defaultUser.IsAdmin = true
	defaultUser.Name = "ece30861defaultadminuser"
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody == (models.AuthenticationRequest{}) || *requestBody.User != defaultUser {
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	// Have not implemented
	// resource - https://mattermost.com/blog/how-to-build-an-authentication-microservice-in-golang-from-scratch/

	w.WriteHeader(http.StatusNotImplemented)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(models.ModelError{
		Code:    501,
		Message: "This system does not support authentication.",
	})
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// done
func PackageByNameDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	packageName := vars["name"]
	if packageName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	filter := bson.M{"metadata.name": packageName}

	cur, err := repoCollection.Find(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	var packageIDs []string
	for cur.Next(context.Background()) {
		var pkg models.PkgResponse
		err := cur.Decode(&pkg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(models.ModelError{
				Code:    500,
				Message: "An error occurred while decoding package data.",
			})
			return
		}
		packageIDs = append(packageIDs, pkg.Data.Content)
	}

	// Check if any packages were found
	if len(packageIDs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    404,
			Message: "Package does not exist.",
		})
		return
	}

	bucket, _ := gridfs.NewBucket(
		contentCollection.Database(),
		options.GridFSBucket().SetName("fs"),
	)

	for _, id := range packageIDs {
		fmt.Println(id)
		id, err := primitive.ObjectIDFromHex(id)
		if err := bucket.Delete(id); err != nil {
			panic(err)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ModelError{
				Code:    http.StatusInternalServerError,
				Message: "An error occurred while deleting the associated GridFS files and chunks.",
			})
			return
		}

	}

	_, err = repoCollection.DeleteMany(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	response := models.ModelError{
		Code:    200,
		Message: "Package is deleted.",
	}
	historyfilter := bson.M{"packageMetadata.name": packageName}
	_, err = historyCollection.DeleteMany(context.Background(), historyfilter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// done... how do i get a 400 error?
func PackageByNameGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	resourceName := vars["name"]

	if resourceName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	packageFilter := bson.M{"metadata.name": resourceName}
	packageFilter2 := bson.D{} //"packagemetadata.name": "axios"

	packageCount, err := repoCollection.CountDocuments(context.Background(), packageFilter)
	numFound, err := historyCollection.CountDocuments(context.Background(), packageFilter2)
	fmt.Println(numFound)

	if err != nil || packageCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusNotFound,
			Message: "No such package.",
		})
		return
	}

	filter := bson.M{"packagemetadata.name": resourceName}
	findOptions := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	cur, err := historyCollection.Find(context.Background(), filter, findOptions)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    0,
			Message: "unexpected error",
		})
		return
	}

	var results []models.PackageHistoryEntry

	// Decode the results into a slice of PackageHistory
	for cur.Next(context.Background()) {
		var result models.PackageHistoryEntry
		err := cur.Decode(&result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ModelError{
				Code:    0,
				Message: "unexpected error",
			})
			return
		}
		results = append(results, result)
	}

	// Return the list of PackageHistory objects as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)

}

type PackageRegExRequest struct {
	RegEx string `json:"regex"`
}

// done
func PackageByRegExGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Read the request body and store it as a regex pattern
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	regexPattern := string(body)
	if regexPattern == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	// Find packages based on regex pattern
	filter := bson.M{"metadata.name": bson.M{"$regex": primitive.Regex{Pattern: regexPattern, Options: ""}}}
	cur, err := repoCollection.Find(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    0,
			Message: "Unexpected error",
		})
		return
	}

	type PackageVersionName struct {
		Name    string `json:"Name"`
		Version string `json:"Version"`
	}

	var results []PackageVersionName

	// Decode the results into a slice of PackageVersionName
	for cur.Next(context.Background()) {
		var pkg models.PkgResponse
		err := cur.Decode(&pkg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(models.ModelError{
				Code:    0,
				Message: "Unexpected error",
			})
			return
		}
		results = append(results, PackageVersionName{
			Name:    pkg.Metadata.Name,
			Version: pkg.Metadata.Version,
		})
	}

	if len(results) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusNotFound,
			Message: "No package found under this regex.",
		})
		return
	}

	// Return the list of PackageVersionName objects as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

// done.. dont need auth?
func PackageCreate(w http.ResponseWriter, r *http.Request) {
	// Get the authentication token from the request header
	// authToken := r.Header.Get("X-Authorization")
	// if authToken == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Println("here1")
	// 	json.NewEncoder(w).Encode(models.ModelError{
	// 		Code:    http.StatusBadRequest,
	// 		Message: " is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
	// 	})
	// 	return
	// }

	// Decode the request body into a ModelPackage struct
	//var modelPackage models.ModelPackage
	var packageData models.PackageData
	err := json.NewDecoder(r.Body).Decode(&packageData)
	if err != nil {
		fmt.Println("here2")
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}
	var newMetadata models.PackageMetadata

	repoPath := strings.TrimPrefix(packageData.URL, "https://github.com/")
	username, repoName := path.Split(repoPath)
	repoPath = strings.TrimSuffix(username+repoName, "/")

	newMetadata.Name = path.Base(repoPath)

	ver, e := extractVersionFromZip(packageData.Content)
	largeString := packageData.Content
	fileID, err := storeLargeString(contentCollection, largeString)
	if err != nil {
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}
	packageData.Content = fileID.Hex()

	if e {
		newMetadata.Version = ver
	} else {
		fmt.Println("here4")
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	modelPackage := models.ModelPackage{
		ID:       primitive.NewObjectID(),
		Data:     &packageData,
		Metadata: &newMetadata,
	}
	newMetadata.ID = modelPackage.ID.Hex()

	if !packageExists(modelPackage.Metadata.Name, modelPackage.Metadata.Version) {

		new_metrics := metrics.GetMetrics("https://github.com", 1, repoPath) //get metrics

		var ndjsonData models.PackageRating

		err := json.Unmarshal([]byte(new_metrics), &ndjsonData)
		if err != nil {
			return
		}

		if ndjsonData.NetScore > 0.5 {
			//insert
			if _, err := repoCollection.InsertOne(context.Background(), modelPackage); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.ModelError{
					Code:    http.StatusInternalServerError,
					Message: "Failed to create package",
				})
				return
			}
			// Return the created package metadata

			ls, err := readLargeString(contentCollection, modelPackage.Data.Content)
			if err != nil {
				return
			}
			modelPackage.Data.Content = ls
			resp := &models.PkgResponse{
				Metadata: modelPackage.Metadata,
				Data:     modelPackage.Data,
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(resp)
			AddPackageHistory(*resp.Metadata, "CREATE")
			return
		} else {
			json.NewEncoder(w).Encode(models.ModelError{
				Code:    424,
				Message: "Package is not uploaded due to the disqualified rating.",
			})
			return
		}
	} else {
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    409,
			Message: "Package already exists",
		})
		return
	}
	// Generate a unique package ID and store the package in the database
}

// done
func packageExists(name string, version string) bool { //other page?
	// Create a filter for the query
	filter := bson.M{
		"metadata.name":    name,
		"metadata.version": version,
	}

	// Find a document in the collection that matches the filter
	result := repoCollection.FindOne(context.Background(), filter)

	// Check if a document was found

	return result.Err() == nil
}

// done
func PackageDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	resourceID := vars["id"]
	objectId, err := primitive.ObjectIDFromHex(resourceID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusBadRequest,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	var temp models.ModelPackage

	err = repoCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&temp)
	if err != nil {
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    404,
			Message: "Package does not exist.",
		})
		return
	}

	bucket, _ := gridfs.NewBucket(
		contentCollection.Database(),
		options.GridFSBucket().SetName("fs"),
	)

	id, err := primitive.ObjectIDFromHex(temp.Data.Content)
	if err := bucket.Delete(id); err != nil {
		panic(err)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusInternalServerError,
			Message: "An error occurred while deleting the associated GridFS files and chunks.",
		})
		return
	}

	result, err := repoCollection.DeleteOne(context.Background(), bson.M{"_id": objectId})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusInternalServerError,
			Message: "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusNotFound,
			Message: "Package does not exist.",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ModelError{
		Code:    http.StatusOK,
		Message: "Package is deleted.",
	})

}

// done... dont need auth?
func PackageRate(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	authToken := r.Header.Get("X-Authorization")
	if authToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusBadRequest,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	vars := mux.Vars(r)
	resourceID := vars["id"]
	objectId, err := primitive.ObjectIDFromHex(resourceID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusBadRequest,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	var result models.ModelPackage

	err = repoCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&result)
	if err != nil {
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    404,
			Message: "Package does not exist.",
		})
		return
	}
	fmt.Printf("%+v\n", result)

	repoPath := strings.TrimPrefix(result.Data.URL, "https://github.com/")
	username, repoName := path.Split(repoPath)
	repoPath = strings.TrimSuffix(username+repoName, "/")
	new_metrics := metrics.GetMetrics("https://github.com", 1, repoPath)

	var ndjsondata models.PackageRating

	err = json.Unmarshal([]byte(new_metrics), &ndjsondata)
	if err != nil {
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    500,
			Message: "The package rating system choked on at least one of the metrics.",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ndjsondata)
	AddPackageHistory(*result.Metadata, "RATE")
}

// done
func PackageRetrieve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	resourceID := vars["id"]
	objectId, err := primitive.ObjectIDFromHex(resourceID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusBadRequest,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	var result models.ModelPackage

	err = repoCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&result)
	if err != nil {
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    404,
			Message: "Package does not exist.",
		})
		return
	}

	ls, _ := readLargeString(contentCollection, result.Data.Content)

	result.Data.Content = ls

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	AddPackageHistory(*result.Metadata, "DOWNLOAD")
}

// done
func PackageUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	resourceID := vars["id"]
	objectId, err := primitive.ObjectIDFromHex(resourceID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusBadRequest,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	var result models.ModelPackage

	err = repoCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&result)
	if err != nil {
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    404,
			Message: "Package does not exist.",
		})
		return
	}

	// Decode the request body into a ModelPackage struct
	var updatedPackage models.ModelPackage
	err = json.NewDecoder(r.Body).Decode(&updatedPackage)
	fmt.Println(updatedPackage.Data.JSProgram)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusBadRequest,
			Message: "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	//if(result.Metadata.Name != updatedPackage.Metadata.Name || result.Metadata.Version != updatedPackage.Metadata.Version){
	if result.Metadata.Name != updatedPackage.Metadata.Name {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusBadRequest,
			Message: "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	largeString := updatedPackage.Data.Content
	fileID, err := storeLargeString(contentCollection, largeString)
	if err != nil {
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    400,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	// Update the fields of the result with the values from the request body
	result.Data.JSProgram = updatedPackage.Data.JSProgram
	fmt.Println(result.Data.JSProgram)
	result.Data.URL = updatedPackage.Data.URL

	oldContentID := result.Data.Content
	result.Data.Content = fileID.Hex()
	// Add other fields as needed

	// Update the package in the MongoDB collection
	updateResult, err := repoCollection.UpdateOne(context.Background(), bson.M{"_id": objectId}, bson.M{"$set": result})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusInternalServerError,
			Message: "An error occurred while updating the package.",
		})
		return
	}

	if updateResult.ModifiedCount == 0 { //could mean that nothing would have gotten updated
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusBadRequest,
			Message: "There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	bucket, _ := gridfs.NewBucket(
		contentCollection.Database(),
		options.GridFSBucket().SetName("fs"),
	)

	id, err := primitive.ObjectIDFromHex(oldContentID)
	if err := bucket.Delete(id); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ModelError{
		Code:    http.StatusOK,
		Message: "Version is updated.",
	})

	AddPackageHistory(*updatedPackage.Metadata, "UPDATE")
}

// Not done the filter for the database might have to be parsed
func PackagesList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	type requestBody struct {
		Version string `json:"Version"`
		Name    string `json:"Name"`
	}

	var search []requestBody
	var results []requestBody

	// Decode the results into a slice of PackageVersionName
	err := json.NewDecoder(r.Body).Decode(&search)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusBadRequest,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	filter := bson.M{
		"metadata.name":    search[0].Name,
		"metadata.version": search[0].Version,
	}

	cur, err := repoCollection.Find(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    0,
			Message: "Unexpected error",
		})
		return
	}
	for cur.Next(context.Background()) {
		var pkg requestBody
		err := cur.Decode(&pkg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(models.ModelError{
				Code:    500,
				Message: "An error occurred while decoding package data.",
			})
			return
		}
		results = append(results, pkg)
	}

	json.NewEncoder(w).Encode(results)
	w.WriteHeader(http.StatusOK)
}

// Check user ? when to return 401
func RegistryReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	authToken := r.Header.Get("X-Authorization")
	if authToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusBadRequest,
			Message: "There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	repoCollection.Drop(context.Background())
	contentCollection.Drop(context.Background())
	historyCollection.Drop(context.Background())
	fsfilesCollection.Drop(context.Background())
	fschunksCollection.Drop(context.Background())

	w.WriteHeader(http.StatusOK)
}

func AddPackageHistory(metadata models.PackageMetadata, action string) error {
	now := time.Now().UTC()
	formattedDate := now.Format("2006-01-02T15:04:05Z")


	hardcodedUser := &models.User{
		Name:    "ece30861defaultadminuser",
		IsAdmin: true,
	}

	history := models.PackageHistoryEntry{
		User:           hardcodedUser,
		Date:           formattedDate,
		PackageMetadata: &metadata,
		Action:          action,
	}

	_, err := historyCollection.InsertOne(context.Background(), history)
	return err
}
