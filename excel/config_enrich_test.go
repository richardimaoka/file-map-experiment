package excel

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonutil"
	"github.com/xuri/excelize/v2"
)

// Regenerate JSON goldens:
//
//	go test ./excel -update -run 'TestEnrichConfig|TestCreateEnrichConfigFixturesInTestdata'
func TestEnrichConfig(t *testing.T) {
	t.Run("ads_report_simple", func(t *testing.T) {
		excelPath := filepath.Join("testdata", "samples", "ads_report_simple.xlsx")
		goldenPath := filepath.Join("testdata", "config", "enrich_ads_report_simple", "config.json")

		cfg := InitConfig("ads_report_simple.xlsx")
		got, err := enrichConfig(cfg, excelPath)
		if err != nil {
			t.Fatalf("enrichConfig: %v", err)
		}
		if updateEnabled() {
			writeConfigGolden(t, goldenPath, got)
		}
		want, err := readConfigGolden(goldenPath)
		if err != nil {
			t.Fatalf("read golden: %v", err)
		}
		fatalIfJSONGoldenMismatch(t, got, want)
	})

	t.Run("multi_sheet", func(t *testing.T) {
		excelPath := filepath.Join("testdata", "enrich_config", "multi_sheet.xlsx")
		goldenPath := filepath.Join("testdata", "config", "enrich_multi_sheet", "config.json")

		cfg := InitConfig("multi_sheet.xlsx")
		got, err := enrichConfig(cfg, excelPath)
		if err != nil {
			t.Fatalf("enrichConfig: %v", err)
		}
		if updateEnabled() {
			writeConfigGolden(t, goldenPath, got)
		}
		want, err := readConfigGolden(goldenPath)
		if err != nil {
			t.Fatalf("read golden: %v", err)
		}
		fatalIfJSONGoldenMismatch(t, got, want)
	})

	t.Run("missing_workbook", func(t *testing.T) {
		cfg := InitConfig("missing.xlsx")
		_, err := enrichConfig(cfg, filepath.Join("testdata", "samples", "no_such_workbook.xlsx"))
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestCreateEnrichConfigFixturesInTestdata(t *testing.T) {
	skipUnlessUpdateEnabled(t)

	destDir := filepath.Join("testdata", "enrich_config")
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	dest := filepath.Join(destDir, "multi_sheet.xlsx")
	f := excelize.NewFile()
	if err := f.SetSheetName("Sheet1", "alpha"); err != nil {
		t.Fatalf("SetSheetName: %v", err)
	}
	if _, err := f.NewSheet("beta"); err != nil {
		t.Fatalf("NewSheet: %v", err)
	}
	if err := f.SaveAs(dest); err != nil {
		t.Fatalf("SaveAs: %v", err)
	}
}

func readConfigGolden(path string) (Config, error) {
	var cfg Config
	if err := jsonutil.Read(path, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func writeConfigGolden(t *testing.T, path string, cfg Config) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := jsonutil.Write(path, cfg); err != nil {
		t.Fatal(err)
	}
}
