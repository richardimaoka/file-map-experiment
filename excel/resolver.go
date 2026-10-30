package excel

import (
	"fmt"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/env"
)

// PathResolver resolves paths for a single Excel case directory under APP_BASE_DIR.
//
// Expected layout:
//
//	$APP_BASE_DIR/<relPath>/
//	  <excelFileName>   # .xlsx/.xlsm; name from config.json
//	  config.json
//	  schemas/
//	    <name>.json
//
// Example: NewPathResolver("cases/client1")
type PathResolver struct {
	baseDir string
	relPath string
}

func NewPathResolver(relPath string) PathResolver {
	baseDir := env.AppBaseDir()

	return PathResolver{
		baseDir: baseDir,
		relPath: relPath,
	}
}

func (p PathResolver) DirPath() string {
	return filepath.Join(p.baseDir, p.relPath)
}

func (p PathResolver) ExcelPath(config ExcelConfig) string {
	return filepath.Join(p.DirPath(), config.ExcelFileName)
}

func (p PathResolver) ConfigPath() string {
	return filepath.Join(p.DirPath(), "config.json")

}

func (p PathResolver) SchemaCollectionPath() string {
	return filepath.Join(p.DirPath(), "schemas")
}

func (p PathResolver) SchemaPath(schemaName string) string {
	return filepath.Join(p.SchemaCollectionPath(), schemaName+".json")
}

func (p PathResolver) ValidateDir() error {
	_, err := ReadExcelConfig(p.ConfigPath())
	if err != nil {
		return fmt.Errorf("ValidateDir: %w", err)
	}
	return nil
}
