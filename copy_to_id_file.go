package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	file, _, filePath, url, err := cfg.idFilePathURL(request, formFileKey, id, supportedMediatypes, cfg.assetsFilePath, cfg.filePathURL)

	if err != nil {
		return url, 0, err
	}

	written, err := copyFormFileTo(file, filePath)

	file.Close()

	return url, written, err
}

func (cfg *apiConfig) copyFormFileToTemp(request *http.Request, formFileKey string, supportedMediatypes []string, id string) (string, string, *string, *os.File, int64, error) {

	file, mediaType, fileName, url, err := cfg.idFilePathURL(request, formFileKey, id, supportedMediatypes, func(s string) string { return s }, cfg.s3FilePathURL)

	if err != nil {
		return mediaType, fileName, url, nil, 0, err
	}

	defer file.Close()

	tempFile, err := os.CreateTemp(cfg.tempRoot, "tubely-upload-*.mp4")

	if err != nil {
		return mediaType, fileName, url, tempFile, 0, err
	}

	written, err := io.Copy(tempFile, file)

	return mediaType, fileName, url, tempFile, written, err
}

func (cfg *apiConfig) copyToS3IdFile(request *http.Request, formFileKey string, supportedMediatypes []string, id string) (*string, int64, error) {

	mediaType, fileName, url, tempFile, written, err := cfg.copyFormFileToTemp(request, formFileKey, supportedMediatypes, id)

	defer os.Remove(tempFile.Name())

	defer tempFile.Close()

	if err != nil {
		return url, written, err
	}

	_, err = tempFile.Seek(0, io.SeekStart)

	if err != nil {
		return url, written, err
	}

	obj, err := cfg.s3Client.PutObject(request.Context(), &s3.PutObjectInput{
		Bucket:      &cfg.s3Bucket,
		Key:         &fileName,
		Body:        tempFile,
		ContentType: &mediaType,
	})

	fmt.Println(*obj)

	if err != nil {
		return url, written, err
	}

	return url, written, err
}

func (cfg *apiConfig) copyToRandomIdFile(request *http.Request, formFileKey string, supportedMediatypes []string, numberBytes int, copyToIdFile func(*http.Request, string, []string, string) (*string, int64, error)) (*string, int64, error) {
	id, err := randomId(numberBytes)

	if err != nil {
		return nil, 0, err
	}

	return copyToIdFile(request, formFileKey, supportedMediatypes, id)
}
