package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func setRawHeaderCell(f *excelize.File, sheet, col string, row int, value string) error {
	cell, err := excelize.JoinCellName(col, row)
	if err != nil {
		return err
	}
	return f.SetCellValue(sheet, cell, value)
}

func setRawHeaderRow(f *excelize.File, sheet string, row int, cols []string, vals []string) error {
	if len(cols) != len(vals) {
		return fmt.Errorf("setRawHeaderRow: cols and vals length mismatch")
	}
	for i, col := range cols {
		if err := setRawHeaderCell(f, sheet, col, row, vals[i]); err != nil {
			return err
		}
	}
	return nil
}

func mergeRawHeader(f *excelize.File, sheet, upperLeftCol string, upperLeftRow int, bottomRightCol string, bottomRightRow int) error {
	from, err := excelize.JoinCellName(upperLeftCol, upperLeftRow)
	if err != nil {
		return err
	}
	to, err := excelize.JoinCellName(bottomRightCol, bottomRightRow)
	if err != nil {
		return err
	}
	return f.MergeCell(sheet, from, to)
}
