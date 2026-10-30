package shinoken

// AssetsConfig is config.json for an assets case (see AssetsResolver layout).
type AssetsConfig struct {
	ParamsExcelFile string `json:"params_excel_file"`
	TextExcelFile   string `json:"text_excel_file"`
}

// func CreateAssetConfig(paramsExcelFile, textExcelFile)

// func ReadAssetConfig(configPath)
