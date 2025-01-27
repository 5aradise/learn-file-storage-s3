package main

import (
	"io"
	"os"
	"path"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func (cfg apiConfig) saveInLocal(name string, mimeType string, src io.Reader) (url string, err error) {
	ext, err := extractExtFromMime(mimeType)
	if err != nil {
		return "", err
	}
	fullName := name + "." + ext

	assetPath := path.Join(cfg.assetsRoot, fullName)
	f, err := os.Create(assetPath)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(f, src)
	if err != nil {
		return "", err
	}

	return "http://localhost:" + cfg.port + assetsPath + fullName, nil
}
