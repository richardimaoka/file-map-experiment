package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func SubDirs(parentDirPath string) ([]string, error) {
	entries, err := os.ReadDir(parentDirPath)
	if err != nil {
		return nil, fmt.Errorf("fs.SubDirs: %v", err)
	}

	dirNames := []string{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirNames = append(dirNames, entry.Name())
	}

	return dirNames, nil
}

// FilesWithExtensions lists regular files directly under parentDir whose extension matches
// one of extensions. At least one extension is required. Each extension must start with ".",
// use lowercase letters only (e.g. ".xlsx"), and must not contain whitespace.
func FilesWithExtensions(parentDir string, extensions ...string) ([]string, error) {
	wantExt, err := validateExtensions(extensions)
	if err != nil {
		return nil, fmt.Errorf("fs.FilesWithExtensions: %w", err)
	}

	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, fmt.Errorf("fs.FilesWithExtensions: %w", err)
	}

	var paths []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		got := strings.ToLower(filepath.Ext(entry.Name()))
		if !extensionAllowed(got, wantExt) {
			continue
		}
		paths = append(paths, filepath.Join(parentDir, entry.Name()))
	}
	return paths, nil
}

func validateExtensions(extensions []string) ([]string, error) {
	if len(extensions) == 0 {
		return nil, fmt.Errorf("at least one extension is required")
	}
	for _, ext := range extensions {
		if ext == "" {
			return nil, fmt.Errorf("extension must not be empty")
		}
		for _, r := range ext {
			if unicode.IsSpace(r) {
				return nil, fmt.Errorf("extension %q must not contain whitespace", ext)
			}
		}
		if !strings.HasPrefix(ext, ".") {
			return nil, fmt.Errorf("extension %q must start with '.'", ext)
		}
		if ext != strings.ToLower(ext) {
			return nil, fmt.Errorf("extension %q must be lowercase", ext)
		}
	}
	return extensions, nil
}

func extensionAllowed(got string, want []string) bool {
	for _, ext := range want {
		if got == ext {
			return true
		}
	}
	return false
}
