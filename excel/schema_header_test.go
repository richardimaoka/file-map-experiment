package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

const schemaHeaderFixturesDir = "schema_header"

// Schema header fixtures under excel/testdata/schema_header/{1row,2row,3row}/<case>/.
// Header layouts: fixture_header_layouts_test.go; data band uses r{row}c{col} literals.
//
// Regenerate workbook fixtures and schema_prep.json:
//
//	go test ./excel -update -run TestCreateSchemaHeaderFixturesInTestdata
//
// Regenerate table_schema.json goldens:
//
//	go test ./excel -update -run TestReadSchemaInternalHeader

type schemaHeaderTrackCase struct {
	tier             string // "1row", "2row", "3row"
	dir              string
	headerEndRow     int
	headerWrite      func(f *excelize.File, sheet string) error
	dataBand         CellRange
	wantSchemaErr    bool
	skipDataColumns map[string]bool
}

func (tc schemaHeaderTrackCase) caseRelPath() string {
	return filepath.Join(tc.tier, tc.dir)
}

func (tc schemaHeaderTrackCase) writeWorkbook(f *excelize.File, sheet string) error {
	if err := tc.headerWrite(f, sheet); err != nil {
		return err
	}
	return fillSchemaHeaderDataBand(f, sheet, tc.dataBand, tc.skipDataColumns)
}

func schemaHeaderCases() []schemaHeaderTrackCase {
	return append(append(schemaHeader1RowCases(), schemaHeader2RowCases()...), schemaHeader3RowCases()...)
}

func schemaHeader1RowCases() []schemaHeaderTrackCase {
	return []schemaHeaderTrackCase{
		{"1row", "1-row-simple", 1, writeFixture1RowSimple, NewCellRange("A", 2, "C", 4), false, nil},
		{"1row", "1-row-empty-header", 1, writeFixture1RowEmptyHeader, NewCellRange("A", 2, "D", 4), true, nil},
		{
			"1row", "1-row-empty-header-valid", 1,
			writeFixture1RowEmptyHeader,
			NewCellRange("A", 2, "D", 4),
			false,
			map[string]bool{"C": true},
		},
		{"1row", "1-row-line-break", 1, writeFixture1RowLineBreak, NewCellRange("A", 2, "B", 4), false, nil},
		{"1row", "1-row-horizontal-merge", 1, writeFixture1RowHorizontalMerge, NewCellRange("A", 2, "D", 4), true, nil},
		{"1row", "1-row-duplicate-header", 1, writeFixture1RowDuplicateHeader, NewCellRange("A", 2, "C", 4), true, nil},
	}
}

func schemaHeader2RowCases() []schemaHeaderTrackCase {
	return []schemaHeaderTrackCase{
		{"2row", "2-row-simple", 2, writeFixture2RowSimple, NewCellRange("A", 3, "C", 5), false, nil},
		{"2row", "2-row-duplicate-header", 2, writeFixture2RowDuplicateHeader, NewCellRange("A", 3, "C", 5), true, nil},
		{"2row", "2-row-empty-cell", 2, writeFixture2RowEmptyCell, NewCellRange("A", 3, "C", 5), false, nil},
		{"2row", "2-row-horizontal-merge", 2, writeFixture2RowHorizontalMerge, NewCellRange("A", 3, "C", 5), false, nil},
		{"2row", "2-row-vertical-merge", 2, writeFixture2RowVerticalMerge, NewCellRange("A", 3, "C", 5), false, nil},
		{"2row", "2-row-bidirectional-merge", 2, writeFixture2RowBidirectionalMerge, NewCellRange("A", 3, "C", 5), true, nil},
		{"2row", "2-row-horizontal-merge-repeated", 2, writeFixture2RowHorizontalMergeRepeated, NewCellRange("A", 3, "G", 5), false, nil},
		{"2row", "2-row-vertical-merge-repeated", 2, writeFixture2RowVerticalMergeRepeated, NewCellRange("A", 3, "E", 5), false, nil},
		{"2row", "2-row-bidirectional-merge-rep", 2, writeFixture2RowBidirectionalMergeRepeated, NewCellRange("A", 3, "I", 5), true, nil},
		{"2row", "2-row-mixed-merge", 2, writeFixture2RowMixedMerge, NewCellRange("A", 3, "L", 5), true, nil},
	}
}

func schemaHeader3RowCases() []schemaHeaderTrackCase {
	return []schemaHeaderTrackCase{
		{"3row", "3-row-simple", 3, writeFixture3RowSimple, NewCellRange("A", 4, "C", 6), false, nil},
		{"3row", "3-row-bidirectional-merge", 3, writeFixture3RowBidirectionalMerge, NewCellRange("A", 4, "C", 6), true, nil},
		{"3row", "3-row-mixed-merge", 3, writeFixture3RowMixedMerge, NewCellRange("A", 4, "H", 6), false, nil},
	}
}

func TestCreateSchemaHeaderFixturesInTestdata(t *testing.T) {
	skipUnlessUpdateEnabled(t)
	testCreateSchemaHeaderFixturesInTestdata(t, schemaHeaderCases())
}

