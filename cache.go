package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"
)

const (
	CachePath = ".cache/tldr"
	tldrURL   = "https://codeload.github.com/tldr-pages/tldr/zip/refs/heads/main"
)

var (
	ErrorCacheNotExists = errors.New("cache doesn't exists")
)

func getCustomPath() (path string) {
	return getCachePath() + "/custom/"
}

func getOfficalPath() (path string) {
	return getCachePath() + "/tldr-main/"
}

func getCachePath() (path string) {

	homeDir, err := getHomeDir()
	if err != nil {
		return
	}

	path = homeDir + "/" + CachePath

	return path
}

func updateCache() (err error) {
	// Check if cache folder exists
	cachePath := getCachePath()
	zipFilePath := cachePath + "/tldr-main.zip"
	upzipFilePath := cachePath + "/tldr-main"
	if err = os.RemoveAll(upzipFilePath); err != nil {
		return
	}
	if err = os.RemoveAll(zipFilePath); err != nil {
		return
	}
	if err = os.MkdirAll(cachePath, os.ModePerm); err != nil {
		return
	}

	if err = downloadZip(tldrURL, zipFilePath); err != nil {
		return
	}

	if err = unzip(zipFilePath, cachePath); err != nil {
		return
	}

	return nil
}

func checkLocalCache(cfg *Config, name string) (tldrFile string, page []string, err error) {
	platforms := []string{PlatformCommon, *cfg.Platform}
	tldrFile = ""

	// Check file for all platforms
	for _, platform := range platforms {
		// Build path to the local pages
		langFile := buildLocalPath(*cfg.Language, platform) + name + ".md"
		enFile := buildLocalPath("", platform) + name + ".md"

		if isFileExists(langFile) {
			tldrFile = langFile
		} else if isFileExists(enFile) {
			tldrFile = enFile
		}

		// If page not exist, just return
		if len(tldrFile) > 0 {
			var data []byte
			// Read page data
			data, err = ioutil.ReadFile(tldrFile)
			if err != nil {
				return
			}
			// Split text to lines
			page = strings.Split(string(data), "\n")
		}
	}
	return
}

// *cfg.Language

func buildLocalPath(lang, platform string) (dir string) {

	if lang != "en" && lang != "" {
		dir = getCachePath() + "/tldr-main/pages." + lang + "/" + platform + "/"
	} else {
		dir = getCachePath() + "/tldr-main/pages/" + platform + "/"
	}
	return
}

func buildCustomPath(name string) (file string) {
	file = getCachePath() + "/custom/" + name + ".md"
	return
}

func checkCustom(name string) (tldrFile string, page []string, err error) {
	tldrFile = ""

	// Build path to the local pages
	langFile := buildCustomPath(name)

	if isFileExists(langFile) {
		tldrFile = langFile
	}

	// If page not exist, just return
	if len(tldrFile) > 0 {
		var data []byte
		// Read page data
		data, err = ioutil.ReadFile(tldrFile)
		if err != nil {
			return
		}
		// Split text to lines
		page = strings.Split(string(data), "\n")
	}

	return
}

func getCacheAge() (bool, time.Time) {
	cachePath := getCachePath()
	zipFilePath := cachePath + "/tldr-main.zip"

	file, err := os.Stat(zipFilePath)

	if err != nil {
		fmt.Printf("Error get cache file: %v\n", err)
		return false, time.Time{}
	}

	modifiedtime := file.ModTime()
	return true, modifiedtime
}

func humanizeDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
		{"second", seconds},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d %ss", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}
