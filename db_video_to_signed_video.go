package main

import (
	"errors"
	"strings"
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/s3"
)

func splitExtendedKey(extendedURL *string) (string, string, error) {
	splitted := strings.Split(*extendedURL, ",")

	if len(splitted) != 2 {
		return "", "", errors.New("invalid video url")
	}

	return splitted[0], splitted[1], nil
}

func (cfg *apiConfig) extendedKey(key string) *string {
	extended := cfg.s3Bucket + "," + key
	return &extended
}
func (cfg *apiConfig) extendedKeyToPresignedURL(extendedURL *string) (*string, error) {

	if extendedURL == nil {
		return nil, nil
	}

	bucket, key, err := splitExtendedKey(extendedURL)

	if err != nil {
		return nil, err
	}

	url, err := s3.GeneratePresignedURL(cfg.s3Client, bucket, key, 15*time.Minute)

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {

	url, err := cfg.extendedKeyToPresignedURL(video.VideoURL)

	if err != nil {
		return video, err
	}

	video.VideoURL = url

	return video, nil
}
