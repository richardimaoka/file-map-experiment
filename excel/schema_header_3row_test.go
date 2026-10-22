package excel

import (
	"testing"
)

// Schema header fixtures (3-row) under excel/testdata/schema_header_3row/<case>/.
// Data band uses r{row}c{col} literals below header layouts from fixture_header_layouts_test.go.
//
// Regenerate workbooks and schema_prep.json:
//
//	go test ./excel -update -run TestCreateSchemaHeader3RowFixturesInTestdata
//
// Regenerate table_schema.json:
//
//	go test ./excel -update -run TestReadSchemaInternal3Row

func schemaHeader3RowCases() []schemaHeaderTrackCase {
	return []schemaHeaderTrackCase{
		{"3-row-simple", writeRawHeader3RowSimple, NewCellRange("A", 4, "C", 6), false, nil},
		{
			"3-row-bidirectional-merge",
			writeRawHeader3RowBidirectionalMerge,
			NewCellRange("A", 4, "C", 6),
			true,
			nil,
		},
		{"3-row-mixed-merge", writeRawHeader3RowMixedMerge, NewCellRange("A", 4, "H", 6), false, nil},
	}
}

func TestCreateSchemaHeader3RowFixturesInTestdata(t *testing.T) {
	skipUnlessUpdateEnabled(t)
	testCreateSchemaHeaderFixturesInTestdata(t, "schema_header_3row", 3, schemaHeader3RowCases())
}

// Golden TableSchema (3-row header track): header_name joins three row labels with " > ".
func TestReadSchemaInternal3Row(t *testing.T) {
	testReadSchemaInternalHeaderTrack(t, "schema_header_3row", schemaHeader3RowCases())
}
