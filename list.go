package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func getCachedCommandList() (commands []string) {
	platforms := []string{PlatformCommon, platform}

	for _, platform := range platforms {
		path := buildLocalPath(language, platform)
		if !isFileExists(path) {
			return
		}
		_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			commands = append(commands, strings.ReplaceAll(info.Name(), ".md", ""))
			return nil
		})
		sort.Strings(commands)
	}
	return
}

func getLocalCommandList() (commands []string) {

	customPath := getCustomPath()

	if !isFileExists(customPath) {
		return
	}
	_ = filepath.Walk(customPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		commands = append(commands, strings.ReplaceAll(info.Name(), ".md", ""))
		return nil
	})
	sort.Strings(commands)

	return
}
