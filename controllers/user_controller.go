package controllers

import (
	"ECE461-Team1-Repository/configs"
	"ECE461-Team1-Repository/metrics"
	models "ECE461-Team1-Repository/models"
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
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
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the AuthenticationRequest or it is formed improperly."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	// Have not implemented
	// resource - https://mattermost.com/blog/how-to-build-an-authentication-microservice-in-golang-from-scratch/

	w.WriteHeader(http.StatusNotImplemented)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte("This system does not support authentication."))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
	return
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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	filter := bson.M{"metadata.name": packageName}

	cur, err := repoCollection.Find(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	var packageIDs []string
	for cur.Next(context.Background()) {
		var pkg models.PkgResponse
		err := cur.Decode(&pkg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("error decoding..."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}
		packageIDs = append(packageIDs, pkg.Data.Content)
	}

	// Check if any packages were found
	if len(packageIDs) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
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
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("error decoding 2..."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}

	}

	_, err = repoCollection.DeleteMany(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	historyfilter := bson.M{"packageMetadata.name": packageName}
	_, err = historyCollection.DeleteMany(context.Background(), historyfilter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageName/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte("Package is deleted."))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}

// done... how do i get a 400 error?
func PackageByNameGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	resourceName := vars["name"]

	if resourceName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageQuery/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	packageFilter := bson.M{"metadata.name": resourceName}
	packageFilter2 := bson.D{} //"packagemetadata.name": "axios"

	packageCount, err := repoCollection.CountDocuments(context.Background(), packageFilter)
	if err != nil || packageCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("No such package."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}
	numFound, err := historyCollection.CountDocuments(context.Background(), packageFilter2)
	fmt.Println(numFound)

	if err != nil || packageCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("No such package."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	filter := bson.M{"packagemetadata.name": resourceName}
	findOptions := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	cur, err := historyCollection.Find(context.Background(), filter, findOptions)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("Unexpected error"))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	var results []models.PackageHistoryEntry

	// Decode the results into a slice of PackageHistory
	for cur.Next(context.Background()) {
		var result models.PackageHistoryEntry
		err := cur.Decode(&result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("Unexpected error"))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}
		results = append(results, result)
	}

	// Return the list of PackageHistory objects as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)

}

type PackageRegExRequest struct {
	RegEx string `json:"RegEx"`
}

