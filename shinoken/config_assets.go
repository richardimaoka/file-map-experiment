package shinoken

import (
	"fmt"

	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

// AssetsConfig is config.json for an assets case (see AssetsResolver layout).
type AssetsConfig struct {
	ParamsExcelFile string `json:"params_excel_file"`
	TextExcelFile   string `json:"text_excel_file"`
}

func CreateAssetsConfig(paramsExcelFile, textExcelFile string) AssetsConfig {
	return AssetsConfig{
		ParamsExcelFile: paramsExcelFile,
		TextExcelFile:   textExcelFile,
	}
}

func (c AssetsConfig) Validate() error {
	if c.ParamsExcelFile == "" {
		return fmt.Errorf("AssetsConfig.Validate: params_excel_file is not set")
	}
	if c.TextExcelFile == "" {
		return fmt.Errorf("AssetsConfig.Validate: text_excel_file is not set")
	}
	return nil
}

func ReadAssetsConfig(configPath string) (AssetsConfig, error) {
	var cfg AssetsConfig
	if err := jsonutil.Read(configPath, &cfg); err != nil {
		return AssetsConfig{}, fmt.Errorf("ReadAssetsConfig: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return AssetsConfig{}, fmt.Errorf("ReadAssetsConfig: %w", err)
	}
	return cfg, nil
}
