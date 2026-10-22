package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Schema data fixtures under excel/testdata/schema_data/<case>/.
// Header matches 1-row-simple everywhere; only the data band varies (default rows 2–4, A–C).
//
// Regenerate workbooks and schema_prep.json:
//
//	go test ./excel -run TestCreateSchemaDataFixturesInTestdata
//
// Regenerate table_schema.json:
//
//	go test ./excel -update -run TestReadSchemaInternalData
//
// Regenerate raw_data.json:
//
//	go test ./excel -update -run TestReadRawDataSchemaData
//
// Regenerate rows_data.json:
//
//	go test ./excel -update -run TestReadRowsSchemaData

type schemaDataCase struct {
	dir           string
	dataWrite     func(f *excelize.File, sheet string) error
	wantSchemaErr bool
	// dataBand overrides default rows 2–4 when set (e.g. multi-group vertical merges).
	dataBand *CellRange
}

func schemaDataBand() CellRange {
	return NewCellRange("A", 2, "C", 4)
}

// Twelve data rows (2–13): supports 2×6-row and 4×3-row vertical merge groups.
func schemaDataBand12Rows() CellRange {
	return NewCellRange("A", 2, "C", 13)
}

func (tc schemaDataCase) effectiveDataBand() CellRange {
	if tc.dataBand != nil {
		return *tc.dataBand
	}
	return schemaDataBand()
}

func schemaDataCases() []schemaDataCase {
	band12 := schemaDataBandPtr(schemaDataBand12Rows())
	return []schemaDataCase{
		{dir: "data-all-string", dataWrite: writeSchemaDataAllString},
		{dir: "data-empty-column", dataWrite: writeSchemaDataEmptyColumn},
		{dir: "data-empty-cell", dataWrite: writeSchemaDataEmptyCell},
		{dir: "data-mixed-types", dataWrite: writeSchemaDataMixedTypes},
		{dir: "data-vertical-merge", dataWrite: writeSchemaDataVerticalMerge},
		{dir: "data-vmerge-a6-b3", dataWrite: writeSchemaDataVerticalMergeA6B3, dataBand: band12},
		{dir: "data-vmerge-b6-c3", dataWrite: writeSchemaDataVerticalMergeB6C3, dataBand: band12},
		{dir: "data-horizontal-merge", dataWrite: writeSchemaDataHorizontalMerge, wantSchemaErr: true},
		{dir: "data-formula", dataWrite: writeSchemaDataFormula},
	}
}

func (tc schemaDataCase) writeWorkbook(f *excelize.File, sheet string) error {
	if err := writeRawHeader1RowSimple(f, sheet); err != nil {
		return err
	}
	return tc.dataWrite(f, sheet)
}

func writeSchemaDataAllString(f *excelize.File, sheet string) error {
	return fillSchemaHeader1RowDataBand(f, sheet, schemaDataBand(), nil)
}

func writeSchemaDataEmptyColumn(f *excelize.File, sheet string) error {
	return fillSchemaHeader1RowDataBand(f, sheet, schemaDataBand(), map[string]bool{"C": true})
}

func writeSchemaDataEmptyCell(f *excelize.File, sheet string) error {
	if err := fillSchemaHeader1RowDataBand(f, sheet, schemaDataBand(), nil); err != nil {
		return err
	}
	return clearRawHeaderCell(f, sheet, "B", 3)
}

// Vertical merge A2:A3 in the data band; column A is_row_grouped true.
func writeSchemaDataVerticalMerge(f *excelize.File, sheet string) error {
	if err := fillSchemaHeader1RowDataBand(f, sheet, schemaDataBand(), nil); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "A", 2, "A", 3)
}

func schemaDataBandPtr(b CellRange) *CellRange {
	return &b
}

