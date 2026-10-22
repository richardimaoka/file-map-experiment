package excel

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Schema header fixtures (2-row) under excel/testdata/schema_header_2row/<case>/.
// Data band uses r{row}c{col} literals below header layouts from fixture_header_layouts_test.go.
//
// Regenerate workbooks and schema_prep.json:
//
//	go test ./excel -run TestCreateSchemaHeader2RowFixturesInTestdata
//
// Regenerate table_schema.json:
//
//	go test ./excel -update -run TestReadSchemaInternal2Row

type schemaHeader2RowCase struct {
	dir              string
	headerWrite      func(f *excelize.File, sheet string) error
	dataBand         CellRange
	wantSchemaErr    bool
	skipDataColumns map[string]bool
}

func schemaHeader2RowCases() []schemaHeader2RowCase {
	return []schemaHeader2RowCase{
		{"2-row-simple", writeRawHeader2RowSimple, NewCellRange("A", 3, "C", 5), false, nil},
		{
			"2-row-duplicate-header",
			writeRawHeader2RowDuplicateHeader,
			NewCellRange("A", 3, "C", 5),
			true,
			nil,
		},
		{"2-row-empty-cell", writeRawHeader2RowEmptyCell, NewCellRange("A", 3, "C", 5), false, nil},
		{"2-row-horizontal-merge", writeRawHeader2RowHorizontalMerge, NewCellRange("A", 3, "C", 5), false, nil},
		{"2-row-vertical-merge", writeRawHeader2RowVerticalMerge, NewCellRange("A", 3, "C", 5), false, nil},
		{
			"2-row-bidirectional-merge",
			writeRawHeader2RowBidirectionalMerge,
			NewCellRange("A", 3, "C", 5),
			true,
			nil,
		},
		{
			"2-row-horizontal-merge-repeated",
			writeRawHeader2RowHorizontalMergeRepeated,
			NewCellRange("A", 3, "G", 5),
			false,
			nil,
		},
		{
			"2-row-vertical-merge-repeated",
			writeRawHeader2RowVerticalMergeRepeated,
			NewCellRange("A", 3, "E", 5),
			false,
			nil,
		},
		{
			"2-row-bidirectional-merge-rep",
			writeRawHeader2RowBidirectionalMergeRepeated,
			NewCellRange("A", 3, "I", 5),
			true,
			nil,
		},
		{
			"2-row-mixed-merge",
			writeRawHeader2RowMixedMerge,
			NewCellRange("A", 3, "L", 5),
			true,
			nil,
		},
	}
}

func (tc schemaHeader2RowCase) writeWorkbook(f *excelize.File, sheet string) error {
	if err := tc.headerWrite(f, sheet); err != nil {
		return err
	}
	return fillSchemaHeader1RowDataBand(f, sheet, tc.dataBand, tc.skipDataColumns)
}

func TestCreateSchemaHeader2RowFixturesInTestdata(t *testing.T) {
	skipUnlessUpdateEnabled(t)
	destBase := filepath.Join("testdata", "schema_header_2row")
	cases := schemaHeader2RowCases()

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			destDir := filepath.Join(destBase, tc.dir)
			if err := os.MkdirAll(destDir, 0o755); err != nil {
				t.Fatal(err)
			}

			if err := writeSchemaPrepFile(destDir, fixtureSchemaPrep(tc.dir, 2, tc.dataBand)); err != nil {
				t.Fatal(err)
			}
			prep, err := readSchemaPrepFixture(filepath.Join(destDir, "schema_prep.json"))
			if err != nil {
				t.Fatal(err)
			}
			workbookPath := filepath.Join(destDir, "workbook.xlsx")
			if err := writeRawHeaderWorkbook(workbookPath, tc.dir, tc.writeWorkbook, prep); err != nil {
				t.Fatal(err)
			}
		})
	}

	sheets := make([]fixtureSheetCase, len(cases))
	for i, tc := range cases {
		prep, err := readSchemaPrepFixture(filepath.Join(destBase, tc.dir, "schema_prep.json"))
		if err != nil {
			t.Fatal(err)
		}
		sheets[i] = fixtureSheetCase{dir: tc.dir, write: tc.writeWorkbook, prep: prep}
	}
	combinedPath := fixtureCombinedWorkbookPath(destBase)
	if err := writeFixtureCombinedWorkbook(combinedPath, sheets); err != nil {
		t.Fatal(err)
	}
}

// Golden TableSchema (2-row header track): header_name is row labels joined with " > "
// (resolvedMultiRowHeaderName); is_row_grouped false until data-band grouping exists.
func TestReadSchemaInternal2Row(t *testing.T) {
	for _, tc := range schemaHeader2RowCases() {
		t.Run(tc.dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "schema_header_2row", tc.dir)
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

