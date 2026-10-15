package excel

import (
	"github.com/xuri/excelize/v2"
)

func isCellInMergeRange(cellAxis string, m excelize.MergeCell) bool {
	col, row, err := excelize.CellNameToCoordinates(cellAxis)
	if err != nil {
		return false
	}

	startCol, startRow, err := excelize.CellNameToCoordinates(m.GetStartAxis())
	if err != nil {
		return false
	}

	endCol, endRow, err := excelize.CellNameToCoordinates(m.GetEndAxis())
	if err != nil {
		return false
	}

	return col >= startCol && col <= endCol && row >= startRow && row <= endRow
}
