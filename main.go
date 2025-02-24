package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
}

func RespondWithError(
	w http.ResponseWriter,
	status int,
	msg string,
) {
	w.WriteHeader(status)
	error := ApiError{Error: msg}
	errBody, _ := json.Marshal(&error)
	w.Write(errBody)
	return
}

func main() {
	// Get and set all env variables.
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
	}

	// Prep handlers with required auth wrapper.
	createRelease, _ := cfg.createReleaseHandler().(http.Handler)
	authdCreateRelease, _ := clerkhttp.
		WithHeaderAuthorization()(createRelease).(http.HandlerFunc)
	clerk.SetKey(clerkKey)

	// Create Router
	r := chi.NewRouter()

	// Add Middleware To Router
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		MaxAge:         300,
	}))

	// Routes/Handlers
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome, but this an api bruh\n"))
	})
	r.Post("/api/webhooks/users", cfg.userWebhookHandler)
	r.Get("/api/releases/{id}", cfg.getReleaseHandler)
	r.Post("/api/releases", authdCreateRelease)
	http.ListenAndServe(":3000", r)
}
