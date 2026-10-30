package excel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/fs"
)

func ensureDirExists(resolver Resolver) error {
	path := resolver.FullDirPath()
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("ensureDirExists: %w", err)
	}
	return nil
}

func copyExcelFile(sourceExcelPath string, resolver Resolver) error {
	excelName := filepath.Base(sourceExcelPath)
	cfg := InitConfig(excelName)
	destExcel := resolver.ExcelPath(cfg)
	if err := fs.CopyFile(sourceExcelPath, destExcel); err != nil {
		return fmt.Errorf("copyExcelFile: %w", err)
	}
	return nil
}

func ensureConfigExists(sourceExcelPath string, resolver Resolver) error {
	configPath := resolver.ConfigPath()

	_, err := ReadConfig(configPath)
	if err != nil {
		excelName := filepath.Base(sourceExcelPath)
		config := InitConfig(excelName)
		data, err := json.Marshal(config)
		if err != nil {
			return fmt.Errorf("ensureConfigExists: %w", err)
		}
		if err := os.WriteFile(configPath, data, 0o644); err != nil {
			return fmt.Errorf("ensureConfigExists: %w", err)
		}
	}

	return nil
}

// BuildDir bootstraps a case directory at relDirPath:
// creates the directory, copies the workbook, and writes config.json if missing.
// It does not create schemas/ or schema JSON files. For the full on-disk layout, see [Resolver].
func BuildDir(sourceExcelPath string, relDirPath string) error {
	excelName := filepath.Base(sourceExcelPath)
	if excelName == "." || excelName == "/" || excelName == "" {
		return fmt.Errorf("BuildExcelDir: invalid excel path %q", sourceExcelPath)
	}

	resolver := NewResolver(relDirPath)

	// ensure the target dir exists
	if err := ensureDirExists(resolver); err != nil {
		return fmt.Errorf("BuildExcelDir: %w", err)
	}

	// copy Excel file
	if err := copyExcelFile(sourceExcelPath, resolver); err != nil {
		return fmt.Errorf("BuildExcelDir: %w", err)
	}

	// ensure config.json exists
	if err := ensureConfigExists(sourceExcelPath, resolver); err != nil {
		return fmt.Errorf("BuildExcelDir: %w", err)
	}

	return nil
}
