package excel

import (
	"testing"
)

// Schema header fixtures (2-row) under excel/testdata/schema_header_2row/<case>/.
// Data band uses r{row}c{col} literals below header layouts from fixture_header_layouts_test.go.
//
// Regenerate workbooks and schema_prep.json:
//
//	go test ./excel -update -run TestCreateSchemaHeader2RowFixturesInTestdata
//
// Regenerate table_schema.json:
//
//	go test ./excel -update -run TestReadSchemaInternal2Row

func schemaHeader2RowCases() []schemaHeaderTrackCase {
	return []schemaHeaderTrackCase{
		{"2-row-simple", writeRawHeader2RowSimple, NewCellRange("A", 3, "C", 5), false, nil},
		{
			"2-row-duplicate-header",
			writeRawHeader2RowDuplicateHeader,
			NewCellRange("A", 3, "C", 5),
			true,
			nil,
		},
		{"2-row-empty-cell", writeRawHeader2RowEmptyCell, NewCellRange("A", 3, "C", 5), false, nil},
		{"2-row-horizontal-merge", writeRawHeader2RowHorizontalMerge, NewCellRange("A", 3, "C", 5), false, nil},
		{"2-row-vertical-merge", writeRawHeader2RowVerticalMerge, NewCellRange("A", 3, "C", 5), false, nil},
		{
			"2-row-bidirectional-merge",
			writeRawHeader2RowBidirectionalMerge,
			NewCellRange("A", 3, "C", 5),
			true,
			nil,
		},
		{
			"2-row-horizontal-merge-repeated",
			writeRawHeader2RowHorizontalMergeRepeated,
			NewCellRange("A", 3, "G", 5),
			false,
			nil,
		},
		{
			"2-row-vertical-merge-repeated",
			writeRawHeader2RowVerticalMergeRepeated,
			NewCellRange("A", 3, "E", 5),
			false,
			nil,
		},
		{
			"2-row-bidirectional-merge-rep",
			writeRawHeader2RowBidirectionalMergeRepeated,
			NewCellRange("A", 3, "I", 5),
			true,
			nil,
		},
		{
			"2-row-mixed-merge",
			writeRawHeader2RowMixedMerge,
			NewCellRange("A", 3, "L", 5),
			true,
			nil,
		},
	}
}

func TestCreateSchemaHeader2RowFixturesInTestdata(t *testing.T) {
	skipUnlessUpdateEnabled(t)
	testCreateSchemaHeaderFixturesInTestdata(t, "schema_header_2row", 2, schemaHeader2RowCases())
}

// Golden TableSchema (2-row header track): header_name is row labels joined with " > "
// (resolvedMultiRowHeaderName); is_row_grouped false until data-band grouping exists.
func TestReadSchemaInternal2Row(t *testing.T) {
	testReadSchemaInternalHeaderTrack(t, "schema_header_2row", schemaHeader2RowCases())
}
