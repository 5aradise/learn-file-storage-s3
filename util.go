package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
)

func randURLName() (string, error) {
	randSl := make([]byte, 32)
	_, err := rand.Read(randSl)
	if err != nil {
		return "", err
	}
	name := base64.RawURLEncoding.EncodeToString(randSl)
	return name, nil
}

var (
	errInvalidMimeType = errors.New("invalid mime type")
)

func extractExtFromMime(mimeType string) (string, error) {
	parts := strings.Split(mimeType, "/")
	if len(parts) != 2 {
		return "", errInvalidMimeType
	}
	return parts[1], nil
}

type aspectRatio string

var (
	LANDSCAPE aspectRatio = "16:9"
	PORTRAIT  aspectRatio = "9:16"
)

type ffprobeData struct {
	Streams []struct {
		AspectRatio aspectRatio `json:"display_aspect_ratio"`
	} `json:"streams"`
}

var (
	errEmptyStreamsSlice = errors.New("empty data streams slice")
)

func getVideoAspectRatio(filePath string) (aspectRatio, error) {
	res := &bytes.Buffer{}

	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	cmd.Stdout = res
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	data := ffprobeData{}
	err = json.NewDecoder(res).Decode(&data)
	if err != nil {
		return "", err
	}
	if len(data.Streams) == 0 {
		return "", errEmptyStreamsSlice
	}

	return data.Streams[0].AspectRatio, nil
}

func processVideoForFastStart(srcPath, dstPath string) error {
	if srcPath != dstPath {
		return exec.Command("ffmpeg", "-i", srcPath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", dstPath).Run()
	}

	dstPath = srcPath + ".processing"

	cmd := exec.Command("ffmpeg", "-i", srcPath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", dstPath)
	err := cmd.Run()
	if err != nil {
		os.Remove(dstPath)
		return err
	}

	err = os.Remove(srcPath)
	if err != nil {
		os.Remove(dstPath)
		return err
	}
	return os.Rename(dstPath, srcPath)
}
