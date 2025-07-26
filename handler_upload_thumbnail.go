package main

import (
	"fmt"
	"net/http"
)

const MAX_MEMORY = 10 << 20

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {

	videoID, id, userID, ok := cfg.getPathIdvalidateUserIDParseMultipartForm(w, r, "videoID", MAX_MEMORY)

	if !ok {
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	url, written, err := cfg.copyToRandomIdFile(r, "thumbnail", []string{"image/jpeg", "image/png"}, 32)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't copy file to id file", err)
		return
	}

	fmt.Println("copied", written, "bytes to", *url)

	video, err := cfg.db.GetVideo(id)
	if err != nil || video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get video", err)
		return
	}

	video.ThumbnailURL = url

	err = cfg.db.UpdateVideo(video)

	if err != nil || video.UserID != userID {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}
