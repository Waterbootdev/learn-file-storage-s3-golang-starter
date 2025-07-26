package main

import (
	"fmt"
	"mime"
	"mime/multipart"
	"slices"
)

func madiaFileExtension(header *multipart.FileHeader, supported []string) ([]string, error) {
	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	if !slices.Contains(supported, mediaType) {
		return nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}

	extensions, err := mime.ExtensionsByType(mediaType)

	if err != nil {
		return nil, err
	}

	return extensions, nil
}

func idFileName(id string, header *multipart.FileHeader, supported []string) (string, error) {
	extension, err := madiaFileExtension(header, supported)

	if err != nil {
		return "", err
	}

	return id + extension[0], nil
}
