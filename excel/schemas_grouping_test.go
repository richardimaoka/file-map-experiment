package excel

import "testing"

func TestDataColumnIsRowGrouped(t *testing.T) {
	verticalMerge := func(col string, row int, ulRow, brRow int) CellValue {
		return CellValue{
			Value:    "x",
			Position: CellPosition{ColumnName: col, RowNumber: row},
			MergedRange: &MergedRange{
				Range: NewCellRange(col, ulRow, col, brRow),
			},
		}
	}

	t.Run("vertical merge in data band", func(t *testing.T) {
		v := RawVertical{
			ExcelColumn: "A",
			CellValues: []CellValue{
				verticalMerge("A", 2, 2, 3),
				verticalMerge("A", 3, 2, 3),
				{Value: "r4cA", Position: CellPosition{ColumnName: "A", RowNumber: 4}},
			},
		}
		if !dataColumnIsRowGrouped(v) {
			t.Fatal("expected row grouped for column A")
		}
	})

	t.Run("no merge", func(t *testing.T) {
		v := RawVertical{
			ExcelColumn: "B",
			CellValues: []CellValue{
				{Value: "r2cB", Position: CellPosition{ColumnName: "B", RowNumber: 2}},
			},
		}
		if dataColumnIsRowGrouped(v) {
			t.Fatal("expected not row grouped")
		}
	})

	t.Run("horizontal merge only", func(t *testing.T) {
		v := RawVertical{
			ExcelColumn: "A",
			CellValues: []CellValue{
				{
					Value:    "x",
					Position: CellPosition{ColumnName: "A", RowNumber: 2},
					MergedRange: &MergedRange{
						Range: NewCellRange("A", 2, "B", 2),
					},
				},
			},
		}
		if dataColumnIsRowGrouped(v) {
			t.Fatal("horizontal data merge should not set row grouped yet")
		}
	})
}
