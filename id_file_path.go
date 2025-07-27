package main

import (
	"mime/multipart"
	"path/filepath"
)

func (cfg *apiConfig) assetsFilePath(fileName string) string {
	return filepath.Join(cfg.assetsRoot, fileName)
}

func (cfg *apiConfig) idFilePath(id string, header *multipart.FileHeader, supported []string, idFilePath func(string) string) (string, string, error) {

	mediaType, fileName, err := idFileName(id, header, supported)

	if err != nil {
		return mediaType, "", err
	}

	return mediaType, idFilePath(fileName), nil
}
