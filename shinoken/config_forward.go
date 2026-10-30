package shinoken

import (
	"fmt"

	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

type ForwardConfig struct {
	FormatName string `json:"format_name"`
	SchemaName string `json:"scema_file"`
	SheetName  string `json:"sheet_name"`
}

func CreateForwardConfig(formatName, schemaName, sheetName string) ForwardConfig {
	return ForwardConfig{
		FormatName: formatName,
		SchemaName: schemaName,
		SheetName:  sheetName,
	}
}

func (c ForwardConfig) Validate() error {
	if c.FormatName == "" {
		return fmt.Errorf("ForwardConfig.Validate: format_name is not set")
	}
	if c.SchemaName == "" {
		return fmt.Errorf("ForwardConfig.Validate: scema_file is not set")
	}
	if c.SheetName == "" {
		return fmt.Errorf("ForwardConfig.Validate: sheet_name is not set")
	}
	return nil
}

func ReadForwardConfig(configPath string) (ForwardConfig, error) {
	var cfg ForwardConfig
	if err := jsonutil.Read(configPath, &cfg); err != nil {
		return ForwardConfig{}, fmt.Errorf("ReadForwardConfig: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return ForwardConfig{}, fmt.Errorf("ReadForwardConfig: %w", err)
	}
	return cfg, nil
}
