package main

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/auth"
	"github.com/google/uuid"
)

const MAX_MEMORY = 10 << 20

func madiaFileExtension(header *multipart.FileHeader, supported []string) (string, error) {
	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}

	if !slices.Contains(supported, mediaType) {
		return "", fmt.Errorf("unsupported media type: %s", mediaType)
	}

	extension, err := mime.ExtensionsByType(mediaType)

	if err != nil {
		return "", err
	}

	return extension[0], nil
}

func thumbnailFileName(videoID string, header *multipart.FileHeader) (string, error) {
	extension, err := madiaFileExtension(header, []string{"image/jpeg", "image/png"})

	if err != nil {
		return "", err
	}

	return videoID + extension, nil
}

func (cfg *apiConfig) videoFilePath(fileName string) string {
	return filepath.Join(cfg.assetsRoot, fileName)
}

func (cfg *apiConfig) portURL() string {
	return fmt.Sprintf("http://localhost:%s", cfg.port)
}
func (cfg *apiConfig) getFilePathURL(filePath string) *string {
	path := filepath.Join(cfg.portURL(), filePath)
	return &path
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

	fileName, err := thumbnailFileName(videoIDString, header)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	filePath := cfg.videoFilePath(fileName)

	assetFile, err := os.Create(filePath)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create file", err)
		return
	}

	defer assetFile.Close()

	_, err = io.Copy(assetFile, file)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't copy file", err)
		return
	}

	video, err := cfg.db.GetVideo(videoID)
	if err != nil || video.UserID != userID {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get video", err)
		return
	}

	video.ThumbnailURL = cfg.getFilePathURL(filePath)

	err = cfg.db.UpdateVideo(video)

	if err != nil || video.UserID != userID {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update video", err)
		return
	}

	respondWithJSON(w, http.StatusOK, video)
}
