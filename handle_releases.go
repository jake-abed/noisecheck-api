package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/jake-abed/noisecheck-api/internal/database"
)

const MAX_IMG_SIZE = 10 << 18

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

		err = r.ParseMultipartForm(MAX_IMG_SIZE)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		body := r.FormValue("data")
		newRelBody := NewReleaseBody{}
		err = json.Unmarshal([]byte(body), &newRelBody)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		file, fileHeader, err := r.FormFile("file")
		defer file.Close()

		contentType := fileHeader.Header.Get("Content-Type")
		fileType, _, err := mime.ParseMediaType(contentType)
		if fileType != "image/jpeg" && fileType != "image/png" {
			RespondWithError(w, http.StatusUnsupportedMediaType, "Wrong file type!")
			return
		}
		fileFormat := strings.Replace(contentType, "image/", "", 1)

		randBytes := make([]byte, 8)
		rand.Read(randBytes)
		imagePathMod := base64.RawURLEncoding.EncodeToString(randBytes)

		fileName := fmt.Sprintf("release/image/%s-%s.%s",
			imagePathMod, newRelBody.Name, fileFormat)

		s3Params := &s3.PutObjectInput{
			Bucket:      &c.S3Bucket,
			Key:         &fileName,
			Body:        file,
			ContentType: &contentType,
		}

		_, err = c.S3Client.PutObject(context.Background(), s3Params)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "AWS Error")
			fmt.Println(err)
			return
		}

		newRelParams := database.CreateReleaseParams{
			Name:     newRelBody.Name,
			UserID:   user.ID,
			Url:      "#",
			Imgurl:   c.CloudfrontUrl + "/" + fileName,
			IsPublic: newRelBody.IsPublic,
			IsSingle: newRelBody.IsSingle,
		}

		newRel, err := c.Db.CreateRelease(context.Background(), newRelParams)
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

func (c *apiConfig) getReleaseHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		RespondWithError(w, http.StatusNotFound, "Release ID required")
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid Release Id")
		return
	}

	rel, err := c.Db.GetReleaseById(context.Background(), idInt)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	convertedRel := convertDbRelease(rel)

	releaseReturn, _ := json.Marshal(&convertedRel)
	w.WriteHeader(http.StatusOK)
	w.Write(releaseReturn)
	return
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
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if user.ID != relBody.UserID {
			RespondWithError(w, http.StatusUnauthorized,
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
			RespondWithError(w, http.StatusInternalServerError, err.Error())
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
