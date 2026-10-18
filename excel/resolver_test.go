package excel

import (
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/env"
)

const readyCaseRelPath = "excel/testdata/resolver/ready/cases/client1"

// Case tree under APP_BASE_DIR/<relPath>/ must match PathResolver's documented layout.
// APP_BASE_DIR comes from .env (see TestMain) or defaults to the module root.
func TestPathResolver_paths(t *testing.T) {
	r := NewPathResolver(readyCaseRelPath)
	caseDir := filepath.Join(env.AppBaseDir(), readyCaseRelPath)

	if r.DirPath() != caseDir {
		t.Fatalf("case root: got %q want %q", r.DirPath(), caseDir)
	}

	config := InitConfig("folder_map.xlsx")

	// Each entry is a path under the case root; absolute must equal Join(caseDir, relative).
	layout := []struct {
		relative string
		absolute string
	}{
		{"folder_map.xlsx", r.ExcelPath(config)},
		{"config.json", r.ConfigPath()},
		{"schemas", r.SchemaCollectionPath()},
		{"schemas/sheet1.json", r.SchemaPath("sheet1")},
		{"schemas/sheet2.json", r.SchemaPath("sheet2")},
	}
	for _, ent := range layout {
		t.Run(ent.relative, func(t *testing.T) {
			want := filepath.Join(caseDir, ent.relative)
			if ent.absolute != want {
				t.Errorf("got %q want %q", ent.absolute, want)
			}
		})
	}
}

// ValidateDir delegates to ReadExcelConfig; config edge cases live in config_test.go.
func TestPathResolver_ValidateDir(t *testing.T) {
	cases := []struct {
		name    string
		relPath string // under APP_BASE_DIR
		wantErr bool
	}{
		// Case tree has valid config.json and workbook.
		{name: "ready", relPath: readyCaseRelPath},
		// Case directory exists but config.json is absent.
		{
			name:    "not_ready",
			relPath: "excel/testdata/resolver/not_ready/cases/client1",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewPathResolver(tc.relPath).ValidateDir()
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ValidateDir: %v", err)
			}
		})
	}
}

// NewPathResolver delegates to env.AppBaseDir, which panics if APP_BASE_DIR is unset.
func TestNewPathResolver_panicsWithoutAppBaseDir(t *testing.T) {
	t.Setenv("APP_BASE_DIR", "")

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when APP_BASE_DIR is unset")
		}
	}()
	NewPathResolver("cases/client1")
}
