package excel

import (
	"path/filepath"
	"testing"
)

// Case tree under APP_BASE_DIR/<dirPath>/ must match PathResolver's documented layout.
func TestPathResolver_paths(t *testing.T) {
	base := filepath.Join("testdata", "resolver", "ready")
	t.Setenv("APP_BASE_DIR", base)

	r := NewPathResolver("cases/client1")
	caseDir := filepath.Join(base, "cases", "client1")

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
	const relPath = "cases/client1"

	cases := []struct {
		name    string
		fixture string // testdata/resolver/<fixture> as APP_BASE_DIR
		wantErr bool
	}{
		// Case tree has valid config.json and workbook.
		{name: "ready", fixture: "ready"},
		// Case directory exists but config.json is absent.
		{name: "not_ready", fixture: "not_ready", wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("APP_BASE_DIR", filepath.Join("testdata", "resolver", tc.fixture))

			err := NewPathResolver(relPath).ValidateDir()
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
