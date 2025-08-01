package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/ff"
)

func copyFormFileTo(formFile multipart.File, filePath string) (int64, error) {

	file, err := os.Create(filePath)

	if err != nil {
		return 0, err
	}

	defer file.Close()

	return io.Copy(file, formFile)
}

func (cfg *apiConfig) copyToIdFile(request *http.Request, formFileKey string, supportedMediatypes []string, id string) (string, int64, error) {

	file, _, filePath, err := cfg.idFilePathURL(request, formFileKey, id, supportedMediatypes, cfg.assetsFilePath)

	if err != nil {
		return filePath, 0, err
	}

	written, err := copyFormFileTo(file, filePath)

	file.Close()

	return filePath, written, err
}

func (cfg *apiConfig) copyToTempFile(request *http.Request, formFileKey string, supportedMediatypes []string, id string) (string, string, string, int64, error) {

	file, mediaType, fileName, err := cfg.idFilePathURL(request, formFileKey, id, supportedMediatypes, func(s string) string { return s })

	if err != nil {
		return mediaType, fileName, "", 0, err
	}

	defer file.Close()

	tempFile, err := os.CreateTemp(cfg.tempRoot, "tubely-upload-*.mp4")

	if err != nil {
		return mediaType, fileName, "", 0, err
	}

	defer tempFile.Close()

	tempFileName := tempFile.Name()

	written, err := io.Copy(tempFile, file)

	if err != nil {
		defer os.Remove(tempFile.Name())
		return mediaType, fileName, tempFileName, written, err
	}

	return mediaType, fileName, tempFileName, written, err
}

func getPrefixSchema(tempFile string) (string, error) {

	aspectRatio, err := ff.GetVideoAspectRatio(tempFile)

	if err != nil {
		return "", err
	}

	schema, err := ff.PrefixSchema(aspectRatio)

	if err != nil {
		return "", err
	}

	return schema, nil
}

func getPrefixFilePathURL(tempFile string, fileName string) (string, error) {

	prefixSchema, err := getPrefixSchema(tempFile)

	if err != nil {
		return "", err
	}

	return prefixSchema + "/" + fileName, nil
}

func (cfg *apiConfig) copyToS3IdFile(request *http.Request, formFileKey string, supportedMediatypes []string, id string) (string, int64, error) {

	mediaType, fileName, tempFileName, written, err := cfg.copyToTempFile(request, formFileKey, supportedMediatypes, id)

	if err != nil {
		return fileName, written, err
	}

	defer os.Remove(tempFileName)

	tempFastStartFileName, err := ff.ProcessVideoForFastStart(tempFileName)

	if err != nil {
		return fileName, written, err
	}

	fileName, err = getPrefixFilePathURL(tempFastStartFileName, fileName)

	if err != nil {
		return fileName, written, err
	}

	tempFile, err := os.Open(tempFastStartFileName)

	if err != nil {
		return fileName, written, err
	}

	defer tempFile.Close()
	defer os.Remove(tempFastStartFileName)

	obj, err := cfg.s3Client.PutObject(request.Context(), &s3.PutObjectInput{
		Bucket:      &cfg.s3Bucket,
		Key:         &fileName,
		Body:        tempFile,
		ContentType: &mediaType,
	})

	fmt.Println(*obj)

	if err != nil {
		return fileName, written, err
	}

	return fileName, written, err
}

func (cfg *apiConfig) copyToRandomIdFile(request *http.Request, formFileKey string, supportedMediatypes []string, numberBytes int, copyToIdFile func(*http.Request, string, []string, string) (string, int64, error)) (string, int64, error) {
	id, err := randomId(numberBytes)

	if err != nil {
		return "", 0, err
	}

	return copyToIdFile(request, formFileKey, supportedMediatypes, id)
}
