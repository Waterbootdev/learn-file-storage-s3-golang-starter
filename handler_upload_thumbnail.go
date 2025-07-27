package main

import (
	"fmt"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {

	videoIDString, videoID, userID, ok := cfg.getPathIdvalidateUserIDParseMultipartForm(w, r, "videoID", 10<<20)

	if !ok {
		return
	}

	fmt.Println("uploading thumbnail for video", videoIDString, "by user", userID)

	filePath, written, err := cfg.copyToRandomIdFile(r, "thumbnail", []string{"image/jpeg", "image/png"}, 32, cfg.copyToIdFile)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't copy file to id file", err)
		return
	}

	url := cfg.filePathURL(filePath)

	fmt.Println("copied", written, "bytes to", *url)

	cfg.updateVideo(w, videoID, userID, func(video *database.Video) {
		video.ThumbnailURL = url
	})

}
