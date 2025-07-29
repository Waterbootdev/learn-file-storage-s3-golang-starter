package main

import (
	"errors"
	"strings"
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/s3"
)

func splitExtendedURL(video database.Video) (string, string, error) {
	splitted := strings.Split(*video.VideoURL, ",")

	if len(splitted) != 2 {
		return "", "", errors.New("invalid video url")
	}

	return splitted[0], splitted[1], nil
}

func (cfg *apiConfig) extendURL(key string) string {
	return cfg.s3Bucket + "," + key
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	bucket, key, err := splitExtendedURL(video)

	if err != nil {
		return video, err
	}

	url, err := s3.GeneratePresignedURL(cfg.s3Client, bucket, key, 15*time.Minute)

	if err != nil {
		return video, err
	}

	video.VideoURL = &url

	return video, nil

}
