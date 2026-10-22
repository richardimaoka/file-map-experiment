package excel

import (
	"fmt"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// Schema header fixtures (1-row) under excel/testdata/schema_header_1row/<case>/.
// Data band is filled with r{row}c{col} literals below the header layout writers in fixture_header_layouts_test.go.
//
// Regenerate workbooks and schema_prep.json:
//
//	go test ./excel -update -run TestCreateSchemaHeader1RowFixturesInTestdata
//
// Regenerate table_schema.json:
//
//	go test ./excel -update -run TestReadSchemaInternal1Row

func schemaHeader1RowCases() []schemaHeaderTrackCase {
	return []schemaHeaderTrackCase{
		{"1-row-simple", writeRawHeader1RowSimple, NewCellRange("A", 2, "C", 4), false, nil},
		{"1-row-empty-header", writeRawHeader1RowEmptyHeader, NewCellRange("A", 2, "D", 4), true, nil},
		{
			"1-row-empty-header-valid",
			writeRawHeader1RowEmptyHeader,
			NewCellRange("A", 2, "D", 4),
			false,
			map[string]bool{"C": true},
		},
		{"1-row-line-break", writeRawHeader1RowLineBreak, NewCellRange("A", 2, "B", 4), false, nil},
		{"1-row-horizontal-merge", writeRawHeader1RowHorizontalMerge, NewCellRange("A", 2, "D", 4), true, nil},
		{"1-row-duplicate-header", writeRawHeader1RowDuplicateHeader, NewCellRange("A", 2, "C", 4), true, nil},
	}
}

func TestCreateSchemaHeader1RowFixturesInTestdata(t *testing.T) {
	skipUnlessUpdateEnabled(t)
	testCreateSchemaHeaderFixturesInTestdata(t, "schema_header_1row", 1, schemaHeader1RowCases())
}

// Golden TableSchema assumptions (1-row header track): sheet_name, data band bounds, and column
// bounds; header_name from the header row; skip empty names; merged headers only on the merge
// anchor column; is_row_grouped from single-column vertical merges in the data band.
func TestReadSchemaInternal1Row(t *testing.T) {
	testReadSchemaInternalHeaderTrack(t, "schema_header_1row", schemaHeader1RowCases())
}

func readTableSchemaGolden(path string) (TableSchema, error) {
	var schema TableSchema
	if err := jsonfile.Read(path, &schema); err != nil {
		return TableSchema{}, err
	}
	return schema, nil
}

func writeTableSchemaGolden(t *testing.T, path string, schema TableSchema) {
	if err := WriteSchemaToJson(path, schema); err != nil {
		t.Fatal(err)
	}
}

func fillSchemaHeader1RowDataBand(f *excelize.File, sheet string, band CellRange, skipColumns map[string]bool) error {
	startRow := band.StartRowNumber()
	endRow := band.EndRowNumber()
	startCol, err := band.StartColumnNumber()
	if err != nil {
		return err
	}
	endCol, err := band.EndColumnNumber()
	if err != nil {
		return err
	}

	for colNum := startCol; colNum <= endCol; colNum++ {
		colName, err := excelize.ColumnNumberToName(colNum)
		if err != nil {
			return err
		}
		if skipColumns != nil && skipColumns[colName] {
			continue
		}
		for row := startRow; row <= endRow; row++ {
			val := fmt.Sprintf("r%dc%s", row, colName)
			if err := setRawHeaderCell(f, sheet, colName, row, val); err != nil {
				return err
			}
		}
	}
	return nil
}
