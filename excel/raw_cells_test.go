package excel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Cell fixtures under excel/testdata/raw_cells/<case>/workbook.xlsx, merged_ranges.json, and <cell>.json (e.g. B1.json).
// Regenerate workbooks: go test ./excel -run TestCreateRawCellsFixturesInTestdata

type rawCellsCase struct {
	dir   string
	write func(f *excelize.File, sheet string) error
	cells []string
}

func rawCellsCases() []rawCellsCase {
	return []rawCellsCase{
		{"none", writeRawCellsNone, []string{"A1"}},
		{"horizontal", writeRawCellsHorizontal, []string{"A1", "B1", "C1"}},
		{"vertical", writeRawCellsVertical, []string{"A1", "A2", "B1"}},
		{"block", writeRawCellsBlock, []string{"A1", "B1", "C1", "B2", "C2"}},
		{"two-horizontal", writeRawCellsTwoHorizontal, []string{"A1", "B1", "C1", "D1", "E1"}},
	}
}

func TestCreateRawCellsFixturesInTestdata(t *testing.T) {
	base := filepath.Join("testdata", "raw_cells")
	cases := rawCellsCases()
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

	sheets := make([]fixtureSheetCase, len(cases))
	for i, tc := range cases {
		sheets[i] = fixtureSheetCase{dir: tc.dir, write: tc.write}
	}
	combinedPath := filepath.Join(base, "test_cases.xlsx")
	if err := writeFixtureCombinedWorkbook(combinedPath, sheets); err != nil {
		t.Fatal(err)
	}
}

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

func TestLoadMergedRanges(t *testing.T) {
	for _, tc := range rawCellsCases() {
		t.Run(tc.dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "raw_cells", tc.dir)
			excelPath := filepath.Join(caseDir, "workbook.xlsx")
			goldenPath := filepath.Join(caseDir, "merged_ranges.json")

			f, err := excelize.OpenFile(excelPath)
			if err != nil {
				t.Fatalf("open workbook: %v", err)
			}
			defer f.Close()

			got, err := loadMergedRanges(f, tc.dir)
			if err != nil {
				t.Fatalf("loadMergedRanges: %v", err)
			}
			sortMergedRanges(got)

			if os.Getenv("UPDATE_RAW_CELLS_GOLDEN") != "" {
				writeMergedRangesGolden(t, goldenPath, got)
				return
			}

			want, err := readMergedRangesGolden(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			sortMergedRanges(want)
			if !reflect.DeepEqual(got, want) {
				gotJSON, _ := json.MarshalIndent(got, "", "  ")
				wantJSON, _ := json.MarshalIndent(want, "", "  ")
				t.Fatalf("merged ranges mismatch:\n--- got\n%s\n--- want\n%s", gotJSON, wantJSON)
			}
		})
	}
}

func TestReadCellValue(t *testing.T) {
	for _, tc := range rawCellsCases() {
		for _, cellName := range tc.cells {
			t.Run(tc.dir+"/"+cellName, func(t *testing.T) {
				caseDir := filepath.Join("testdata", "raw_cells", tc.dir)
				excelPath := filepath.Join(caseDir, "workbook.xlsx")
				goldenPath := filepath.Join(caseDir, cellName+".json")

				col, row, err := excelize.SplitCellName(cellName)
				if err != nil {
					t.Fatalf("split cell name: %v", err)
				}

				f, err := excelize.OpenFile(excelPath)
				if err != nil {
					t.Fatalf("open workbook: %v", err)
				}
				defer f.Close()

				merges, err := loadMergedRanges(f, tc.dir)
				if err != nil {
					t.Fatalf("loadMergedRanges: %v", err)
				}

				got, err := readCellValue(f, tc.dir, col, row, merges)
				if err != nil {
					t.Fatalf("readCellValue: %v", err)
				}

				if os.Getenv("UPDATE_RAW_CELLS_GOLDEN") != "" {
					writeCellValueGolden(t, goldenPath, got)
					return
				}

				want, err := readCellValueGolden(goldenPath)
				if err != nil {
					t.Fatalf("read golden: %v", err)
				}
				if !reflect.DeepEqual(got, want) {
					gotJSON, _ := json.MarshalIndent(got, "", "  ")
					wantJSON, _ := json.MarshalIndent(want, "", "  ")
					t.Fatalf("cell value mismatch:\n--- got\n%s\n--- want\n%s", gotJSON, wantJSON)
				}
			})
		}
	}
}

func writeRawCellsNone(f *excelize.File, sheet string) error {
	return setRawHeaderCell(f, sheet, "A", 1, "only")
}

func writeRawCellsHorizontal(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "B", 1, "headerBC"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "B", 1, "C", 1)
}

func writeRawCellsVertical(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "headerA"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "A", 1, "A", 2)
}

func writeRawCellsBlock(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "B", 1, "block"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "B", 1, "C", 2)
}

func writeRawCellsTwoHorizontal(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "left"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "D", 1, "right"); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "A", 1, "B", 1); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "D", 1, "E", 1)
}

func sortMergedRanges(merges []MergedRange) {
	sort.Slice(merges, func(i, j int) bool {
		return mergedRangeSortKey(merges[i]) < mergedRangeSortKey(merges[j])
	})
}

func mergedRangeSortKey(m MergedRange) string {
	ul := m.Range.UpperLeft
	br := m.Range.BottomRight
	return fmt.Sprintf("%04d-%s-%04d-%s", ul.RowNumber, ul.ColumnName, br.RowNumber, br.ColumnName)
}

func readMergedRangesGolden(path string) ([]MergedRange, error) {
	var merges []MergedRange
	if err := jsonfile.Read(path, &merges); err != nil {
		return nil, err
	}
	return merges, nil
}

func writeMergedRangesGolden(t *testing.T, path string, merges []MergedRange) {
	data, err := json.MarshalIndent(merges, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func readCellValueGolden(path string) (CellValue, error) {
	var cv CellValue
	if err := jsonfile.Read(path, &cv); err != nil {
		return CellValue{}, err
	}
	return cv, nil
}

func writeCellValueGolden(t *testing.T, path string, cv CellValue) {
	data, err := json.MarshalIndent(cv, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
