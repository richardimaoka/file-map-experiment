# Schema fixtures: 3-row headers

Part of the schema fixture plan; see [../schema/README.md](../schema/README.md). **Initial pass** (3 cases).

Cases for `readSchemaInternal` (header track, three header rows). Header cells match the same-named case under `testdata/raw_header/`; the data band (`data_start_row`–`data_end_row`) is filled with values like `r4cA`.

Per case:

- `workbook.xlsx` — sheet name matches the directory name
- `schema_prep.json` — `SchemaPrep` (header + data row bounds)
- `table_schema.json` — golden `TableSchema` from `readSchemaInternal` (**passing cases only**; error cases omit this file)
- `raw_headers.json` — optional copy from `raw_header` for comparison only

Combined workbook: `schema_header_3row_tests.xlsx` (one sheet per case).

`header_name` joins each header row’s label top-to-bottom with ` > `. Example: `group1 > subgroup1 > headerA`.

**Merge rule:** any merge spanning more than one column that touches the **leaf** header row (`header_end_row` = 3) fails schema (e.g. `3-row-bidirectional-merge`). Group/subgroup rows may use horizontal merges (not covered in this initial pass).

| Directory | Description |
|-----------|-------------|
| `3-row-simple` | Three header rows on A–C |
| `3-row-bidirectional-merge` | 2D merge A1:B3 — **errors** (multi-column merge on leaf row) |
| `3-row-mixed-merge` | Group merges A–C, partial D–E tree, vertical leaf merges F–H; data band A–H rows 4–6 — **passes** |

Regenerate workbooks:

```bash
go test ./excel -run TestCreateSchemaHeader3RowFixturesInTestdata
```

Regenerate `table_schema.json`:

```bash
go test ./excel -update -run TestReadSchemaInternal3Row
```
