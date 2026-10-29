package excel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonutil"
	"github.com/xuri/excelize/v2"
)

// Vertical range fixtures under excel/testdata/raw_vertical/<case>/.
// Cases are grouped by directory prefix; each group has its own combined workbook at
// excel/testdata/raw_vertical/<combined_workbook> (see rawVerticalGroups).
//
// Regenerate per-case workbooks and combined workbooks:
//
//	go test ./excel -run TestCreateRawVerticalFixturesInTestdata
//
// Regenerate raw_verticals.json:
//
//	go test ./excel -update -run TestReadRawVerticals

type rawVerticalCase struct {
	dir       string
	readRange CellRange
	write     func(f *excelize.File, sheet string) error
	// merge crosses read range boundary; readRawVerticals must error.
	wantMergeCrossesRangeErr bool
}

type rawVerticalGroup struct {
	combinedWorkbook string
	cases            []rawVerticalCase
}

func rawVerticalGroups() []rawVerticalGroup {
	return []rawVerticalGroup{
		{
			combinedWorkbook: "combined_extent.xlsx",
			cases: []rawVerticalCase{
				{"ext-1-row-simple", NewCellRange("A", 1, "C", 1), writeRawVerticalExt1RowSimple, false},
				{"ext-3-row-simple", NewCellRange("A", 1, "C", 3), writeRawVerticalExt3RowSimple, false},
				{"ext-1x1-single-cell", NewCellRange("B", 2, "B", 2), writeRawVerticalExt1x1, false},
			},
		},
		{
			combinedWorkbook: "combined_sparse.xlsx",
			cases: []rawVerticalCase{
				{"sparse-3-row-empty-cell", NewCellRange("A", 1, "C", 3), writeRawVerticalSparseEmptyCell, false},
				{"sparse-3-row-empty-column", NewCellRange("A", 1, "C", 3), writeRawVerticalSparseEmptyColumn, false},
				{"sparse-1-row-gap-column", NewCellRange("A", 1, "D", 1), writeRawVerticalSparseGapColumn, false},
			},
		},
		{
			combinedWorkbook: "combined_merge.xlsx",
			cases: []rawVerticalCase{
				{"merge-1-row-horizontal", NewCellRange("A", 1, "C", 1), writeRawVerticalMerge1RowHorizontal, false},
				{"merge-2-row-vertical", NewCellRange("A", 1, "B", 2), writeRawVerticalMerge2RowVertical, false},
				{"merge-2-row-block", NewCellRange("A", 1, "C", 2), writeRawVerticalMerge2RowBlock, false},
				{"merge-2-row-bidirectional", NewCellRange("A", 1, "C", 2), writeRawVerticalMerge2RowBidirectional, false},
				{"merge-1-row-two-horizontal", NewCellRange("A", 1, "E", 1), writeRawVerticalMerge1RowTwoHorizontal, false},
			},
		},
		{
			combinedWorkbook: "combined_merge_boundary.xlsx",
			cases: []rawVerticalCase{
				{"bound-merge-above", NewCellRange("A", 2, "A", 3), writeRawVerticalBoundMergeAbove, true},
				{"bound-merge-below", NewCellRange("A", 1, "A", 2), writeRawVerticalBoundMergeBelow, true},
				{"bound-merge-left", NewCellRange("B", 1, "C", 1), writeRawVerticalBoundMergeLeft, true},
				{"bound-merge-right", NewCellRange("A", 1, "B", 1), writeRawVerticalBoundMergeRight, true},
			},
		},
		{
			combinedWorkbook: "combined_hidden.xlsx",
			cases: []rawVerticalCase{
				{"hidden-3-row-hidden-row", NewCellRange("A", 1, "C", 3), writeRawVerticalHiddenRow, false},
				{"hidden-3-row-hidden-column", NewCellRange("A", 1, "C", 3), writeRawVerticalHiddenColumn, false},
			},
		},
		{
			combinedWorkbook: "combined_column.xlsx",
			cases: []rawVerticalCase{
				{"col-1-row-wide", NewCellRange("A", 1, "G", 1), writeRawVerticalCol1RowWide, false},
				{"col-1-row-column-aa", NewCellRange("Z", 1, "AA", 1), writeRawVerticalCol1RowColumnAA, false},
			},
		},
		{
			combinedWorkbook: "combined_range.xlsx",
			cases: []rawVerticalCase{
				{"range-3-row-all-empty", NewCellRange("A", 1, "C", 3), writeRawVerticalRange3RowAllEmpty, false},
				{"range-3-row-below-used", NewCellRange("A", 10, "C", 12), writeRawVerticalRange3RowBelowUsed, false},
			},
		},
		{
			combinedWorkbook: "combined_stress.xlsx",
			cases: []rawVerticalCase{
				{"stress-2-row-many-merge", NewCellRange("A", 1, "H", 2), writeRawVerticalStress2RowManyMerge, false},
			},
		},
	}
}

