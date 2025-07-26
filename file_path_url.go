package main

import (
	"fmt"
	"mime/multipart"
)

func (cfg *apiConfig) getPortURL() string {
	return fmt.Sprintf("http://localhost:%s/", cfg.port)
}

func (cfg *apiConfig) filePathURL(filePath string) *string {
	path := cfg.getPortURL() + filePath
	return &path
}

func (cfg *apiConfig) idFilePathURL(id string, header *multipart.FileHeader, supported []string) (string, *string, error) {
	filePath, err := cfg.idFilePath(id, header, supported)

	if err != nil {
		return "", nil, err
	}

	return filePath, cfg.filePathURL(filePath), nil
}
