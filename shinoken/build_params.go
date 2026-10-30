package shinoken

import (
	"fmt"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/excel"
	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

func updateParamsConfig(excelPath string, resolver ParamsResolver) error {
	configPath := resolver.ConfigPath()
	excelFileName := filepath.Base(excelPath)
	config := CreateParamsConfig(excelFileName, "更新用シート", "更新用シート")

	if err := config.Validate(); err != nil {
		return fmt.Errorf("updateParamsConfig: %w", err)
	}

	if err := jsonutil.Write(configPath, config); err != nil {
		return fmt.Errorf("updateParamsConfig: %w", err)
	}

	return nil
}

func writeParamsData(resolver ParamsResolver) error {
	configPath := resolver.ConfigPath()
	config, err := ReadParamsConfig(configPath)
	if err != nil {
		return fmt.Errorf("writeParamsData: %w", err)
	}

	rows, err := excel.ReadRowsFromExcel(resolver.Resolver, config.SchemaName)
	if err != nil {
		return fmt.Errorf("writeParamsData: %w", err)
	}

	if err := jsonutil.Write(resolver.ParamsAllPath(config), rows); err != nil {
		return fmt.Errorf("writeParamsData: %w", err)
	}

	return nil
}

// BuildParamsDir bootstraps a parameters case directory for [ParamsResolver]:
// - run [excel.BuildDir] (case dir, workbook copy, config schema_preps, schemas/)
// - write config.json as [ParamsConfig] for 更新用シート
// - run [excel.UpdateSchemasFromConfig] from the params config
// - write data/<schemaName>.json from the workbook
//
// For the full on-disk layout, see [ParamsResolver].
func BuildParamsDir(excelPath string, resolver ParamsResolver) error {
	if err := excel.BuildDir(excelPath, resolver.RelDirPath()); err != nil {
		return fmt.Errorf("BuildParamsDir: %w", err)
	}

	if err := updateParamsConfig(excelPath, resolver); err != nil {
		return fmt.Errorf("BuildParamsDir: %w", err)
	}

	if err := excel.UpdateSchemasFromConfig(resolver.Resolver); err != nil {
		return fmt.Errorf("BuildParamsDir: %w", err)
	}

	if err := writeParamsData(resolver); err != nil {
		return fmt.Errorf("BuildParamsDir: %w", err)
	}

	return nil
}
