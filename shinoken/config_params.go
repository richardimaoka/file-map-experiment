package shinoken

import "github.com/richardimaoka/file-map-experiment/excel"

type ParamsConfig struct {
	SchemaName string `json:"scema_file"`
	SheetName  string `json:"sheet_name"`
	excel.Config
}

// func CreateParamsConfig(schemaName, sheetName)

// func ReadParamsConfig(configPath)
// validate if the fields are non-empty
// if len(errors) > 0, return error
