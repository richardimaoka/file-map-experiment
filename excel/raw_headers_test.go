package excel

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

// readRawHeaders golden files: testdata/raw_header/<case>/raw_headers.json
func TestReadRawHeaders(t *testing.T) {
	for _, dir := range rawHeaderCaseDirs() {
		if !strings.HasPrefix(dir, "1-row-") {
			continue
		}
		t.Run(dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "raw_header", dir)
			prepPath := filepath.Join(caseDir, "schema_prep.json")
			excelPath := filepath.Join(caseDir, "workbook.xlsx")
			goldenPath := filepath.Join(caseDir, "raw_headers.json")

			var prep SchemaPrep
			if err := jsonfile.Read(prepPath, &prep); err != nil {
				t.Fatalf("read schema_prep: %v", err)
			}

			got, err := readRawHeaders(excelPath, prep)
			if err != nil {
				t.Fatalf("readRawHeaders: %v", err)
			}

			if os.Getenv("UPDATE_RAW_HEADER_GOLDEN") != "" {
				writeRawHeadersGolden(t, goldenPath, got)
				return
			}

			want, err := readRawHeadersGolden(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			if !reflect.DeepEqual(got, want) {
				gotJSON, _ := json.MarshalIndent(got, "", "  ")
				wantJSON, _ := json.MarshalIndent(want, "", "  ")
				t.Fatalf("raw headers mismatch:\n--- got\n%s\n--- want\n%s", gotJSON, wantJSON)
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

func readRawHeadersGolden(path string) ([]RawHeader, error) {
	var headers []RawHeader
	if err := jsonfile.Read(path, &headers); err != nil {
		return nil, err
	}
	return headers, nil
}

func writeRawHeadersGolden(t *testing.T, path string, headers []RawHeader) {
	data, err := json.MarshalIndent(headers, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
