package instruction

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/richardimaoka/file-map-experiment/excel"
	"github.com/richardimaoka/file-map-experiment/fs"
)

// BuildDir rewrites defaults/<schemaName>.json for ready schema_preps after excel.BuildDir.
func TestBuildDir_writesDefaultData(t *testing.T) {
	const rel = "instruction_build_case"
	fixture := filepath.Join("..", "excel", "testdata", "read_api", "data_all_string")
	sourceExcel := filepath.Join(fixture, "workbook.xlsx")
	goldenRows := filepath.Join("..", "excel", "testdata", "schema_data", "data-all-string", "rows_data.json")

	appBase := t.TempDir()
	t.Setenv("APP_BASE_DIR", appBase)
	caseDir := filepath.Join(appBase, rel)
	if err := fs.CopyDir(fixture, caseDir); err != nil {
		t.Fatalf("CopyDir: %v", err)
	}

	if err := BuildDir(sourceExcel, rel); err != nil {
		t.Fatalf("BuildDir: %v", err)
	}

	want, err := excel.ReadRowsFromJson(goldenRows)
	if err != nil {
		t.Fatalf("ReadRowsFromJson golden: %v", err)
	}
	got, err := excel.ReadRowsFromJson(NewResolver(rel).DefaultDataPath("data-all-string"))
	if err != nil {
		t.Fatalf("ReadRowsFromJson defaults: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("defaults rows: got %v want %v", got, want)
	}
}
