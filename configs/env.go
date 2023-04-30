package configs

import (
	// "log"
	"context"
	"fmt"
	//"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
	// "github.com/joho/godotenv"
)


func getSecret(projectID, secretName string) (string, error) {
    ctx := context.Background()
    client, err := secretmanager.NewClient(ctx, option.WithUserAgent("my-app/0.1"))
    if err != nil {
        return "", err
    }
    defer client.Close()

    req := &secretmanagerpb.AccessSecretVersionRequest{
        Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName),
    }

    result, err := client.AccessSecretVersion(ctx, req)
    if err != nil {
        return "", err
    }

    return string(result.Payload.Data), nil
}

func EnvMongoURI() string {
	// err := godotenv.Load()
	// if err != nil {
	//     log.Fatal("Error loading .env file")
	// }
	//projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	//secretName := "MONGOURI"
	//mongouri, _ := getSecret(projectID, secretName)
	// return mongouri
	return "mongodb+srv://ece461team7_new:uQaUIvLWcMbgujtt@cluster0.a747pma.mongodb.net/?retryWrites=true&w=majority"
}
