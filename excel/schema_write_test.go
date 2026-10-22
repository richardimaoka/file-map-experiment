package excel

import (
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/fs"
)

func TestWriteSchemaToJson_roundTrip(t *testing.T) {
	resolver := NewPathResolver(readAPICaseRelPath)
	want, err := ReadSchemaFromJson(resolver.SchemaPath("data-all-string"))
	if err != nil {
		t.Fatalf("ReadSchemaFromJson: %v", err)
	}
	path := filepath.Join(t.TempDir(), "schema.json")
	if err := WriteSchemaToJson(path, want); err != nil {
		t.Fatalf("WriteSchemaToJson: %v", err)
	}
	got, err := ReadSchemaFromJson(path)
	if err != nil {
		t.Fatalf("ReadSchemaFromJson: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, got, want)
}

func TestWriteSchemaFromExcel_readAPIFixture(t *testing.T) {
	resolver := NewPathResolver(readAPICaseRelPath)
	want, err := ReadSchemaFromExcel(*resolver, "data-all-string")
	if err != nil {
		t.Fatalf("ReadSchemaFromExcel: %v", err)
	}
	dir := t.TempDir()
	rel := "write_schema_case"
	caseDir := filepath.Join(dir, rel)
	if err := fs.CopyDir(resolver.DirPath(), caseDir); err != nil {
		t.Fatalf("CopyDir: %v", err)
	}
	t.Setenv("APP_BASE_DIR", dir)
	tmpResolver := NewPathResolver(rel)
	if err := WriteSchemaFromExcel(*tmpResolver, "data-all-string"); err != nil {
		t.Fatalf("WriteSchemaFromExcel: %v", err)
	}
	got, err := ReadSchemaFromJson(tmpResolver.SchemaPath("data-all-string"))
	if err != nil {
		t.Fatalf("ReadSchemaFromJson: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, got, want)
}

func TestWriteSchemaFromExcel_notReady(t *testing.T) {
	resolver := NewPathResolver("excel/testdata/build/existing_dir_with_config")
	err := WriteSchemaFromExcel(*resolver, "preserved_marker")
	if err == nil {
		t.Fatal("expected error when is_ready is false")
	}
}
