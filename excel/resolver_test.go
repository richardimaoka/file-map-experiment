package excel

import (
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/env"
)

const readyCaseRelPath = "excel/testdata/resolver/ready/cases/client1"

// Case tree under APP_BASE_DIR/<relDirPath>/ must match PathResolver's documented layout.
// APP_BASE_DIR comes from .env (see TestMain) or defaults to the module root.
func TestResolver_paths(t *testing.T) {
	r := NewResolver(readyCaseRelPath)
	caseDir := filepath.Join(env.AppBaseDir(), readyCaseRelPath)

	if r.FullDirPath() != caseDir {
		t.Fatalf("case root: got %q want %q", r.FullDirPath(), caseDir)
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

// ValidateDir delegates to ReadConfig; config edge cases live in config_test.go.
func TestResolver_ValidateDir(t *testing.T) {
	cases := []struct {
		name       string
		relDirPath string // under APP_BASE_DIR
		wantErr    bool
	}{
		// Case tree has valid config.json and workbook.
		{name: "ready", relDirPath: readyCaseRelPath},
		// Case directory exists but config.json is absent.
		{
			name:       "not_ready",
			relDirPath: "excel/testdata/resolver/not_ready/cases/client1",
			wantErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewResolver(tc.relDirPath).ValidateDir()
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

// NewResolver delegates to env.AppBaseDir, which panics if APP_BASE_DIR is unset.
func TestNewResolver_panicsWithoutAppBaseDir(t *testing.T) {
	t.Setenv("APP_BASE_DIR", "")

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when APP_BASE_DIR is unset")
		}
	}()
	NewResolver("cases/client1")
}
