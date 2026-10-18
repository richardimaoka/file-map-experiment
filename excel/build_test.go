package excel

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/richardimaoka/file-map-experiment/fs"
)

// BuildExcelDir prepares a working directory from an existing Excel file.
func TestBuildExcelDir(t *testing.T) {
	excelName := "ads_report_simple.xlsx"
	sourceExcel := filepath.Join("testdata", "samples", excelName)
	sourceBytes := mustReadFile(t, sourceExcel)

	preservedConfig := mustParseConfigFixture(t, filepath.Join("testdata", "build", "existing_dir_with_config", "config.json"))

	cases := []buildExcelDirCase{
		// Target path does not exist yet: create it, copy the source Excel file, write config.json.
		{
			name:      "creates_target_dir_when_missing",
			wantExcel: sourceBytes,
		},
		// testdata/build/existing_dir_only: empty directory.
		{
			name:         "existing_dir_only",
			buildFixture: "existing_dir_only",
			wantExcel:    sourceBytes,
		},
		// testdata/build/existing_dir_with_config: config.json only.
		{
			name:         "existing_dir_with_config",
			buildFixture: "existing_dir_with_config",
			wantExcel:    sourceBytes,
			wantConfig:   preservedConfig,
		},
		// testdata/build/existing_dir_with_excel: Excel file only (overwritten from source).
		{
			name:         "existing_dir_with_excel",
			buildFixture: "existing_dir_with_excel",
			wantExcel:    sourceBytes,
		},
		// testdata/build/existing_dir_with_excel_and_config: Excel overwritten; config.json preserved.
		{
			name:         "existing_dir_with_excel_and_config",
			buildFixture: "existing_dir_with_excel_and_config",
			wantExcel:    sourceBytes,
			wantConfig:   preservedConfig,
		},
	}

	caseRelPath := "cases/client1"

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			appBase := t.TempDir()
			t.Setenv("APP_BASE_DIR", appBase)
			targetDir := filepath.Join(appBase, caseRelPath)
			if tc.buildFixture != "" {
				copyBuildFixtureIntoTarget(t, tc.buildFixture, targetDir)
			}

			excelPath := filepath.Join(targetDir, excelName)
			configPath := filepath.Join(targetDir, "config.json")

			if err := BuildExcelDir(sourceExcel, caseRelPath); err != nil {
				t.Fatalf("BuildExcelDir: %v", err)
			}

			if info, err := os.Stat(targetDir); err != nil || !info.IsDir() {
				t.Fatalf("target dir: %v", err)
			}

			assertFileBytes(t, excelPath, tc.wantExcel)
			assertConfigAfter(t, configPath, tc.wantConfig, excelName)
		})
	}
}

type buildExcelDirCase struct {
	name         string
	buildFixture string // testdata/build/<name>; empty = target directory does not exist yet
	wantExcel    []byte
	wantConfig   *ExcelConfig // nil: expect InitConfig-shaped config after build
}

// copyBuildFixtureIntoTarget copies testdata/build/<fixture> into targetDir (read-only fixtures).
func copyBuildFixtureIntoTarget(t *testing.T, fixture, targetDir string) {
	srcDir := filepath.Join("testdata", "build", fixture)
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Fatalf("read build fixture %s: %v", fixture, err)
	}
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	for _, ent := range entries {
		if ent.IsDir() || ent.Name() == ".gitkeep" {
			continue
		}
		src := filepath.Join(srcDir, ent.Name())
		dst := filepath.Join(targetDir, ent.Name())
		if err := fs.CopyFile(src, dst); err != nil {
			t.Fatalf("copy %s -> %s: %v", src, dst, err)
		}
	}
}

func assertConfigAfter(t *testing.T, configPath string, want *ExcelConfig, excelName string) {
	cfg, err := ReadExcelConfig(configPath)
	if err != nil {
		t.Fatalf("ReadExcelConfig: %v", err)
	}
	if want != nil {
		if !reflect.DeepEqual(cfg, want) {
			t.Errorf("config: got %+v want %+v", cfg, want)
		}
		return
	}
	if cfg.ExcelFileName != excelName {
		t.Errorf("ExcelFileName: got %q want %q", cfg.ExcelFileName, excelName)
	}
	if len(cfg.SchemaPreps) != 0 {
		t.Errorf("SchemaPreps: got %v want empty", cfg.SchemaPreps)
	}
}

func mustParseConfigFixture(t *testing.T, path string) *ExcelConfig {
	data := mustReadFile(t, path)
	var cfg ExcelConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("parse config fixture %s: %v", path, err)
	}
	return &cfg
}

func mustReadFile(t *testing.T, path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return b
}

func assertFileBytes(t *testing.T, path string, want []byte) {
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !bytes.Equal(want, got) {
		t.Errorf("%s: got %q want %q", path, got, want)
	}
}
