# Schema fixtures: header row tracks

Schema header track for `readSchemaInternal`; global rules: [../README.md#schema-fixtures-readschemainternal](../README.md#schema-fixtures-readschemainternal). **Initial pass complete** (19 cases). Data band track: [../schema_data/README.md](../schema_data/README.md).

Layout:

```
schema_header/
  1row/<case>/   — single header row (6 cases)
  2row/<case>/   — two header rows (10 cases)
  3row/<case>/   — three header rows (3 cases)
  schema_header_1row_tests.xlsx   — combined workbook (one sheet per 1row case)
  schema_header_2row_tests.xlsx
  schema_header_3row_tests.xlsx
```

Per case:

- `workbook.xlsx` — sheet name matches the case directory name
- `schema_prep.json` — `SchemaPrep` (header + data row bounds)
- `table_schema.json` — JSON golden from `readSchemaInternal` (passing cases only)
- `rows_data.json` — JSON golden row read (`TestReadRowsSchemaHeader`)

Header cells are built in `excel/fixture_header_layouts_test.go`; the data band is filled with values like `r2cA`.

## 1-row cases (`1row/`)

| Directory | Description |
|-----------|-------------|
| `1-row-simple` | Three plain header cells |
| `1-row-empty-header` | Empty header in column C with data in C — **errors** |
| `1-row-empty-header-valid` | Same header; column C has no data — **passes** |
| `1-row-line-break` | Multiline header text in B |
| `1-row-horizontal-merge` | Horizontal merge B–C in header — **errors** |
| `1-row-duplicate-header` | Duplicate header names — **errors** |

Data band: rows 2–4, all in-range columns.

## 2-row cases (`2row/`)

Data band: rows 3–5 (`r3cA`, …).

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

## 3-row cases (`3row/`)

Data band: rows 4–6 (`r4cA`, …).

`header_name` joins each header row’s label top-to-bottom with ` > `. Example: `group1 > subgroup1 > headerA`.

**Merge rule:** any merge spanning more than one column that touches the **leaf** header row (`header_end_row` = 3) fails schema (e.g. `3-row-bidirectional-merge`). Group/subgroup rows may use horizontal merges.

| Directory | Description |
|-----------|-------------|
| `3-row-simple` | Three header rows on A–C |
| `3-row-bidirectional-merge` | 2D merge A1:B3 — **errors** (multi-column merge on leaf row) |
| `3-row-mixed-merge` | Group merges A–C, partial D–E tree, vertical leaf merges F–H; data band A–H rows 4–6 — **passes** |

## Regenerate

```bash
go test ./excel -update -run TestCreateSchemaHeaderFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternalHeader
go test ./excel -update -run TestReadRowsSchemaHeader
```
