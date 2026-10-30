package excel

import (
	"fmt"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/env"
)

// Resolver resolves paths for a single Excel case directory under APP_BASE_DIR.
//
// Expected layout:
//
//	$APP_BASE_DIR/<relDirPath>/
//	  <excelFileName>   # .xlsx/.xlsm; name from config.json
//	  config.json
//	  schemas/
//	    <name>.json
//
// Example: NewResolver("cases/client1")
type Resolver struct {
	appBaseDir string
	relDirPath string
}

func NewResolver(relDirPath string) Resolver {
	baseDir := env.AppBaseDir()

	return Resolver{
		appBaseDir: baseDir,
		relDirPath: relDirPath,
	}
}

func (p Resolver) FullDirPath() string {
	return filepath.Join(p.appBaseDir, p.relDirPath)
}

func (p Resolver) ExcelPath(config Config) string {
	return filepath.Join(p.FullDirPath(), config.ExcelFileName)
}

func (p Resolver) ConfigPath() string {
	return filepath.Join(p.FullDirPath(), "config.json")
}

func (p Resolver) SchemaCollectionPath() string {
	return filepath.Join(p.FullDirPath(), "schemas")
}

func (p Resolver) SchemaPath(schemaName string) string {
	return filepath.Join(p.SchemaCollectionPath(), schemaName+".json")
}

func (p Resolver) ValidateDir() error {
	_, err := ReadConfig(p.ConfigPath())
	if err != nil {
		return fmt.Errorf("ValidateDir: %w", err)
	}
	return nil
}
