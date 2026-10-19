# Schema fixtures: 1-row headers

Cases for `readSchemaInternal` (header track, single header row). Each directory mirrors the same-named case under `testdata/raw_header/` until we consolidate or remove `raw_header`.

Per case:

- `workbook.xlsx` — sheet name matches the directory name
- `schema_prep.json` — `SchemaPrep` (header + data row bounds)
- `table_schema.json` — golden `TableSchema` from `readSchemaInternal`
- `raw_headers.json` — optional copy from `raw_header` for comparison only (not tested here)

Combined workbook: `test_cases.xlsx` (one sheet per case).

| Directory | Description |
|-----------|-------------|
| `1-row-simple` | Three plain header cells |
| `1-row-empty-header` | Empty header in column C; range A–D |
| `1-row-line-break` | Multiline header text in B |
| `1-row-horizontal-merge` | Horizontal merge B–C |

Regenerate from `raw_header` writers and JSON copies:

```bash
go test ./excel -run TestCreateSchemaHeader1RowFixturesInTestdata
```

Regenerate `table_schema.json`:

```bash
go test ./excel -update -run TestReadSchemaInternal1Row
```
