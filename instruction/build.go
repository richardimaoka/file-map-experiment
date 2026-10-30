package instruction

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardimaoka/file-map-experiment/excel"
	"github.com/richardimaoka/file-map-experiment/fs"
	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

func updateDefaultDataFromConfig(resolver Resolver) error {
	if err := os.MkdirAll(resolver.DefaultsDir(), 0o755); err != nil {
		return fmt.Errorf("updateDefaultDataFromConfig: %w", err)
	}

	schemaPaths, err := fs.FilesWithExtensions(resolver.SchemasDir(), ".json")
	if err != nil {
		return fmt.Errorf("updateDefaultDataFromConfig: %w", err)
	}

	for _, schemaPath := range schemaPaths {
		fileName := filepath.Base(schemaPath)
		schemaName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		rows, err := excel.ReadRowsFromExcel(resolver.Resolver, schemaName)
		if err != nil {
			return fmt.Errorf("updateDefaultDataFromConfig: %w", err)
		}
		if err := jsonutil.Write(resolver.DefaultDataPath(schemaName), rows); err != nil {
			return fmt.Errorf("updateDefaultDataFromConfig: %w", err)
		}
	}
	return nil
}

func syncInstructionConfig(resolver Resolver) error {
	cfg, err := ReadConfig(resolver.ConfigPath())
	if err != nil {
		return fmt.Errorf("syncInstructionConfig: %w", err)
	}
	excelPath := resolver.ExcelPath(cfg.Config)
	cfg, err = filterSchemaPreps(cfg, excelPath)
	if err != nil {
		return fmt.Errorf("syncInstructionConfig: %w", err)
	}
	if err := jsonutil.Write(resolver.ConfigPath(), cfg.Config); err != nil {
		return fmt.Errorf("syncInstructionConfig: %w", err)
	}
	return nil
}

// BuildDir bootstraps an instruction case at relDirPath:
// - run [excel.BuildDir]
// - rewrite config.json from filterSchemaPreps (management sheets, comments)
// - rewrite defaults/<schemaName>.json for each schemas/<schemaName>.json
//
// For the full on-disk layout, see [Resolver].
func BuildDir(sourceExcelPath string, relDirPath string) error {
	excelFileName := filepath.Base(sourceExcelPath)
	relDirBase := filepath.Base(relDirPath)
	slog.Info("instruction.BuildDir: starting", "rel_dir", relDirBase, "excelFile", excelFileName)

	// run excel.BuildDir
	if err := excel.BuildDir(sourceExcelPath, relDirPath); err != nil {
		return fmt.Errorf("instruction.BuildDir: %w", err)
	}

	resolver := NewResolver(relDirPath)

	// rewrite config.json from filterSchemaPreps
	if err := syncInstructionConfig(resolver); err != nil {
		return fmt.Errorf("instruction.BuildDir: %w", err)
	}

	// rewrite defaults/<schemaName>.json for each schemas/<schemaName>.json
	if err := updateDefaultDataFromConfig(resolver); err != nil {
		return fmt.Errorf("instruction.BuildDir: %w", err)
	}

	return nil
}
