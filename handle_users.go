package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *apiConfig) getUserHandler(w http.ResponseWriter, r *http.Request) {

}

func (c *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type userBody struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	decoder := json.NewDecoder(r.Body)
	newBody := userBody{}
	if err := decoder.Decode(&newBody); err != nil {
		w.WriteHeader(401)
		error := ApiError{Error: "Malformed body!"}
		errorBody, _ := json.Marshal(&error)
		w.Write(errorBody)
		return
	}

	fmt.Printf("%s added to the app I guess.", newBody.Email)
	w.WriteHeader(204)
	return
}
