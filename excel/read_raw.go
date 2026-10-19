package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

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
	CellValues   []CellValue `json:"cell_values"`
}

func (c *RawDataColumn) Type() (string, error) {
	return "", fmt.Errorf("RawDataColumn.Type: not implemented")
}

type RawColumn struct {
	Header *RawHeader     `json:"header"`
	Data   *RawDataColumn `json:"data"`
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

func readRawHeaders(excelPath string, prep SchemaPrep) ([]RawHeader, error) {
	if err := prep.Validate(); err != nil {
		return nil, fmt.Errorf("readRawHeaders: %w", err)
	}

	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, fmt.Errorf("readRawHeaders: %w", err)
	}
	defer f.Close()

	sheet := prep.SheetName
	if idx, err := f.GetSheetIndex(sheet); err != nil || idx < 0 {
		return nil, fmt.Errorf("readRawHeaders: sheet %q not found", sheet)
	}

	merges, err := loadMergedRanges(f, sheet)
	if err != nil {
		return nil, fmt.Errorf("readRawHeaders: %w", err)
	}

	startCol, err := excelize.ColumnNameToNumber(prep.StartColumn)
	if err != nil {
		return nil, fmt.Errorf("readRawHeaders: start_column: %w", err)
	}
	endCol, err := excelize.ColumnNameToNumber(prep.EndColumn)
	if err != nil {
		return nil, fmt.Errorf("readRawHeaders: end_column: %w", err)
	}

	headers := make([]RawHeader, 0, endCol-startCol+1)
	for colNum := startCol; colNum <= endCol; colNum++ {
		colName, err := excelize.ColumnNumberToName(colNum)
		if err != nil {
			return nil, fmt.Errorf("readRawHeaders: %w", err)
		}

		cellValues := make([]CellValue, 0, prep.HeaderEndRow-prep.HeaderStartRow+1)
		for row := prep.HeaderStartRow; row <= prep.HeaderEndRow; row++ {
			cv, err := readCellValue(f, sheet, colName, row, merges)
			if err != nil {
				return nil, fmt.Errorf("readRawHeaders: %w", err)
			}
			cellValues = append(cellValues, cv)
		}

		headers = append(headers, RawHeader{
			ExcelColumn:    colName,
			HeaderStartRow: prep.HeaderStartRow,
			HeaderEndRow:   prep.HeaderEndRow,
			CellValues:     cellValues,
		})
	}

	return headers, nil
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
