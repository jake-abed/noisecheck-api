package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jake-abed/noisecheck-api/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type apiConfig struct {
	Db *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	tursoUrl := os.Getenv("TURSO_DATABASE_URL")
	tursoToken := os.Getenv("TURSO_AUTH_TOKEN")
	clerkKey := os.Getenv("CLERK_KEY")

	dbURL := fmt.Sprintf("%s?authToken=%s", tursoUrl, tursoToken)

	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %s: %s", dbURL, err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)

	cfg := apiConfig{Db: dbQueries}

	createRelease, _ := cfg.createReleaseHandler().(http.Handler)
	authdCreateRelease, _ := clerkhttp.
		WithHeaderAuthorization()(createRelease).(http.HandlerFunc)
	clerk.SetKey(clerkKey)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome\n"))
	})
	r.Post("/api/webhooks/users", cfg.userWebhookHandler)
	r.Post("/api/releases", authdCreateRelease)
	http.ListenAndServe(":3000", r)
}
