package excel

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Schema coupling fixtures under excel/testdata/schema_coupling/<case>/.
// Combines multi-row header layouts with data-band patterns from schema_data.
//
// Regenerate workbooks and schema_prep.json:
//
//	go test ./excel -run TestCreateSchemaCouplingFixturesInTestdata
//
// Regenerate table_schema.json:
//
//	go test ./excel -update -run TestReadSchemaInternalCoupling

type schemaCouplingCase struct {
	dir             string
	headerWrite     func(f *excelize.File, sheet string) error
	dataBand        CellRange
	dataAfterFill   func(f *excelize.File, sheet string) error
	wantSchemaErr   bool
	skipDataColumns map[string]bool
	// raw_header case for schema_prep.json and optional raw_headers.json copy; empty skips raw_headers.
	rawHeaderRef string
}

func schemaCouplingCases() []schemaCouplingCase {
	band2Row := NewCellRange("A", 3, "C", 5)
	return []schemaCouplingCase{
		{
			dir:           "coupling-2row-data-vmerge",
			headerWrite:   writeRawHeader2RowSimple,
			dataBand:      band2Row,
			dataAfterFill: mergeDataVerticalMergeA,
			rawHeaderRef:  "2-row-simple",
		},
		{
			dir:           "coupling-2row-hdr-data-vmerge",
			headerWrite:   writeRawHeader2RowVerticalMerge,
			dataBand:      band2Row,
			dataAfterFill: mergeDataVerticalMergeA,
			rawHeaderRef:  "2-row-vertical-merge",
		},
		{
			dir:           "coupling-2row-empty-header-data",
			headerWrite:   writeCoupling2RowEmptyHeaderColumnC,
			dataBand:      band2Row,
			wantSchemaErr: true,
			rawHeaderRef:  "2-row-simple",
		},
	}
}

// Vertical merge A3:A4 in the standard 2-row data band (rows 3–5).
func mergeDataVerticalMergeA(f *excelize.File, sheet string) error {
	return mergeRawHeader(f, sheet, "A", 3, "A", 4)
}

// No resolved header_name for column C (both header rows empty) but data band still fills C.
func writeCoupling2RowEmptyHeaderColumnC(f *excelize.File, sheet string) error {
	if err := writeRawHeader2RowSimple(f, sheet); err != nil {
		return err
	}
	if err := clearRawHeaderCell(f, sheet, "C", 1); err != nil {
		return err
	}
	return clearRawHeaderCell(f, sheet, "C", 2)
}

func (tc schemaCouplingCase) writeWorkbook(f *excelize.File, sheet string) error {
	if err := tc.headerWrite(f, sheet); err != nil {
		return err
	}
	if err := fillSchemaHeader1RowDataBand(f, sheet, tc.dataBand, tc.skipDataColumns); err != nil {
		return err
	}
	if tc.dataAfterFill != nil {
		return tc.dataAfterFill(f, sheet)
	}
	return nil
}

func TestCreateSchemaCouplingFixturesInTestdata(t *testing.T) {
	destBase := filepath.Join("testdata", "schema_coupling")
	srcBase := filepath.Join("testdata", "raw_header")
	cases := schemaCouplingCases()

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			destDir := filepath.Join(destBase, tc.dir)
			if err := os.MkdirAll(destDir, 0o755); err != nil {
				t.Fatal(err)
			}

			if tc.rawHeaderRef != "" && tc.dir != "coupling-2row-empty-header-data" { // no raw_header match for cleared C headers
				srcDir := filepath.Join(srcBase, tc.rawHeaderRef)
				if err := copySchemaHeaderFixtureFile(filepath.Join(srcDir, "raw_headers.json"), filepath.Join(destDir, "raw_headers.json")); err != nil {
					t.Fatal(err)
				}
			}
			if err := writeSchemaPrepForCouplingCase(destDir, tc); err != nil {
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

// Golden TableSchema where header resolution and data-band grouping interact.
func TestReadSchemaInternalCoupling(t *testing.T) {
	for _, tc := range schemaCouplingCases() {
		t.Run(tc.dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "schema_coupling", tc.dir)
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

func writeSchemaPrepForCouplingCase(destDir string, tc schemaCouplingCase) error {
	ref := tc.rawHeaderRef
	if ref == "" {
		ref = "2-row-simple"
	}
	srcPath := filepath.Join("testdata", "raw_header", ref, "schema_prep.json")
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
