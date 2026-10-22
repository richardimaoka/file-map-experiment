package excel

import (
	"fmt"
	"strings"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

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

// ReadTableSchema builds TableSchema from a workbook and prep.
func ReadTableSchema(excelPath string, prep SchemaPrep) (TableSchema, error) {
	return readSchemaInternal(excelPath, prep)
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

	if err := rejectMultiColumnMergesIntersectingDataBand(raw.DataVerticals, prep.DataStartRow, prep.DataEndRow); err != nil {
		return TableSchema{}, fmt.Errorf("readSchemaInternal: %w", err)
	}

	columns, err := columnSchemasFromHeader(prep, raw.HeaderVerticals)
	if err != nil {
		return TableSchema{}, fmt.Errorf("readSchemaInternal: %w", err)
	}

	if err := rejectDuplicateHeaderNames(columns); err != nil {
		return TableSchema{}, fmt.Errorf("readSchemaInternal: %w", err)
	}

	if err := validateDataColumnsHaveRowDataKeys(columns, raw.DataVerticals); err != nil {
		return TableSchema{}, fmt.Errorf("readSchemaInternal: %w", err)
	}

	columns = applyIsRowGroupedFromData(columns, raw.DataVerticals)

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

// returned []ColumnSchema is in an intermediate state, which will be enriched later.
func columnSchemasFromSingleRowHeader(verticals []RawVertical, headerRow int) ([]ColumnSchema, error) {
	if err := rejectMultiCellMergesInSingleRowHeader(verticals, headerRow); err != nil {
		return nil, err
	}

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
			HeaderName:  cv.Value,
			ExcelColumn: v.ExcelColumn,
			//IsRowGrouped will be calculated in a separate function.
		})
	}

	return columns, nil
}

func rejectMultiCellMergesInSingleRowHeader(verticals []RawVertical, headerRow int) error {
	for _, v := range verticals {
		cv, err := headerCellValueAtRow(v, headerRow)
		if err != nil {
			return err
		}
		if cv.MergedRange == nil {
			continue
		}
		if !mergeSpansMultipleCells(*cv.MergedRange) {
			continue
		}
		return fmt.Errorf(
			"merged header cells are not allowed for single-row headers (merge %s)",
			mergeRangeLabel(cv.MergedRange.Range),
		)
	}
	return nil
}

func mergeSpansMultipleCells(m MergedRange) bool {
	ul := m.Range.UpperLeft
	br := m.Range.BottomRight
	return ul.ColumnName != br.ColumnName || ul.RowNumber != br.RowNumber
}

func mergeSpansMultipleColumns(m MergedRange) bool {
	ul := m.Range.UpperLeft
	br := m.Range.BottomRight
	return ul.ColumnName != br.ColumnName
}

func mergeIntersectsRow(m MergedRange, row int) bool {
	ul := m.Range.UpperLeft
	br := m.Range.BottomRight
	return ul.RowNumber <= row && row <= br.RowNumber
}

// Each data row must map to one column at the finest granularity; multi-column merges in the data band are ambiguous.
func rejectMultiColumnMergesIntersectingDataBand(verticals []RawVertical, dataStartRow, dataEndRow int) error {
	seen := make(map[string]bool)
	for _, v := range verticals {
		for _, cv := range v.CellValues {
			if cv.MergedRange == nil {
				continue
			}
			m := *cv.MergedRange
			if !mergeSpansMultipleColumns(m) {
				continue
			}
			if !mergeIntersectsRowRange(m, dataStartRow, dataEndRow) {
				continue
			}
			label := mergeRangeLabel(m.Range)
			if seen[label] {
				continue
			}
			seen[label] = true
			return fmt.Errorf(
				"multi-column merge intersects the data band (merge %s)",
				label,
			)
		}
	}
	return nil
}

func mergeIntersectsRowRange(m MergedRange, startRow, endRow int) bool {
	ul := m.Range.UpperLeft
	br := m.Range.BottomRight
	return ul.RowNumber <= endRow && br.RowNumber >= startRow
}

// Leaf header row must name each column; a merge across columns on that row (or a 2D merge
// that includes it) makes RowData keys ambiguous.
func rejectMultiColumnMergesIntersectingHeaderEndRow(verticals []RawVertical, headerEndRow int) error {
	seen := make(map[string]bool)
	for _, v := range verticals {
		for _, cv := range v.CellValues {
			if cv.MergedRange == nil {
				continue
			}
			m := *cv.MergedRange
			if !mergeSpansMultipleColumns(m) {
				continue
			}
			if !mergeIntersectsRow(m, headerEndRow) {
				continue
			}
			label := mergeRangeLabel(m.Range)
			if seen[label] {
				continue
			}
			seen[label] = true
			return fmt.Errorf(
				"multi-column header merge intersects the leaf header row (merge %s)",
				label,
			)
		}
	}
	return nil
}

func columnSchemasFromMultiRowHeader(prep SchemaPrep, verticals []RawVertical) ([]ColumnSchema, error) {
	if err := rejectMultiColumnMergesIntersectingHeaderEndRow(verticals, prep.HeaderEndRow); err != nil {
		return nil, err
	}

	columns := make([]ColumnSchema, 0, len(verticals))
	for _, v := range verticals {
		name := resolvedMultiRowHeaderName(v, prep.HeaderStartRow, prep.HeaderEndRow)
		if name == "" {
			continue
		}
		columns = append(columns, ColumnSchema{
			HeaderName:  name,
			ExcelColumn: v.ExcelColumn,
		})
	}
	return columns, nil
}

