package excel

import (
	"fmt"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

// value is in string, because excelize treats all cell values as string, and we don't want to lose precision
type RowData map[string]string

// readRowsFromExcelInternal builds one RowData per data row from data_start_row through the
// last row with data in schema columns (see detectDataEndRow).
func readRowsFromExcelInternal(excelPath string, schema TableSchema) ([]RowData, error) {
	if err := schema.ValidateForDataRead(); err != nil {
		return nil, fmt.Errorf("readRowsFromExcelInternal: %w", err)
	}
	verticals, err := readRawDataForSchema(excelPath, schema)
	if err != nil {
		return nil, fmt.Errorf("readRowsFromExcelInternal: %w", err)
	}
	return rowDataFromRawVerticals(verticals, schema)
}

func rowDataFromRawVerticals(verticals []RawVertical, schema TableSchema) ([]RowData, error) {
	byCol := make(map[string]map[int]string, len(verticals))
	for _, v := range verticals {
		rowMap := make(map[int]string, len(v.CellValues))
		for _, cv := range v.CellValues {
			rowMap[cv.Position.RowNumber] = cv.Value
		}
		byCol[v.ExcelColumn] = rowMap
	}

	dataStart := schema.DataStartRow
	dataEnd := dataStart - 1
	if len(verticals) > 0 {
		dataStart = verticals[0].StartRow
		dataEnd = verticals[0].EndRow
	}

	n := dataEnd - dataStart + 1
	if n < 0 {
		n = 0
	}
	rows := make([]RowData, 0, n)
	for row := dataStart; row <= dataEnd; row++ {
		rd := RowData{}
		for _, c := range schema.Columns {
			val := ""
			if rowMap, ok := byCol[c.ExcelColumn]; ok {
				if v, ok := rowMap[row]; ok {
					val = v
				}
			}
			rd[c.HeaderName] = val
		}
		rows = append(rows, rd)
	}
	return rows, nil
}

func ReadRowsFromExcel(resolver PathResolver, schemaName string) ([]RowData, error) {
	config, err := ReadExcelConfig(resolver.ConfigPath())
	if err != nil {
		return nil, fmt.Errorf("ReadRowsFromExcel: %w", err)
	}
	schema, err := ReadSchemaFromJson(resolver.SchemaPath(schemaName))
	if err != nil {
		return nil, fmt.Errorf("ReadRowsFromExcel: %w", err)
	}
	if schema.SchemaName != schemaName {
		return nil, fmt.Errorf("ReadRowsFromExcel: schema_name %q does not match requested %q", schema.SchemaName, schemaName)
	}
	if err := schema.ValidateForDataRead(); err != nil {
		return nil, fmt.Errorf("ReadRowsFromExcel: %w", err)
	}
	excelPath := resolver.ExcelPath(config)
	return readRowsFromExcelInternal(excelPath, schema)
}

func ReadRowsFromJson(jsonPath string) ([]RowData, error) {
	var rows []RowData
	if err := jsonfile.Read(jsonPath, &rows); err != nil {
		return nil, fmt.Errorf("ReadRowsFromJson: %w", err)
	}
	return rows, nil
}

func ReadSingleRowFromJson(jsonPath string) ([]RowData, error) {
	return nil, fmt.Errorf("ReadSingleRowFromJson: not implemented")
}
