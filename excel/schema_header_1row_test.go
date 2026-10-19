package excel

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Schema header fixtures (1-row) under excel/testdata/schema_header_1row/<case>/.
// Header layout matches raw_header for the same case name; data band is filled with r{row}c{col} literals.
//
// Regenerate workbooks and sync schema_prep.json / raw_headers.json from raw_header:
//
//	go test ./excel -run TestCreateSchemaHeader1RowFixturesInTestdata
//
// Regenerate table_schema.json (once readSchemaInternal is implemented):
//
//	go test ./excel -update -run TestReadSchemaInternal1Row

type schemaHeader1RowCase struct {
	dir         string
	headerWrite func(f *excelize.File, sheet string) error
	dataBand    CellRange
	// data in a column with no schema header name (RowData key) must fail readSchemaInternal.
	wantSchemaErr bool
	// omit data band cells for these columns (e.g. empty header column with no RowData key).
	skipDataColumns map[string]bool
	// copy schema_prep.json / raw_headers.json from this raw_header case (default: same as dir).
	rawHeaderFixture string
}

func schemaHeader1RowCases() []schemaHeader1RowCase {
	return []schemaHeader1RowCase{
		{"1-row-simple", writeRawHeader1RowSimple, NewCellRange("A", 2, "C", 4), false, nil, ""},
		{"1-row-empty-header", writeRawHeader1RowEmptyHeader, NewCellRange("A", 2, "D", 4), true, nil, ""},
		{
			"1-row-empty-header-valid",
			writeRawHeader1RowEmptyHeader,
			NewCellRange("A", 2, "D", 4),
			false,
			map[string]bool{"C": true},
			"1-row-empty-header",
		},
		{"1-row-line-break", writeRawHeader1RowLineBreak, NewCellRange("A", 2, "B", 4), false, nil, ""},
		{"1-row-horizontal-merge", writeRawHeader1RowHorizontalMerge, NewCellRange("A", 2, "D", 4), true, nil, ""},
	}
}

func (tc schemaHeader1RowCase) writeWorkbook(f *excelize.File, sheet string) error {
	if err := tc.headerWrite(f, sheet); err != nil {
		return err
	}
	return fillSchemaHeader1RowDataBand(f, sheet, tc.dataBand, tc.skipDataColumns)
}

func (tc schemaHeader1RowCase) rawHeaderSourceDir() string {
	if tc.rawHeaderFixture != "" {
		return tc.rawHeaderFixture
	}
	return tc.dir
}

func TestCreateSchemaHeader1RowFixturesInTestdata(t *testing.T) {
	destBase := filepath.Join("testdata", "schema_header_1row")
	srcBase := filepath.Join("testdata", "raw_header")
	cases := schemaHeader1RowCases()

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
			if err := writeSchemaPrepForCase(destDir, tc); err != nil {
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

// Golden TableSchema assumptions (1-row header track): prep supplies schema_name, data_start_row,
// and column bounds; header_name from the header row; skip empty names; merged headers only on
// the merge anchor column; is_row_grouped is false until data-range grouping is implemented.
func TestReadSchemaInternal1Row(t *testing.T) {
	for _, tc := range schemaHeader1RowCases() {
		t.Run(tc.dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "schema_header_1row", tc.dir)
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

func readTableSchemaGolden(path string) (TableSchema, error) {
	var schema TableSchema
	if err := jsonfile.Read(path, &schema); err != nil {
		return TableSchema{}, err
	}
	return schema, nil
}

func writeTableSchemaGolden(t *testing.T, path string, schema TableSchema) {
	data, err := marshalJSONGoldenIndent(schema)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeSchemaPrepForCase(destDir string, tc schemaHeader1RowCase) error {
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

func fillSchemaHeader1RowDataBand(f *excelize.File, sheet string, band CellRange, skipColumns map[string]bool) error {
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

func copySchemaHeaderFixtureFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
