package instruction

import (
	"path/filepath"
	"testing"
)

func TestTopLevelResolver_formatSchemaPaths(t *testing.T) {
	const formatName = "client-a"
	moduleRoot := filepath.Join("..")
	t.Setenv("APP_BASE_DIR", moduleRoot)

	top := NewTopLevelResolver()
	rel := top.FormatSchemaRelDir(formatName)
	wantRel := filepath.Join("assets", "入稿指示書", "format_schemas", formatName)
	if rel != wantRel {
		t.Fatalf("FormatSchemaRelDir: got %q want %q", rel, wantRel)
	}

	abs := top.FormatSchemaDir(formatName)
	wantAbs := filepath.Join(moduleRoot, rel)
	if abs != wantAbs {
		t.Fatalf("FormatSchemaDir: got %q want %q", abs, wantAbs)
	}
}