// done
func PackageByRegExGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Read the request body and store it as a regex pattern
	body, err := ioutil.ReadAll(r.Body)
	body = bytes.ReplaceAll(body, []byte(`\`), []byte(`\\`))
	// fmt.Println(string(body))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	var inputBody PackageRegExRequest
	if err := json.Unmarshal(body, &inputBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	regexPattern := inputBody.RegEx

	if regexPattern == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	// Find packages based on regex pattern (name)
	firstFilterNames := []string{}
	firstFilterVersions := []string{}

	filter := bson.M{"metadata.name": bson.M{"$regex": primitive.Regex{Pattern: regexPattern, Options: ""}}}
	cur, err := repoCollection.Find(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
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
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("Unexpected error"))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}
		results = append(results, PackageVersionName{
			Name:    pkg.Metadata.Name,
			Version: pkg.Metadata.Version,
		})
		firstFilterNames = append(firstFilterNames, pkg.Metadata.Name)
		firstFilterVersions = append(firstFilterVersions, pkg.Metadata.Version)
	}

	defer cur.Close(context.Background())

	// Find packages based on regex pattern (readme)
	secondFilter := bson.M{
		"$and": []bson.M{
			{"metadata.name": bson.M{"$nin": firstFilterNames}},
			{"metadata.version": bson.M{"$nin": firstFilterVersions}},
		},
	}

	cur, err = repoCollection.Find(context.Background(), secondFilter)
	if err != nil {
		fmt.Printf("Error occurred while querying: %v\n", err)
		// Handle the error accordingly
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	for cur.Next(context.Background()) {
		var pkg models.PkgResponse
		err := cur.Decode(&pkg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}

		fsFileID := pkg.Data.Content

		content, err := readLargeString(contentCollection, fsFileID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}

		// Decode the base64 encoded zip file
		decodedZipFile, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}

		// Create a bytes reader for the zip file
		zipReader, err := zip.NewReader(bytes.NewReader(decodedZipFile), int64(len(decodedZipFile)))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}

		folderName := zipReader.File[0].Name

		// Iterate through the files in the zip archive
		for _, file := range zipReader.File {
			// Check if the file is a README (case insensitive) within the inner folder
			if strings.HasPrefix(file.Name, folderName) && (strings.EqualFold(file.Name, folderName+"readme") || strings.EqualFold(file.Name, folderName+"readme.txt") || strings.EqualFold(file.Name, folderName+"readme.md") || strings.EqualFold(file.Name, folderName+"README.md")) {
				// Open the README file
				readmeFile, err := file.Open()
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
					if err != nil {
						fmt.Println("Error writing response:", err)
					}
					return
				}
				defer readmeFile.Close()

				// Read the README file content
				readmeContent, err := ioutil.ReadAll(readmeFile)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					_, err := w.Write([]byte("There is missing field(s) in the PackageRegEx/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
					if err != nil {
						fmt.Println("Error writing response:", err)
					}
					return
				}

				// Check if the regex pattern matches the README content
				if regex.Match(readmeContent) {
					// The regex pattern was found in the README content
					// Add this package to the results
					results = append(results, PackageVersionName{
						Name:    pkg.Metadata.Name,
						Version: pkg.Metadata.Version,
					})
				}

				// Since we found the README file, no need to check other files in the zip
				break
			}
		}

	}

	if len(results) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("No package found under this regex."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	// Return the list of PackageVersionName objects as JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

// done.. dont need auth?
func PackageCreate(w http.ResponseWriter, r *http.Request) {
	var packageData models.PackageData
	err := json.NewDecoder(r.Body).Decode(&packageData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
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
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}
	packageData.Content = fileID.Hex()

	if e {
		newMetadata.Version = ver
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageData/AuthenticationToken or it is formed improperly (e.g. Content and URL are both set), or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
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
		ndjsonData.NetScore = 0.6; //hardcoded for now
		if ndjsonData.NetScore > 0.5 {
			//insert
			if _, err := repoCollection.InsertOne(context.Background(), modelPackage); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				_, err := w.Write([]byte("Unexpected error"))
				if err != nil {
					fmt.Println("Error writing response:", err)
				}
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
			if err := AddPackageHistory(*resp.Metadata, "CREATE"); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				_, err := w.Write([]byte("Unexpected error"))
				if err != nil {
					fmt.Println("Error writing response:", err)
				}
				return
			}
			w.WriteHeader(http.StatusCreated)
			// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			// _, err = w.Write([]byte("Success. Check the ID in the returned metadata for the official ID."))
			// if err != nil {
			// 	fmt.Println("Error writing response:", err)
			// }
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			w.WriteHeader(http.StatusFailedDependency)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("Package is not uploaded due to the disqualified rating."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("Package exists already."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	var temp models.ModelPackage

	err = repoCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&temp)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err = w.Write([]byte("Package does not exist."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
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
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	if result.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err = w.Write([]byte("Package does not exist."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte("Package is deleted."))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}

// done... dont need auth?
func PackageRate(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// authToken := r.Header.Get("X-Authorization")
	// if authToken == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// 	_, err := w.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	// 	if err != nil {
	// 		fmt.Println("Error writing response:", err)
	// 	}
	// 	return
	// }

	vars := mux.Vars(r)
	resourceID := vars["id"]
	objectId, err := primitive.ObjectIDFromHex(resourceID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	var result models.ModelPackage

	err = repoCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("Package does not exist."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("The package rating system choked on at least one of the metrics."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	if err := AddPackageHistory(*result.Metadata, "RATE"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("The package rating system choked on at least one of the metrics."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ndjsondata)
}

// done
func PackageRetrieve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	resourceID := vars["id"]
	objectId, err := primitive.ObjectIDFromHex(resourceID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	var result models.PkgResponse

	err = repoCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("Package does not exist."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	ls, _ := readLargeString(contentCollection, result.Data.Content)

	result.Data.Content = ls

	if err := AddPackageHistory(*result.Metadata, "DOWNLOAD"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("Failed to upload history..."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// done
func PackageUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	resourceID := vars["id"]
	objectId, err := primitive.ObjectIDFromHex(resourceID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	var result models.ModelPackage

	err = repoCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("Package does not exist."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	// Decode the request body into a ModelPackage struct
	var updatedPackage models.ModelPackage
	err = json.NewDecoder(r.Body).Decode(&updatedPackage)
	fmt.Println(updatedPackage.Data.JSProgram)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	//if(result.Metadata.Name != updatedPackage.Metadata.Name || result.Metadata.Version != updatedPackage.Metadata.Version){
	if result.Metadata.Name != updatedPackage.Metadata.Name {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	largeString := updatedPackage.Data.Content
	fileID, err := storeLargeString(contentCollection, largeString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
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
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("failed to update..."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	if updateResult.ModifiedCount == 0 { //could mean that nothing would have gotten updated
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageID/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	bucket, _ := gridfs.NewBucket(
		contentCollection.Database(),
		options.GridFSBucket().SetName("fs"),
	)

	id, err := primitive.ObjectIDFromHex(oldContentID)
	if err != nil {
		panic(err)
	}
	if err := bucket.Delete(id); err != nil {
		panic(err)
	}

	if err := AddPackageHistory(*updatedPackage.Metadata, "UPDATE"); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("failed to add pacakge history..."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err = w.Write([]byte("Version is updated."))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}

}

// where to return error 413??
func PackagesList(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	offset := r.Header.Get("offset")
	if offset == "" || offset == "0" {
		offset = "1"
	}

	offsetNum, _ := strconv.Atoi(offset)
	offsetNum = (offsetNum - 1) * 10

	type packageResponse struct {
		Version string `json:"Version"`
		Name    string `json:"Name"`
		ID      string `json:"ID"`
	}
	var search []models.PackageQuery
	var results []packageResponse

	// Decode the results into a slice of PackageVersionName
	err := json.NewDecoder(r.Body).Decode(&search)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageQuery/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	stringmap := extractVersionRanges(search[0].Version)

	var filter, filter1, filter2, filter3 bson.M
	if search[0].Name != "*" {
		filter = bson.M{
			"metadata.name":    search[0].Name,
			"metadata.version": stringmap.Exact,
		}
		filter1 = bson.M{
			"metadata.name": search[0].Name,
			"metadata.version": bson.M{
				"$gte": stringmap.BoundedRange[0],
				"$lte": stringmap.BoundedRange[1],
			},
		}
		filter2 = bson.M{
			"metadata.name": search[0].Name,
			"metadata.version": bson.M{
				"$gte": stringmap.Caret,
				"$lte": getCaretUpperBound(stringmap.Caret),
			},
		}
		filter3 = bson.M{
			"metadata.name": search[0].Name,
			"metadata.version": bson.M{
				"$gte": stringmap.Tilde,
				"$lte": geTildeUpperBound(stringmap.Tilde),
			},
		}
	} else {
		filter = bson.M{
			"metadata.version": stringmap.Exact,
		}
		filter1 = bson.M{
			"metadata.version": bson.M{
				"$gte": stringmap.BoundedRange[0],
				"$lte": stringmap.BoundedRange[1],
			},
		}
		filter2 = bson.M{
			"metadata.version": bson.M{
				"$gte": stringmap.Caret,
				"$lte": getCaretUpperBound(stringmap.Caret),
			},
		}
		filter3 = bson.M{
			"metadata.version": bson.M{
				"$gte": stringmap.Tilde,
				"$lte": geTildeUpperBound(stringmap.Tilde),
			},
		}
	}

	cur, err := repoCollection.Find(context.Background(), filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("unexpected cur"))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}
	cur1, err := repoCollection.Find(context.Background(), filter1)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("unexpected cur1"))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}
	cur2, err := repoCollection.Find(context.Background(), filter2)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("unexpected cur2"))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}
	cur3, err := repoCollection.Find(context.Background(), filter3)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("unexpected cur3"))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	for cur.Next(context.Background()) {
		var pkg packageResponse
		var myMap map[string]interface{}
		err := cur.Decode(&myMap)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("error decoding"))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}

		pkg.Version = stringmap.Exact
		pkg.Name = myMap["metadata"].(map[string]interface{})["name"].(string)
		pkg.ID = myMap["metadata"].(map[string]interface{})["id"].(string)
		if pkg.Name == "" || pkg.Version == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("error decoding 2."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}

		results = append(results, pkg)
	}
	for cur1.Next(context.Background()) {

		var pkg packageResponse
		var myMap map[string]interface{}
		err := cur1.Decode(&myMap)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("error decoding 3."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}

		pkg.Version = stringmap.BoundedRange[0] + "-" + stringmap.BoundedRange[1]
		pkg.Name = myMap["metadata"].(map[string]interface{})["name"].(string)
		pkg.ID = myMap["metadata"].(map[string]interface{})["id"].(string)
		if pkg.Name == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("error decoding 4."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}

		results = append(results, pkg)
	}

	for cur2.Next(context.Background()) {
		var pkg packageResponse
		var myMap map[string]interface{}
		err := cur2.Decode(&myMap)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("error decoding 5."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}

		pkg.Version = "^" + stringmap.Caret
		pkg.Name = myMap["metadata"].(map[string]interface{})["name"].(string)
		pkg.ID = myMap["metadata"].(map[string]interface{})["id"].(string)
		if pkg.Name == "" || pkg.Version == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("error decoding 6."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}
		results = append(results, pkg)
	}
	for cur3.Next(context.Background()) {
		var pkg packageResponse
		var myMap map[string]interface{}
		err := cur3.Decode(&myMap)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("error decoding 7."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}

		pkg.Version = "~" + stringmap.Tilde
		pkg.Name = myMap["metadata"].(map[string]interface{})["name"].(string)
		pkg.ID = myMap["metadata"].(map[string]interface{})["id"].(string)
		if pkg.Name == "" || pkg.Version == "" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte("error decoding 8."))
			if err != nil {
				fmt.Println("Error writing response:", err)
			}
			return
		}
		results = append(results, pkg)
	}

	if len(results) == 0 || len(results) < offsetNum*10 {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, err := w.Write([]byte("There is missing field(s) in the PackageQuery/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}
	if len(results) < 10 || (len(results) < (offsetNum+1)*10) && (len(results) > offsetNum*10) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results[offsetNum:])
		return
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results[offsetNum : offsetNum+10])
		return
	}
}

func geTildeUpperBound(version string) string {
	segments := strings.Split(version, ".")

	major, err := strconv.Atoi(segments[0])
	if err != nil {
		return ""
	}

	minor, err := strconv.Atoi(segments[1])
	if err != nil {
		return ""
	}

	// Increment minor version and reset patch version
	upperBound := fmt.Sprintf("%d.%d.0", major, minor+1)

	return upperBound
}

func getCaretUpperBound(version string) string {
	segments := strings.Split(version, ".")

	major, err := strconv.Atoi(segments[0])
	if err != nil {
		return ""
	}

	// Increment major version and reset minor and patch versions
	upperBound := fmt.Sprintf("%d.0.0", major+1)

	return upperBound
}

type VersionRanges struct {
	Exact        string
	BoundedRange []string
	Caret        string
	Tilde        string
}

// extracts the strings given a request body that contains Exact, BoundedRange, Caret, and Tilde version ranges
func extractVersionRanges(versionString string) VersionRanges {
	versionPattern := `\d+\.\d+\.\d+`
	exactPattern := `Exact \((` + versionPattern + `)\)?`
	boundedRangePattern := `Bounded range \((` + versionPattern + `)-(` + versionPattern + `)\)?`
	caretPattern := `Carat \(\^(` + versionPattern + `)\)?`
	tildePattern := `Tilde \(~(` + versionPattern + `)\)?`

	exactRegexp := regexp.MustCompile(exactPattern)
	boundedRangeRegexp := regexp.MustCompile(boundedRangePattern)
	caretRegexp := regexp.MustCompile(caretPattern)
	tildeRegexp := regexp.MustCompile(tildePattern)

	exactMatches := exactRegexp.FindStringSubmatch(versionString)
	boundedRangeMatches := boundedRangeRegexp.FindStringSubmatch(versionString)
	caretMatches := caretRegexp.FindStringSubmatch(versionString)
	tildeMatches := tildeRegexp.FindStringSubmatch(versionString)

	versionRanges := VersionRanges{
		Exact:        "",
		BoundedRange: []string{"", ""},
		Caret:        "",
		Tilde:        "",
	}

	if len(exactMatches) > 0 {
		versionRanges.Exact = string(exactMatches[1])
	}
	if len(boundedRangeMatches) > 0 {
		versionRanges.BoundedRange = boundedRangeMatches[1:]
	}
	if len(caretMatches) > 0 {
		versionRanges.Caret = caretMatches[1]
	}
	if len(tildeMatches) > 0 {
		versionRanges.Tilde = tildeMatches[1]
	}

	return versionRanges
}

// Check user ? when to return 401
func RegistryReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// authToken := r.Header.Get("X-Authorization")
	// if authToken == "" {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// 	_, err := w.Write([]byte("There is missing field(s) in the AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
	// 	if err != nil {
	// 		fmt.Println("Error writing response:", err)
	// 	}
	// 	return
	// }

	if err := repoCollection.Drop(context.Background()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusInternalServerError,
			Message: "An error occurred while dropping the repo collection.",
		})
	}

	if err := contentCollection.Drop(context.Background()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusInternalServerError,
			Message: "An error occurred while dropping the content collection.",
		})
	}
	if err := historyCollection.Drop(context.Background()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusInternalServerError,
			Message: "An error occurred while dropping the history collection.",
		})
	}
	if err := fsfilesCollection.Drop(context.Background()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusInternalServerError,
			Message: "An error occurred while dropping the fs.files collection.",
		})
	}
	if err := fschunksCollection.Drop(context.Background()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ModelError{
			Code:    http.StatusInternalServerError,
			Message: "An error occurred while dropping the fs.chunks collection.",
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := w.Write([]byte("Registry is reset."))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
	return
}

func AddPackageHistory(metadata models.PackageMetadata, action string) error {
	now := time.Now().UTC()
	formattedDate := now.Format("2006-01-02T15:04:05Z")

	hardcodedUser := &models.User{
		Name:    "ece30861defaultadminuser",
		IsAdmin: true,
	}

	history := models.PackageHistoryEntry{
		User:            hardcodedUser,
		Date:            formattedDate,
		PackageMetadata: &metadata,
		Action:          action,
	}

	_, err := historyCollection.InsertOne(context.Background(), history)
	return err
}
