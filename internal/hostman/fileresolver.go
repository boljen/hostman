package hostman

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ResolveConfigFilePath resolves a configuration file path based on user input.
// If input is just a filename (no path separators), it searches the current
// directory and then ascends parent directories until the root to find the file.
// If input contains a path, it resolves that specific path only (relative to cwd if relative).
// Returns the absolute path to the existing file or an error if not found.
func ResolveConfigFilePath(input string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return resolveFileFrom(input, cwd)
}

func resolveFileFrom(filename string, cwd string) (string, error) {
	if filename == "" {
		return "", errors.New("empty config path")
	}

	if isFilename(filename) {
		return resolveFileUpwards(filename, cwd)
	} else {
		return resolveSpecificPath(filename, cwd)
	}
}

func resolveFileUpwards(filename string, searchDir string) (string, error) {
	for {
		candidate := filepath.Join(searchDir, filename)
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			abs, err := filepath.Abs(candidate)
			if err != nil {
				return "", err
			}
			return abs, nil
		}

		parent := parentDir(searchDir)
		if parent == searchDir {
			break
		}
		searchDir = parent
	}
	return "", os.ErrNotExist
}

func resolveSpecificPath(location string, startDir string) (string, error) {
	if !filepath.IsAbs(location) {
		location = filepath.Join(startDir, location)
	}

	abs, err := filepath.Abs(location)
	if err != nil {
		return "", err
	}

	if info, err := os.Stat(abs); err == nil && !info.IsDir() {
		return abs, nil
	}
	return "", os.ErrNotExist
}

func isFilename(path string) bool {
	if path == "" {
		return false
	} else if path == "." || path == ".." {
		return false
	}
	return !strings.ContainsRune(path, '/') && !strings.ContainsRune(path, '\\')
}

func parentDir(p string) string {
	if p == "" {
		return ""
	}
	return filepath.Dir(p)
}
