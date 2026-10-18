package excel

import "strings"

// include1RowRawHeaderColumn skips empty columns and horizontal-merge continuations (e.g. C when B1:C1).
func include1RowRawHeaderColumn(cv CellValue) bool {
	if cv.IsMerged() {
		// Merged cells often have empty Value except on the anchor cell; require both Value and merged_range.value.
		return strings.TrimSpace(cv.Value) != "" && strings.TrimSpace(cv.MergedRange.Value) != ""
	}
	return strings.TrimSpace(cv.Value) != ""
}

// includeMultiRowRawHeaderColumn skips columns with no content in any header row (merge rules TBD).
func includeMultiRowRawHeaderColumn(cellValues []CellValue) bool {
	for _, cv := range cellValues {
		if cv.IsMerged() || strings.TrimSpace(cv.Value) != "" {
			return true
		}
	}
	return false
}
