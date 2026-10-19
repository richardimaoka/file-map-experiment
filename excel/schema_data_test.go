package excel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Schema data fixtures under excel/testdata/schema_data/<case>/.
// Header matches 1-row-simple everywhere; only the data band (rows 2–4, A–C) varies.
//
// Regenerate workbooks and schema_prep.json:
//
//	go test ./excel -run TestCreateSchemaDataFixturesInTestdata
//
// Regenerate table_schema.json:
//
//	go test ./excel -update -run TestReadSchemaInternalData

type schemaDataCase struct {
	dir       string
	dataWrite func(f *excelize.File, sheet string) error
}

func schemaDataBand() CellRange {
	return NewCellRange("A", 2, "C", 4)
}

func schemaDataCases() []schemaDataCase {
	return []schemaDataCase{
		{"data-all-string", writeSchemaDataAllString},
		{"data-empty-column", writeSchemaDataEmptyColumn},
		{"data-empty-cell", writeSchemaDataEmptyCell},
		{"data-mixed-types", writeSchemaDataMixedTypes},
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
	headerRef := filepath.Join("testdata", "raw_header", "1-row-simple", "raw_headers.json")
	cases := schemaDataCases()

	for _, tc := range cases {
		t.Run(tc.dir, func(t *testing.T) {
			destDir := filepath.Join(destBase, tc.dir)
			if err := os.MkdirAll(destDir, 0o755); err != nil {
				t.Fatal(err)
			}

			workbookPath := filepath.Join(destDir, "workbook.xlsx")
			if err := writeRawHeaderWorkbook(workbookPath, tc.dir, tc.writeWorkbook); err != nil {
				t.Fatal(err)
			}

			if err := copySchemaHeaderFixtureFile(headerRef, filepath.Join(destDir, "raw_headers.json")); err != nil {
				t.Fatal(err)
			}
			if err := writeSchemaPrepForDataCase(destDir, tc); err != nil {
				t.Fatal(err)
			}
		})
	}

	sheets := make([]fixtureSheetCase, len(cases))
	for i, tc := range cases {
		sheets[i] = fixtureSheetCase{dir: tc.dir, write: tc.writeWorkbook}
	}
	combinedPath := fixtureCombinedWorkbookPath(destBase)
	if err := writeFixtureCombinedWorkbook(combinedPath, sheets); err != nil {
		t.Fatal(err)
	}
}

// Golden TableSchema: same 1-row-simple header; data band does not change column list today.
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

func writeSchemaPrepForDataCase(destDir string, tc schemaDataCase) error {
	prep := SchemaPrep{
		SchemaName:     tc.dir,
		SheetName:      tc.dir,
		IsReady:        true,
		HeaderStartRow: 1,
		HeaderEndRow:   1,
		DataStartRow:   2,
		DataEndRow:     4,
		StartColumn:    "A",
		EndColumn:      "C",
	}
	data, err := json.MarshalIndent(prep, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(destDir, "schema_prep.json"), data, 0o644)
}
