package excel

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/fs"
)

// BuildDir prepares a working directory from an existing Excel file.
func TestBuildDir(t *testing.T) {
	excelName := "ads_report_simple.xlsx"
	sourceExcel := filepath.Join("testdata", "samples", excelName)
	sourceBytes := mustReadFile(t, sourceExcel)

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
		// testdata/build/existing_dir_with_config: config.json only; schema_preps resynced from workbook.
		{
			name:         "existing_dir_with_config",
			buildFixture: "existing_dir_with_config",
			wantExcel:    sourceBytes,
		},
		// testdata/build/existing_dir_with_excel: Excel file only (overwritten from source).
		{
			name:         "existing_dir_with_excel",
			buildFixture: "existing_dir_with_excel",
			wantExcel:    sourceBytes,
		},
		// testdata/build/existing_dir_with_excel_and_config: Excel overwritten; schema_preps resynced.
		{
			name:         "existing_dir_with_excel_and_config",
			buildFixture: "existing_dir_with_excel_and_config",
			wantExcel:    sourceBytes,
		},
	}

	caseRelPath := "cases/client1"

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			appBase := t.TempDir()
			t.Setenv("APP_BASE_DIR", appBase)
			targetDir := filepath.Join(appBase, caseRelPath)
			if tc.buildFixture != "" {
				srcDir := filepath.Join("testdata", "build", tc.buildFixture)
				if err := fs.CopyDir(srcDir, targetDir); err != nil {
					t.Fatalf("CopyDir: %v", err)
				}
			}

			excelPath := filepath.Join(targetDir, excelName)
			configPath := filepath.Join(targetDir, "config.json")

			if err := BuildDir(sourceExcel, caseRelPath); err != nil {
				t.Fatalf("BuildExcelDir: %v", err)
			}

			if info, err := os.Stat(targetDir); err != nil || !info.IsDir() {
				t.Fatalf("target dir: %v", err)
			}

			assertFileBytes(t, excelPath, tc.wantExcel)
			assertConfigAfter(t, configPath)
			assertSchemasDirAfter(t, filepath.Join(targetDir, "schemas"))
		})
	}
}

func assertSchemasDirAfter(t *testing.T, schemasDir string) {
	info, err := os.Stat(schemasDir)
	if err != nil || !info.IsDir() {
		t.Fatalf("schemas dir: %v", err)
	}
}

type buildExcelDirCase struct {
	name         string
	buildFixture string // testdata/build/<name>; empty = target directory does not exist yet
	wantExcel    []byte
}

func assertConfigAfter(t *testing.T, configPath string) {
	cfg, err := ReadConfig(configPath)
	if err != nil {
		t.Fatalf("ReadConfig: %v", err)
	}
	wantCfg, err := ReadConfig(filepath.Join("testdata", "config", "enrich_ads_report_simple", "config.json"))
	if err != nil {
		t.Fatalf("ReadConfig golden: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, cfg, wantCfg)
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
