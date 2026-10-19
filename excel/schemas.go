package excel

import "fmt"

type ColumnSchema struct {
	HeaderName   string `json:"header_name"`
	ExcelColumn  string `json:"excel_column"`
	IsRowGrouped bool   `json:"is_row_grouped"`
}

type TableSchema struct {
	SchemaName   string         `json:"schema_name"`
	DataStartRow int            `json:"data_start_row"`
	StartColumn  string         `json:"start_column"`
	EndColumn    string         `json:"end_column"`
	Columns      []ColumnSchema `json:"columns"`
}

func readSchemaInternal(excelPath string, prep SchemaPrep) (TableSchema, error) {
	hRange := prep.HeaderRange()
	dRange := prep.DataRange()

	_, err := readRawTable(excelPath, prep.SheetName, hRange, dRange)
	if err != nil {
		return TableSchema{}, fmt.Errorf("readSchemaInternal: %w", err)
	}

	return TableSchema{}, fmt.Errorf("readSchemaInternal: not implemented")
}

func validateExcelAgainstSchema(resolver PathResolver, schemaName string) (TableSchema, error) {
	return TableSchema{}, fmt.Errorf("validateExcelAgainstSchema: not implemented")
}

func ReadSchemaFromExcel(resolver PathResolver, schemaName string) (TableSchema, error) {
	return TableSchema{}, fmt.Errorf("ReadSchemaFromExcel: not implemented")
}

func ReadSchemaFromJson(jsonPath string) (TableSchema, error) {
	return TableSchema{}, fmt.Errorf("ReadSchemaFromJson: not implemented")
}
