package excel

import (
	"path/filepath"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

const (
	readAPICaseRelPath                     = "excel/testdata/read_api/data_all_string"
	readAPIAdsReportSimpleRelPath          = "excel/testdata/read_api/ads_report_simple"
	adsReportSimpleSchemaName              = "ads-report-simple"
	readAPIAdsReportMultiRowHeadersRelPath = "excel/testdata/read_api/ads_report_multi_row_headers"
	adsReportMultiRowHeadersSchemaName     = "ads-report-multi-row-headers"
)

// Public read API: resolver layout with config.json, cached schema JSON, and workbook.
func TestReadSchemaFromJson(t *testing.T) {
	schemaPath := filepath.Join("testdata", "schema_data", "data-all-string", "table_schema.json")
	got, err := ReadSchemaFromJson(schemaPath)
	if err != nil {
		t.Fatalf("ReadSchemaFromJson: %v", err)
	}
	if got.SchemaName != "data-all-string" || len(got.Columns) != 3 {
		t.Fatalf("unexpected schema: name=%q columns=%d", got.SchemaName, len(got.Columns))
	}
}

func TestReadSchemaFromExcel_readAPIFixture(t *testing.T) {
	resolver := NewResolver(readAPICaseRelPath)
	got, err := ReadSchemaFromExcel(resolver, "data-all-string")
	if err != nil {
		t.Fatalf("ReadSchemaFromExcel: %v", err)
	}
	wantPath := resolver.SchemaPath("data-all-string")
	want, err := ReadSchemaFromJson(wantPath)
	if err != nil {
		t.Fatalf("read cached schema: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, got, want)
}

func TestReadRowsFromExcel_readAPIFixture(t *testing.T) {
	resolver := NewResolver(readAPICaseRelPath)
	got, err := ReadRowsFromExcel(resolver, "data-all-string")
	if err != nil {
		t.Fatalf("ReadRowsFromExcel: %v", err)
	}
	goldenPath := filepath.Join("testdata", "schema_data", "data-all-string", "rows_data.json")
	var want []RowData
	if err := jsonutil.Read(goldenPath, &want); err != nil {
		t.Fatalf("read rows golden: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, got, want)
}

// Hand-styled sample workbook (excel/testdata/samples/ads_report_simple.xlsx) in resolver layout.
func TestReadSchemaFromExcel_adsReportSimple(t *testing.T) {
	resolver := NewResolver(readAPIAdsReportSimpleRelPath)
	got, err := ReadSchemaFromExcel(resolver, adsReportSimpleSchemaName)
	if err != nil {
		t.Fatalf("ReadSchemaFromExcel: %v", err)
	}
	schemaPath := resolver.SchemaPath(adsReportSimpleSchemaName)
	if updateEnabled() {
		writeTableSchemaGolden(t, schemaPath, got)
		return
	}
	want, err := ReadSchemaFromJson(schemaPath)
	if err != nil {
		t.Fatalf("read cached schema: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, got, want)
}

func TestReadRowsFromExcel_adsReportSimple(t *testing.T) {
	resolver := NewResolver(readAPIAdsReportSimpleRelPath)
	got, err := ReadRowsFromExcel(resolver, adsReportSimpleSchemaName)
	if err != nil {
		t.Fatalf("ReadRowsFromExcel: %v", err)
	}
	goldenPath := filepath.Join("testdata", "read_api", "ads_report_simple", "rows_data.json")
	if updateEnabled() {
		writeRowsDataGolden(t, goldenPath, got)
		return
	}
	want, err := readRowsDataGolden(goldenPath)
	if err != nil {
		t.Fatalf("read rows golden: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, got, want)
}

// Hand-styled sample with two header rows (group row + leaf labels).
func TestReadSchemaFromExcel_adsReportMultiRowHeaders(t *testing.T) {
	resolver := NewResolver(readAPIAdsReportMultiRowHeadersRelPath)
	got, err := ReadSchemaFromExcel(resolver, adsReportMultiRowHeadersSchemaName)
	if err != nil {
		t.Fatalf("ReadSchemaFromExcel: %v", err)
	}
	schemaPath := resolver.SchemaPath(adsReportMultiRowHeadersSchemaName)
	if updateEnabled() {
		writeTableSchemaGolden(t, schemaPath, got)
		return
	}
	want, err := ReadSchemaFromJson(schemaPath)
	if err != nil {
		t.Fatalf("read cached schema: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, got, want)
}

func TestReadRowsFromExcel_adsReportMultiRowHeaders(t *testing.T) {
	resolver := NewResolver(readAPIAdsReportMultiRowHeadersRelPath)
	got, err := ReadRowsFromExcel(resolver, adsReportMultiRowHeadersSchemaName)
	if err != nil {
		t.Fatalf("ReadRowsFromExcel: %v", err)
	}
	goldenPath := filepath.Join("testdata", "read_api", "ads_report_multi_row_headers", "rows_data.json")
	if updateEnabled() {
		writeRowsDataGolden(t, goldenPath, got)
		return
	}
	want, err := readRowsDataGolden(goldenPath)
	if err != nil {
		t.Fatalf("read rows golden: %v", err)
	}
	fatalIfJSONGoldenMismatch(t, got, want)
}
