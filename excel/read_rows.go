package excel

import (
	"fmt"
)

// value is in string, because excelize treats all cell values as string, and we don't want to lose precision
type RowData map[string]string

func readRowsFromExcelInternal(excelPath string, schema TableSchema) ([]RowData, error) {
	return nil, fmt.Errorf("readRowsFromExcelInternal: not implemented")
}

func ReadRowsFromExcel(resolver PathResolver, schemaName string) ([]RowData, error) {
	return nil, fmt.Errorf("ReadRowsFromExcel: not implemented")
}

func ReadRowsFromJson(jsonPath string) ([]RowData, error) {
	return nil, fmt.Errorf("ReadRowsFromJson: not implemented")
}

func ReadSingleRowFromJson(jsonPath string) ([]RowData, error) {
	return nil, fmt.Errorf("ReadRowsFromJson: not implemented")
}
