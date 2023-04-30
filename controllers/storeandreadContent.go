package controllers

import (
	"bytes"
	"fmt"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func storeLargeString(contentCollection *mongo.Collection, largeString string) (primitive.ObjectID, error) {
	// Convert the string to a byte array
	content := []byte(largeString)

	// Create a GridFS bucket with the default name ("fs")
	bucket, err := gridfs.NewBucket(
		contentCollection.Database(),
	)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// Create a bytes.Reader from the byte array
	contentReader := bytes.NewReader(content)

	// Upload the content to GridFS
	fileID := primitive.NewObjectID()
	uploadStream, err := bucket.OpenUploadStreamWithID(fileID, "")
	if err != nil {
		return primitive.NilObjectID, err
	}
	defer uploadStream.Close()

	_, err = io.Copy(uploadStream, contentReader)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return fileID, nil
}

func readLargeString(contentCollection *mongo.Collection, fileIDStr string) (string, error) {
	// Convert the string to a primitive.ObjectID
	fileID, err := primitive.ObjectIDFromHex(fileIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid fileID: %v", err)
	}

	// Create a GridFS bucket with the default name ("fs")
	bucket, err := gridfs.NewBucket(
		contentCollection.Database(),
	)
	if err != nil {
		return "", err
	}

	// Open the download stream for the given fileID
	downloadStream, err := bucket.OpenDownloadStream(fileID)
	if err != nil {
		return "", err
	}
	defer downloadStream.Close()

	// Create a buffer to store the content
	var contentBuffer bytes.Buffer

	// Read the content from GridFS
	_, err = io.Copy(&contentBuffer, downloadStream)
	if err != nil {
		return "", err
	}

	// Convert the content buffer to a string
	largeString := contentBuffer.String()

	return largeString, nil
}
