package main

import (
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

var (
	errInvalidMimeType = errors.New("invalid mime type")
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func (cfg apiConfig) saveAsset(name string, mimeType string, src io.Reader) (url string, err error) {
	parts := strings.Split(mimeType, "/")
	if len(parts) != 2 {
		return "", errInvalidMimeType
	}
	assetName := name + "." + parts[1]

	assetPath := path.Join(cfg.assetsRoot, assetName)
	f, err := os.Create(assetPath)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(f, src)
	if err != nil {
		return "", err
	}

	return "http://localhost:" + cfg.port + assetsPath + assetName, nil
}
