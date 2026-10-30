package shinoken

import (
	"fmt"

	"github.com/richardimaoka/file-map-experiment/excel"
	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

type ParamsConfig struct {
	SchemaName string `json:"scema_file"`
	SheetName  string `json:"sheet_name"`
	excel.Config
}

func CreateParamsConfig(excelFileName, schemaName, sheetName string) ParamsConfig {
	excelConfig := excel.InitConfig(excelFileName)
	excelConfig.SchemaPreps = []excel.SchemaPrep{
		{
			SchemaName:     schemaName,
			SheetName:      sheetName,
			IsReady:        true,
			HeaderStartRow: 4,
			HeaderEndRow:   4,
			DataStartRow:   205,
			DataEndRow:     5000,
			StartColumn:    "B",
			EndColumn:      "P",
		},
	}

	return ParamsConfig{
		SchemaName: schemaName,
		SheetName:  sheetName,
		Config:     excelConfig,
	}
}

func (c ParamsConfig) Validate() error {
	if c.SchemaName == "" {
		return fmt.Errorf("ParamsConfig.Validate: scema_file is not set")
	}
	if c.SheetName == "" {
		return fmt.Errorf("ParamsConfig.Validate: sheet_name is not set")
	}
	if c.ExcelFileName == "" {
		return fmt.Errorf("ParamsConfig.Validate: excel_file_name is not set")
	}
	return nil
}

func ReadParamsConfig(configPath string) (ParamsConfig, error) {
	var cfg ParamsConfig
	if err := jsonutil.Read(configPath, &cfg); err != nil {
		return ParamsConfig{}, fmt.Errorf("ReadParamsConfig: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return ParamsConfig{}, fmt.Errorf("ReadParamsConfig: %w", err)
	}
	return cfg, nil
}
