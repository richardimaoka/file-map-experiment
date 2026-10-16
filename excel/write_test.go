package excel

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

// Writes excel/testdata/empty.xlsx (default single sheet, no cell data) for read/build tests.
func TestCreateEmptyExcelInTestdata(t *testing.T) {
	dir := filepath.Join("testdata")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "empty.xlsx")
	f := excelize.NewFile()
	defer f.Close()

	if err := f.SaveAs(path); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Fatal("empty.xlsx was not written")
	}
}
