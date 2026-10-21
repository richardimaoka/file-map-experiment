package excel

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

// Golden rows_data.json: one RowData map per data-band row (schema_data track first).
func TestReadRowsSchemaData(t *testing.T) {
	for _, tc := range schemaDataCases() {
		if tc.wantSchemaErr {
			continue
		}
		t.Run(tc.dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "schema_data", tc.dir)
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
