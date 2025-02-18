package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/jake-abed/noisecheck-api/internal/database"
)

type release struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserID   string `json:"user_id"`
	Url      string `json:"url"`
	Imgurl   string `json:"image_url"`
	IsPublic bool   `json:"is_public"`
	IsSingle bool   `json:"is_single"`
}

func (c *apiConfig) createReleaseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := clerk.SessionClaimsFromContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"access": "unauthorized"}`))
			return
		}

		usr, err := user.Get(ctx, claims.Subject)
		if err != nil {
			fmt.Printf("error: %s", err.Error())
		}

		if usr == nil {
			w.Write([]byte(`{"error": "User does not exist"}`))
			return
		}

		decoder := json.NewDecoder(r.Body)
		relBody := release{}
		if err := decoder.Decode(&relBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			error := ApiError{Error: "Malformed body!"}
			errorBody, _ := json.Marshal(&error)
			w.Write(errorBody)
			return
		}

		params := database.CreateReleaseParams{
			Name:     relBody.Name,
			Url:      "",
			UserID:   usr.ID,
			Imgurl:   relBody.Imgurl,
			IsPublic: relBody.IsPublic,
			IsSingle: relBody.IsSingle,
		}

		newRel, err := c.Db.CreateRelease(context.Background(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			error := ApiError{Error: err.Error()}
			errorBody, _ := json.Marshal(&error)
			w.Write(errorBody)
			return
		}

		rel := release{
			ID:       int(newRel.ID),
			Name:     newRel.Name,
			UserID:   newRel.UserID,
			Url:      newRel.Url,
			Imgurl:   newRel.Imgurl,
			IsPublic: newRel.IsPublic,
			IsSingle: newRel.IsSingle,
		}
		releaseReturn, _ := json.Marshal(&rel)
		w.WriteHeader(http.StatusOK)
		w.Write(releaseReturn)
		return
	})
}
