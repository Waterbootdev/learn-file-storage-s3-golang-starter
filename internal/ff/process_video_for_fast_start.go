package ff

import (
	"os/exec"
)

func ProcessVideoForFastStart(filePath string) (string, error) {

	file := filePath + ".processing"

	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", file)

	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return file, nil
}
