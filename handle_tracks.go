package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jake-abed/noisecheck-api/internal/database"
)

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

func convertDbTrack(dbTrack database.Track) Track {
	return Track{
		ID:        int(dbTrack.ID),
		Name:      dbTrack.Name,
		ReleaseId: int(dbTrack.ReleaseID),
		CreatedAt: dbTrack.CreatedAt,
		UpdatedAt: dbTrack.UpdatedAt,
	}
}
