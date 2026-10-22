# Next work (memo)

_Last updated: 2026-10-23._

## Done (read path)

1. **Row goldens** — `rows_data.json` for `schema_data`, `schema_header_*`, and `schema_coupling` (`TestReadRowsSchema*`).
2. **Public read API** — `ReadSchemaFromJson`, `ReadSchemaFromExcel`, `WriteSchemaToJson`, `WriteSchemaFromExcel`, `ReadRowsFromExcel` (resolver + `config.json` + persisted `schemas/<name>.json` + workbook). Fixture: `excel/testdata/read_api/data_all_string/`.

## Done (fixture architecture & excel package)

1. **Schema prep in Go** — `fixtureSchemaPrep` + `TestCreateSchemaHeader*FixturesInTestdata`; no copies from `testdata/raw_header`.
2. **Removed `testdata/raw_header/`** — `schema_header_*` trees are canonical.
3. **`RawDataColumn` removed** — data band and row read use `RawVertical` only.
4. **Row read cleanup** — single `readRowsFromExcelInternal`; validate row span; comments on `byCol` pivot.
5. **File layout** — `vertical.go`, `merge.go`, `schema.go`, `rows_read.go`, `rows_write.go`, `rows_ops.go`; fixture helpers in `fixture_*_test.go` (styling only in tests).
6. **excelize** — `v2.11.1-0.20261003002531-6258dcebc4e2` (2026 security advisories including GO-2026-6695).

## Done (read API smoke)

1. **`read_api/ads_report_simple`** — Hand-styled [samples/ads_report_simple.xlsx](excel/testdata/samples/ads_report_simple.xlsx) in resolver layout (`config.json`, `schemas/ads-report-simple.json`, `rows_data.json`). Tests: `TestReadSchemaFromExcel_adsReportSimple`, `TestReadRowsFromExcel_adsReportSimple`.
2. **`read_api/ads_report_multi_row_headers`** — [samples/ads_report_multi_row_headers.xlsx](excel/testdata/samples/ads_report_multi_row_headers.xlsx) (2 header rows, grouped column titles). Tests: `TestReadSchemaFromExcel_adsReportMultiRowHeaders`, `TestReadRowsFromExcel_adsReportMultiRowHeaders`. See [excel/testdata/read_api/README.md](excel/testdata/read_api/README.md).

**Scope:** `read_api` is resolver-layout smoke (`config.json` + persisted schema + workbook), not a second copy of `schema_header_*` / `schema_data` / `schema_coupling`. Edge cases (row grouping, coupling, etc.) stay on internal `ReadTableSchema` / row-read tests.

3. **Persisted schema** — `WriteSchemaFromExcel` (reference workbook + `is_ready` prep → `schemas/<name>.json`); `WriteSchemaToJson` for file-level serde. `ReadRowsFromExcel` loads that JSON; `validateExcelAgainstSchema` (future) checks instance workbooks against it.

## Optional next
- **Fixture test dedupe** — shared `testCreateSchemaHeaderFixturesInTestdata` / `testReadSchemaInternalHeaderTrack` (keep three `schema_header_*` testdata trees and case files).
- **Rename helpers** — `writeRawHeader*` / `writeRawHeaderWorkbook` → `writeFixture*` (mechanical).
- **Toolchain** — Go **1.26.6+** for **GO-2026-6088** (`encoding/xml` via excelize open/read).

## Defer for now

- `validateExcelAgainstSchema`, `rows_ops` stubs  
- Excel **write** path / formula preservation  
- Expanding **`instruction/`** (workbook-level policy, production layout for ~30 books)

## Commands

```bash
export APP_BASE_DIR="$(pwd)"   # or module root per .env

# Row goldens (all tracks)
go test ./excel -update -run 'TestReadRowsSchema'

# Public read API smoke
go test ./excel -run 'TestReadSchemaFrom|TestReadRowsFromExcel'

# Vulnerability scan (excelize + import paths)
govulncheck ./...
```

See [excel/testdata/schema/README.md](excel/testdata/schema/README.md) for full fixture regenerate commands.
