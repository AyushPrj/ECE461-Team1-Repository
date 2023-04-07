package controllers

import (
	"ECE461-Team1-Repository/configs"
	"ECE461-Team1-Repository/metrics"
	models "ECE461-Team1-Repository/models"
	"context"
	"encoding/json"
	"fmt"

	//"fmt"
	"regexp"
	"strings"

	//"io"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson" //add this
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var repoCollection *mongo.Collection = configs.GetCollection(configs.DB, "repos")
var contentCollection *mongo.Collection = configs.GetCollection(configs.DB, "largeStrings")

func CreateAuthToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func PackageByNameDelete(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("X-Authorization")
	if authToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Missing X-Authorization header",
		})
		return
	}

	// Implement your authentication and authorization logic based on authToken

	packageName := r.URL.Query().Get("name")
	if packageName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "There is missing field(s) in the PackageName/AuthenticationToken\\ or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	filter := bson.M{"name": packageName}
	result, err := repoCollection.DeleteMany(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "There is missing field(s) in the PackageName/AuthenticationToken\\ or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(Response{
			Status:  "error",
			Message: "Package does not exist.",
		})
		return
	}

	response := Response{
		Status:  "success",
		Message: "Package is deleted.",
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func PackageByNameGet(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("X-Authorization")
	if authToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing X-Authorization header",
		})
		return
	}

	// Implement your authentication and authorization logic based on authToken

	packageName := r.URL.Query().Get("name")
	if packageName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing package name in the path",
		})
		return
	}

	filter := bson.M{"name": packageName}
	findOptions := options.Find()
	var packageHistory []models.PackageHistoryEntry

	cursor, err := repoCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error retrieving package history",
		})
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &packageHistory); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error decoding package history",
		})
		return
	}

	if len(packageHistory) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Package not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(packageHistory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type PackageRegExRequest struct {
	RegEx string `json:"regex"`
}

func PackageByRegExGet(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("X-Authorization")
	if authToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing X-Authorization header",
		})
		return
	}

	// Implement your authentication and authorization logic based on authToken

	var packageRegExRequest PackageRegExRequest
	err := json.NewDecoder(r.Body).Decode(&packageRegExRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid JSON in request body",
		})
		return
	}

	findOptions := options.Find()
	var packages []models.PackageMetadata
	cursor, err := repoCollection.Find(context.Background(), bson.D{}, findOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error retrieving packages",
		})
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &packages); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error decoding packages",
		})
		return
	}

	regex, err := regexp.Compile(packageRegExRequest.RegEx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid regular expression",
		})
		return
	}

	matchingPackages := []models.PackageMetadata{}
	for _, pkg := range packages {
		if regex.MatchString(pkg.Name) && regex.MatchString(pkg.Version) {
			matchingPackages = append(matchingPackages, pkg)
		}
	}

	if len(matchingPackages) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(map[string]string{
			"error": "No packages found for the provided regex",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(matchingPackages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func PackageCreate(w http.ResponseWriter, r *http.Request) {
	// Get the authentication token from the request header
	authToken := r.Header.Get("X-Authorization")
	if authToken == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusBadRequest,
			Message: " is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid.",
		})
		return
	}

	// Decode the request body into a ModelPackage struct
	//var modelPackage models.ModelPackage
	var packageData models.PackageData
	err := json.NewDecoder(r.Body).Decode(&packageData)
	if err != nil {
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

func PackageDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

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
	if err != nil{
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
}

func PackageRetrieve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PackageUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PackagesList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func RegistryReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
