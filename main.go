package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

func main() {
	godotenv.Load()
	_ = os.Getenv("TURSO_DATABASE_URL")
	_ = os.Getenv("TURSO_AUTH_TOKEN")
	clerkKey := os.Getenv("CLERK_KEY")

	createRelease, _ := createReleaseHandler().(http.Handler)
	authdCreateRelease, _ := clerkhttp.WithHeaderAuthorization()(createRelease).(http.HandlerFunc)
	clerk.SetKey(clerkKey)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome\n"))
	})
	r.Post("/users", createUserHandler)
	r.Post("/releases", authdCreateRelease)
	http.ListenAndServe(":3000", r)
}
