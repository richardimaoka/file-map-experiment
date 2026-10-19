# Schema fixtures: data band

Part of the schema fixture plan; see [../schema/README.md](../schema/README.md).

Cases for `readSchemaInternal` (data track). **Header is always `1-row-simple`** (A–C, row 1); `schema_prep` bounds match that layout. Only rows 2–4 in A–C vary.

Per case:

- `workbook.xlsx` — sheet name matches the directory name
- `schema_prep.json` — `SchemaPrep` (header + data row bounds)
- `table_schema.json` — golden `TableSchema` (today identical across cases: header-driven columns)
- `raw_headers.json` — copy from `raw_header/1-row-simple` for reference only

Combined workbook: `schema_data_tests.xlsx` (one sheet per case).

| Directory | Description |
|-----------|-------------|
| `data-all-string` | Full data band with `r{row}c{col}` literals |
| `data-empty-column` | Header C present; no data cells in column C |
| `data-empty-cell` | Hole at B3; other cells filled |
| `data-mixed-types` | Column A strings, B numbers, C dates (display strings when read later) |

Regenerate workbooks:

```bash
go test ./excel -run TestCreateSchemaDataFixturesInTestdata
```

Regenerate `table_schema.json`:

```bash
go test ./excel -update -run TestReadSchemaInternalData
```

Row-level goldens (`RowData`) belong here once `readRowsFromExcelInternal` exists.
