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
	if err := prep.Validate(); err != nil {
		return TableSchema{}, fmt.Errorf("readSchemaInternal: %w", err)
	}

	hRange := prep.HeaderRange()
	dRange := prep.DataRange()

	raw, err := readRawTable(excelPath, prep.SheetName, hRange, dRange)
	if err != nil {
		return TableSchema{}, fmt.Errorf("readSchemaInternal: %w", err)
	}

	columns, err := columnSchemasFromHeader(prep, raw.HeaderVerticals)
	if err != nil {
		return TableSchema{}, fmt.Errorf("readSchemaInternal: %w", err)
	}

	return TableSchema{
		SchemaName:   prep.SchemaName,
		DataStartRow: prep.DataStartRow,
		StartColumn:  prep.StartColumn,
		EndColumn:    prep.EndColumn,
		Columns:      columns,
	}, nil
}

func columnSchemasFromHeader(prep SchemaPrep, verticals []RawVertical) ([]ColumnSchema, error) {
	if prep.HeaderStartRow == prep.HeaderEndRow {
		return columnSchemasFromSingleRowHeader(verticals, prep.HeaderEndRow)
	}
	return columnSchemasFromMultiRowHeader(prep, verticals)
}

func columnSchemasFromSingleRowHeader(verticals []RawVertical, headerRow int) ([]ColumnSchema, error) {
	columns := make([]ColumnSchema, 0, len(verticals))
	for _, v := range verticals {
		cv, err := headerCellValueAtRow(v, headerRow)
		if err != nil {
			return nil, err
		}
		if cv.Value == "" {
			continue
		}
		if cv.MergedRange != nil && !isMergeAnchorCell(cv) {
			continue
		}

		columns = append(columns, ColumnSchema{
			HeaderName:   cv.Value,
			ExcelColumn:  v.ExcelColumn,
			IsRowGrouped: false, // data-range row merge grouping; set when data schema is implemented
		})
	}

	return columns, nil
}

func columnSchemasFromMultiRowHeader(prep SchemaPrep, verticals []RawVertical) ([]ColumnSchema, error) {
	return nil, fmt.Errorf("columnSchemasFromMultiRowHeader: not implemented")
}

func headerCellValueAtRow(v RawVertical, row int) (CellValue, error) {
	for _, cv := range v.CellValues {
		if cv.Position.RowNumber == row {
			return cv, nil
		}
	}
	return CellValue{}, fmt.Errorf("headerCellValueAtRow: row %d not in vertical %s", row, v.ExcelColumn)
}

func isMergeAnchorCell(cv CellValue) bool {
	if cv.MergedRange == nil {
		return true
	}
	ul := cv.MergedRange.Range.UpperLeft
	return cv.Position.ColumnName == ul.ColumnName && cv.Position.RowNumber == ul.RowNumber
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
