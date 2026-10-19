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

func (c *CellRange) StartColumnNumber() (int, error) {
	n, err := excelize.ColumnNameToNumber(c.StartColumnName())
	if err != nil {
		return 0, fmt.Errorf("CellRange.StartColumnNumber: %w", err)
	}
	return n, nil
}

func (c *CellRange) EndColumnNumber() (int, error) {
	n, err := excelize.ColumnNameToNumber(c.EndColumnName())
	if err != nil {
		return 0, fmt.Errorf("CellRange.EndColumnNumber: %w", err)
	}
	return n, nil
}

// Intersects reports whether any cell lies in both ranges (inclusive edges).
func (c CellRange) Intersects(other CellRange) (bool, error) {
	a0, err := c.StartColumnNumber()
	if err != nil {
		return false, err
	}
	a1, err := c.EndColumnNumber()
	if err != nil {
		return false, err
	}
	b0, err := other.StartColumnNumber()
	if err != nil {
		return false, err
	}
	b1, err := other.EndColumnNumber()
	if err != nil {
		return false, err
	}

	if a1 < b0 || b1 < a0 {
		return false, nil
	}
	if c.EndRowNumber() < other.StartRowNumber() || other.EndRowNumber() < c.StartRowNumber() {
		return false, nil
	}
	return true, nil
}

// Contains reports whether every cell in other lies in c (inclusive edges).
func (c CellRange) Contains(other CellRange) (bool, error) {
	ok, err := c.Intersects(other)
	if err != nil || !ok {
		return false, err
	}
	a0, err := c.StartColumnNumber()
	if err != nil {
		return false, err
	}
	a1, err := c.EndColumnNumber()
	if err != nil {
		return false, err
	}
	b0, err := other.StartColumnNumber()
	if err != nil {
		return false, err
	}
	b1, err := other.EndColumnNumber()
	if err != nil {
		return false, err
	}

	return b0 >= a0 && b1 <= a1 &&
		other.StartRowNumber() >= c.StartRowNumber() &&
		other.EndRowNumber() <= c.EndRowNumber(), nil
}

type MergedRange struct {
	Value string    `json:"value"`
	Range CellRange `json:"range"`
}
