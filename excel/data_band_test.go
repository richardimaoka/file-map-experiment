package excel

import (
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonutil"
	"github.com/xuri/excelize/v2"
)

// detectDataEndRow matches prep data_end_row on schema_data goldens (no trailing empty band).
func TestDetectDataEndRow_schemaDataFixtures(t *testing.T) {
	caseDir := filepath.Join("testdata", "schema_data", "data-all-string")
	var prep SchemaPrep
	if err := jsonutil.Read(filepath.Join(caseDir, "schema_prep.json"), &prep); err != nil {
		t.Fatalf("read prep: %v", err)
	}
	var schema TableSchema
	if err := jsonutil.Read(filepath.Join(caseDir, "table_schema.json"), &schema); err != nil {
		t.Fatalf("read schema: %v", err)
	}

	f, err := excelize.OpenFile(filepath.Join(caseDir, "workbook.xlsx"))
	if err != nil {
		t.Fatalf("open workbook: %v", err)
	}
	defer f.Close()

	merges, err := loadMergedRanges(f, schema.SheetName)
	if err != nil {
		t.Fatalf("loadMergedRanges: %v", err)
	}
	got, err := detectDataEndRow(f, schema, merges)
	if err != nil {
		t.Fatalf("detectDataEndRow: %v", err)
	}
	if got != prep.DataEndRow {
		t.Fatalf("data end row: got %d want %d (prep band)", got, prep.DataEndRow)
	}
}
