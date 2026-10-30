package excel

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/fs"
	"github.com/richardimaoka/file-map-experiment/jsonutil"
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

// syncConfigFile loads config.json (or InitConfig), merges schema_preps from the case workbook, writes config.json.
func syncConfigFile(sourceExcelPath string, resolver Resolver) error {
	configPath := resolver.ConfigPath()
	excelName := filepath.Base(sourceExcelPath)

	config, err := ReadConfig(configPath)
	if err != nil {
		if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
			config = InitConfig(excelName)
		} else {
			return fmt.Errorf("syncConfigFile: %w", err)
		}
	}

	workbookPath := resolver.ExcelPath(config)
	config, err = enrichConfig(config, workbookPath)
	if err != nil {
		return fmt.Errorf("syncConfigFile: %w", err)
	}
	if err := jsonutil.Write(configPath, config); err != nil {
		return fmt.Errorf("syncConfigFile: %w", err)
	}
	return nil
}

func UpdateSchemasFromConfig(resolver Resolver) error {
	config, err := ReadConfig(resolver.ConfigPath())
	if err != nil {
		return fmt.Errorf("updateSchemasFromConfig: %w", err)
	}

	if err := os.MkdirAll(resolver.SchemasDir(), 0o755); err != nil {
		return fmt.Errorf("updateSchemasFromConfig: %w", err)
	}

	for _, prep := range config.SchemaPreps {
		if !prep.IsReady {
			continue
		}
		if err := prep.Validate(); err != nil {
			return fmt.Errorf("updateSchemasFromConfig: %w", err)
		}
		if err := WriteSchemaFromExcel(resolver, prep.SchemaName); err != nil {
			return fmt.Errorf("updateSchemasFromConfig: %w", err)
		}
	}
	return nil
}

// BuildDir bootstraps a case directory at relDirPath:
// - ensure the target dir exists
// - copy Excel file
// - sync config.json schema_preps with the copied workbook
// - rewrite schemas/ from config schema_preps (workbook was just overwritten)
//
// For the full on-disk layout, see [Resolver].
func BuildDir(sourceExcelPath string, relDirPath string) error {
	excelName := filepath.Base(sourceExcelPath)
	if excelName == "." || excelName == "/" || excelName == "" {
		return fmt.Errorf("excel.BuildDir: invalid excel path %q", sourceExcelPath)
	}

	resolver := NewResolver(relDirPath)

	// ensure the target dir exists
	if err := ensureDirExists(resolver); err != nil {
		return fmt.Errorf("excel.BuildDir: %w", err)
	}

	// copy Excel file
	if err := copyExcelFile(sourceExcelPath, resolver); err != nil {
		return fmt.Errorf("excel.BuildDir: %w", err)
	}

	// sync config.json schema_preps with the copied workbook
	if err := syncConfigFile(sourceExcelPath, resolver); err != nil {
		return fmt.Errorf("excel.BuildDir: %w", err)
	}

	// rewrite schemas/ from config schema_preps (workbook was just overwritten)
	if err := UpdateSchemasFromConfig(resolver); err != nil {
		return fmt.Errorf("excel.BuildDir: %w", err)
	}

	return nil
}
