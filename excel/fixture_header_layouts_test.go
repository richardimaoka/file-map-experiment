package excel

import (
	"github.com/xuri/excelize/v2"
)

// Header layout writers for schema_header_* and related fixture tests.
// Workbook save/load helpers live in fixture_workbook_test.go; cell helpers in fixture_cells_test.go.

func writeFixture1RowSimple(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, "headerB"); err != nil {
		return err
	}
	return setRawHeaderCell(f, sheet, "C", 1, "headerC")
}

// Column C is intentionally left blank (empty header).
func writeFixture1RowEmptyHeader(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, "headerB"); err != nil {
		return err
	}
	return setRawHeaderCell(f, sheet, "D", 1, "headerC")
}

func writeFixture1RowLineBreak(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	return setRawHeaderCell(f, sheet, "B", 1, "this\nis a long\nheader with\nline breaks")
}

func writeFixture1RowHorizontalMerge(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, "headerBC"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "D", 1, "headerD"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "B", 1, "C", 1)
}

// Columns B and C share the same header text (duplicate RowData keys).
func writeFixture1RowDuplicateHeader(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, "headerB"); err != nil {
		return err
	}
	return setRawHeaderCell(f, sheet, "C", 1, "headerB")
}

func writeFixture2RowSimple(f *excelize.File, sheet string) error {
	if err := setRawHeaderRow(f, sheet, 1, []string{"A", "B", "C"}, []string{"group1", "group2", "group3"}); err != nil {
		return err
	}
	return setRawHeaderRow(f, sheet, 2, []string{"A", "B", "C"}, []string{"headerA", "headerB", "headerC"})
}

// Columns B and C both resolve to group2 > headerB (duplicate RowData keys).
func writeFixture2RowDuplicateHeader(f *excelize.File, sheet string) error {
	if err := setRawHeaderRow(f, sheet, 1, []string{"A", "B", "C"}, []string{"group1", "group2", "group2"}); err != nil {
		return err
	}
	return setRawHeaderRow(f, sheet, 2, []string{"A", "B", "C"}, []string{"headerA", "headerB", "headerB"})
}

func writeFixture2RowEmptyCell(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "group1"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 1, "group3"); err != nil {
		return err
	}
	return setRawHeaderRow(f, sheet, 2, []string{"A", "B", "C"}, []string{"headerA", "headerB", "headerC"})
}

func writeFixture2RowHorizontalMerge(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "group1"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 1, "group3"); err != nil {
		return err
	}
	if err := setRawHeaderRow(f, sheet, 2, []string{"A", "B", "C"}, []string{"headerA", "headerB", "headerC"}); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "A", 1, "B", 1)
}

func writeFixture2RowHorizontalMergeRepeated(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "group1"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 1, "group2"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "D", 1, "group 3"); err != nil {
		return err
	}
	headers := []string{"headerA", "headerB", "headerC", "headerD", "headerE", "headerF", "headerG"}
	for i, h := range headers {
		col, err := excelize.ColumnNumberToName(i + 1)
		if err != nil {
			return err
		}
		if err := setRawHeaderCell(f, sheet, col, 2, h); err != nil {
			return err
		}
	}
	if err := mergeRawHeader(f, sheet, "A", 1, "B", 1); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "D", 1, "F", 1)
}

func writeFixture2RowVerticalMerge(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, "group2"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 1, "group3"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 2, "headerB"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 2, "headerC"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "A", 1, "A", 2)
}

func writeFixture2RowVerticalMergeRepeated(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, "group2"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 1, "headerC"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "D", 1, "headerD"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 2, "headerB"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "E", 2, "headerE"); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "A", 1, "A", 2); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "C", 1, "C", 2); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "D", 1, "D", 2)
}

func writeFixture2RowBidirectionalMerge(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 1, "group3"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 2, "headerC"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "A", 1, "B", 2)
}

func writeFixture2RowBidirectionalMergeRepeated(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 1, "headerB"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "F", 1, "headerD"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "E", 2, "headerC"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "I", 2, "headerE"); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "A", 1, "B", 2); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "C", 1, "D", 2); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "F", 1, "H", 2)
}

func writeFixture2RowMixedMerge(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "B", 1, "headerB"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "E", 1, "groupA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "G", 1, "groupB"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "I", 1, "groupC"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "J", 1, "headerJ"); err != nil {
		return err
	}
	if err := setRawHeaderRow(f, sheet, 2,
		[]string{"A", "C", "D", "E", "F", "G", "H", "I", "L"},
		[]string{"headerA", "headerC", "headerD", "headerE", "headerF", "headerG", "headerH", "headerI", "headerK"},
	); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "J", 1, "K", 2); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "B", 1, "B", 2); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "E", 1, "F", 1); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "G", 1, "H", 1)
}

func writeFixture3RowSimple(f *excelize.File, sheet string) error {
	if err := setRawHeaderRow(f, sheet, 1, []string{"A", "B", "C"}, []string{"group1", "group2", "group3"}); err != nil {
		return err
	}
	if err := setRawHeaderRow(f, sheet, 2, []string{"A", "B", "C"}, []string{"subgroup1", "subgroup2", "subgroup3"}); err != nil {
		return err
	}
	return setRawHeaderRow(f, sheet, 3, []string{"A", "B", "C"}, []string{"headerA", "headerB", "headerC"})
}

// 2D merge A1:B3 intersects the leaf header row across multiple columns.
func writeFixture3RowBidirectionalMerge(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "tierA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 1, "group3"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 2, "subgroup3"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 3, "headerC"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "A", 1, "B", 3)
}

// Horizontal group/subgroup merges (A–C, D–E) plus vertical leaf merges (F–H); data band only in A–C.
func writeFixture3RowMixedMerge(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "group1"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "D", 1, "group2"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "A", 2, "subgroup1"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 2, "subgroup2"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "D", 2, "subgroup1"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "E", 2, "subgroup2"); err != nil {
		return err
	}
	if err := setRawHeaderRow(f, sheet, 2, []string{"F", "G", "H"}, []string{"headerA", "headerB", "headerC"}); err != nil {
		return err
	}
	if err := setRawHeaderRow(f, sheet, 3, []string{"A", "B", "C", "D", "E"}, []string{"headerA", "headerB", "headerC", "headerA", "headerB"}); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "A", 1, "C", 1); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "B", 2, "C", 2); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "F", 2, "F", 3); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "G", 2, "G", 3); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "H", 2, "H", 3)
}
