package excel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Header fixtures under excel/testdata/raw_header/<case>/ (workbook.xlsx, raw_headers.json).
// Regenerate workbooks: go test ./excel -run TestCreateRawHeaderFixturesInTestdata
// Regenerate raw_headers.json: go test ./excel -update -run TestReadRawHeaders

type rawHeaderCase struct {
	dir   string
	write func(f *excelize.File, sheet string) error
}

func rawHeaderCases() []rawHeaderCase {
	return []rawHeaderCase{
		{"1-row-simple", writeRawHeader1RowSimple},
		{"1-row-empty-header", writeRawHeader1RowEmptyHeader},
		{"1-row-line-break", writeRawHeader1RowLineBreak},
		{"1-row-horizontal-merge", writeRawHeader1RowHorizontalMerge},
		{"2-row-simple", writeRawHeader2RowSimple},
		{"2-row-empty-cell", writeRawHeader2RowEmptyCell},
		{"2-row-horizontal-merge", writeRawHeader2RowHorizontalMerge},
		{"2-row-horizontal-merge-repeated", writeRawHeader2RowHorizontalMergeRepeated},
		{"2-row-vertical-merge", writeRawHeader2RowVerticalMerge},
		{"2-row-vertical-merge-repeated", writeRawHeader2RowVerticalMergeRepeated},
		{"2-row-bidirectional-merge", writeRawHeader2RowBidirectionalMerge},
		{"2-row-bidirectional-merge-rep", writeRawHeader2RowBidirectionalMergeRepeated},
		{"2-row-mixed-merge", writeRawHeader2RowMixedMerge},
	}
}

func TestCreateRawHeaderFixturesInTestdata(t *testing.T) {
	base := filepath.Join("testdata", "raw_header")
	cases := rawHeaderCases()

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			dir := filepath.Join(base, tc.dir)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(dir, "workbook.xlsx")
			if err := writeRawHeaderWorkbook(path, tc.dir, tc.write); err != nil {
				t.Fatal(err)
			}
		})
	}

	combinedPath := fixtureCombinedWorkbookPath(base)
	if err := writeRawHeaderCombinedWorkbook(combinedPath, cases); err != nil {
		t.Fatal(err)
	}
}

func TestReadRawHeaders(t *testing.T) {
	for _, dir := range rawHeaderCaseDirs() {
		t.Run(dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "raw_header", dir)
			prepPath := filepath.Join(caseDir, "schema_prep.json")
			excelPath := filepath.Join(caseDir, "workbook.xlsx")
			goldenPath := filepath.Join(caseDir, "raw_headers.json")

			var prep SchemaPrep
			if err := jsonfile.Read(prepPath, &prep); err != nil {
				t.Fatalf("read schema_prep: %v", err)
			}

			got, err := readRawHeaders(excelPath, prep)
			if err != nil {
				t.Fatalf("readRawHeaders: %v", err)
			}

			if updateGoldens() {
				writeRawHeadersGolden(t, goldenPath, got)
				return
			}

			want, err := readRawHeadersGolden(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			fatalIfJSONGoldenMismatch(t, got, want)
		})
	}
}

func rawHeaderCaseDirs() []string {
	cases := rawHeaderCases()
	dirs := make([]string, len(cases))
	for i, c := range cases {
		dirs[i] = c.dir
	}
	return dirs
}

func readRawHeadersGolden(path string) ([]RawHeader, error) {
	var headers []RawHeader
	if err := jsonfile.Read(path, &headers); err != nil {
		return nil, err
	}
	return headers, nil
}

func writeRawHeadersGolden(t *testing.T, path string, headers []RawHeader) {
	data, err := json.MarshalIndent(headers, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

type fixtureSheetCase struct {
	dir   string
	write func(f *excelize.File, sheet string) error
}

// Combined workbook for a testdata tree: <dir_name>_tests.xlsx (e.g. raw_header_tests.xlsx).
func fixtureCombinedWorkbookPath(testdataDir string) string {
	return filepath.Join(testdataDir, filepath.Base(testdataDir)+"_tests.xlsx")
}

func writeRawHeaderWorkbook(path, sheetName string, write func(f *excelize.File, sheet string) error) error {
	f := excelize.NewFile()
	defer f.Close()

	if err := renameDefaultSheet(f, sheetName); err != nil {
		return err
	}
	if err := write(f, sheetName); err != nil {
		return err
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

func writeRawHeaderCombinedWorkbook(path string, cases []rawHeaderCase) error {
	sheets := make([]fixtureSheetCase, len(cases))
	for i, tc := range cases {
		sheets[i] = fixtureSheetCase{dir: tc.dir, write: tc.write}
	}
	return writeFixtureCombinedWorkbook(path, sheets)
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

func writeRawHeader2RowSimple(f *excelize.File, sheet string) error {
	if err := setRawHeaderRow(f, sheet, 1, []string{"A", "B", "C"}, []string{"group1", "group2", "group3"}); err != nil {
		return err
	}
	return setRawHeaderRow(f, sheet, 2, []string{"A", "B", "C"}, []string{"headerA", "headerB", "headerC"})
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
