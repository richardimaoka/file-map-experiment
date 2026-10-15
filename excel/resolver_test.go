package excel

import (
	"path/filepath"
	"testing"
)

// Verifies path joins for the case directory layout (see PathResolver doc); no files created.
func TestPathResolver_paths(t *testing.T) {
	base := t.TempDir()
	t.Setenv("APP_BASE_DIR", base)

	r := NewPathResolver("cases/client1", "folder_map.xlsx")
	caseDir := filepath.Join(base, "cases", "client1")

	if r.DirPath() != caseDir {
		t.Fatalf("DirPath: got %q want %q", r.DirPath(), caseDir)
	}

	checks := []struct {
		name string
		got  string
		want string
	}{
		{"ExcelPath", r.ExcelPath(), filepath.Join(caseDir, "folder_map.xlsx")},
		{"ConfigPath", r.ConfigPath(), filepath.Join(caseDir, "config.json")},
		{"SchemaDirPath", r.SchemaDirPath(), filepath.Join(caseDir, "schemas")},
		{"SchemaPath", r.SchemaPath("sheet1"), filepath.Join(caseDir, "schemas", "sheet1.json")},
	}
	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			if c.got != c.want {
				t.Errorf("got %q want %q", c.got, c.want)
			}
		})
	}
}

// NewPathResolver delegates to env.AppBaseDir, which panics if APP_BASE_DIR is empty.
func TestNewPathResolver_panicsWithoutAppBaseDir(t *testing.T) {
	t.Setenv("APP_BASE_DIR", "")

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when APP_BASE_DIR is unset")
		}
	}()
	NewPathResolver("cases/client1", "folder_map.xlsx")
}
