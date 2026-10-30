package shinoken

import (
	"fmt"

	"github.com/richardimaoka/file-map-experiment/excel"
	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

type TextConfig struct {
	SchemaName string `json:"schema_name"`
	SheetName  string `json:"sheet_name"`
	excel.Config
}

func CreateTextConfig(excelFileName, schemaName, sheetName string) TextConfig {
	excelConfig := excel.InitConfig(excelFileName)
	excelConfig.SchemaPreps = []excel.SchemaPrep{
		{
			SchemaName:     schemaName,
			SheetName:      sheetName,
			IsReady:        true,
			HeaderStartRow: 2,
			HeaderEndRow:   2,
			DataStartRow:   3,
			DataEndRow:     33,
			StartColumn:    "A",
			EndColumn:      "Z",
		},
	}

	return TextConfig{
		SchemaName: schemaName,
		SheetName:  sheetName,
		Config:     excelConfig,
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

func (c TextConfig) SchemaPrep() (excel.SchemaPrep, error) {
	for _, p := range c.SchemaPreps {
		if p.SchemaName == c.SchemaName {
			return p, nil
		}
	}

	return excel.SchemaPrep{}, fmt.Errorf("TextConfig.SchemaPrep: schema %q not found in config.json", c.SchemaName)
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
