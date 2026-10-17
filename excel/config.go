package excel

import "fmt"

// SchemaPrep is needed before reading Excel table schema
type SchemaPrep struct {
	SchemaName     string `json:"schema_name"`
	IsReady        bool   `json:"is_ready"`
	HeaderStartRow int    `json:"header_start_row"`
	HeaderEndRow   int    `json:"header_end_row"`
	DataStartRow   int    `json:"data_start_row"`
	StartColumn    string `json:"start_column"`
	EndColumn      string `json:"end_column"`
}

func (c *SchemaPrep) HeaderRange() CellRange {
	return NewCellRange(
		c.StartColumn,
		c.HeaderStartRow,
		c.EndColumn,
		c.HeaderEndRow,
	)
}

type ExcelConfig struct {
	ExcelFileName string       `json:"excel_file_name"`
	SchemaPreps   []SchemaPrep `json:"schema_preps"` // manually updated in config.json
}

func (c *ExcelConfig) GetSchemaPrep(schemaName string) (*SchemaPrep, error) {
	for _, prep := range c.SchemaPreps {
		if prep.SchemaName == schemaName {
			return &prep, nil
		}
	}

	return nil, fmt.Errorf("ExcelConfig.GetSchemaPrep: schema %q not found in config", schemaName)
}

// The excelFileName argument is needed, assuming you already know
// which Excel file to read upon creating config.json
func InitConfig(excelFileName string) *ExcelConfig {
	return &ExcelConfig{
		ExcelFileName: excelFileName,
	}
}

func ReadConfigFromJson(jsonPath string) (*ExcelConfig, error) {
	return nil, fmt.Errorf("ReadConfigFromJson: not implemented")
}
