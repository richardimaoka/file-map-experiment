package excel

import "strings"

// includeRawHeaderColumn applies 1-row or multi-row inclusion rules for readRawColumns.
func includeRawHeaderColumn(prep SchemaPrep, cellValues []CellValue) bool {
	if prep.HeaderStartRow == prep.HeaderEndRow {
		return include1RowRawHeaderColumn(cellValues)
	}
	return includeMultiRowRawHeaderColumn(cellValues)
}

// include1RowRawHeaderColumn skips empty columns and horizontal-merge continuations (e.g. C when B1:C1).
func include1RowRawHeaderColumn(cellValues []CellValue) bool {
	hasContent := false
	for _, cv := range cellValues {
		if cv.IsMerged || strings.TrimSpace(cv.Value) != "" {
			hasContent = true
		}
		// Skip columns in a merged range except the left-most (e.g. C when B1:C1 is merged).
		if cv.IsMerged && cv.MergedRange != nil && cv.Position.ColumnName != cv.MergedRange.Range.UpperLeft.ColumnName {
			return false
		}
	}
	return hasContent
}

// includeMultiRowRawHeaderColumn skips columns with no content in any header row (merge rules TBD).
func includeMultiRowRawHeaderColumn(cellValues []CellValue) bool {
	for _, cv := range cellValues {
		if cv.IsMerged || strings.TrimSpace(cv.Value) != "" {
			return true
		}
	}
	return false
}
