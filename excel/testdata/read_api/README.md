# Public read API fixtures

Each case is a full resolver tree under `APP_BASE_DIR`:

```
<case>/
  workbook.xlsx        # Excel fixture (hand-styled or copied)
  config.json          # schema_preps (bounds + is_ready)
  schemas/
    <schema_name>.json # JSON golden: persisted TableSchema
  rows_data.json       # JSON golden: []RowData (ReadRowsFromExcel)
```

| Case | Source | Tests |
|------|--------|--------|
| `data_all_string` | Generated `schema_data/data-all-string` workbook | `TestReadSchemaFromExcel_readAPIFixture`, `TestReadRowsFromExcel_readAPIFixture` |
| `ads_report_simple` | Hand-styled [samples/ads_report_simple.xlsx](../samples/ads_report_simple.xlsx) | `TestReadSchemaFromExcel_adsReportSimple`, `TestReadRowsFromExcel_adsReportSimple` |
| `ads_report_multi_row_headers` | Hand-styled [samples/ads_report_multi_row_headers.xlsx](../samples/ads_report_multi_row_headers.xlsx) (2 header rows, data rows 3–27) | `TestReadSchemaFromExcel_adsReportMultiRowHeaders`, `TestReadRowsFromExcel_adsReportMultiRowHeaders` |

Regenerate JSON goldens (`schemas/*.json`, `rows_data.json`) for sample cases (`workbook.xlsx` is not rewritten by these tests):

```bash
export APP_BASE_DIR="$(pwd)"
go test ./excel -update -run 'TestReadSchemaFromExcel_adsReportSimple|TestReadRowsFromExcel_adsReportSimple'
go test ./excel -update -run 'TestReadSchemaFromExcel_adsReportMultiRowHeaders|TestReadRowsFromExcel_adsReportMultiRowHeaders'
```
