package instruction

import (
	"os"
	"path/filepath"
	"testing"
)

func dataFormulaCaseDir(t *testing.T) string {
	t.Helper()
	return filepath.Join("testdata", "data-formula")
}

func dataFormulaWorkbookPath(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "excel", "testdata", "schema_data", "data-formula", "workbook.xlsx")
}

func TestReadSheetPolicyDataFormulaGolden(t *testing.T) {
	caseDir := dataFormulaCaseDir(t)
	goldenPath := filepath.Join(caseDir, PolicyFileName)

	want := SheetPolicy{
		SchemaName:           "data-formula",
		PreserveFormulaCells: true,
		ColumnTypes:          map[string]string{"B": ColumnTypeFormula},
	}

	if os.Getenv("UPDATE_GOLDENS") == "1" {
		writePolicyGolden(t, goldenPath, want)
		return
	}

	got, err := ReadSheetPolicyFromJSON(goldenPath)
	if err != nil {
		t.Fatalf("ReadSheetPolicyFromJSON: %v", err)
	}
	if got.SchemaName != want.SchemaName || !got.PreserveFormulaCells {
		t.Fatalf("got %+v", got)
	}
	if got.ColumnType("B") != ColumnTypeFormula {
		t.Fatalf("column B: %q", got.ColumnType("B"))
	}
}

func TestLoadSheetDataFormula(t *testing.T) {
	caseDir := dataFormulaCaseDir(t)
	excelPath := dataFormulaWorkbookPath(t)

	sheet, err := LoadSheet(caseDir, excelPath)
	if err != nil {
		t.Fatal(err)
	}
	if sheet.Policy == nil {
		t.Fatal("expected sheet_policy.json")
	}
	if !sheet.Policy.PreserveFormulaCells {
		t.Fatal("expected preserve_formula_cells")
	}
	if sheet.Policy.ColumnType("B") != ColumnTypeFormula {
		t.Fatalf("column B: %q", sheet.Policy.ColumnType("B"))
	}
	if sheet.Schema.Columns[0].HeaderName != "headerA" {
		t.Fatalf("schema: %+v", sheet.Schema.Columns)
	}
}

func writePolicyGolden(t *testing.T, path string, policy SheetPolicy) {
	t.Helper()
	data, err := jsonMarshalGolden(policy)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