// Column A: two vertical merges of six rows each; column B: four merges of three rows each.
func writeSchemaDataVerticalMergeA6B3(f *excelize.File, sheet string) error {
	band := schemaDataBand12Rows()
	if err := fillSchemaHeader1RowDataBand(f, sheet, band, nil); err != nil {
		return err
	}
	if err := mergeDataColumnVerticalGroups(f, sheet, "A", band.StartRowNumber(), 6, 2); err != nil {
		return err
	}
	return mergeDataColumnVerticalGroups(f, sheet, "B", band.StartRowNumber(), 3, 4)
}

// Same grouping pattern on columns B (6×2) and C (3×4).
func writeSchemaDataVerticalMergeB6C3(f *excelize.File, sheet string) error {
	band := schemaDataBand12Rows()
	if err := fillSchemaHeader1RowDataBand(f, sheet, band, nil); err != nil {
		return err
	}
	if err := mergeDataColumnVerticalGroups(f, sheet, "B", band.StartRowNumber(), 6, 2); err != nil {
		return err
	}
	return mergeDataColumnVerticalGroups(f, sheet, "C", band.StartRowNumber(), 3, 4)
}

func mergeDataColumnVerticalGroups(f *excelize.File, sheet, col string, startRow, rowsPerGroup, groupCount int) error {
	for g := 0; g < groupCount; g++ {
		top := startRow + g*rowsPerGroup
		bottom := top + rowsPerGroup - 1
		if err := mergeRawHeader(f, sheet, col, top, col, bottom); err != nil {
			return err
		}
	}
	return nil
}

// Horizontal merge B2:C2 in the data band — readSchemaInternal errors (multi-column data merge).
func writeSchemaDataHorizontalMerge(f *excelize.File, sheet string) error {
	if err := fillSchemaHeader1RowDataBand(f, sheet, schemaDataBand(), nil); err != nil {
		return err
	}
	return mergeRawHeader(f, sheet, "B", 2, "C", 2)
}

// Formula in B3; schema columns unchanged (formula semantics live in raw_data / raw_cells).
func writeSchemaDataFormula(f *excelize.File, sheet string) error {
	if err := fillSchemaHeader1RowDataBand(f, sheet, schemaDataBand(), nil); err != nil {
		return err
	}
	cell, err := excelize.JoinCellName("B", 3)
	if err != nil {
		return err
	}
	return f.SetCellFormula(sheet, cell, "B2&\"-calc\"")
}

func writeSchemaDataMixedTypes(f *excelize.File, sheet string) error {
	band := schemaDataBand()
	startRow := band.StartRowNumber()
	endRow := band.EndRowNumber()

	for row := startRow; row <= endRow; row++ {
		if err := setRawHeaderCell(f, sheet, "A", row, fmtSchemaDataString(row, "A")); err != nil {
			return err
		}
		num := float64(row) * 10.5
		if err := setSchemaDataNumberCell(f, sheet, "B", row, num); err != nil {
			return err
		}
		d := time.Date(2024, 3, row-1, 0, 0, 0, 0, time.UTC)
		if err := setSchemaDataDateCell(f, sheet, "C", row, d); err != nil {
			return err
		}
	}
	return nil
}

func fmtSchemaDataString(row int, col string) string {
	return fmt.Sprintf("r%dc%s", row, col)
}

func clearRawHeaderCell(f *excelize.File, sheet, col string, row int) error {
	cell, err := excelize.JoinCellName(col, row)
	if err != nil {
		return err
	}
	return f.SetCellValue(sheet, cell, nil)
}

func setSchemaDataNumberCell(f *excelize.File, sheet, col string, row int, n float64) error {
	cell, err := excelize.JoinCellName(col, row)
	if err != nil {
		return err
	}
	return f.SetCellValue(sheet, cell, n)
}

func setSchemaDataDateCell(f *excelize.File, sheet, col string, row int, d time.Time) error {
	cell, err := excelize.JoinCellName(col, row)
	if err != nil {
		return err
	}
	customFmt := "yyyy-mm-dd"
	styleID, err := f.NewStyle(&excelize.Style{CustomNumFmt: &customFmt})
	if err != nil {
		return err
	}
	if err := f.SetCellStyle(sheet, cell, cell, styleID); err != nil {
		return err
	}
	return f.SetCellValue(sheet, cell, d)
}

