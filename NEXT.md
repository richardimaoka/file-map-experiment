# Next work (memo)

_Last updated: 2026-10-22._

## Done (read path)

1. **Row goldens** — `rows_data.json` for `schema_data`, `schema_header_1row`, `schema_header_2row`, `schema_header_3row`, and `schema_coupling` (`TestReadRowsSchema*`).
2. **Public read API** — `ReadSchemaFromJson`, `ReadSchemaFromExcel`, `ReadRowsFromExcel` (resolver + `config.json` + cached `schemas/<name>.json` + workbook). Fixture: `excel/testdata/read_api/data_all_string/`.

## Optional next

- **Smoke** — One or two `excel/testdata/samples/ads_report_*.xlsx` cases with hand-written prep/config + `schemas/*.json`.
- **`ReadSchemaFromExcel` in production flow** — Materialize `schemas/*.json` after prep is ready (no writer yet).

## Defer for now

- `validateExcelAgainstSchema`, `rows_operations` helpers  
- Excel **write** path / formula preservation  
- Expanding **`instruction/`** (workbook-level policy, production layout for ~30 books)

## Commands

```bash
# Row goldens (all tracks)
go test ./excel -update -run 'TestReadRowsSchema'

# Public read API smoke
go test ./excel -run 'TestReadSchemaFrom|TestReadRowsFromExcel'
```

See [excel/testdata/schema/README.md](excel/testdata/schema/README.md) for full fixture regenerate commands.
