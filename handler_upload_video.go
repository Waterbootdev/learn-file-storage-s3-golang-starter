package main

import (
	"fmt"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {
	videoIDString, videoID, userID, ok := cfg.getPathIdvalidateUserIDParseMultipartForm(w, r, "videoID", 10<<30)

	if !ok {
		return
	}

	fmt.Println("uploading video", videoIDString, "by user", userID)

	url, written, err := cfg.copyToRandomIdFile(r, "video", []string{"video/mp4", ""}, 32, cfg.copyToS3IdFile)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't copy file to s3 id file", err)
		return
	}

	fmt.Println("copied", written, "bytes to", *url)

	cfg.updateVideo(w, videoID, userID, func(video *database.Video) {
		video.VideoURL = url
	})
}
