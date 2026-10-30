package shinoken

type ForwardConfig struct {
	FormatName string `json:"format_name"`
	SchemaName string `json:"scema_file"`
	SheetName  string `json:"sheet_name"`
}

// func CreateForwardConfig(formatName, schemaName, sheetName)

// func ReadForwardConfig(configPath)
// validate if the fields are non-empty
// if len(errors) > 0, return error
