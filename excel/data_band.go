package excel

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// DataRange returns the column band and row span for reading table data.
// dataEndRow is usually from detectDataEndRow (last row with data in schema columns).
func (s TableSchema) DataRange(dataEndRow int) CellRange {
	return NewCellRange(
		s.StartColumn,
		s.DataStartRow,
		s.EndColumn,
		dataEndRow,
	)
}

// detectDataEndRow returns the last row at or below data_start_row where any schema
// column has a non-empty display value. If there is no data, returns data_start_row - 1.
func detectDataEndRow(f *excelize.File, schema TableSchema, merges []MergedRange) (int, error) {
	if len(schema.Columns) == 0 {
		return schema.DataStartRow - 1, nil
	}

	maxRow, err := sheetLastRow(f, schema.SheetName)
	if err != nil {
		return 0, fmt.Errorf("detectDataEndRow: %w", err)
	}
	if maxRow < schema.DataStartRow {
		return schema.DataStartRow - 1, nil
	}

	last := schema.DataStartRow - 1
	for row := schema.DataStartRow; row <= maxRow; row++ {
		for _, col := range schema.Columns {
			cv, err := readCellValue(f, schema.SheetName, col.ExcelColumn, row, merges)
			if err != nil {
				return 0, fmt.Errorf("detectDataEndRow: %w", err)
			}
			if cv.Value != "" {
				last = row
				break
			}
		}
	}
	return last, nil
}

func sheetLastRow(f *excelize.File, sheet string) (int, error) {
	maxRow := 0

	dim, err := f.GetSheetDimension(sheet)
	if err != nil {
		return 0, err
	}
	dim = strings.TrimSpace(dim)
	if dim != "" {
		parts := strings.Split(dim, ":")
		br := strings.TrimSpace(parts[len(parts)-1])
		_, row, err := excelize.SplitCellName(br)
		if err != nil {
			return 0, fmt.Errorf("sheetLastRow: parse dimension %q: %w", dim, err)
		}
		if row > maxRow {
			maxRow = row
		}
	}

	// GetSheetDimension is often stale (e.g. "A1" on small fixtures); GetRows reflects used rows.
	rows, err := f.GetRows(sheet)
	if err != nil {
		return 0, fmt.Errorf("sheetLastRow: %w", err)
	}
	if len(rows) > maxRow {
		maxRow = len(rows)
	}
	return maxRow, nil
}
