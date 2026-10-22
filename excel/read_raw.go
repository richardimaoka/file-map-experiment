package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type RawDataColumn struct {
	ExcelColumn  string      `json:"excel_column"`
	DataStartRow int         `json:"data_start_row"`
	DataEndRow   int         `json:"data_end_row"`
	CellValues   []CellValue `json:"cell_values"`
}

func (c *RawDataColumn) Type() (string, error) {
	return "", fmt.Errorf("RawDataColumn.Type: not implemented")
}

type RawTable struct {
	HeaderVerticals []RawVertical `json:"header_verticals"`
	DataVerticals   []RawVertical `json:"data_verticals"`
}

func readRawTable(excelPath string, sheetName string, hRange, dRange CellRange) (RawTable, error) {
	h, err := readRawVerticals(excelPath, sheetName, hRange)
	if err != nil {
		return RawTable{}, fmt.Errorf("readRawTable: failed to read header verticals: %w", err)
	}

	d, err := readRawVerticals(excelPath, sheetName, dRange)
	if err != nil {
		return RawTable{}, fmt.Errorf("readRawTable: failed to read data verticals: %w", err)
	}

	return RawTable{
		HeaderVerticals: h,
		DataVerticals:   d,
	}, nil

}

type RawVertical struct {
	ExcelColumn string      `json:"excel_column"`
	StartRow    int         `json:"start_row"`
	EndRow      int         `json:"end_row"`
	CellValues  []CellValue `json:"cell_values"`
}

func readRawVerticals(excelPath string, sheetName string, r CellRange) ([]RawVertical, error) {
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, fmt.Errorf("readRawVerticals: %w", err)
	}
	defer f.Close()

	merges, err := loadMergedRanges(f, sheetName)
	if err != nil {
		return nil, fmt.Errorf("readRawVerticals: %w", err)
	}
	if err := rejectMergesCrossingCellRange(merges, r); err != nil {
		return nil, err
	}

	startRow := r.StartRowNumber()
	endRow := r.EndRowNumber()
	startCol, err := r.StartColumnNumber()
	if err != nil {
		return nil, fmt.Errorf("readRawVerticals: start_column: %w", err)
	}
	endCol, err := r.EndColumnNumber()
	if err != nil {
		return nil, fmt.Errorf("readRawVerticals: end_column: %w", err)
	}

	verticals := make([]RawVertical, 0, endCol-startCol+1)
	for colNum := startCol; colNum <= endCol; colNum++ {
		colName, err := excelize.ColumnNumberToName(colNum)
		if err != nil {
			return nil, fmt.Errorf("readRawVerticals: %w", err)
		}

		cellValues := make([]CellValue, 0, endRow-startRow+1)
		for row := startRow; row <= endRow; row++ {
			cv, err := readCellValue(f, sheetName, colName, row, merges)
			if err != nil {
				return nil, fmt.Errorf("readRawVerticals: %w", err)
			}
			cellValues = append(cellValues, cv)
		}

		verticals = append(verticals, RawVertical{
			ExcelColumn: colName,
			StartRow:    startRow,
			EndRow:      endRow,
			CellValues:  cellValues,
		})
	}

	return verticals, nil
}

func rejectMergesCrossingCellRange(merges []MergedRange, r CellRange) error {
	for _, m := range merges {
		hits, err := r.Intersects(m.Range)
		if err != nil {
			return fmt.Errorf("readRawVerticals: %w", err)
		}
		if !hits {
			continue
		}
		inside, err := r.Contains(m.Range)
		if err != nil {
			return fmt.Errorf("readRawVerticals: %w", err)
		}
		if inside {
			continue
		}
		return fmt.Errorf(
			"readRawVerticals: merged range %s must lie entirely inside read range %s",
			mergeRangeLabel(m.Range),
			mergeRangeLabel(r),
		)
	}
	return nil
}

func mergeRangeLabel(r CellRange) string {
	ul := r.UpperLeft
	br := r.BottomRight
	return fmt.Sprintf("%s%d:%s%d", ul.ColumnName, ul.RowNumber, br.ColumnName, br.RowNumber)
}

func readRawDataForSchema(excelPath string, schema TableSchema) ([]RawDataColumn, error) {
	if err := schema.ValidateForDataRead(); err != nil {
		return nil, fmt.Errorf("readRawDataForSchema: %w", err)
	}

	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, fmt.Errorf("readRawDataForSchema: %w", err)
	}
	defer f.Close()

	merges, err := loadMergedRanges(f, schema.SheetName)
	if err != nil {
		return nil, fmt.Errorf("readRawDataForSchema: %w", err)
	}
	dataEndRow, err := detectDataEndRow(f, schema, merges)
	if err != nil {
		return nil, fmt.Errorf("readRawDataForSchema: %w", err)
	}

	verticals, err := readRawVerticals(excelPath, schema.SheetName, schema.DataRange(dataEndRow))
	if err != nil {
		return nil, fmt.Errorf("readRawDataForSchema: %w", err)
	}

	columns := make([]RawDataColumn, len(verticals))
	for i, v := range verticals {
		columns[i] = RawDataColumn{
			ExcelColumn:  v.ExcelColumn,
			DataStartRow: schema.DataStartRow,
			DataEndRow:   dataEndRow,
			CellValues:   v.CellValues,
		}
	}
	return columns, nil
}

func readCellValue(f *excelize.File, sheet, col string, row int, merges []MergedRange) (CellValue, error) {
	pos := CellPosition{ColumnName: col, RowNumber: row}
	cellName, err := excelize.JoinCellName(col, row)
	if err != nil {
		return CellValue{}, err
	}

	value, err := f.GetCellValue(sheet, cellName)
	if err != nil {
		return CellValue{}, err
	}
	formula, err := f.GetCellFormula(sheet, cellName)
	if err != nil {
		return CellValue{}, err
	}

	var merged *MergedRange
	for _, m := range merges {
		if !isCellInMergeRange(pos, m) {
			continue
		}
		// Copy before taking a pointer so we don't alias the range loop variable.
		mr := m
		merged = &mr
		break
	}

	return CellValue{
		Value:       value,
		Formula:     formula,
		Position:    pos,
		MergedRange: merged,
	}, nil
}
