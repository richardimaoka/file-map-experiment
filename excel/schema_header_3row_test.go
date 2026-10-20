package excel

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Schema header fixtures (3-row) under excel/testdata/schema_header_3row/<case>/.
// Header layout matches raw_header for the same case name; data band uses r{row}c{col} literals.
//
// Regenerate workbooks and sync schema_prep.json / raw_headers.json from raw_header:
//
//	go test ./excel -run TestCreateSchemaHeader3RowFixturesInTestdata
//
// Regenerate table_schema.json:
//
//	go test ./excel -update -run TestReadSchemaInternal3Row

type schemaHeader3RowCase struct {
	dir              string
	headerWrite      func(f *excelize.File, sheet string) error
	dataBand         CellRange
	wantSchemaErr    bool
	skipDataColumns  map[string]bool
	rawHeaderFixture string
}

func schemaHeader3RowCases() []schemaHeader3RowCase {
	return []schemaHeader3RowCase{
		{"3-row-simple", writeRawHeader3RowSimple, NewCellRange("A", 4, "C", 6), false, nil, ""},
		{
			"3-row-bidirectional-merge",
			writeRawHeader3RowBidirectionalMerge,
			NewCellRange("A", 4, "C", 6),
			true,
			nil,
			"",
		},
		{"3-row-mixed-merge", writeRawHeader3RowMixedMerge, NewCellRange("A", 4, "H", 6), false, nil, ""},
	}
}

func (tc schemaHeader3RowCase) writeWorkbook(f *excelize.File, sheet string) error {
	if err := tc.headerWrite(f, sheet); err != nil {
		return err
	}
	return fillSchemaHeader1RowDataBand(f, sheet, tc.dataBand, tc.skipDataColumns)
}

func (tc schemaHeader3RowCase) rawHeaderSourceDir() string {
	if tc.rawHeaderFixture != "" {
		return tc.rawHeaderFixture
	}
	return tc.dir
}

func TestCreateSchemaHeader3RowFixturesInTestdata(t *testing.T) {
	destBase := filepath.Join("testdata", "schema_header_3row")
	srcBase := filepath.Join("testdata", "raw_header")
	cases := schemaHeader3RowCases()

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
			if err := writeSchemaPrepFor3RowCase(destDir, tc); err != nil {
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

// Golden TableSchema (3-row header track): header_name joins three row labels with " > ".
func TestReadSchemaInternal3Row(t *testing.T) {
	for _, tc := range schemaHeader3RowCases() {
		t.Run(tc.dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "schema_header_3row", tc.dir)
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

func writeSchemaPrepFor3RowCase(destDir string, tc schemaHeader3RowCase) error {
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
