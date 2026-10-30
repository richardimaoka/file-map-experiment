package excel

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/richardimaoka/file-map-experiment/fs"
)

func TestWriteSchemaToJson_roundTrip(t *testing.T) {
	resolver := NewResolver(readAPICaseRelPath)
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
	resolver := NewResolver(readAPICaseRelPath)
	want, err := ReadSchemaFromExcel(resolver, "data-all-string")
	if err != nil {
		t.Fatalf("ReadSchemaFromExcel: %v", err)
	}
	dir := t.TempDir()
	rel := "write_schema_case"
	caseDir := filepath.Join(dir, rel)
	if err := fs.CopyDir(resolver.FullDirPath(), caseDir); err != nil {
		t.Fatalf("CopyDir: %v", err)
	}
	t.Setenv("APP_BASE_DIR", dir)
	tmpResolver := NewResolver(rel)
	if err := WriteSchemaFromExcel(tmpResolver, "data-all-string"); err != nil {
		t.Fatalf("WriteSchemaFromExcel: %v", err)
	}
	got, err := ReadSchemaFromJson(tmpResolver.SchemaPath("data-all-string"))
	if err != nil {
		t.Fatalf("ReadSchemaFromJson: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, got, want)
}

func TestWriteSchemaFromExcel_notReady(t *testing.T) {
	resolver := NewResolver("excel/testdata/build/existing_dir_with_config")
	err := WriteSchemaFromExcel(resolver, "preserved_marker")
	if err == nil {
		t.Fatal("expected error when is_ready is false")
	}
}

func TestValidateExcelAgainstSchema_readAPIFixture_matches(t *testing.T) {
	resolver := NewResolver(readAPICaseRelPath)
	var got TableSchema
	stderr := captureStderr(func() {
		got = ValidateExcelAgainstSchema(resolver, "data-all-string")
	})
	if stderr != "" {
		t.Fatalf("unexpected stderr:\n%s", stderr)
	}
	if got.SchemaName != "data-all-string" || len(got.Columns) != 3 {
		t.Fatalf("unexpected schema: name=%q columns=%d", got.SchemaName, len(got.Columns))
	}
}

func TestWarnExcelSchemaMismatch_shortMessage(t *testing.T) {
	persisted := TableSchema{SchemaName: "a", SheetName: "Sheet1"}
	fromExcel := TableSchema{SchemaName: "b", SheetName: "Sheet1"}
	stderr := captureStderr(func() {
		warnExcelSchemaMismatch("test-schema", persisted, fromExcel)
	})
	if !strings.Contains(stderr, "does not match persisted schema") {
		t.Fatalf("stderr: %q", stderr)
	}
	if strings.Contains(stderr, "- ") || strings.Contains(stderr, "+ ") {
		t.Fatalf("expected short message without JSON diff, got:\n%s", stderr)
	}
}

func captureStderr(fn func()) string {
	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stderr = w
	fn()
	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	r.Close()
	return buf.String()
}
