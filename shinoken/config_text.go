package shinoken

import "github.com/richardimaoka/file-map-experiment/excel"

type TextConfig struct {
	SchemaName string `json:"scema_file"`
	SheetName  string `json:"sheet_name"`
	excel.Config
}

// func CreateTextConfig(schemaName, sheetName)

// func ReadTextConfig(configPath)
// validate if the fields are non-empty
// if len(errors) > 0, return error
