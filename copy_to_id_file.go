package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func copyFormFileTo(formFile multipart.File, filePath string) (int64, error) {

	file, err := os.Create(filePath)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	return io.Copy(file, formFile)
}

func (cfg *apiConfig) copyToIdFile(request *http.Request, formFileKey string, supportedMediatypes []string, id string) (*string, int64, error) {

	file, header, err := request.FormFile(formFileKey)

	if err != nil {
		return nil, 0, err
	}

	defer file.Close()

	filePath, url, err := cfg.idFilePathURL(id, header, supportedMediatypes)

	if err != nil {
		return url, 0, err
	}

	written, err := copyFormFileTo(file, filePath)

	return url, written, err
}

func (cfg *apiConfig) copyToRandomIdFile(request *http.Request, formFileKey string, supportedMediatypes []string, numberBytes int) (*string, int64, error) {
	id, err := randomId(numberBytes)

	if err != nil {
		return nil, 0, err
	}

	return cfg.copyToIdFile(request, formFileKey, supportedMediatypes, id)
}
