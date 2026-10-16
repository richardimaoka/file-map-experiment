package excel

import "fmt"

type SingleRowHeaderRaw struct {
	ExcelColumn string `json:"excel_column"`
	HeaderName  string `json:"header_name"`
}

type MultiRowHeaderRaw struct {
	ExcelColumn string   `json:"excel_column"`
	CellValues  []string `json:"cell_values"`
}

func (h *MultiRowHeaderRaw) ExtractHeaderName() (string, error) {
	return "", fmt.Errorf("extractHeaderName: not implemented")
}

type DataColumnRaw struct {
	ExcelColumn  string `json:"excel_column"`
	DataStartRow int    `json:"data_start_row"`
	IsMerged     bool   `json:"is_merged"`
}

func (c *DataColumnRaw) DetectColumnType() (string, error) {
	return "", fmt.Errorf("extractHeaderName: not implemented")
}

// isCellInMergeRange reports whether pos lies inside m.Range (inclusive on all edges).
func isCellInMergeRange(pos CellPosition, m MergedRange) bool {
	col, err := pos.ColumnNumber()
	if err != nil || pos.RowNumber < 1 {
		return false
	}

	startCol, err := m.Range.UpperLeft.ColumnNumber()
	if err != nil {
		return false
	}
	endCol, err := m.Range.BottomRight.ColumnNumber()
	if err != nil {
		return false
	}
	startRow := m.Range.UpperLeft.RowNumber
	endRow := m.Range.BottomRight.RowNumber
	if startRow < 1 || endRow < 1 {
		return false
	}

	row := pos.RowNumber
	return col >= startCol && col <= endCol && row >= startRow && row <= endRow
}
