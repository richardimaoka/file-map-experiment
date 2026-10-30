package shinoken

import (
	"fmt"

	"github.com/richardimaoka/file-map-experiment/excel"
	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

type TextConfig struct {
	SchemaName string `json:"scema_file"`
	SheetName  string `json:"sheet_name"`
	excel.Config
}

func CreateTextConfig(excelFileName, schemaName, sheetName string) TextConfig {
	return TextConfig{
		SchemaName: schemaName,
		SheetName:  sheetName,
		Config:     excel.InitConfig(excelFileName),
	}
}

func (c TextConfig) Validate() error {
	if c.SchemaName == "" {
		return fmt.Errorf("TextConfig.Validate: scema_file is not set")
	}
	if c.SheetName == "" {
		return fmt.Errorf("TextConfig.Validate: sheet_name is not set")
	}
	if c.ExcelFileName == "" {
		return fmt.Errorf("TextConfig.Validate: excel_file_name is not set")
	}
	return nil
}

func ReadTextConfig(configPath string) (TextConfig, error) {
	var cfg TextConfig
	if err := jsonutil.Read(configPath, &cfg); err != nil {
		return TextConfig{}, fmt.Errorf("ReadTextConfig: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return TextConfig{}, fmt.Errorf("ReadTextConfig: %w", err)
	}
	return cfg, nil
}
