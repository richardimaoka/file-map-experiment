package excel

import "fmt"

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

type MergedRange struct {
	Value string    `json:"value"`
	Range CellRange `json:"range"`
}