// RowData keys join each header row's label top-to-bottom with " > ".
func resolvedMultiRowHeaderName(v RawVertical, headerStartRow, headerEndRow int) string {
	parts := make([]string, 0, headerEndRow-headerStartRow+1)
	for row := headerStartRow; row <= headerEndRow; row++ {
		cv, err := headerCellValueAtRow(v, row)
		if err != nil {
			continue
		}
		if !headerRowLabelIncluded(cv) {
			continue
		}
		parts = append(parts, cv.Value)
	}
	return strings.Join(parts, " > ")
}

// Unmerged cells and merge anchors always count. Non-anchor cells count only for horizontal
// merges on that row (group labels spanning columns). Vertical / 2D merge fill cells are skipped.
func headerRowLabelIncluded(cv CellValue) bool {
	if cv.Value == "" {
		return false
	}
	if cv.MergedRange == nil || isMergeAnchorCell(cv) {
		return true
	}
	return mergeIsHorizontalOnRow(*cv.MergedRange, cv.Position.RowNumber)
}

func mergeIsHorizontalOnRow(m MergedRange, row int) bool {
	ul := m.Range.UpperLeft
	br := m.Range.BottomRight
	return ul.RowNumber == row && br.RowNumber == row && ul.ColumnName != br.ColumnName
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

// RowData keys are header names; two columns with the same non-empty name are ambiguous.
func rejectDuplicateHeaderNames(columns []ColumnSchema) error {
	byName := make(map[string]string, len(columns))
	for _, c := range columns {
		if c.HeaderName == "" {
			continue
		}
		if firstCol, ok := byName[c.HeaderName]; ok {
			return fmt.Errorf(
				"duplicate header name %q (columns %s and %s)",
				c.HeaderName, firstCol, c.ExcelColumn,
			)
		}
		byName[c.HeaderName] = c.ExcelColumn
	}
	return nil
}

// Every Excel column with non-empty data cells must map to a ColumnSchema with a non-empty
// header_name so rows can be encoded as RowData (flat string map keyed by header).
func validateDataColumnsHaveRowDataKeys(columns []ColumnSchema, dataVerticals []RawVertical) error {
	byExcelColumn := make(map[string]ColumnSchema, len(columns))
	for _, c := range columns {
		byExcelColumn[c.ExcelColumn] = c
	}

	for _, dv := range dataVerticals {
		if !rawVerticalHasNonEmptyCell(dv) {
			continue
		}
		col, ok := byExcelColumn[dv.ExcelColumn]
		if !ok || col.HeaderName == "" {
			return fmt.Errorf(
				"column %s has data cells but no header name for RowData key",
				dv.ExcelColumn,
			)
		}
	}
	return nil
}

func rawVerticalHasNonEmptyCell(v RawVertical) bool {
	for _, cv := range v.CellValues {
		if cv.Value != "" {
			return true
		}
	}
	return false
}

// IsRowGrouped is true when the data band has a single-column merge spanning multiple rows
// (Excel row grouping / repeated values under one anchor cell).
func applyIsRowGroupedFromData(columns []ColumnSchema, dataVerticals []RawVertical) []ColumnSchema {
	byCol := make(map[string]RawVertical, len(dataVerticals))
	for _, dv := range dataVerticals {
		byCol[dv.ExcelColumn] = dv
	}
	out := make([]ColumnSchema, len(columns))
	for i, c := range columns {
		c.IsRowGrouped = false
		if dv, ok := byCol[c.ExcelColumn]; ok {
			c.IsRowGrouped = dataColumnIsRowGrouped(dv)
		}
		out[i] = c
	}
	return out
}

func dataColumnIsRowGrouped(v RawVertical) bool {
	for _, cv := range v.CellValues {
		if cv.MergedRange == nil {
			continue
		}
		m := *cv.MergedRange
		if !mergeSpansMultipleRows(m) || !mergeIsSingleColumn(m) {
			continue
		}
		if m.Range.UpperLeft.ColumnName == v.ExcelColumn {
			return true
		}
	}
	return false
}

func mergeSpansMultipleRows(m MergedRange) bool {
	ul := m.Range.UpperLeft
	br := m.Range.BottomRight
	return ul.RowNumber != br.RowNumber
}

func mergeIsSingleColumn(m MergedRange) bool {
	ul := m.Range.UpperLeft
	br := m.Range.BottomRight
	return ul.ColumnName == br.ColumnName
}

func validateExcelAgainstSchema(resolver PathResolver, schemaName string) (TableSchema, error) {
	return TableSchema{}, fmt.Errorf("validateExcelAgainstSchema: not implemented")
}

func ReadSchemaFromExcel(resolver PathResolver, schemaName string) (TableSchema, error) {
	config, err := ReadExcelConfig(resolver.ConfigPath())
	if err != nil {
		return TableSchema{}, fmt.Errorf("ReadSchemaFromExcel: %w", err)
	}
	prep, err := config.GetSchemaPrep(schemaName)
	if err != nil {
		return TableSchema{}, fmt.Errorf("ReadSchemaFromExcel: %w", err)
	}
	if err := prep.Validate(); err != nil {
		return TableSchema{}, fmt.Errorf("ReadSchemaFromExcel: %w", err)
	}
	excelPath := resolver.ExcelPath(config)
	return ReadTableSchema(excelPath, prep)
}

func ReadSchemaFromJson(jsonPath string) (TableSchema, error) {
	var schema TableSchema
	if err := jsonfile.Read(jsonPath, &schema); err != nil {
		return TableSchema{}, fmt.Errorf("ReadSchemaFromJson: %w", err)
	}
	return schema, nil
}
