package excel

import "testing"

func TestIsCellInMergeRange(t *testing.T) {
	merged := MergedRange{
		Value: "title",
		Range: NewCellRange("B", 2, "D", 4),
	}

	tests := []struct {
		name string
		pos  CellPosition
		want bool
	}{
		{"upper_left_corner", CellPosition{ColumnName: "B", RowNumber: 2}, true},
		{"lower_right_corner", CellPosition{ColumnName: "D", RowNumber: 4}, true},
		{"interior", CellPosition{ColumnName: "C", RowNumber: 3}, true},
		{"left_of_range", CellPosition{ColumnName: "A", RowNumber: 3}, false},
		{"right_of_range", CellPosition{ColumnName: "E", RowNumber: 3}, false},
		{"above_range", CellPosition{ColumnName: "C", RowNumber: 1}, false},
		{"below_range", CellPosition{ColumnName: "C", RowNumber: 5}, false},
		{"invalid_column", CellPosition{ColumnName: "", RowNumber: 3}, false},
		{"invalid_row", CellPosition{ColumnName: "C", RowNumber: 0}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isCellInMergeRange(tc.pos, merged); got != tc.want {
				t.Errorf("isCellInMergeRange(%v): got %v want %v", tc.pos, got, tc.want)
			}
		})
	}
}

func TestIsCellInMergeRange_invalidMergeRange(t *testing.T) {
	pos := CellPosition{ColumnName: "A", RowNumber: 1}
	bad := MergedRange{
		Range: NewCellRange("???", 1, "B", 1),
	}
	if isCellInMergeRange(pos, bad) {
		t.Fatal("expected false when merge range has invalid column names")
	}
}
