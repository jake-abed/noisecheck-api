package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jake-abed/noisecheck-api/internal/database"
)

type ClerkUser struct {
	Data struct {
		EmailAddresses []struct {
			EmailAddress string `json:"email_address"`
			ID           string `json:"id"`
			LinkedTo     []any  `json:"linked_to"`
			Object       string `json:"object"`
		} `json:"email_addresses"`
		FirstName       string `json:"first_name"`
		ID              string `json:"id"`
		ImageURL        string `json:"image_url"`
		LastName        string `json:"last_name"`
		ProfileImageURL string `json:"profile_image_url"`
		Username        string `json:"username"`
	} `json:"data"`
	Object    string `json:"object"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"`
}

func (c *apiConfig) getUserHandler(w http.ResponseWriter, r *http.Request) {
	type GetUserBody struct {
		Username string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	userBody := GetUserBody{}
	if err := decoder.Decode(&userBody); err != nil {
		w.WriteHeader(401)
		error := ApiError{Error: "Malformed body!"}
		errorBody, _ := json.Marshal(&error)
		w.Write(errorBody)
		return
	}

}

func (c *apiConfig) userWebhookHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	newBody := ClerkUser{}
	if err := decoder.Decode(&newBody); err != nil {
		w.WriteHeader(401)
		error := ApiError{Error: "Malformed body!"}
		errorBody, _ := json.Marshal(&error)
		w.Write(errorBody)
		return
	}

	if len(newBody.Data.EmailAddresses) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		error := ApiError{Error: "No email in payload!"}
		errorBody, _ := json.Marshal(&error)
		w.Write(errorBody)
		return
	}

	switch newBody.Type {
	case "user.created":
		username := newBody.Data.Username

		if username == "" {
			username = newBody.Data.FirstName + newBody.Data.LastName
		}

		params := database.CreateUserParams{
			ID:       newBody.Data.ID,
			Email:    newBody.Data.EmailAddresses[0].EmailAddress,
			Username: username,
		}

		newUser, err := c.Db.CreateUser(context.Background(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			error := ApiError{Error: err.Error()}
			errorBody, _ := json.Marshal(&error)
			w.Write(errorBody)
			return
		}

		fmt.Printf("%s added to the app I guess.\n", newUser.Email)
		userResp, _ := json.Marshal(&newUser)
		w.WriteHeader(http.StatusOK)
		w.Write(userResp)
		return
	case "user.updated":
		params := database.UpdateUserByIdParams{
			Email:    newBody.Data.EmailAddresses[0].EmailAddress,
			Username: newBody.Data.Username,
			ID:       newBody.Data.ID,
		}

		updatedUser, err := c.Db.UpdateUserById(context.Background(), params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			error := ApiError{Error: err.Error()}
			errorBody, _ := json.Marshal(&error)
			w.Write(errorBody)
			return
		}

		fmt.Printf("%s updated\n", updatedUser.Username)
		userResp, _ := json.Marshal(&updatedUser)
		w.WriteHeader(http.StatusOK)
		w.Write(userResp)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		error := ApiError{Error: newBody.Type + " not a valid event type."}
		errorBody, _ := json.Marshal(&error)
		w.Write(errorBody)
		return
	}
}
