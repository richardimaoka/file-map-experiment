package excel

import "fmt"

type RawHeader struct {
	ExcelColumn    string      `json:"excel_column"`
	HeaderStartRow int         `json:"header_start_row"`
	HeaderEndRow   int         `json:"header_end_row"`
	CellValues     []CellValue `json:"cell_values"`
}

func (h *RawHeader) Type() (string, error) {
	return "", fmt.Errorf("RawHeader.Type: not implemented")
}

func (h *RawHeader) HeaderName() (string, error) {
	return "", fmt.Errorf("HeaderName: not implemented")
}

type RawDataColumn struct {
	ExcelColumn  string      `json:"excel_column"`
	DataStartRow int         `json:"data_start_row"`
	DataEndRow   int         `json:"data_end_row"`
	IsMerged     bool        `json:"is_merged"`
	CellValues   []CellValue `json:"cell_values"`
}

func (c *RawDataColumn) Type() (string, error) {
	return "", fmt.Errorf("RawDataColumn.Type: not implemented")
}

type RawColumn struct {
	Header *RawHeader     `json:"header"`
	Data   *RawDataColumn `json:"data"`
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

func readRawColumns(excelPath string, prep *SchemaPrep) ([]RawColumn, error) {
	return nil, fmt.Errorf("readRawColumns: not implemented")
}
