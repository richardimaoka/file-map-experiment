package instruction

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/richardimaoka/file-map-experiment/excel"
	"github.com/richardimaoka/file-map-experiment/jsonutil"
	"github.com/xuri/excelize/v2"
)

func TestReadConfig(t *testing.T) {
	path := filepath.Join("testdata", "config", "read_config", "config.json")

	var want excel.Config
	if err := jsonutil.Read(path, &want); err != nil {
		t.Fatalf("read golden: %v", err)
	}

	got, err := ReadConfig(path)
	if err != nil {
		t.Fatalf("ReadConfig: %v", err)
	}
	if !reflect.DeepEqual(got, Config{Config: want}) {
		t.Fatalf("config: got %+v want %+v", got, Config{Config: want})
	}
}

func TestInstructionEnrichConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "workbook.xlsx")
	f := excelize.NewFile()
	if err := f.SetSheetName("Sheet1", "表紙"); err != nil {
		t.Fatalf("SetSheetName: %v", err)
	}
	if _, err := f.NewSheet("data_ok"); err != nil {
		t.Fatalf("NewSheet data_ok: %v", err)
	}
	if err := f.SetCellValue("data_ok", "A1", "入力規則"); err != nil {
		t.Fatalf("SetCellValue: %v", err)
	}
	if _, err := f.NewSheet("needs_review"); err != nil {
		t.Fatalf("NewSheet needs_review: %v", err)
	}
	if err := f.SaveAs(path); err != nil {
		t.Fatalf("SaveAs: %v", err)
	}

	cfg := Config{
		Config: excel.Config{
			ExcelFileName: "workbook.xlsx",
			SchemaPreps: []excel.SchemaPrep{
				excelPrep("表紙"),
				excelPrep("data_ok"),
				excelPrep("needs_review"),
			},
		},
	}

	got, err := enrichConfig(cfg, path)
	if err != nil {
		t.Fatalf("enrichConfig: %v", err)
	}

	want := []excel.SchemaPrep{
		excelPrep("data_ok"),
		{
			SchemaName:     "needs_review",
			SheetName:      "needs_review",
			Comment:        prepCommentNeedsReview,
			IsReady:        false,
			HeaderStartRow: 1,
			HeaderEndRow:   1,
			StartColumn:    "A",
			EndColumn:      "A",
		},
	}
	if !reflect.DeepEqual(got.SchemaPreps, want) {
		t.Fatalf("SchemaPreps: got %+v want %+v", got.SchemaPreps, want)
	}
}

func TestInstructionEnrichConfig_dataValidationRetainsWithoutComment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "workbook.xlsx")
	f := excelize.NewFile()
	sheet := f.GetSheetName(f.GetActiveSheetIndex())
	dv := excelize.NewDataValidation(true)
	dv.Sqref = "A1:A10"
	dv.SetDropList([]string{"a", "b"})
	if err := f.AddDataValidation(sheet, dv); err != nil {
		t.Fatalf("AddDataValidation: %v", err)
	}
	if err := f.SaveAs(path); err != nil {
		t.Fatalf("SaveAs: %v", err)
	}

	cfg := Config{
		Config: excel.Config{
			ExcelFileName: "workbook.xlsx",
			SchemaPreps:   []excel.SchemaPrep{excelPrep(sheet)},
		},
	}

	got, err := enrichConfig(cfg, path)
	if err != nil {
		t.Fatalf("enrichConfig: %v", err)
	}
	if len(got.SchemaPreps) != 1 {
		t.Fatalf("SchemaPreps len: got %d want 1", len(got.SchemaPreps))
	}
	if got.SchemaPreps[0].Comment != "" {
		t.Errorf("Comment: got %q want empty", got.SchemaPreps[0].Comment)
	}
}

func TestInstructionEnrichConfig_emptySheetName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "workbook.xlsx")
	f := excelize.NewFile()
	if err := f.SaveAs(path); err != nil {
		t.Fatalf("SaveAs: %v", err)
	}

	cfg := Config{
		Config: excel.Config{
			ExcelFileName: "workbook.xlsx",
			SchemaPreps: []excel.SchemaPrep{{
				SchemaName: "only-schema",
				SheetName:  "",
			}},
		},
	}
	_, err := enrichConfig(cfg, path)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestInstructionEnrichConfig_missingSheet(t *testing.T) {
	path := filepath.Join(t.TempDir(), "workbook.xlsx")
	f := excelize.NewFile()
	if err := f.SaveAs(path); err != nil {
		t.Fatalf("SaveAs: %v", err)
	}

	cfg := Config{
		Config: excel.Config{
			ExcelFileName: "workbook.xlsx",
			SchemaPreps:   []excel.SchemaPrep{excelPrep("missing")},
		},
	}
	_, err := enrichConfig(cfg, path)
	if err == nil {
		t.Fatal("expected error")
	}
}

func excelPrep(sheetName string) excel.SchemaPrep {
	return excel.SchemaPrep{
		SchemaName:     sheetName,
		SheetName:      sheetName,
		Comment:        "",
		IsReady:        false,
		HeaderStartRow: 1,
		HeaderEndRow:   1,
		StartColumn:    "A",
		EndColumn:      "A",
	}
}
