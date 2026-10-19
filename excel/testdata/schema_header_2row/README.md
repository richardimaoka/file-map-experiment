# Schema fixtures: 2-row headers

Part of the schema fixture plan; see [../schema/README.md](../schema/README.md).

Cases for `readSchemaInternal` (header track, two header rows). Header cells match the same-named case under `testdata/raw_header/`; the data band (`data_start_row`–`data_end_row`) is filled with values like `r3cA`.

Per case:

- `workbook.xlsx` — sheet name matches the directory name
- `schema_prep.json` — `SchemaPrep` (header + data row bounds)
- `table_schema.json` — golden `TableSchema` from `readSchemaInternal`
- `raw_headers.json` — optional copy from `raw_header` for comparison only

Combined workbook: `schema_header_2row_tests.xlsx` (one sheet per case).

`header_name` joins each header row’s label top-to-bottom with ` > ` (merge **anchor** cells only per row). Example: `group1 > headerA`. Columns with no resolvable labels are omitted (e.g. `2-row-bidirectional-merge` column B inside a 2D merge); that case leaves column B empty in the data band.

| Directory | Description |
|-----------|-------------|
| `2-row-simple` | Group row + leaf headers on A–C |
| `2-row-empty-cell` | Empty cell in row 1, column B |
| `2-row-horizontal-merge` | Group merge A–B on row 1 |
| `2-row-vertical-merge` | Vertical merge in column A |
| `2-row-bidirectional-merge` | 2D merge A–B; schema columns A and C only |
| `2-row-horizontal-merge-repeated` | Multiple horizontal group merges, A–G |

Regenerate workbooks:

```bash
go test ./excel -run TestCreateSchemaHeader2RowFixturesInTestdata
```

Regenerate `table_schema.json`:

```bash
go test ./excel -update -run TestReadSchemaInternal2Row
```
