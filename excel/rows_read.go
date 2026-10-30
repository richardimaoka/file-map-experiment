package excel

import (
	"fmt"

	"github.com/richardimaoka/file-map-experiment/jsonutil"
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

	// byCol[excelColumn][rowNumber] = cell display value in the data band.
	byCol := make(map[string]map[int]string, len(verticals))
	for _, v := range verticals {
		rowMap := make(map[int]string, len(v.CellValues))
		for _, cv := range v.CellValues {
			rowMap[cv.Position.RowNumber] = cv.Value
		}
		byCol[v.ExcelColumn] = rowMap
	}

	// Verticals share the same start_row/end_row (the data band actually read, after detectDataEndRow).
	dataStart := schema.DataStartRow
	dataEnd := dataStart - 1
	if len(verticals) > 0 {
		dataStart = verticals[0].StartRow
		dataEnd = verticals[0].EndRow
	}

	// For each sheet row from dataStart through dataEnd, build one RowData (map header_name → value).
	// Every schema column gets a key; if the cell was empty or absent in byCol, the value is "".
	n := dataEnd - dataStart + 1
	if n < 0 {
		return nil, fmt.Errorf(
			"readRowsFromExcelInternal: data end row %d is before start row %d",
			dataEnd, dataStart,
		)
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

func ReadRowsFromExcel(resolver Resolver, schemaName string) ([]RowData, error) {
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
	excelPath := resolver.ExcelPath(config)
	return readRowsFromExcelInternal(excelPath, schema)
}

func ReadRowsFromJson(jsonPath string) ([]RowData, error) {
	var rows []RowData
	if err := jsonutil.Read(jsonPath, &rows); err != nil {
		return nil, fmt.Errorf("ReadRowsFromJson: %w", err)
	}
	return rows, nil
}

func ReadSingleRowFromJson(jsonPath string) ([]RowData, error) {
	return nil, fmt.Errorf("ReadSingleRowFromJson: not implemented")
}
