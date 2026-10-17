package excel

import (
	"path/filepath"
	"testing"
)

func TestInitConfig(t *testing.T) {
	cfg := InitConfig("folder_map.xlsx")
	if cfg.ExcelFileName != "folder_map.xlsx" {
		t.Errorf("ExcelFileName: got %q", cfg.ExcelFileName)
	}
	if cfg.SchemaPreps != nil {
		t.Errorf("SchemaPreps: got %v want nil", cfg.SchemaPreps)
	}
}

func TestSchemaPrep_HeaderRange(t *testing.T) {
	prep := SchemaPrep{
		StartColumn:    "A",
		HeaderStartRow: 1,
		EndColumn:      "H",
		HeaderEndRow:   2,
	}
	want := NewCellRange("A", 1, "H", 2)
	if got := prep.HeaderRange(); got != want {
		t.Errorf("HeaderRange: got %v want %v", got, want)
	}
}

func TestExcelConfig_GetSchemaPrep(t *testing.T) {
	cfg := &ExcelConfig{
		SchemaPreps: []SchemaPrep{
			{SchemaName: "sheet1", DataStartRow: 3},
			{SchemaName: "sheet2"},
		},
	}

	prep, err := cfg.GetSchemaPrep("sheet1")
	if err != nil {
		t.Fatal(err)
	}
	if prep.SchemaName != "sheet1" || prep.DataStartRow != 3 {
		t.Fatalf("got %+v", prep)
	}

	_, err = cfg.GetSchemaPrep("missing")
	if err == nil {
		t.Fatal("expected error for missing schema")
	}
}

func TestReadExcelConfig(t *testing.T) {
	path := filepath.Join("testdata", "config", "read_config", "config.json")

	cfg, err := ReadExcelConfig(path)
	if err != nil {
		t.Fatalf("ReadExcelConfig: %v", err)
	}
	if cfg.ExcelFileName != "folder_map.xlsx" {
		t.Errorf("ExcelFileName: got %q", cfg.ExcelFileName)
	}
	if len(cfg.SchemaPreps) != 0 {
		t.Errorf("SchemaPreps: got %d want 0", len(cfg.SchemaPreps))
	}
}

func TestReadExcelConfig_with_schema_preps(t *testing.T) {
	path := filepath.Join("testdata", "config", "read_config_with_schema_preps", "config.json")

	cfg, err := ReadExcelConfig(path)
	if err != nil {
		t.Fatalf("ReadExcelConfig: %v", err)
	}
	if cfg.ExcelFileName != "ads_report_simple.xlsx" {
		t.Fatalf("ExcelFileName: got %q", cfg.ExcelFileName)
	}

	prep, err := cfg.GetSchemaPrep("main")
	if err != nil {
		t.Fatal(err)
	}
	if !prep.IsReady || prep.EndColumn != "H" {
		t.Fatalf("got %+v", prep)
	}
}

func TestReadExcelConfig_file_missing(t *testing.T) {
	_, err := ReadExcelConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReadExcelConfig_invalid_json(t *testing.T) {
	path := filepath.Join("testdata", "config", "read_config_invalid_json", "config.json")

	_, err := ReadExcelConfig(path)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReadExcelConfig_missing_excel_file_name(t *testing.T) {
	path := filepath.Join("testdata", "config", "read_config_missing_excel_file_name", "config.json")

	_, err := ReadExcelConfig(path)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestReadExcelConfig_excel_file_missing(t *testing.T) {
	path := filepath.Join("testdata", "config", "read_config_excel_file_missing", "config.json")

	_, err := ReadExcelConfig(path)
	if err == nil {
		t.Fatal("expected error")
	}
}
