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

func TestSchemaPrep_DataRange(t *testing.T) {
	prep := SchemaPrep{
		StartColumn:  "B",
		DataStartRow: 3,
		EndColumn:    "G",
		DataEndRow:   10,
	}
	want := NewCellRange("B", 3, "G", 10)
	if got := prep.DataRange(); got != want {
		t.Errorf("DataRange: got %v want %v", got, want)
	}
}

func TestSchemaPrep_Validate(t *testing.T) {
	valid := SchemaPrep{
		SchemaName:     "1-row-simple",
		SheetName:      "1-row-simple",
		HeaderStartRow: 1,
		HeaderEndRow:   1,
		StartColumn:    "A",
		EndColumn:      "C",
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid prep: %v", err)
	}

	if err := (SchemaPrep{SheetName: "x", StartColumn: "A", EndColumn: "B", HeaderStartRow: 1, HeaderEndRow: 1}).Validate(); err == nil {
		t.Fatal("expected error when schema_name is missing")
	}
	if err := (SchemaPrep{SchemaName: "x", StartColumn: "A", EndColumn: "B", HeaderStartRow: 1, HeaderEndRow: 1}).Validate(); err == nil {
		t.Fatal("expected error when sheet_name is missing")
	}
	if err := (SchemaPrep{
		SchemaName: "x", SheetName: "x",
		StartColumn: "C", EndColumn: "A",
		HeaderStartRow: 1, HeaderEndRow: 1,
	}).Validate(); err == nil {
		t.Fatal("expected error when start_column after end_column")
	}
	if err := (SchemaPrep{
		SchemaName: "x", SheetName: "x",
		StartColumn: "A", EndColumn: "B",
		HeaderStartRow: 3, HeaderEndRow: 1,
	}).Validate(); err == nil {
		t.Fatal("expected error when header_start_row after header_end_row")
	}
	if err := (SchemaPrep{
		SchemaName: "x", SheetName: "x",
		StartColumn: "A", EndColumn: "B",
		HeaderStartRow: 1, HeaderEndRow: 1,
		DataStartRow: 5, DataEndRow: 3,
	}).Validate(); err == nil {
		t.Fatal("expected error when data_start_row after data_end_row")
	}
	if err := (SchemaPrep{
		SchemaName: "x", SheetName: "x",
		StartColumn: "A", EndColumn: "B",
		HeaderStartRow: 1, HeaderEndRow: 1,
		DataStartRow: 2,
	}).Validate(); err == nil {
		t.Fatal("expected error when data_end_row is unset but data_start_row is set")
	}
}

func TestConfig_GetSchemaPrep(t *testing.T) {
	cfg := Config{
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

// ReadConfig loads config.json and runs Config.Validate.
func TestReadConfig(t *testing.T) {
	cases := []struct {
		name           string
		fixture        string // testdata/config/<fixture>/config.json
		wantErr        bool
		hasSchemaPreps bool // when wantErr is false, whether schema_preps is non-empty
	}{
		{
			// Valid config with Excel file on disk; no schema_preps.
			name:    "read_config_without_schema_preps",
			fixture: "read_config_without_schema_preps",
		},
		{
			// Valid config including schema_preps.
			name:           "read_config_with_schema_preps",
			fixture:        "read_config_with_schema_preps",
			hasSchemaPreps: true,
		},
		{
			// testdata/config/config_file_missing: no config.json in the directory.
			name:    "config_file_missing",
			fixture: "config_file_missing",
			wantErr: true,
		},
		{
			// Config file is not valid JSON.
			name:    "read_config_invalid_json",
			fixture: "read_config_invalid_json",
			wantErr: true,
		},
		{
			// excel_file_name missing or empty in config.
			name:    "read_config_missing_excel_file_name",
			fixture: "read_config_missing_excel_file_name",
			wantErr: true,
		},
		{
			// excel_file_name does not end with .xlsx or .xlsm.
			name:    "read_config_invalid_excel_file_name",
			fixture: "read_config_invalid_excel_file_name",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join("testdata", "config", tc.fixture, "config.json")

			cfg, err := ReadConfig(path)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ReadConfig: %v", err)
			}
			if tc.hasSchemaPreps {
				if len(cfg.SchemaPreps) == 0 {
					t.Errorf("SchemaPreps: got 0 entries want at least one")
				}
			} else if len(cfg.SchemaPreps) != 0 {
				t.Errorf("SchemaPreps: got %d want 0", len(cfg.SchemaPreps))
			}
		})
	}
}
