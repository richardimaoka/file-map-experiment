# Schema fixtures: 1-row headers

Part of the schema fixture plan; see [../schema/README.md](../schema/README.md). **Initial pass complete** (6 cases). Data-band work continues under [../schema_data/README.md](../schema_data/README.md).

Cases for `readSchemaInternal` (header track, single header row). Header cells match the same-named case under `testdata/raw_header/`; the data band (`data_start_row`–`data_end_row`, every column in `start_column`–`end_column`) is filled with values like `r2cA` for use when header and data verticals are read together.

Per case:

- `workbook.xlsx` — sheet name matches the directory name
- `schema_prep.json` — `SchemaPrep` (header + data row bounds)
- `table_schema.json` — golden `TableSchema` from `readSchemaInternal` (passing cases only)
- `raw_headers.json` — optional copy from `raw_header` for comparison only (not tested here)

Combined workbook: `schema_header_1row_tests.xlsx` (one sheet per case).

| Directory | Description |
|-----------|-------------|
| `1-row-simple` | Three plain header cells |
| `1-row-empty-header` | Empty header in column C with data in C — `readSchemaInternal` **errors** (no RowData key) |
| `1-row-empty-header-valid` | Same header; column C has no data — **passes** (A, B, D only in schema) |
| `1-row-line-break` | Multiline header text in B |
| `1-row-horizontal-merge` | Horizontal merge B–C in header — **errors** (merged headers disallowed for single-row) |
| `1-row-duplicate-header` | Columns B and C both `headerB` — **errors** (duplicate RowData key) |

Data band: rows 2–4, all in-range columns.

Regenerate workbooks (header + data) and sync JSON from `raw_header`:

```bash
go test ./excel -run TestCreateSchemaHeader1RowFixturesInTestdata
```

Regenerate `table_schema.json`:

```bash
go test ./excel -update -run TestReadSchemaInternal1Row
```
