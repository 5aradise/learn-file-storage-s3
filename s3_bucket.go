package main

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (cfg apiConfig) saveInBucket(ctx context.Context, name string, mimeType string, data io.Reader) (url string, err error) {
	ext, err := extractExtFromMime(mimeType)
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("", "tubely-upload-*."+ext)
	if err != nil {
		return "", err
	}
	tempFilePath := tempFile.Name()
	defer os.Remove(tempFilePath)
	defer tempFile.Close()

	_, err = io.Copy(tempFile, data)
	if err != nil {
		return "", err
	}
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	err = processVideoForFastStart(tempFilePath, tempFilePath)
	if err != nil {
		return "", err
	}
	tempFile, err = os.Open(tempFilePath)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	ratio, err := getVideoAspectRatio(tempFilePath)
	if err != nil {
		return "", err
	}
	var prefix string
	switch ratio {
	case LANDSCAPE:
		prefix = "landscape/"
	case PORTRAIT:
		prefix = "portrait/"
	default:
		prefix = "other/"
	}

	fullName := name + "." + ext
	path := prefix + fullName
	return cfg.s3PutObject(ctx, path, mimeType, tempFile)
}

func (cfg apiConfig) s3PutObject(ctx context.Context, path, mimeType string, data io.Reader) (url string, err error) {
	_, err = cfg.s3Client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(cfg.s3Bucket),
			Key:         aws.String(path),
			Body:        data,
			ContentType: aws.String(mimeType),
		},
	)
	if err != nil {
		return "", err
	}

	return cfg.s3CfDistribution + path, nil
}
