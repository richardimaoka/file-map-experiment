package excel

import (
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

const readAPICaseRelPath = "excel/testdata/read_api/data_all_string"

// Public read API: resolver layout with config.json, cached schema JSON, and workbook.
func TestReadSchemaFromJson(t *testing.T) {
	schemaPath := filepath.Join("testdata", "schema_data", "data-all-string", "table_schema.json")
	got, err := ReadSchemaFromJson(schemaPath)
	if err != nil {
		t.Fatalf("ReadSchemaFromJson: %v", err)
	}
	if got.SchemaName != "data-all-string" || len(got.Columns) != 3 {
		t.Fatalf("unexpected schema: name=%q columns=%d", got.SchemaName, len(got.Columns))
	}
}

func TestReadSchemaFromExcel_readAPIFixture(t *testing.T) {
	resolver := NewPathResolver(readAPICaseRelPath)
	got, err := ReadSchemaFromExcel(*resolver, "data-all-string")
	if err != nil {
		t.Fatalf("ReadSchemaFromExcel: %v", err)
	}
	wantPath := resolver.SchemaPath("data-all-string")
	want, err := ReadSchemaFromJson(wantPath)
	if err != nil {
		t.Fatalf("read cached schema: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, got, want)
}

func TestReadRowsFromExcel_readAPIFixture(t *testing.T) {
	resolver := NewPathResolver(readAPICaseRelPath)
	got, err := ReadRowsFromExcel(*resolver, "data-all-string")
	if err != nil {
		t.Fatalf("ReadRowsFromExcel: %v", err)
	}
	goldenPath := filepath.Join("testdata", "schema_data", "data-all-string", "rows_data.json")
	var want []RowData
	if err := jsonfile.Read(goldenPath, &want); err != nil {
		t.Fatalf("read rows golden: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, got, want)
}
