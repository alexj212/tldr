package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func getCachedCommandList(cfg *Config) (commands []string) {
	platforms := []string{PlatformCommon, *cfg.Platform}

	for _, platform := range platforms {
		path := buildLocalPath(*cfg.Language, platform)
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

	customPath := getCachePath() + "/custom/"

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
