package shinoken

import (
	"fmt"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/fs"
)

func manualFileExists(path string) error {
	if !fs.FileExists(path) {
		baseName := filepath.Base(path)
		return fmt.Errorf("%q does not exist, make it manually", baseName)
	}
	return nil
}
