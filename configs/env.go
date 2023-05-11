package configs

import (
	// "log"
	"context"
	"fmt"

	// "log"
	"os"

	//"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"

	// "github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var GITHUB_TOKEN string

func GetGithubUserToken() string {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	secretName := "GITHUB_TOKEN"
	GITHUB_TOKEN, _ := GetSecret(projectID, secretName)
	if GITHUB_TOKEN == "" {
		GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")
	}
	return GITHUB_TOKEN
}

func GetSecret(projectID, secretName string) (string, error) {
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
	
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	secretName := "MONGOURI"
	mongouri, _ := GetSecret(projectID, secretName)
	if mongouri == "" {
		mongouri = os.Getenv("MONGOURI")
	}
	return mongouri
}
