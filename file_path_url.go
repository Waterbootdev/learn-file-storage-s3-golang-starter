package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
)

func (cfg *apiConfig) getPortURL() string {
	return fmt.Sprintf("http://localhost:%s/", cfg.port)
}

func (cfg *apiConfig) filePathURL(filePath string) *string {
	path := cfg.getPortURL() + filePath
	return &path
}

func (cfg *apiConfig) idFilePathURL(request *http.Request, formFileKey string, id string, supported []string, idFilePath func(string) string) (multipart.File, string, string, error) {

	file, header, err := request.FormFile(formFileKey)

	if err != nil {
		return file, "", "", err
	}

	mediaType, filePath, err := cfg.idFilePath(id, header, supported, idFilePath)

	if err != nil {
		file.Close()
		return file, mediaType, filePath, err
	}

	return file, mediaType, filePath, nil
}
