package excel

import (
	"fmt"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

// value is in string, because excelize treats all cell values as string, and we don't want to lose precision
type RowData map[string]string

// readRowsFromExcelInternal builds one RowData per physical data row (data_start_row–data_end_row).
// Cell values come from readRawData (merge anchors and excelize display strings).
func readRowsFromExcelInternal(excelPath string, prep SchemaPrep, schema TableSchema) ([]RowData, error) {
	rawCols, err := readRawData(excelPath, prep)
	if err != nil {
		return nil, fmt.Errorf("readRowsFromExcelInternal: %w", err)
	}
	return rowDataFromRawColumns(rawCols, prep, schema)
}

func rowDataFromRawColumns(rawCols []RawDataColumn, prep SchemaPrep, schema TableSchema) ([]RowData, error) {
	byCol := make(map[string]map[int]string, len(rawCols))
	for _, col := range rawCols {
		rowMap := make(map[int]string, len(col.CellValues))
		for _, cv := range col.CellValues {
			rowMap[cv.Position.RowNumber] = cv.Value
		}
		byCol[col.ExcelColumn] = rowMap
	}

	n := prep.DataEndRow - prep.DataStartRow + 1
	rows := make([]RowData, 0, n)
	for row := prep.DataStartRow; row <= prep.DataEndRow; row++ {
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
	return nil, fmt.Errorf("ReadRowsFromExcel: not implemented")
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
