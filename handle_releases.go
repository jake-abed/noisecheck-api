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
	"github.com/jake-abed/noisecheck-api/internal/utils"
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

		sanitizedRelName := utils.SanitizeReleaseName(newRelBody.Name)

		fileName := fmt.Sprintf("release/image/%s-%s.%s",
			imagePathMod, sanitizedRelName, fileFormat)

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
			Imgurl:   c.CloudfrontUrl + "/" + fileName,
			IsPublic: newRelBody.IsPublic,
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

func (c *apiConfig) getAllReleasesHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	releases, err := c.Db.GetAllPublicReleases(context.Background())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type allReleasesResp struct {
		Data []Release `json:"data"`
	}

	allReleases := []Release{}

	for _, rel := range releases {
		allReleases = append(allReleases, convertDbRelease(rel))
	}

	allResp := allReleasesResp{Data: allReleases}
	allRespBody, _ := json.Marshal(&allResp)
	w.WriteHeader(http.StatusOK)
	w.Write(allRespBody)
	return
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

	tracks, err := c.Db.GetTracksByRelease(context.Background(), idInt)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	convertedRel := convertDbRelease(rel)
	convertedTracks := []Track{}

	for _, t := range tracks {
		convertedTracks = append(convertedTracks, convertDbTrack(t))
	}

	releaseWithTracks := ReleaseWithTracks{
		Release: convertedRel,
		Tracks:  convertedTracks,
	}
	releaseReturn, _ := json.Marshal(&releaseWithTracks)
	w.WriteHeader(http.StatusOK)
	w.Write(releaseReturn)
	return
}

func (c *apiConfig) getUserReleasesHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		claims, ok := clerk.SessionClaimsFromContext(ctx)
		if !ok {
			RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		user, err := clerkUser.Get(ctx, claims.Subject)
		if err != nil {
			fmt.Printf("error: %s", err.Error())
		}

		if user == nil {
			RespondWithError(w, http.StatusNotFound, "user does not exist")
			return
		}

		userId := r.PathValue("userId")
		if userId == "" {
			RespondWithError(w, http.StatusBadRequest, "must provide user in path")
			return
		}

		if userId != user.ID {
			RespondWithError(w, http.StatusUnauthorized,
				"you may not access another user's releases")
			return
		}

		releases, err := c.Db.GetAllReleasesByUser(context.Background(), userId)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		type allReleasesResp struct {
			Data []Release `json:"data"`
		}

		allReleases := []Release{}

		for _, rel := range releases {
			allReleases = append(allReleases, convertDbRelease(rel))
		}

		allResp := allReleasesResp{Data: allReleases}
		allRespBody, _ := json.Marshal(&allResp)
		w.WriteHeader(http.StatusOK)
		w.Write(allRespBody)
		return
	})
}

func (c *apiConfig) deleteReleaseHandler() http.Handler {
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

		id := r.PathValue("id")
		if id == "" {
			RespondWithError(w, http.StatusBadRequest, "must provide user in path")
			return
		}

		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid Release Id")
			return
		}

		err = c.Db.DeleteReleaseById(context.Background(), idInt)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusAccepted)
		return
	})
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
			Imgurl:   relBody.Imgurl,
			IsPublic: relBody.IsPublic,
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
		ID:        int(rel.ID),
		Name:      rel.Name,
		UserID:    rel.UserID,
		Imgurl:    rel.Imgurl,
		IsPublic:  rel.IsPublic,
		CreatedAt: rel.CreatedAt,
		UpdatedAt: rel.UpdatedAt,
	}
}
