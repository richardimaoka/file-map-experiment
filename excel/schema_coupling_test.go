package excel

import (
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
}

func schemaCouplingCases() []schemaCouplingCase {
	band2Row := NewCellRange("A", 3, "C", 5)
	return []schemaCouplingCase{
		{
			dir:           "coupling-2row-data-vmerge",
			headerWrite:   writeFixture2RowSimple,
			dataBand:      band2Row,
			dataAfterFill: mergeDataVerticalMergeA,
		},
		{
			dir:           "coupling-2row-hdr-data-vmerge",
			headerWrite:   writeFixture2RowVerticalMerge,
			dataBand:      band2Row,
			dataAfterFill: mergeDataVerticalMergeA,
		},
		{
			dir:           "coupling-2row-empty-header-data",
			headerWrite:   writeCoupling2RowEmptyHeaderColumnC,
			dataBand:      band2Row,
			wantSchemaErr: true,
		},
	}
}

// Vertical merge A3:A4 in the standard 2-row data band (rows 3–5).
func mergeDataVerticalMergeA(f *excelize.File, sheet string) error {
	return mergeRawHeader(f, sheet, "A", 3, "A", 4)
}

// No resolved header_name for column C (both header rows empty) but data band still fills C.
func writeCoupling2RowEmptyHeaderColumnC(f *excelize.File, sheet string) error {
	if err := writeFixture2RowSimple(f, sheet); err != nil {
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
	skipUnlessUpdateEnabled(t)
	destBase := filepath.Join("testdata", "schema_coupling")
	cases := schemaCouplingCases()

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
			if err := writeFixtureWorkbook(workbookPath, tc.dir, tc.writeWorkbook, prep); err != nil {
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

