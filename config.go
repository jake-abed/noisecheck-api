package main

import (
	"context"
	"database/sql"
	"fmt"
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

	dbURL := fmt.Sprintf("%s?authToken=%s", tursoUrl, tursoToken)

	// Try to open DB and exit if it won't work.
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %s: %s", dbURL, err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)

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
		Db:            dbQueries,
		S3Bucket:      s3Bucket,
		S3Client:      s3Client,
		S3Region:      s3Region,
		CloudfrontUrl: cloudfrontUrl,
		ClerkKey:      clerkKey,
	}

	return cfg
}
