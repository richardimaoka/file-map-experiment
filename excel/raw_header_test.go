package excel

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Shared header layout writers and fixture workbook helpers for schema/raw_* tests.
// Checked-in workbooks under testdata/raw_header/ are maintained via schema_header_* fixture generators.

func readSchemaPrepFixture(path string) (*SchemaPrep, error) {
	var prep SchemaPrep
	if err := jsonfile.Read(path, &prep); err != nil {
		return nil, err
	}
	return &prep, nil
}

type fixtureSheetCase struct {
	dir   string
	write func(f *excelize.File, sheet string) error
	prep  *SchemaPrep
}

// Combined workbook for a testdata tree: <dir_name>_tests.xlsx (e.g. raw_header_tests.xlsx).
func fixtureCombinedWorkbookPath(testdataDir string) string {
	return filepath.Join(testdataDir, filepath.Base(testdataDir)+"_tests.xlsx")
}

func writeRawHeaderWorkbook(path, sheetName string, write func(f *excelize.File, sheet string) error, prep *SchemaPrep) error {
	f := excelize.NewFile()
	defer f.Close()

	if err := renameDefaultSheet(f, sheetName); err != nil {
		return err
	}
	if err := write(f, sheetName); err != nil {
		return err
	}
	if prep != nil {
		if err := applyFixtureWorkbookStyle(f, sheetName, *prep); err != nil {
			return err
		}
	}
	if err := f.SaveAs(path); err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return fmt.Errorf("%s was not written", path)
	}
	return nil
}

func writeFixtureCombinedWorkbook(path string, cases []fixtureSheetCase) error {
	f := excelize.NewFile()
	defer f.Close()

	for i, tc := range cases {
		sheet := tc.dir
		if i == 0 {
			if err := renameDefaultSheet(f, sheet); err != nil {
				return err
			}
		} else if _, err := f.NewSheet(sheet); err != nil {
			return err
		}
		if err := tc.write(f, sheet); err != nil {
			return err
		}
		if tc.prep != nil {
			if err := applyFixtureWorkbookStyle(f, sheet, *tc.prep); err != nil {
				return err
			}
		}
	}

	if err := f.SaveAs(path); err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return fmt.Errorf("%s was not written", path)
	}
	return nil
}

func renameDefaultSheet(f *excelize.File, name string) error {
	defaultSheet := f.GetSheetName(0)
	if defaultSheet == name {
		return nil
	}
	return f.SetSheetName(defaultSheet, name)
}

func setRawHeaderCell(f *excelize.File, sheet, col string, row int, value string) error {
	cell, err := excelize.JoinCellName(col, row)
	if err != nil {
		return err
	}
	return f.SetCellValue(sheet, cell, value)
}

func setRawHeaderRow(f *excelize.File, sheet string, row int, cols []string, vals []string) error {
	if len(cols) != len(vals) {
		return fmt.Errorf("setRawHeaderRow: cols and vals length mismatch")
	}
	for i, col := range cols {
		if err := setRawHeaderCell(f, sheet, col, row, vals[i]); err != nil {
			return err
		}
	}
	return nil
}

func mergeRawHeader(f *excelize.File, sheet, upperLeftCol string, upperLeftRow int, bottomRightCol string, bottomRightRow int) error {
	from, err := excelize.JoinCellName(upperLeftCol, upperLeftRow)
	if err != nil {
		return err
	}
	to, err := excelize.JoinCellName(bottomRightCol, bottomRightRow)
	if err != nil {
		return err
	}
	return f.MergeCell(sheet, from, to)
}

func writeRawHeader1RowSimple(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, "headerB"); err != nil {
		return err
	}
	return setRawHeaderCell(f, sheet, "C", 1, "headerC")
}

// Column C is intentionally left blank (empty header).
func writeRawHeader1RowEmptyHeader(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, "headerB"); err != nil {
		return err
	}
	return setRawHeaderCell(f, sheet, "D", 1, "headerC")
}

func writeRawHeader1RowLineBreak(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	return setRawHeaderCell(f, sheet, "B", 1, "this\nis a long\nheader with\nline breaks")
}

func writeRawHeader1RowHorizontalMerge(f *excelize.File, sheet string) error {
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
func writeRawHeader1RowDuplicateHeader(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, "headerB"); err != nil {
		return err
	}
	return setRawHeaderCell(f, sheet, "C", 1, "headerB")
}

func writeRawHeader2RowSimple(f *excelize.File, sheet string) error {
	if err := setRawHeaderRow(f, sheet, 1, []string{"A", "B", "C"}, []string{"group1", "group2", "group3"}); err != nil {
		return err
	}
	return setRawHeaderRow(f, sheet, 2, []string{"A", "B", "C"}, []string{"headerA", "headerB", "headerC"})
}

// Columns B and C both resolve to group2 > headerB (duplicate RowData keys).
func writeRawHeader2RowDuplicateHeader(f *excelize.File, sheet string) error {
	if err := setRawHeaderRow(f, sheet, 1, []string{"A", "B", "C"}, []string{"group1", "group2", "group2"}); err != nil {
		return err
	}
	return setRawHeaderRow(f, sheet, 2, []string{"A", "B", "C"}, []string{"headerA", "headerB", "headerB"})
}

func writeRawHeader2RowEmptyCell(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "group1"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 1, "group3"); err != nil {
		return err
	}
	return setRawHeaderRow(f, sheet, 2, []string{"A", "B", "C"}, []string{"headerA", "headerB", "headerC"})
}

func writeRawHeader2RowHorizontalMerge(f *excelize.File, sheet string) error {
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

func writeRawHeader2RowHorizontalMergeRepeated(f *excelize.File, sheet string) error {
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

func writeRawHeader2RowVerticalMerge(f *excelize.File, sheet string) error {
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

func writeRawHeader2RowVerticalMergeRepeated(f *excelize.File, sheet string) error {
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

func writeRawHeader2RowBidirectionalMerge(f *excelize.File, sheet string) error {
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

func writeRawHeader2RowBidirectionalMergeRepeated(f *excelize.File, sheet string) error {
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

func writeRawHeader2RowMixedMerge(f *excelize.File, sheet string) error {
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

func writeRawHeader3RowSimple(f *excelize.File, sheet string) error {
	if err := setRawHeaderRow(f, sheet, 1, []string{"A", "B", "C"}, []string{"group1", "group2", "group3"}); err != nil {
		return err
	}
	if err := setRawHeaderRow(f, sheet, 2, []string{"A", "B", "C"}, []string{"subgroup1", "subgroup2", "subgroup3"}); err != nil {
		return err
	}
	return setRawHeaderRow(f, sheet, 3, []string{"A", "B", "C"}, []string{"headerA", "headerB", "headerC"})
}

// 2D merge A1:B3 intersects the leaf header row across multiple columns.
func writeRawHeader3RowBidirectionalMerge(f *excelize.File, sheet string) error {
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
func writeRawHeader3RowMixedMerge(f *excelize.File, sheet string) error {
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
