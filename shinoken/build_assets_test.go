package shinoken

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/fs"
	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

func TestBuildAssetsDir(t *testing.T) {
	const rel = "data/shinoken/assets/case-a"
	sample := filepath.Join("..", "excel", "testdata", "samples", "ads_report_simple.xlsx")

	appBase := t.TempDir()
	t.Setenv("APP_BASE_DIR", appBase)
	resolver := NewAssetsResolver(rel)

	cfg := CreateAssetsConfig("ads_report_simple.xlsx", "ads_report_simple.xlsx")
	if err := os.MkdirAll(resolver.FullDirPath(), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := jsonutil.Write(resolver.ConfigPath(), cfg); err != nil {
		t.Fatalf("Write config: %v", err)
	}
	if err := fs.CopyFile(sample, resolver.ParamsExcelPath(cfg)); err != nil {
		t.Fatalf("Copy params excel: %v", err)
	}
	if err := fs.CopyFile(sample, resolver.TextExcelPath(cfg)); err != nil {
		t.Fatalf("Copy text excel: %v", err)
	}

	if err := BuildAssetsDir(resolver); err != nil {
		t.Fatalf("BuildAssetsDir: %v", err)
	}
}
