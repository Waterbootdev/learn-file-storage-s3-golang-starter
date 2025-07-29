package main

import (
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) updateVideo(w http.ResponseWriter, videoId uuid.UUID, userID uuid.UUID, update func(*database.Video)) {

	video, err := cfg.db.GetVideo(videoId)
	if err != nil || video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get video", err)
		return
	}

	update(&video)

	err = cfg.db.UpdateVideo(video)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	cfg.responsVideo(w, video)
}

func (cfg *apiConfig) responsVideo(w http.ResponseWriter, video database.Video) {
	video, err := cfg.dbVideoToSignedVideo(video)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't sign video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}

func (cfg *apiConfig) responsVideos(w http.ResponseWriter, videos []database.Video) {

	responseVideos := make([]database.Video, len(videos))

	for i, v := range videos {
		currentVideo, err := cfg.dbVideoToSignedVideo(v)

		if err != nil {
			return
		}
		responseVideos[i] = currentVideo
	}

	respondWithJSON(w, http.StatusOK, responseVideos)
}
