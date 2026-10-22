package excel

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestValidateExcelAgainstSchema_readAPIFixture_matches(t *testing.T) {
	resolver := NewPathResolver(readAPICaseRelPath)
	var got TableSchema
	stderr := captureStderr(func() {
		got = ValidateExcelAgainstSchema(*resolver, "data-all-string")
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
