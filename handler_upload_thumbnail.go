package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

const MAX_MEMORY = 10 << 20

func getDataUrl(data []byte, mediaType string) *string {
	encodedData := base64.StdEncoding.EncodeToString(data)
	url := fmt.Sprintf("data:%s;base64,%s", mediaType, encodedData)
	return &url
}

func getMapUrl(data []byte, mediaType string, port string, videoID uuid.UUID) *string {

	videoThumbnail := thumbnail{
		data:      data,
		mediaType: mediaType,
	}

	videoThumbnails[videoID] = videoThumbnail

	url := fmt.Sprintf("http://localhost:%s/api/thumbnails/{%s}", port, videoID)
	return &url
}

func (cfg *apiConfig) handlerUploadThumbnail(w http.ResponseWriter, r *http.Request) {
	videoIDString := r.PathValue("videoID")
	videoID, err := uuid.Parse(videoIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid ID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	fmt.Println("uploading thumbnail for video", videoID, "by user", userID)

	// TODO: implement the upload here

	err = r.ParseMultipartForm(MAX_MEMORY)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse multipart form", err)
		return
	}

	file, header, err := r.FormFile("thumbnail")

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse form file", err)
		return
	}

	defer file.Close()

	mediaType := header.Header.Get("Content-Type")

	data, err := io.ReadAll(file)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "cant reead file", err)
		return
	}

	video, err := cfg.db.GetVideo(videoID)
	if err != nil || video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get video", err)
		return
	}

	video.ThumbnailURL = getDataUrl(data, mediaType)

	//video.ThumbnailURL = getMapUrl(data, mediaType, cfg.port, videoID)

	err = cfg.db.UpdateVideo(video)

	if err != nil || video.UserID != userID {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}
