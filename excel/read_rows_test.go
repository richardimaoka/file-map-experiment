package excel

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

type readRowsGoldenCase struct {
	dir           string
	wantSchemaErr bool
}

func testReadRowsGoldenTrack(t *testing.T, trackDir string, cases []readRowsGoldenCase) {
	for _, tc := range cases {
		if tc.wantSchemaErr {
			continue
		}
		t.Run(tc.dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", trackDir, tc.dir)
			prepPath := filepath.Join(caseDir, "schema_prep.json")
			excelPath := filepath.Join(caseDir, "workbook.xlsx")
			schemaPath := filepath.Join(caseDir, "table_schema.json")
			goldenPath := filepath.Join(caseDir, "rows_data.json")

			var prep SchemaPrep
			if err := jsonfile.Read(prepPath, &prep); err != nil {
				t.Fatalf("read schema_prep: %v", err)
			}
			var schema TableSchema
			if err := jsonfile.Read(schemaPath, &schema); err != nil {
				t.Fatalf("read table_schema: %v", err)
			}

			got, err := readRowsFromExcelInternal(excelPath, prep, schema)
			if err != nil {
				t.Fatalf("readRowsFromExcelInternal: %v", err)
			}

			if updateGoldens() {
				writeRowsDataGolden(t, goldenPath, got)
				return
			}

			want, err := readRowsDataGolden(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			fatalIfJSONGoldenMismatch(t, got, want)
		})
	}
}

// Golden rows_data.json: one RowData map per data-band row (schema_data track).
func TestReadRowsSchemaData(t *testing.T) {
	cases := make([]readRowsGoldenCase, 0, len(schemaDataCases()))
	for _, tc := range schemaDataCases() {
		cases = append(cases, readRowsGoldenCase{dir: tc.dir, wantSchemaErr: tc.wantSchemaErr})
	}
	testReadRowsGoldenTrack(t, "schema_data", cases)
}

func TestReadRowsSchemaHeader1Row(t *testing.T) {
	cases := make([]readRowsGoldenCase, 0, len(schemaHeader1RowCases()))
	for _, tc := range schemaHeader1RowCases() {
		cases = append(cases, readRowsGoldenCase{dir: tc.dir, wantSchemaErr: tc.wantSchemaErr})
	}
	testReadRowsGoldenTrack(t, "schema_header_1row", cases)
}

func TestReadRowsSchemaHeader2Row(t *testing.T) {
	cases := make([]readRowsGoldenCase, 0, len(schemaHeader2RowCases()))
	for _, tc := range schemaHeader2RowCases() {
		cases = append(cases, readRowsGoldenCase{dir: tc.dir, wantSchemaErr: tc.wantSchemaErr})
	}
	testReadRowsGoldenTrack(t, "schema_header_2row", cases)
}

func TestReadRowsSchemaHeader3Row(t *testing.T) {
	cases := make([]readRowsGoldenCase, 0, len(schemaHeader3RowCases()))
	for _, tc := range schemaHeader3RowCases() {
		cases = append(cases, readRowsGoldenCase{dir: tc.dir, wantSchemaErr: tc.wantSchemaErr})
	}
	testReadRowsGoldenTrack(t, "schema_header_3row", cases)
}

func TestReadRowsSchemaCoupling(t *testing.T) {
	cases := make([]readRowsGoldenCase, 0, len(schemaCouplingCases()))
	for _, tc := range schemaCouplingCases() {
		cases = append(cases, readRowsGoldenCase{dir: tc.dir, wantSchemaErr: tc.wantSchemaErr})
	}
	testReadRowsGoldenTrack(t, "schema_coupling", cases)
}

func readRowsDataGolden(path string) ([]RowData, error) {
	var rows []RowData
	if err := jsonfile.Read(path, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func writeRowsDataGolden(t *testing.T, path string, rows []RowData) {
	data, err := marshalJSONGoldenIndent(rows)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
