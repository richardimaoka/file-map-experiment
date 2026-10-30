package shinoken

type ReverseConfig struct {
	FormatName string `json:"format_name"`
	SchemaName string `json:"scema_file"`
	SheetName  string `json:"sheet_name"`
}

// func CreateReverseConfig(formatName, schemaName, sheetName)

// func ReadReverseConfig(configPath)
// validate if the fields are non-empty
// if len(errors) > 0, return error
