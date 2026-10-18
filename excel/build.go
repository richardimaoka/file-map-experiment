package excel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/fs"
)

// BuildExcelDir creates targetDirPath, copies the workbook if absent, and writes config.json from InitConfig when missing.
func BuildExcelDir(excelPath string, targetDirPath string) error {
	excelName := filepath.Base(excelPath)
	if excelName == "." || excelName == "/" || excelName == "" {
		return fmt.Errorf("BuildExcelDir: invalid excel path %q", excelPath)
	}

	if err := os.MkdirAll(targetDirPath, 0o755); err != nil {
		return fmt.Errorf("BuildExcelDir: %w", err)
	}

	destExcel := filepath.Join(targetDirPath, excelName)
	if _, err := os.Stat(destExcel); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("BuildExcelDir: %w", err)
		}
		if err := fs.CopyFile(excelPath, destExcel); err != nil {
			return fmt.Errorf("BuildExcelDir: %w", err)
		}
	}

	configPath := filepath.Join(targetDirPath, "config.json")
	if _, err := os.Stat(configPath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("BuildExcelDir: %w", err)
		}
		cfg := InitConfig(excelName)
		data, err := json.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("BuildExcelDir: %w", err)
		}
		if err := os.WriteFile(configPath, data, 0o644); err != nil {
			return fmt.Errorf("BuildExcelDir: %w", err)
		}
	}

	return nil
}
