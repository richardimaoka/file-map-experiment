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

type CellValue struct {
	Value       string       `json:"value"`
	Formula     string       `json:"formula"`
	Position    CellPosition `json:"position"`
	IsMerged    bool         `json:"is_merged"`
	MergedRange *MergedRange `json:"merged_range,omitempty"` //nil if not merged
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
