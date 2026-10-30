package fs

import (
	"errors"
	"os"
)

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if err == nil {
		return true // file exists
	}
	if errors.Is(err, os.ErrNotExist) {
		return false // file does not exist
	}

	return false // Other error (e.g., permission issues)
}
