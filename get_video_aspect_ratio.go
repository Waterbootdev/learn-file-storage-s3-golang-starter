package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type ffprobe struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func (ff *ffprobe) aspectRatio(more bool) string {
	width := ff.Streams[0].Width
	height := ff.Streams[0].Height

	if width > height {
		if height*16/10 == width*9/10 {
			return "16:9"
		} else if more && height*4/10 == width*3/10 {
			return "4:3"
		}
	} else if width < height {
		if width*16/10 == height*9/10 {
			return "9:16"
		} else if more && width*4/10 == height*3/10 {
			return "3:4"
		}
	}

	return "other"
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	buffer := new(bytes.Buffer)

	cmd.Stdout = buffer

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	var probe ffprobe

	err = json.Unmarshal(buffer.Bytes(), &probe)

	if err != nil {
		return "", err
	}

	return probe.aspectRatio(true), nil
}

func prefixSchema(aspectRatio string) (string, error) {
	switch aspectRatio {
	case "16:9":
		return "landscape", nil
	case "9:16":
		return "portrait", nil
	case "3:4":
		return "portrait", nil
	case "4:3":
		return "landscape", nil
	case "other":
		return "other", nil
	default:
		return "", fmt.Errorf("unsupported aspect ratio (%s)", aspectRatio)
	}
}
