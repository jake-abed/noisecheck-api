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

const MAX_TRACK_SIZE = 8 << 29

func (c *apiConfig) getTrackHandler(w http.ResponseWriter, r *http.Request) {
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

	dbTrack, err := c.Db.GetTrackById(context.Background(), idInt)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	track := convertDbTrack(dbTrack)
	trackBody, err := json.Marshal(&track)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(trackBody)
	return
}

func (c *apiConfig) createTrackHandler() http.Handler {
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

		err = r.ParseMultipartForm(MAX_TRACK_SIZE)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		body := r.FormValue("data")
		newTrackBody := NewTrackBody{}
		err = json.Unmarshal([]byte(body), &newTrackBody)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		file, fileHeader, err := r.FormFile("file")
		defer file.Close()

		contentType := fileHeader.Header.Get("Content-Type")
		fileType, _, err := mime.ParseMediaType(contentType)
		if fileType != "audio/wav" && fileType != "audio/flac" {
			RespondWithError(w, http.StatusUnsupportedMediaType, "Wrong file type!")
			return
		}
		fileFormat := strings.Replace(contentType, "audio/", "", 1)

		randBytes := make([]byte, 8)
		rand.Read(randBytes)
		imagePathMod := base64.RawURLEncoding.EncodeToString(randBytes)

		sanitizedRelName := utils.SanitizeReleaseName(newTrackBody.Name)

		fileName := fmt.Sprintf("track/audio/%s-%s.%s",
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

		newRelParams := database.CreateTrackParams{
			Name:      newTrackBody.Name,
			ReleaseID: int64(newTrackBody.ReleaseId),
			TrackUrl:  c.CloudfrontUrl + "/" + fileName,
		}

		addedTrack, err := c.Db.CreateTrack(context.Background(), newRelParams)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			error := ApiError{Error: err.Error()}
			errorBody, _ := json.Marshal(&error)
			w.Write(errorBody)
			return
		}

		track := convertDbTrack(addedTrack)

		trackReturn, _ := json.Marshal(&track)
		w.WriteHeader(http.StatusOK)
		w.Write(trackReturn)
		return

	})
}

func convertDbTrack(dbTrack database.Track) Track {
	return Track{
		ID:        int(dbTrack.ID),
		Name:      dbTrack.Name,
		ReleaseId: int(dbTrack.ReleaseID),
		CreatedAt: dbTrack.CreatedAt,
		UpdatedAt: dbTrack.UpdatedAt,
	}
}