func TestCreateSchemaDataFixturesInTestdata(t *testing.T) {
	destBase := filepath.Join("testdata", "schema_data")
	cases := schemaDataCases()

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			destDir := filepath.Join(destBase, tc.dir)
			if err := os.MkdirAll(destDir, 0o755); err != nil {
				t.Fatal(err)
			}

			if err := writeSchemaPrepFile(destDir, fixtureSchemaPrep(tc.dir, 1, tc.effectiveDataBand())); err != nil {
				t.Fatal(err)
			}
			prep, err := readSchemaPrepFixture(filepath.Join(destDir, "schema_prep.json"))
			if err != nil {
				t.Fatal(err)
			}
			workbookPath := filepath.Join(destDir, "workbook.xlsx")
			if err := writeRawHeaderWorkbook(workbookPath, tc.dir, tc.writeWorkbook, prep); err != nil {
				t.Fatal(err)
			}
		})
	}

	sheets := make([]fixtureSheetCase, len(cases))
	for i, tc := range cases {
		prep, err := readSchemaPrepFixture(filepath.Join(destBase, tc.dir, "schema_prep.json"))
		if err != nil {
			t.Fatal(err)
		}
		sheets[i] = fixtureSheetCase{dir: tc.dir, write: tc.writeWorkbook, prep: prep}
	}
	combinedPath := fixtureCombinedWorkbookPath(destBase)
	if err := writeFixtureCombinedWorkbook(combinedPath, sheets); err != nil {
		t.Fatal(err)
	}
}

// Golden raw_data.json: prep data band via readRawVerticals (reference only).
func TestReadRawDataSchemaData(t *testing.T) {
	for _, tc := range schemaDataCases() {
		t.Run(tc.dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "schema_data", tc.dir)
			prepPath := filepath.Join(caseDir, "schema_prep.json")
			excelPath := filepath.Join(caseDir, "workbook.xlsx")
			goldenPath := filepath.Join(caseDir, "raw_data.json")

			var prep SchemaPrep
			if err := jsonfile.Read(prepPath, &prep); err != nil {
				t.Fatalf("read schema_prep: %v", err)
			}

			got, err := readRawVerticals(excelPath, prep.SheetName, prep.DataRange())
			if err != nil {
				t.Fatalf("readRawVerticals: %v", err)
			}

			if updateGoldens() {
				writeRawVerticalsGolden(t, goldenPath, got)
				return
			}

			want, err := readRawVerticalsGolden(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			fatalIfJSONGoldenMismatch(t, got, want)
		})
	}
}

// Golden TableSchema: 1-row-simple header; is_row_grouped reflects vertical merges in the data band.
func TestReadSchemaInternalData(t *testing.T) {
	for _, tc := range schemaDataCases() {
		t.Run(tc.dir, func(t *testing.T) {
			caseDir := filepath.Join("testdata", "schema_data", tc.dir)
			prepPath := filepath.Join(caseDir, "schema_prep.json")
			excelPath := filepath.Join(caseDir, "workbook.xlsx")
			goldenPath := filepath.Join(caseDir, "table_schema.json")

			var prep SchemaPrep
			if err := jsonfile.Read(prepPath, &prep); err != nil {
				t.Fatalf("read schema_prep: %v", err)
			}

			got, err := readSchemaInternal(excelPath, prep)
			if tc.wantSchemaErr {
				if err == nil {
					t.Fatal("readSchemaInternal: expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("readSchemaInternal: %v", err)
			}

			if updateGoldens() {
				writeTableSchemaGolden(t, goldenPath, got)
				return
			}

			want, err := readTableSchemaGolden(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}
			fatalIfJSONGoldenMismatch(t, got, want)
		})
	}
}

