package excel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Standard checked-in fixture: sheet name matches directory name, header rows 1..headerEndRow.
func fixtureSchemaPrep(sheetName string, headerEndRow int, band CellRange) SchemaPrep {
	return SchemaPrep{
		SchemaName:     sheetName,
		SheetName:      sheetName,
		IsReady:        true,
		HeaderStartRow: 1,
		HeaderEndRow:   headerEndRow,
		DataStartRow:   band.StartRowNumber(),
		DataEndRow:     band.EndRowNumber(),
		StartColumn:    band.StartColumnName(),
		EndColumn:      band.EndColumnName(),
	}
}

func writeSchemaPrepFile(destDir string, prep SchemaPrep) error {
	data, err := json.MarshalIndent(prep, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(destDir, "schema_prep.json"), data, 0o644)
}

type fixtureSheetCase struct {
	dir   string
	write func(f *excelize.File, sheet string) error
	prep  *SchemaPrep
}

func readSchemaPrepFixture(path string) (*SchemaPrep, error) {
	var prep SchemaPrep
	if err := jsonfile.Read(path, &prep); err != nil {
		return nil, err
	}
	return &prep, nil
}

// Combined workbook for a testdata tree: <dir_name>_tests.xlsx (e.g. schema_data_tests.xlsx).
func fixtureCombinedWorkbookPath(testdataDir string) string {
	return filepath.Join(testdataDir, filepath.Base(testdataDir)+"_tests.xlsx")
}

func writeFixtureWorkbook(path, sheetName string, write func(f *excelize.File, sheet string) error, prep *SchemaPrep) error {
	f := excelize.NewFile()
	defer f.Close()

	if err := renameDefaultSheet(f, sheetName); err != nil {
		return err
	}
	if err := write(f, sheetName); err != nil {
		return err
	}
	if prep != nil {
		if err := applyFixtureWorkbookStyle(f, sheetName, *prep); err != nil {
			return err
		}
	}
	if err := f.SaveAs(path); err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return fmt.Errorf("%s was not written", path)
	}
	return nil
}

func writeFixtureCombinedWorkbook(path string, cases []fixtureSheetCase) error {
	f := excelize.NewFile()
	defer f.Close()

	for i, tc := range cases {
		sheet := tc.dir
		if i == 0 {
			if err := renameDefaultSheet(f, sheet); err != nil {
				return err
			}
		} else if _, err := f.NewSheet(sheet); err != nil {
			return err
		}
		if err := tc.write(f, sheet); err != nil {
			return err
		}
		if tc.prep != nil {
			if err := applyFixtureWorkbookStyle(f, sheet, *tc.prep); err != nil {
				return err
			}
		}
	}

	if err := f.SaveAs(path); err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return fmt.Errorf("%s was not written", path)
	}
	return nil
}

func renameDefaultSheet(f *excelize.File, name string) error {
	defaultSheet := f.GetSheetName(0)
	if defaultSheet == name {
		return nil
	}
	return f.SetSheetName(defaultSheet, name)
}

// Light header fill and darker grid lines for generated testdata workbooks (not black).
const (
	fixtureHeaderFillColor = "#EDEDED"
	fixtureGridBorderColor = "#9E9E9E"
	fixtureBorderLineStyle = 1
)

func applyFixtureWorkbookStyle(f *excelize.File, sheet string, prep SchemaPrep) error {
	headerStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{fixtureHeaderFillColor},
			Pattern: 1,
		},
		Border: fixtureGrayBorders(),
	})
	if err != nil {
		return fmt.Errorf("applyFixtureWorkbookStyle: header style: %w", err)
	}
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: fixtureGrayBorders(),
	})
	if err != nil {
		return fmt.Errorf("applyFixtureWorkbookStyle: data style: %w", err)
	}

	hTL, hBR, err := cellRangeAxis(prep.HeaderRange())
	if err != nil {
		return err
	}
	if err := f.SetCellStyle(sheet, hTL, hBR, headerStyle); err != nil {
		return fmt.Errorf("applyFixtureWorkbookStyle: header range: %w", err)
	}

	dTL, dBR, err := cellRangeAxis(prep.DataRange())
	if err != nil {
		return err
	}
	if err := f.SetCellStyle(sheet, dTL, dBR, dataStyle); err != nil {
		return fmt.Errorf("applyFixtureWorkbookStyle: data range: %w", err)
	}
	return nil
}

func fixtureGrayBorders() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: fixtureGridBorderColor, Style: fixtureBorderLineStyle},
		{Type: "right", Color: fixtureGridBorderColor, Style: fixtureBorderLineStyle},
		{Type: "top", Color: fixtureGridBorderColor, Style: fixtureBorderLineStyle},
		{Type: "bottom", Color: fixtureGridBorderColor, Style: fixtureBorderLineStyle},
	}
}

func cellRangeAxis(r CellRange) (topLeft, bottomRight string, err error) {
	topLeft, err = excelize.JoinCellName(r.StartColumnName(), r.StartRowNumber())
	if err != nil {
		return "", "", fmt.Errorf("cellRangeAxis: %w", err)
	}
	bottomRight, err = excelize.JoinCellName(r.EndColumnName(), r.EndRowNumber())
	if err != nil {
		return "", "", fmt.Errorf("cellRangeAxis: %w", err)
	}
	return topLeft, bottomRight, nil
}
