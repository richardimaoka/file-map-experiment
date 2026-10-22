package excel

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Shared case shape for schema_header_{1,2,3}row tracks.
type schemaHeaderTrackCase struct {
	dir              string
	headerWrite      func(f *excelize.File, sheet string) error
	dataBand         CellRange
	wantSchemaErr    bool
	skipDataColumns map[string]bool
}

func (tc schemaHeaderTrackCase) writeWorkbook(f *excelize.File, sheet string) error {
	if err := tc.headerWrite(f, sheet); err != nil {
		return err
	}
	return fillSchemaHeader1RowDataBand(f, sheet, tc.dataBand, tc.skipDataColumns)
}

func testCreateSchemaHeaderFixturesInTestdata(t *testing.T, trackDir string, headerEndRow int, cases []schemaHeaderTrackCase) {
	destBase := filepath.Join("testdata", trackDir)

	// One directory per case: schema_prep.json + styled workbook.xlsx.
	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			destDir := filepath.Join(destBase, tc.dir)
			if err := os.MkdirAll(destDir, 0o755); err != nil {
				t.Fatal(err)
			}

			// Bounds and is_ready for readSchemaInternal / ReadTableSchema.
			if err := writeSchemaPrepFile(destDir, fixtureSchemaPrep(tc.dir, headerEndRow, tc.dataBand)); err != nil {
				t.Fatal(err)
			}
			prep, err := readSchemaPrepFixture(filepath.Join(destDir, "schema_prep.json"))
			if err != nil {
				t.Fatal(err)
			}
			workbookPath := filepath.Join(destDir, "workbook.xlsx")
			// Header layout + r{row}c{col} data band; styling from prep.
			if err := writeRawHeaderWorkbook(workbookPath, tc.dir, tc.writeWorkbook, prep); err != nil {
				t.Fatal(err)
			}
		})
	}

	// Single multi-sheet workbook at <track>/<track>_tests.xlsx (one sheet per case).
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

func testReadSchemaInternalHeaderTrack(t *testing.T, trackDir string, cases []schemaHeaderTrackCase) {
	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", trackDir, tc.dir)
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
