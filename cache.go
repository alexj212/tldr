package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"
)

const (
	// CachePath location of cache of tldr archive and pages.
	CachePath = ".cache/tldr"
	// ConfigPath location of local tldr config and local pages.
	ConfigPath = ".config/tldr"
	// tldrURL location to download main archive of tldr pages.
	tldrURL = "https://codeload.github.com/tldr-pages/tldr/zip/refs/heads/main"
)

func getCustomPath() (path string) {
	return getConfigPath() + "/custom/"
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

func getConfigPath() (path string) {

	homeDir, err := getHomeDir()
	if err != nil {
		return
	}

	path = homeDir + "/" + ConfigPath

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

func checkLocalCache(name string) (tldrFile string, page []string, err error) {
	platforms := []string{PlatformCommon, platform}
	tldrFile = ""

	// Check file for all platforms
	for _, platf := range platforms {
		// Build path to the local pages
		langFile := buildLocalPath(language, platf) + name + ".md"
		enFile := buildLocalPath("", platf) + name + ".md"

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

func buildLocalPath(lang, platf string) (dir string) {

	if lang != "en" && lang != "" {
		dir = getCachePath() + "/tldr-main/pages." + lang + "/" + platf + "/"
	} else {
		dir = getCachePath() + "/tldr-main/pages/" + platf + "/"
	}
	return
}

func buildCustomPath(name string) (file string) {
	file = getCustomPath() + name + ".md"
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
