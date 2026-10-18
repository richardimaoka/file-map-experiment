package excel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/fs"
)

func ensureDirExists(resolver *PathResolver) error {
	path := resolver.DirPath()
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("ensureDirExists: %w", err)
	}
	return nil
}

func copyExcelFile(sourceExcelPath string, resolver *PathResolver) error {
	excelName := filepath.Base(sourceExcelPath)
	cfg := InitConfig(excelName)
	destExcel := resolver.ExcelPath(cfg)
	if err := fs.CopyFile(sourceExcelPath, destExcel); err != nil {
		return fmt.Errorf("copyExcelFile: %w", err)
	}
	return nil
}

func ensureConfigExists(sourceExcelPath string, resolver *PathResolver) error {
	configPath := resolver.ConfigPath()

	config, err := ReadExcelConfig(configPath)
	if err != nil {
		excelName := filepath.Base(sourceExcelPath)
		config = InitConfig(excelName)
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

// BuildExcelDir creates the case directory under APP_BASE_DIR at relPath
// - ensure the target dir exists
// - copy Excel file
// - ensure config.json exists
func BuildExcelDir(sourceExcelPath string, relPath string) error {
	excelName := filepath.Base(sourceExcelPath)
	if excelName == "." || excelName == "/" || excelName == "" {
		return fmt.Errorf("BuildExcelDir: invalid excel path %q", sourceExcelPath)
	}

	resolver := NewPathResolver(relPath)

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
