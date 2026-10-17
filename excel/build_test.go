package excel

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/fs"
)

// BuildExcelDir prepares a working directory from an existing Excel file.
func TestBuildExcelDir(t *testing.T) {
	excelName := "ads_report_simple.xlsx"
	sourceExcel := filepath.Join("testdata", "samples", excelName)
	sourceBytes := mustReadFile(t, sourceExcel)

	existingExcelBytes := mustReadFile(t, filepath.Join("testdata", "build", "existing_dir_with_excel", excelName))
	existingConfigBytes := mustReadFile(t, filepath.Join("testdata", "build", "existing_dir_with_config", "config.json"))

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
			name:            "existing_dir_with_config",
			buildFixture:    "existing_dir_with_config",
			wantExcel:       sourceBytes,
			wantConfigExact: existingConfigBytes,
		},
		// testdata/build/existing_dir_with_excel: Excel file only.
		{
			name:         "existing_dir_with_excel",
			buildFixture: "existing_dir_with_excel",
			wantExcel:    existingExcelBytes,
		},
		// testdata/build/existing_dir_with_excel_and_config: Excel file and config.json.
		{
			name:            "existing_dir_with_excel_and_config",
			buildFixture:    "existing_dir_with_excel_and_config",
			wantExcel:       existingExcelBytes,
			wantConfigExact: existingConfigBytes,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			targetDir := filepath.Join(t.TempDir(), "cases", "client1")
			if tc.buildFixture != "" {
				copyBuildFixtureIntoTarget(t, tc.buildFixture, targetDir)
			}

			excelPath := filepath.Join(targetDir, excelName)
			configPath := filepath.Join(targetDir, "config.json")

			if err := BuildExcelDir(sourceExcel, targetDir); err != nil {
				t.Fatalf("BuildExcelDir: %v", err)
			}

			if info, err := os.Stat(targetDir); err != nil || !info.IsDir() {
				t.Fatalf("target dir: %v", err)
			}

			assertFileBytes(t, excelPath, tc.wantExcel)
			assertConfigAfter(t, configPath, tc.wantConfigExact)
		})
	}
}

type buildExcelDirCase struct {
	name            string
	buildFixture    string // testdata/build/<name>; empty = target directory does not exist yet
	wantExcel       []byte
	wantConfigExact []byte // nil: config.json must exist; non-nil: exact bytes
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

func assertConfigAfter(t *testing.T, configPath string, wantExact []byte) {
	if wantExact != nil {
		assertFileBytes(t, configPath, wantExact)
		return
	}
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("config.json: %v", err)
	}
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
