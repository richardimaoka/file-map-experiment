package excel

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

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

func writeRawHeaderWorkbook(path, sheetName string, write func(f *excelize.File, sheet string) error, prep *SchemaPrep) error {
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
