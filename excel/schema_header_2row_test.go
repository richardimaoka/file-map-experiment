package excel

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Schema header fixtures (2-row) under excel/testdata/schema_header_2row/<case>/.
// Header layout matches raw_header for the same case name; data band uses r{row}c{col} literals.
//
// Regenerate workbooks and sync schema_prep.json / raw_headers.json from raw_header:
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
	skipDataColumns  map[string]bool
	rawHeaderFixture string
}

func schemaHeader2RowCases() []schemaHeader2RowCase {
	return []schemaHeader2RowCase{
		{"2-row-simple", writeRawHeader2RowSimple, NewCellRange("A", 3, "C", 5), false, nil, ""},
		{"2-row-empty-cell", writeRawHeader2RowEmptyCell, NewCellRange("A", 3, "C", 5), false, nil, ""},
		{"2-row-horizontal-merge", writeRawHeader2RowHorizontalMerge, NewCellRange("A", 3, "C", 5), false, nil, ""},
		{"2-row-vertical-merge", writeRawHeader2RowVerticalMerge, NewCellRange("A", 3, "C", 5), false, nil, ""},
		{
			"2-row-bidirectional-merge",
			writeRawHeader2RowBidirectionalMerge,
			NewCellRange("A", 3, "C", 5),
			true,
			nil,
			"",
		},
		{
			"2-row-horizontal-merge-repeated",
			writeRawHeader2RowHorizontalMergeRepeated,
			NewCellRange("A", 3, "G", 5),
			false,
			nil,
			"",
		},
	}
}

func (tc schemaHeader2RowCase) writeWorkbook(f *excelize.File, sheet string) error {
	if err := tc.headerWrite(f, sheet); err != nil {
		return err
	}
	return fillSchemaHeader1RowDataBand(f, sheet, tc.dataBand, tc.skipDataColumns)
}

func (tc schemaHeader2RowCase) rawHeaderSourceDir() string {
	if tc.rawHeaderFixture != "" {
		return tc.rawHeaderFixture
	}
	return tc.dir
}

func TestCreateSchemaHeader2RowFixturesInTestdata(t *testing.T) {
	destBase := filepath.Join("testdata", "schema_header_2row")
	srcBase := filepath.Join("testdata", "raw_header")
	cases := schemaHeader2RowCases()

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			destDir := filepath.Join(destBase, tc.dir)
			if err := os.MkdirAll(destDir, 0o755); err != nil {
				t.Fatal(err)
			}

			workbookPath := filepath.Join(destDir, "workbook.xlsx")
			if err := writeRawHeaderWorkbook(workbookPath, tc.dir, tc.writeWorkbook); err != nil {
				t.Fatal(err)
			}

			srcDir := filepath.Join(srcBase, tc.rawHeaderSourceDir())
			if err := copySchemaHeaderFixtureFile(filepath.Join(srcDir, "raw_headers.json"), filepath.Join(destDir, "raw_headers.json")); err != nil {
				t.Fatal(err)
			}
			if err := writeSchemaPrepFor2RowCase(destDir, tc); err != nil {
				t.Fatal(err)
			}
		})
	}

	sheets := make([]fixtureSheetCase, len(cases))
	for i, tc := range cases {
		sheets[i] = fixtureSheetCase{dir: tc.dir, write: tc.writeWorkbook}
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

			if updateGoldens() {
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

func writeSchemaPrepFor2RowCase(destDir string, tc schemaHeader2RowCase) error {
	srcPath := filepath.Join("testdata", "raw_header", tc.rawHeaderSourceDir(), "schema_prep.json")
	var prep SchemaPrep
	if err := jsonfile.Read(srcPath, &prep); err != nil {
		return err
	}
	prep.SchemaName = tc.dir
	prep.SheetName = tc.dir
	data, err := json.MarshalIndent(prep, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(destDir, "schema_prep.json"), data, 0o644)
}
