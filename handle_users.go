package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jake-abed/noisecheck-api/internal/database"
)

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

		userResp, _ := json.Marshal(&updatedUser)
		w.WriteHeader(http.StatusOK)
		w.Write(userResp)
		return
	case "user.deleted":
		err := c.Db.DeleteUserById(context.Background(), newBody.Data.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			error := ApiError{Error: err.Error()}
			errorBody, _ := json.Marshal(&error)
			w.Write(errorBody)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		error := ApiError{Error: newBody.Type + " not a valid event type."}
		errorBody, _ := json.Marshal(&error)
		w.Write(errorBody)
		return
	}
}
