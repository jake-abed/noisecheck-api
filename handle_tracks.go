package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/jake-abed/noisecheck-api/internal/database"
	"github.com/jake-abed/noisecheck-api/internal/utils"
)

// Around 600 MB
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
		if validAudioType(fileType) {
			fmt.Println(fileType)
			RespondWithError(w, http.StatusUnsupportedMediaType, "Wrong file type!")
			return
		}

		var fileFormat string

		if strings.Contains(fileType, "wav") {
			fileFormat = "wav"
		} else {
			fileFormat = "flac"
		}

		audioFile, err := os.CreateTemp("", "noisecheck-upload."+fileFormat)
		defer os.Remove(os.TempDir() + "noisecheck-upload." + fileFormat)
		defer file.Close()
		io.Copy(audioFile, file)
		audioFile.Seek(0, io.SeekStart)

		length, _, err := getAudioFileInfo(audioFile.Name())
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		mp3FilePath, err := convertSourceAudioToMp3(audioFile.Name(), fileFormat)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		mp3File, err := os.Open(mp3FilePath)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer os.Remove(mp3FilePath)
		defer mp3File.Close()
		mp3File.Seek(0, io.SeekStart)

		randBytes := make([]byte, 8)
		rand.Read(randBytes)
		losslessPathMod := base64.RawURLEncoding.EncodeToString(randBytes)

		sanitizedName := utils.SanitizeReleaseName(newTrackBody.Name)

		losslessFileName := fmt.Sprintf("track/audio/lossless/%s-%s.%s",
			losslessPathMod, sanitizedName, fileFormat)

		s3Params := &s3.PutObjectInput{
			Bucket:      &c.S3Bucket,
			Key:         &losslessFileName,
			Body:        audioFile,
			ContentType: &contentType,
		}

		_, err = c.S3Client.PutObject(context.Background(), s3Params)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "AWS Error")
			fmt.Println(err)
			return
		}

		randBytes2 := make([]byte, 8)
		rand.Read(randBytes2)
		mp3PathMod := base64.RawURLEncoding.EncodeToString(randBytes2)

		mp3FileName := fmt.Sprintf("track/audio/mp3/%s-%s.%s", mp3PathMod,
			sanitizedName, "mp3")

		s3ParamsMp3 := &s3.PutObjectInput{
			Bucket: &c.S3Bucket,
			Key:    &mp3FileName,
			Body:   mp3File,
		}

		_, err = c.S3Client.PutObject(context.Background(), s3ParamsMp3)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "AWS Error")
			fmt.Println(err)
			return
		}

		newRelParams := database.CreateTrackParams{
			Name:            newTrackBody.Name,
			Length:          int64(length),
			ReleaseID:       int64(newTrackBody.ReleaseId),
			OriginalFileUrl: c.CloudfrontUrl + "/" + losslessFileName,
			Mp3FileUrl:      c.CloudfrontUrl + "/" + mp3FileName,
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
		ID:              int(dbTrack.ID),
		Name:            dbTrack.Name,
		Length:          int(dbTrack.Length),
		OriginalFileUrl: dbTrack.OriginalFileUrl,
		Mp3FileUrl:      dbTrack.Mp3FileUrl,
		ReleaseId:       int(dbTrack.ReleaseID),
		CreatedAt:       dbTrack.CreatedAt,
		UpdatedAt:       dbTrack.UpdatedAt,
	}
}

type ffprobeResult struct {
	Streams []ffprobeStream `json:"streams"`
	Format  ffprobeFormat   `json:"format"`
}

type ffprobeStream struct {
	Duration string `json:"duration"`
}

type ffprobeFormat struct {
	FormatName string `json:"format_name"`
}

func getAudioFileInfo(filePath string) (float64, string, error) {
	ffprobeCmd := exec.Command("ffprobe", "-sexagesimal", "-v", "error",
		"-print_format", "json", "-show_format", "-show_streams", filePath)

	buf := &bytes.Buffer{}
	ffprobeCmd.Stdout = buf
	err := ffprobeCmd.Run()

	if err != nil {
		fmt.Println("Error occurred in getting file info: ", filePath)
		return 0, "", err
	}

	res := ffprobeResult{}
	err = json.Unmarshal(buf.Bytes(), &res)
	if err != nil {
		return 0, "", err
	}

	seconds, err := parseDurationString(res.Streams[0].Duration)
	fmt.Println(err)

	return seconds, res.Format.FormatName, nil
}

func convertSourceAudioToMp3(filePath string, format string) (string, error) {
	outputFilePath := strings.ReplaceAll(filePath, format, "mp3")
	ffmpegCmd := exec.Command("ffmpeg", "-i", filePath, "-vn",
		"-ar", "44100", "-ac", "2", "-b:a", "192k", "-f", "mp3", outputFilePath)
	err := ffmpegCmd.Run()
	if err != nil {
		fmt.Println(outputFilePath, " - ", filePath)
		return "", err
	}
	return outputFilePath, nil
}

func parseDurationString(duration string) (float64, error) {
	times := strings.Split(duration, ":")

	fmt.Println(times[2])

	if len(times) < 3 {
		return 0, fmt.Errorf("invalid duration string: %s", duration)
	}

	hourString := times[0]
	minString := times[1]
	secString := times[2]

	hours, err := strconv.ParseFloat(hourString, 10)
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.ParseFloat(minString, 10)
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.ParseFloat(secString, 10)
	if err != nil {
		return 0, err
	}

	hoursSeconds := hours * 60 * 60
	minutesSeconds := minutes * 60

	return hoursSeconds + minutesSeconds + seconds, nil
}

func validAudioType(audioType string) bool {
	return audioType != "audio/vnd.wave" && audioType != "audio/wav" &&
		audioType != "audio/flac" && audioType != "audio/wave" &&
		audioType != "audio/vnd.wav" && audioType != "audio/x-wav" &&
		audioType != "audio/x-pn-wav" && audioType != "audio/x-flac"
}
