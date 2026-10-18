package excel

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

// readRawColumns golden files: testdata/raw_header/<case>/raw_columns.json
func TestReadRawColumns(t *testing.T) {
	for _, dir := range rawHeaderCaseDirs() {
		t.Run(dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "raw_header", dir)
			prepPath := filepath.Join(caseDir, "schema_prep.json")
			excelPath := filepath.Join(caseDir, "workbook.xlsx")
			goldenPath := filepath.Join(caseDir, "raw_columns.json")

			var prep SchemaPrep
			if err := jsonfile.Read(prepPath, &prep); err != nil {
				t.Fatalf("read schema_prep: %v", err)
			}

			got, err := readRawColumns(excelPath, prep)
			if err != nil {
				t.Fatalf("readRawColumns: %v", err)
			}

			if os.Getenv("UPDATE_RAW_HEADER_GOLDEN") != "" {
				writeRawColumnsGolden(t, goldenPath, got)
				return
			}

			want, err := readRawColumnsGolden(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			if !reflect.DeepEqual(got, want) {
				gotJSON, _ := json.MarshalIndent(got, "", "  ")
				wantJSON, _ := json.MarshalIndent(want, "", "  ")
				t.Fatalf("raw columns mismatch:\n--- got\n%s\n--- want\n%s", gotJSON, wantJSON)
			}
		})
	}
}

func rawHeaderCaseDirs() []string {
	cases := rawHeaderCases()
	dirs := make([]string, len(cases))
	for i, c := range cases {
		dirs[i] = c.dir
	}
	return dirs
}

func readRawColumnsGolden(path string) ([]RawColumn, error) {
	var cols []RawColumn
	if err := jsonfile.Read(path, &cols); err != nil {
		return nil, err
	}
	return cols, nil
}

func writeRawColumnsGolden(t *testing.T, path string, cols []RawColumn) {
	data, err := json.MarshalIndent(cols, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
