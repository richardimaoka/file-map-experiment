# Schema fixtures: 2-row headers

Part of the schema fixture plan; see [../schema/README.md](../schema/README.md).

Cases for `readSchemaInternal` (header track, two header rows). Header cells match the same-named case under `testdata/raw_header/`; the data band (`data_start_row`–`data_end_row`) is filled with values like `r3cA`.

Per case:

- `workbook.xlsx` — sheet name matches the directory name
- `schema_prep.json` — `SchemaPrep` (header + data row bounds)
- `table_schema.json` — golden `TableSchema` from `readSchemaInternal`
- `raw_headers.json` — optional copy from `raw_header` for comparison only

Combined workbook: `schema_header_2row_tests.xlsx` (one sheet per case).

`header_name` joins each header row’s label top-to-bottom with ` > `. Per row: unmerged cells and merge anchors count; non-anchor cells count only for **horizontal** merges on that row (so `group1` applies to every column under A1:B1). Example: `group1 > headerA`.

**Merge rule:** any merge spanning more than one column that touches the **leaf** header row (`header_end_row`) fails schema (e.g. `2-row-bidirectional-merge`). Group-row-only merges such as `2-row-horizontal-merge` are allowed.

| Directory | Description |
|-----------|-------------|
| `2-row-simple` | Group row + leaf headers on A–C |
| `2-row-duplicate-header` | B and C both `group2 > headerB` — **errors** (duplicate RowData key) |
| `2-row-empty-cell` | Empty cell in row 1, column B |
| `2-row-horizontal-merge` | Group merge A–B on row 1 |
| `2-row-vertical-merge` | Vertical merge in column A |
| `2-row-bidirectional-merge` | 2D merge A–B — **errors** (multi-column merge on leaf row) |
| `2-row-horizontal-merge-repeated` | Multiple horizontal group merges, A–G |
| `2-row-vertical-merge-repeated` | Vertical merges on columns A, C, D |
| `2-row-bidirectional-merge-rep` | Repeated A1:B2-style blocks — **errors** |
| `2-row-mixed-merge` | Horizontal group row + vertical + bidirectional J:K — **errors** |

Regenerate workbooks:

```bash
go test ./excel -run TestCreateSchemaHeader2RowFixturesInTestdata
```

Regenerate `table_schema.json`:

```bash
go test ./excel -update -run TestReadSchemaInternal2Row
```
