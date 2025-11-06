package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CreateDirs creates directories with proper error handling
func CreateDirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// WriteFile writes content to a file with error handling
func WriteFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// ReplacePlaceholders replaces template placeholders in content
func ReplacePlaceholders(content string, replacements map[string]string) string {
	for key, value := range replacements {
		content = strings.ReplaceAll(content, key, value)
	}
	return content
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
