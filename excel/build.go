package excel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/fs"
)

// BuildExcelDir creates the case directory under APP_BASE_DIR at relPath, copies (overwrites) the workbook,
// and writes config.json from InitConfig only when config.json is missing.
func BuildExcelDir(sourceExcelPath string, relPath string) error {
	resolver := NewPathResolver(relPath)

	excelName := filepath.Base(sourceExcelPath)
	if excelName == "." || excelName == "/" || excelName == "" {
		return fmt.Errorf("BuildExcelDir: invalid excel path %q", sourceExcelPath)
	}

	if err := os.MkdirAll(resolver.DirPath(), 0o755); err != nil {
		return fmt.Errorf("BuildExcelDir: %w", err)
	}

	cfg := InitConfig(excelName)
	destExcel := resolver.ExcelPath(cfg)
	if err := fs.CopyFile(sourceExcelPath, destExcel); err != nil {
		return fmt.Errorf("BuildExcelDir: %w", err)
	}

	configPath := resolver.ConfigPath()
	if _, err := os.Stat(configPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("BuildExcelDir: %w", err)
		}
		data, err := json.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("BuildExcelDir: %w", err)
		}
		if err := os.WriteFile(configPath, data, 0o644); err != nil {
			return fmt.Errorf("BuildExcelDir: %w", err)
		}
	}

	if _, err := ReadExcelConfig(configPath); err != nil {
		return fmt.Errorf("BuildExcelDir: %w", err)
	}

	return nil
}
