package main

import (
	"fmt"
	"mime"
	"mime/multipart"
	"slices"
)

func madiaFileExtension(header *multipart.FileHeader, supported []string) (string, []string, error) {
	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		return mediaType, nil, err
	}

	if !slices.Contains(supported, mediaType) {
		return mediaType, nil, fmt.Errorf("unsupported media type: %s", mediaType)
	}

	extensions, err := mime.ExtensionsByType(mediaType)

	if err != nil {
		return mediaType, nil, err
	}

	return mediaType, extensions, nil
}

func idFileName(id string, header *multipart.FileHeader, supported []string) (string, string, error) {
	mediatype, extension, err := madiaFileExtension(header, supported)

	if err != nil {
		return mediatype, "", err
	}

	return mediatype, id + extension[0], nil
}
