package excel

import "github.com/xuri/excelize/v2"

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

func loadMergedRanges(f *excelize.File, sheet string) ([]MergedRange, error) {
	cells, err := f.GetMergeCells(sheet)
	if err != nil {
		return nil, err
	}

	merges := make([]MergedRange, 0, len(cells))
	for _, mc := range cells {
		startCol, startRow, err := excelize.SplitCellName(mc.GetStartAxis())
		if err != nil {
			return nil, err
		}
		endCol, endRow, err := excelize.SplitCellName(mc.GetEndAxis())
		if err != nil {
			return nil, err
		}
		merges = append(merges, MergedRange{
			Value: mc.GetCellValue(),
			Range: NewCellRange(startCol, startRow, endCol, endRow),
		})
	}
	return merges, nil
}
