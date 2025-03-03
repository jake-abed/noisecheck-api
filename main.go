package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jake-abed/noisecheck-api/internal/database"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

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
	// Grabs all documented keys
	cfg := createApiConfig()
	dbURL := fmt.Sprintf("%s?authToken=%s", cfg.TursoUrl, cfg.TursoToken)

	// Try to open DB and exit if it won't work.
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %s: %s", dbURL, err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)
	cfg.Db = dbQueries

	// Prep handlers with required auth wrapper.
	createRelease, _ := cfg.createReleaseHandler().(http.Handler)
	getUserReleases, _ := cfg.getUserReleasesHandler().(http.Handler)
	authdCreateRelease, _ := clerkhttp.
		WithHeaderAuthorization()(createRelease).(http.HandlerFunc)
	clerk.SetKey(cfg.ClerkKey)
	authdGetUserReleases, _ := clerkhttp.
		WithHeaderAuthorization()(getUserReleases).(http.HandlerFunc)
	createTrack, _ := cfg.createTrackHandler().(http.Handler)
	authdCreateTrack, _ := clerkhttp.
		WithHeaderAuthorization()(createTrack).(http.HandlerFunc)

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
	r.Get("/api/users/{userId}/releases", authdGetUserReleases)
	r.Get("/api/releases/{id}", cfg.getReleaseHandler)
	r.Post("/api/releases", authdCreateRelease)
	r.Get("/api/tracks/{id}", cfg.getTrackHandler)
	r.Post("/api/tracks", authdCreateTrack)
	http.ListenAndServe(":3000", r)
}