func TestReadSchemaInternalHeader(t *testing.T) {
	testReadSchemaInternalHeaderTrack(t, schemaHeaderCases())
}

func testCreateSchemaHeaderFixturesInTestdata(t *testing.T, cases []schemaHeaderTrackCase) {
	destBase := filepath.Join("testdata", schemaHeaderFixturesDir)

	// One directory per case: schema_prep.json + styled workbook.xlsx.
	for _, tc := range cases {
		t.Run(tc.caseRelPath(), func(t *testing.T) {
			destDir := filepath.Join(destBase, tc.tier, tc.dir)
			if err := os.MkdirAll(destDir, 0o755); err != nil {
				t.Fatal(err)
			}

			// Bounds and is_ready for readSchemaInternal / ReadTableSchema.
			if err := writeSchemaPrepFile(destDir, fixtureSchemaPrep(tc.dir, tc.headerEndRow, tc.dataBand)); err != nil {
				t.Fatal(err)
			}
			prep, err := readSchemaPrepFixture(filepath.Join(destDir, "schema_prep.json"))
			if err != nil {
				t.Fatal(err)
			}
			workbookPath := filepath.Join(destDir, "workbook.xlsx")
			// Header layout + r{row}c{col} data band; styling from prep.
			if err := writeFixtureWorkbook(workbookPath, tc.dir, tc.writeWorkbook, prep); err != nil {
				t.Fatal(err)
			}
		})
	}

	// Per tier: multi-sheet workbook at schema_header/schema_header_<tier>_tests.xlsx.
	for _, tier := range []string{"1row", "2row", "3row"} {
		var tierCases []schemaHeaderTrackCase
		for _, tc := range cases {
			if tc.tier == tier {
				tierCases = append(tierCases, tc)
			}
		}
		if len(tierCases) == 0 {
			continue
		}
		sheets := make([]fixtureSheetCase, len(tierCases))
		for i, tc := range tierCases {
			destDir := filepath.Join(destBase, tc.tier, tc.dir)
			prep, err := readSchemaPrepFixture(filepath.Join(destDir, "schema_prep.json"))
			if err != nil {
				t.Fatal(err)
			}
			sheets[i] = fixtureSheetCase{dir: tc.dir, write: tc.writeWorkbook, prep: prep}
		}
		combinedPath := filepath.Join(destBase, "schema_header_"+tier+"_tests.xlsx")
		if err := writeFixtureCombinedWorkbook(combinedPath, sheets); err != nil {
			t.Fatal(err)
		}
	}
}

func testReadSchemaInternalHeaderTrack(t *testing.T, cases []schemaHeaderTrackCase) {
	for _, tc := range cases {
		t.Run(tc.caseRelPath(), func(t *testing.T) {
			caseDir := filepath.Join("testdata", schemaHeaderFixturesDir, tc.tier, tc.dir)
			prepPath := filepath.Join(caseDir, "schema_prep.json")
			excelPath := filepath.Join(caseDir, "workbook.xlsx")
			goldenPath := filepath.Join(caseDir, "table_schema.json")

			var prep SchemaPrep
			if err := jsonfile.Read(prepPath, &prep); err != nil {
				t.Fatalf("read schema_prep: %v", err)
			}

			got, err := readSchemaInternal(excelPath, prep)
			if tc.wantSchemaErr {
				if err == nil {
					t.Fatal("readSchemaInternal: expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("readSchemaInternal: %v", err)
			}

			if updateEnabled() {
				writeTableSchemaGolden(t, goldenPath, got)
				return
			}

			want, err := readTableSchemaGolden(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			fatalIfJSONGoldenMismatch(t, got, want)
		})
	}
}

func readTableSchemaGolden(path string) (TableSchema, error) {
	var schema TableSchema
	if err := jsonfile.Read(path, &schema); err != nil {
		return TableSchema{}, err
	}
	return schema, nil
}

func writeTableSchemaGolden(t *testing.T, path string, schema TableSchema) {
	if err := WriteSchemaToJson(path, schema); err != nil {
		t.Fatal(err)
	}
}

func fillSchemaHeaderDataBand(f *excelize.File, sheet string, band CellRange, skipColumns map[string]bool) error {
	startRow := band.StartRowNumber()
	endRow := band.EndRowNumber()
	startCol, err := band.StartColumnNumber()
	if err != nil {
		return err
	}
	endCol, err := band.EndColumnNumber()
	if err != nil {
		return err
	}

	for colNum := startCol; colNum <= endCol; colNum++ {
		colName, err := excelize.ColumnNumberToName(colNum)
		if err != nil {
			return err
		}
		if skipColumns != nil && skipColumns[colName] {
			continue
		}
		for row := startRow; row <= endRow; row++ {
			val := fmt.Sprintf("r%dc%s", row, colName)
			if err := setRawHeaderCell(f, sheet, colName, row, val); err != nil {
				return err
			}
		}
	}
	return nil
}

// Used by schema_data and schema_coupling fixture writers (same r{row}c{col} data band).
func fillSchemaHeader1RowDataBand(f *excelize.File, sheet string, band CellRange, skipColumns map[string]bool) error {
	return fillSchemaHeaderDataBand(f, sheet, band, skipColumns)
}
