package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type CellPosition struct {
	ColumnName string `json:"column_name"`
	RowNumber  int    `json:"row_number"`
}

func (c *CellPosition) CellName() string {
	return fmt.Sprintf("%s%ds", c.ColumnName, c.RowNumber)
}

func (c *CellPosition) ColumnNumber() (int, error) {
	n, err := excelize.ColumnNameToNumber(c.ColumnName)
	if err != nil {
		return 0, fmt.Errorf("CellPosition.ColumnNumber: %w", err)
	}

	return n, nil
}

type CellValue struct {
	Value       string       `json:"value"`
	Formula     string       `json:"formula"`
	Position    CellPosition `json:"position"`
	MergedRange *MergedRange `json:"merged_range,omitempty"`
}

func (cv CellValue) IsMerged() bool {
	return cv.MergedRange != nil
}

type CellRange struct {
	UpperLeft   CellPosition `json:"upper_left"`
	BottomRight CellPosition `json:"bottom_right"`
}

func NewCellRange(startColumnName string, startRowNumber int, endColumnName string, endRowNumber int) CellRange {
	return CellRange{
		UpperLeft: CellPosition{
			ColumnName: startColumnName,
			RowNumber:  startRowNumber,
		},
		BottomRight: CellPosition{
			ColumnName: endColumnName,
			RowNumber:  endRowNumber,
		},
	}
}

func (c *CellRange) StartRowNumber() int {
	return c.UpperLeft.RowNumber
}

func (c *CellRange) EndRowNumber() int {
	return c.BottomRight.RowNumber
}

func (c *CellRange) StartColumnName() string {
	return c.UpperLeft.ColumnName
}

func (c *CellRange) EndColumnName() string {
	return c.BottomRight.ColumnName
}

func (c *CellRange) StartColumnNumber() int {
	n, err := excelize.ColumnNameToNumber(c.StartColumnName())
	if err != nil {
		panic("CellRange.StartColumnNumber: " + err.Error())
	}
	return n
}

func (c *CellRange) EndColumnNumber() int {
	n, err := excelize.ColumnNameToNumber(c.EndColumnName())
	if err != nil {
		panic("CellRange.StartColumnNumber: " + err.Error())
	}
	return n
}

type MergedRange struct {
	Value string    `json:"value"`
	Range CellRange `json:"range"`
}