func rawVerticalCases() []rawVerticalCase {
	var all []rawVerticalCase
	for _, g := range rawVerticalGroups() {
		all = append(all, g.cases...)
	}
	return all
}

func TestCreateRawVerticalFixturesInTestdata(t *testing.T) {
	skipUnlessUpdateEnabled(t)
	base := filepath.Join("testdata", "raw_vertical")
	for _, g := range rawVerticalGroups() {
		sheets := make([]fixtureSheetCase, len(g.cases))
		for i, tc := range g.cases {
			t.Run(tc.dir, func(t *testing.T) {
				dir := filepath.Join(base, tc.dir)
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
				path := filepath.Join(dir, "workbook.xlsx")
				if err := writeFixtureWorkbook(path, tc.dir, tc.write, nil); err != nil {
					t.Fatal(err)
				}
				if err := writeRawVerticalReadRangeFile(dir, tc.dir, tc.readRange); err != nil {
					t.Fatal(err)
				}
			})
			sheets[i] = fixtureSheetCase{dir: tc.dir, write: tc.write}
		}

		combinedPath := filepath.Join(base, g.combinedWorkbook)
		if err := writeFixtureCombinedWorkbook(combinedPath, sheets); err != nil {
			t.Fatalf("%s: %v", g.combinedWorkbook, err)
		}
	}
}

// Exercises readRawVerticals against each case (read_range.json bounds, raw_verticals.json golden).
func TestReadRawVerticals(t *testing.T) {
	for _, tc := range rawVerticalCases() {
		t.Run(tc.dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "raw_vertical", tc.dir)
			excelPath := filepath.Join(caseDir, "workbook.xlsx")
			specPath := filepath.Join(caseDir, "read_range.json")
			goldenPath := filepath.Join(caseDir, "raw_verticals.json")

			spec, err := readRawVerticalReadRangeFile(specPath)
			if err != nil {
				t.Fatalf("read read_range: %v", err)
			}

			got, err := readRawVerticals(excelPath, spec.SheetName, spec.Range)
			if tc.wantMergeCrossesRangeErr {
				if err == nil {
					t.Fatal("readRawVerticals: expected error when merge crosses read range")
				}
				return
			}
			if err != nil {
				t.Fatalf("readRawVerticals: %v", err)
			}

			if updateEnabled() {
				writeRawVerticalsGolden(t, goldenPath, got)
				return
			}

			want, err := readRawVerticalsGolden(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			fatalIfJSONGoldenMismatch(t, got, want)
		})
	}
}

type rawVerticalReadRangeFile struct {
	SheetName string    `json:"sheet_name"`
	Range     CellRange `json:"range"`
}

func writeRawVerticalReadRangeFile(caseDir, sheetName string, r CellRange) error {
	spec := rawVerticalReadRangeFile{SheetName: sheetName, Range: r}
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(caseDir, "read_range.json"), data, 0o644)
}

func readRawVerticalReadRangeFile(path string) (rawVerticalReadRangeFile, error) {
	var spec rawVerticalReadRangeFile
	if err := jsonutil.Read(path, &spec); err != nil {
		return rawVerticalReadRangeFile{}, err
	}
	return spec, nil
}

func readRawVerticalsGolden(path string) ([]RawVertical, error) {
	var verticals []RawVertical
	if err := jsonutil.Read(path, &verticals); err != nil {
		return nil, err
	}
	return verticals, nil
}

