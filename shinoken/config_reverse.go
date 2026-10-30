package shinoken

import (
	"fmt"

	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

type ReverseConfig struct {
	FormatName string `json:"format_name"`
	SchemaName string `json:"scema_file"`
	SheetName  string `json:"sheet_name"`
}

func CreateReverseConfig(formatName, schemaName, sheetName string) ReverseConfig {
	return ReverseConfig{
		FormatName: formatName,
		SchemaName: schemaName,
		SheetName:  sheetName,
	}
}

func (c ReverseConfig) Validate() error {
	if c.FormatName == "" {
		return fmt.Errorf("ReverseConfig.Validate: format_name is not set")
	}
	if c.SchemaName == "" {
		return fmt.Errorf("ReverseConfig.Validate: scema_file is not set")
	}
	if c.SheetName == "" {
		return fmt.Errorf("ReverseConfig.Validate: sheet_name is not set")
	}
	return nil
}

func ReadReverseConfig(configPath string) (ReverseConfig, error) {
	var cfg ReverseConfig
	if err := jsonutil.Read(configPath, &cfg); err != nil {
		return ReverseConfig{}, fmt.Errorf("ReadReverseConfig: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return ReverseConfig{}, fmt.Errorf("ReadReverseConfig: %w", err)
	}
	return cfg, nil
}
