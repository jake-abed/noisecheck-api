package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/jake-abed/noisecheck-api/internal/database"
)

func (c *apiConfig) createReleaseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := clerk.SessionClaimsFromContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"access": "unauthorized"}`))
			return
		}

		user, err := clerkUser.Get(ctx, claims.Subject)
		if err != nil {
			fmt.Printf("error: %s", err.Error())
		}

		if user == nil {
			w.Write([]byte(`{"error": "User does not exist"}`))
			return
		}

		decoder := json.NewDecoder(r.Body)
		relBody := Release{}
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
			UserID:   user.ID,
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

		rel := convertDbRelease(newRel)

		releaseReturn, _ := json.Marshal(&rel)
		w.WriteHeader(http.StatusOK)
		w.Write(releaseReturn)
		return
	})
}

func (c *apiConfig) GetReleaseHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	relBody := Release{}
	err := decoder.Decode(&relBody)

	if err != nil {
		RespondWithError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if !relBody.IsPublic {
		RespondWithError(w, r, http.StatusUnauthorized, "This Release is Private")
		return
	}
}

func (c *apiConfig) UpdateReleaseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := clerk.SessionClaimsFromContext(ctx)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"access": "unauthorized"}`))
			return
		}

		user, err := clerkUser.Get(ctx, claims.Subject)
		if err != nil {
			fmt.Printf("error: %s", err.Error())
		}

		if user == nil {
			w.Write([]byte(`{"error": "User does not exist"}`))
			return
		}

		decoder := json.NewDecoder(r.Body)
		relBody := Release{}
		err = decoder.Decode(&relBody)
		if err != nil {
			RespondWithError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		if user.ID != relBody.UserID {
			RespondWithError(w, r, http.StatusUnauthorized,
				"you may not edit another user's track")
			return
		}

		params := database.UpdateReleaseParams{
			Name:     relBody.Name,
			Url:      relBody.Url,
			Imgurl:   relBody.Imgurl,
			IsPublic: relBody.IsPublic,
			IsSingle: relBody.IsSingle,
			ID:       int64(relBody.ID),
		}

		updatedRel, err := c.Db.UpdateRelease(context.Background(), params)
		if err != nil {
			RespondWithError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		relRes := convertDbRelease(updatedRel)
		resBody, _ := json.Marshal(&relRes)
		w.WriteHeader(http.StatusOK)
		w.Write(resBody)
		return
	})
}

func convertDbRelease(rel database.Release) Release {
	return Release{
		ID:       int(rel.ID),
		Name:     rel.Name,
		Url:      rel.Url,
		Imgurl:   rel.Imgurl,
		IsPublic: rel.IsPublic,
		IsSingle: rel.IsSingle,
	}
}
