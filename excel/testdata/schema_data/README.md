# Schema fixtures: data band

Part of the schema fixture plan; see [../schema/README.md](../schema/README.md). **Initial pass complete** (9 cases).

Cases for `readSchemaInternal` (data track). **Header is always `1-row-simple`** (A–C, row 1). Most cases use data rows 2–4; nested vertical-merge cases (`data-vmerge-*`) use rows 2–13.

Per case:

- `workbook.xlsx` — sheet name matches the directory name
- `schema_prep.json` — `SchemaPrep` (header + data row bounds)
- `table_schema.json` — golden `TableSchema` (columns from header; `is_row_grouped` from single-column vertical merges in the data band)
- `raw_data.json` — golden data-band columns from prep bounds (`TestReadRawDataSchemaData`; reference only)
- `rows_data.json` — golden `[]RowData` from `readRowsFromExcelInternal` (passing cases only)

Combined workbook: `schema_data_tests.xlsx` (one sheet per case).

| Directory | Description |
|-----------|-------------|
| `data-all-string` | Full data band with `r{row}c{col}` literals |
| `data-empty-column` | Header C present; no data cells in column C |
| `data-empty-cell` | Hole at B3; other cells filled |
| `data-mixed-types` | Column A strings, B numbers, C dates (display strings when read later) |
| `data-vertical-merge` | Vertical merge A2:A3 — column A `is_row_grouped` true |
| `data-vmerge-a6-b3` | Rows 2–13: A two×6-row merges, B four×3-row merges — A and B `is_row_grouped` |
| `data-vmerge-b6-c3` | Same band: B two×6-row merges, C four×3-row merges — B and C `is_row_grouped` |
| `data-horizontal-merge` | Horizontal merge B2:C2 — **errors** (multi-column merge in data band) |
| `data-formula` | Formula in B3 (`B2&"-calc"`); `table_schema` unchanged, see `raw_data.json` |

Regenerate workbooks:

```bash
go test ./excel -run TestCreateSchemaDataFixturesInTestdata
```

Regenerate `table_schema.json`:

```bash
go test ./excel -update -run TestReadSchemaInternalData
```

Regenerate `raw_data.json`:

```bash
go test ./excel -update -run TestReadRawDataSchemaData
```

Regenerate `rows_data.json`:

```bash
go test ./excel -update -run TestReadRowsSchemaData
```
