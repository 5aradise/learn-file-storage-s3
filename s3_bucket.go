package main

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (cfg apiConfig) saveInBucket(ctx context.Context, name string, mimeType string, src io.Reader) (url string, err error) {
	ext, err := extractExtFromMime(mimeType)
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("", "tubely-upload-*."+ext)
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = io.Copy(tempFile, src)
	if err != nil {
		return "", err
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	fullName := name + "." + ext
	_, err = cfg.s3Client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(cfg.s3Bucket),
			Key:         aws.String(fullName),
			Body:        tempFile,
			ContentType: aws.String(mimeType),
		},
	)
	if err != nil {
		return "", err
	}

	return "https://" + cfg.s3Bucket + ".s3." + cfg.s3Region + ".amazonaws.com/" + fullName, nil
}
