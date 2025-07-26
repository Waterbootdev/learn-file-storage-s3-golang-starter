package main

import (
	"fmt"
	"net/http"
)

const MAX_VIDEO_MEMORY = 10 << 30

func (cfg *apiConfig) handlerUploadVideo(w http.ResponseWriter, r *http.Request) {
	idString, _, userID, ok := cfg.getPathIdvalidateUserIDParseMultipartForm(w, r, "videoID", MAX_VIDEO_MEMORY)

	if !ok {
		return
	}

	fmt.Println("uploading video", idString, "by user", userID)

	panic("not implemented")
}
