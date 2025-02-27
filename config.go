package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jake-abed/noisecheck-api/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
	Db            *database.Queries
	S3Bucket      string
	S3Client      *s3.Client
	S3Region      string
	CloudfrontUrl string
	ClerkKey      string
	TursoToken    string
	TursoUrl      string
}

func createApiConfig() apiConfig {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	tursoUrl := os.Getenv("TURSO_DATABASE_URL")
	if tursoUrl == "" {
		log.Fatal("TURSO_DATABASE_URL environment variable is not set")
	}

	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
	if tursoToken == "" {
		log.Fatal("TURSO_AUTH_TOKEN environment variable is not set")
	}

	clerkKey := os.Getenv("CLERK_KEY")
	if clerkKey == "" {
		log.Fatal("CLERK_KEY environment variable is not set")
	}

	s3Config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("S3Config is incorrect")
	}

	s3Client := s3.NewFromConfig(s3Config)

	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatal("S3_BUCKET environment variable is not set")
	}

	s3Region := os.Getenv("S3_REGION")
	if s3Region == "" {
		log.Fatal("S3_REGION environment variable is not set")
	}

	cloudfrontUrl := os.Getenv("CLOUDFRONT_URL")
	if cloudfrontUrl == "" {
		log.Fatal("CLOUDFRONT_URL environment variable not set")
	}

	cfg := apiConfig{
		S3Bucket:      s3Bucket,
		S3Client:      s3Client,
		S3Region:      s3Region,
		CloudfrontUrl: cloudfrontUrl,
		ClerkKey:      clerkKey,
		TursoToken:    tursoToken,
		TursoUrl:      tursoUrl,
	}

	return cfg
}