func writeRawVerticalsGolden(t *testing.T, path string, verticals []RawVertical) {
	data, err := json.MarshalIndent(verticals, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func verticalCellLabel(col string, row int) string {
	return fmt.Sprintf("r%dc%s", row, col)
}

func fillVerticalGrid(f *excelize.File, sheet string, r CellRange, skip func(col string, row int) bool) error {
	startRow := r.StartRowNumber()
	endRow := r.EndRowNumber()
	startCol, err := r.StartColumnNumber()
	if err != nil {
		return err
	}
	endCol, err := r.EndColumnNumber()
	if err != nil {
		return err
	}
	for colNum := startCol; colNum <= endCol; colNum++ {
		colName, err := excelize.ColumnNumberToName(colNum)
		if err != nil {
			return err
		}
		for row := startRow; row <= endRow; row++ {
			if skip != nil && skip(colName, row) {
				continue
			}
			if err := setRawHeaderCell(f, sheet, colName, row, verticalCellLabel(colName, row)); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeRawVerticalExt1RowSimple(f *excelize.File, sheet string) error {
	return fillVerticalGrid(f, sheet, NewCellRange("A", 1, "C", 1), nil)
}

func writeRawVerticalExt3RowSimple(f *excelize.File, sheet string) error {
	return fillVerticalGrid(f, sheet, NewCellRange("A", 1, "C", 3), nil)
}

func writeRawVerticalExt1x1(f *excelize.File, sheet string) error {
	return setRawHeaderCell(f, sheet, "B", 2, "only")
}

func writeRawVerticalSparseEmptyCell(f *excelize.File, sheet string) error {
	if err := fillVerticalGrid(f, sheet, NewCellRange("A", 1, "C", 3), nil); err != nil {
		return err
	}
	return setRawHeaderCell(f, sheet, "B", 2, "")
}

func writeRawVerticalSparseEmptyColumn(f *excelize.File, sheet string) error {
	return fillVerticalGrid(f, sheet, NewCellRange("A", 1, "C", 3), func(col string, row int) bool {
		return col == "B"
	})
}

func writeRawVerticalSparseGapColumn(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, verticalCellLabel("A", 1)); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, verticalCellLabel("B", 1)); err != nil {
		return err
	}
	return setRawHeaderCell(f, sheet, "D", 1, verticalCellLabel("D", 1))
}

func writeRawVerticalMerge1RowHorizontal(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, verticalCellLabel("A", 1)); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, "mergedBC"); err != nil {
		return err
	}
	if err := mergeRawHeader(f, sheet, "B", 1, "C", 1); err != nil {
		return err
	}
	return nil
}

func writeRawVerticalMerge2RowVertical(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "mergedA"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, verticalCellLabel("B", 1)); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 2, verticalCellLabel("B", 2)); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "A", 1, "A", 2)
}

func writeRawVerticalMerge2RowBlock(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, verticalCellLabel("A", 1)); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "A", 2, verticalCellLabel("A", 2)); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "B", 1, "block"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "B", 1, "C", 2)
}

func writeRawVerticalMerge2RowBidirectional(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "bidir"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 1, verticalCellLabel("C", 1)); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 2, verticalCellLabel("C", 2)); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "A", 1, "B", 2)
}

func writeRawVerticalMerge1RowTwoHorizontal(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "left"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "C", 1, verticalCellLabel("C", 1)); err != nil {
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

func writeRawVerticalBoundMergeAbove(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "mergedA"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "A", 1, "A", 3)
}

func writeRawVerticalBoundMergeBelow(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "mergedA"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "A", 1, "A", 3)
}

func writeRawVerticalBoundMergeLeft(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "mergedABC"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "A", 1, "C", 1)
}

func writeRawVerticalBoundMergeRight(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "B", 1, "mergedBCD"); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "B", 1, "D", 1)
}

func writeRawVerticalHiddenRow(f *excelize.File, sheet string) error {
	if err := fillVerticalGrid(f, sheet, NewCellRange("A", 1, "C", 3), nil); err != nil {
		return err
	}
	return f.SetRowVisible(sheet, 2, false)
}

func writeRawVerticalHiddenColumn(f *excelize.File, sheet string) error {
	if err := fillVerticalGrid(f, sheet, NewCellRange("A", 1, "C", 3), nil); err != nil {
		return err
	}
	return f.SetColVisible(sheet, "B", false)
}

func writeRawVerticalCol1RowWide(f *excelize.File, sheet string) error {
	return fillVerticalGrid(f, sheet, NewCellRange("A", 1, "G", 1), nil)
}

func writeRawVerticalCol1RowColumnAA(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "Z", 1, verticalCellLabel("Z", 1)); err != nil {
		return err
	}
	return setRawHeaderCell(f, sheet, "AA", 1, verticalCellLabel("AA", 1))
}

func writeRawVerticalRange3RowAllEmpty(f *excelize.File, sheet string) error {
	return nil
}

func writeRawVerticalRange3RowBelowUsed(f *excelize.File, sheet string) error {
	return setRawHeaderCell(f, sheet, "A", 1, "used")
}

func writeRawVerticalStress2RowManyMerge(f *excelize.File, sheet string) error {
	if err := setRawHeaderCell(f, sheet, "A", 1, "g1"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "D", 1, "g2"); err != nil {
		return err
	}
	if err := setRawHeaderCell(f, sheet, "G", 1, "g3"); err != nil {
		return err
	}
	headers := []string{"hA", "hB", "hC", "hD", "hE", "hF", "hG", "hH"}
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
	if err := mergeRawHeader(f, sheet, "D", 1, "F", 1); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "G", 1, "H", 1)
}
