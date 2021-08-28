package main

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
)

const (
	// PlatformCommon is the name of common pages available across platforms.
	PlatformCommon = "common"
)

func isFileExists(path string) (isExist bool) {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()
			return err
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			fileReader.Close()
			targetFile.Close()
			return err
		}
		fileReader.Close()
		targetFile.Close()
	}

	return nil
}
func downloadZip(url string, path string) (err error) {
	// Download the ZIP file
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	log.Println(path)
	out, err := os.Create(path)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	return
}

func getHomeDir() (homeDir string, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	homeDir = usr.HomeDir
	return
}
