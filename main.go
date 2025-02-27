package main

import (
	"encoding/json"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

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
	cfg := createApiConfig()

	// Prep handlers with required auth wrapper.
	createRelease, _ := cfg.createReleaseHandler().(http.Handler)
	getUserReleases, _ := cfg.getUserReleasesHandler().(http.Handler)
	authdCreateRelease, _ := clerkhttp.
		WithHeaderAuthorization()(createRelease).(http.HandlerFunc)
	clerk.SetKey(cfg.ClerkKey)
	authdGetUserReleases, _ := clerkhttp.
		WithHeaderAuthorization()(getUserReleases).(http.HandlerFunc)

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
	http.ListenAndServe(":3000", r)
}
