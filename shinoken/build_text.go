package shinoken

import (
	"fmt"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/excel"
	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

func updateTextConfig(excelPath string, resolver TextResolver) error {
	configPath := resolver.ConfigPath()
	excelFileName := filepath.Base(excelPath)
	config := CreateTextConfig(excelFileName, "配信設計", "配信設計")

	if err := config.Validate(); err != nil {
		return fmt.Errorf("updateTextConfig: %w", err)
	}

	if err := jsonutil.Write(configPath, config); err != nil {
		return fmt.Errorf("updateTextConfig: %w", err)
	}

	return nil
}

func writeTextData(resolver TextResolver) error {
	configPath := resolver.ConfigPath()
	config, err := ReadTextConfig(configPath)
	if err != nil {
		return fmt.Errorf("writeTextData: %w", err)
	}

	rows, err := excel.ReadRowsFromExcel(resolver.Resolver, config.SchemaName)
	if err != nil {
		return fmt.Errorf("writeTextData: %w", err)
	}

	if err := jsonutil.Write(resolver.TextDataPath(config), rows); err != nil {
		return fmt.Errorf("writeTextData: %w", err)
	}

	return nil
}

// BuildTextDir bootstraps a text case directory for [TextResolver]:
// - run [excel.BuildDir] (case dir, workbook copy, config schema_preps, schemas/)
// - write config.json as [TextConfig] for 配信設計
// - run [excel.UpdateSchemasFromConfig] from the text config
// - write data/<schemaName>.json from the workbook
//
// For the full on-disk layout, see [TextResolver].
func BuildTextDir(excelPath string, resolver TextResolver) error {
	if err := excel.BuildDir(excelPath, resolver.RelDirPath()); err != nil {
		return fmt.Errorf("BuildTextDir: %w", err)
	}

	if err := updateTextConfig(excelPath, resolver); err != nil {
		return fmt.Errorf("BuildTextDir: %w", err)
	}

	if err := excel.UpdateSchemasFromConfig(resolver.Resolver); err != nil {
		return fmt.Errorf("BuildTextDir: %w", err)
	}

	if err := writeTextData(resolver); err != nil {
		return fmt.Errorf("BuildTextDir: %w", err)
	}

	return nil
}
