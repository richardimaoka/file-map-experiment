package excel

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

// Schema header fixtures (1-row) under excel/testdata/schema_header_1row/<case>/.
// Layouts match excel/testdata/raw_header for the same case names (copied reference goldens).
//
// Regenerate workbooks and sync schema_prep.json / raw_headers.json from raw_header:
//
//	go test ./excel -run TestCreateSchemaHeader1RowFixturesInTestdata
//
// Regenerate table_schema.json (once readSchemaInternal is implemented):
//
//	go test ./excel -update -run TestReadSchemaInternal1Row

func schemaHeader1RowCases() []rawHeaderCase {
	return []rawHeaderCase{
		{"1-row-simple", writeRawHeader1RowSimple},
		{"1-row-empty-header", writeRawHeader1RowEmptyHeader},
		{"1-row-line-break", writeRawHeader1RowLineBreak},
		{"1-row-horizontal-merge", writeRawHeader1RowHorizontalMerge},
	}
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
			if err := writeRawHeaderWorkbook(workbookPath, tc.dir, tc.write); err != nil {
				t.Fatal(err)
			}

			srcDir := filepath.Join(srcBase, tc.dir)
			for _, name := range []string{"schema_prep.json", "raw_headers.json"} {
				if err := copySchemaHeaderFixtureFile(filepath.Join(srcDir, name), filepath.Join(destDir, name)); err != nil {
					t.Fatal(err)
				}
			}
		})
	}

	sheets := make([]fixtureSheetCase, len(cases))
	for i, tc := range cases {
		sheets[i] = fixtureSheetCase{dir: tc.dir, write: tc.write}
	}
	combinedPath := filepath.Join(destBase, "test_cases.xlsx")
	if err := writeFixtureCombinedWorkbook(combinedPath, sheets); err != nil {
		t.Fatal(err)
	}
}

// Golden TableSchema assumptions (1-row header track): prep supplies schema_name, data_start_row,
// and column bounds; one ColumnSchema per in-range column; header_name is the header-row cell
// value; is_row_grouped is false when header_end_row == header_start_row.
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
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
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
