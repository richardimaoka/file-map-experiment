# Excel package test fixtures

JSON **golden** files (expected read output) and Excel **workbook fixtures** (inputs) for `excel` live under this directory. Each tree has its own `README.md` with case tables and regenerate commands.

Requires `APP_BASE_DIR` (see module `.env` / `excel/testmain_test.go`).

## Regenerating fixtures

Use **`go test -update`** to rewrite checked-in testdata: Excel workbook fixtures (`TestCreate*`) and JSON goldens (read-layer tests). Same flag, different artifact types.

```bash
export APP_BASE_DIR="$(pwd)"   # repo root
./scripts/regenerate-testdata.sh          # .xlsx + JSON
./scripts/regenerate-excel-testdata.sh  # .xlsx only
```

With git hooks installed (`./scripts/install-git-hooks.sh`), **pre-commit** regenerates workbooks and **fails the commit** if `excel/testdata` has unstaged changes afterward (stage `.xlsx` and commit again). Plain `go test ./excel` does not rewrite testdata.

## Fixture trees

| Directory | Tests | Status |
|-----------|--------|--------|
| [`raw_cells/`](raw_cells/README.md) | `readCellValue`, merged ranges | Stable |
| [`raw_vertical/`](raw_vertical/README.md) | `readRawVerticals` | Stable |
| [`schema_header/`](schema_header/README.md) | `readSchemaInternal` — header row tracks (`1row/`, `2row/`, `3row/`) | **Complete** (initial pass, 19 cases) |
| [`schema_data/`](schema_data/README.md) | `readSchemaInternal` — data band track | **Complete** (initial pass, 9 cases) |
| [`schema_coupling/`](schema_coupling/README.md) | `readSchemaInternal` — header + data interaction | **Complete** (initial pass, 3 cases) |
| [`read_api/`](read_api/README.md) | `ReadSchemaFromExcel`, `ReadRowsFromExcel` (resolver layout) | Smoke (`data_all_string`, `ads_report_simple`, `ads_report_multi_row_headers`) |
| [`samples/`](samples/) | Realistic `.xlsx` samples (not golden-driven) | Ad hoc |

Lower-level cell and vertical reads (`raw_cells`, `raw_vertical`) support schema fixture design; see [schema fixtures](#schema-fixtures-readschemainternal) below.

## Schema fixtures (`readSchemaInternal`)

Golden tests call `readSchemaInternal` with `schema_prep.json` and compare to `table_schema.json` (passing cases only). Each case directory has `workbook.xlsx` and data cells in the band defined by `data_start_row` / `data_end_row` (values like `r2cA` unless noted).

Row reading uses `RowData` (`[]map[string]string`): each row is a flat JSON object keyed by **header name**. Schema rules are chosen so every non-empty data column has a stable key.

**Clean cases first:** supported layouts only; ambiguous Excel (empty header + data, merged 1-row header, data under non-anchor merge column) are separate **error** fixtures or deferred until real files require them.

Header layouts are defined in Go (`excel/fixture_header_layouts_test.go`) and written by `TestCreateSchemaHeaderFixturesInTestdata`; cases add a filled data band and `table_schema.json` (**passing cases only**; error fixtures omit it).

Case tables: [schema_header/](schema_header/README.md), [schema_data/](schema_data/README.md), [schema_coupling/](schema_coupling/README.md).

### Design rules

| Rule | Where enforced |
|------|----------------|
| Any Excel column with **non-empty data** must appear in `TableSchema.Columns` with a **non-empty** `header_name` | `validateDataColumnsHaveRowDataKeys` |
| **Non-empty `header_name` values must be unique** across columns (RowData keys) | `rejectDuplicateHeaderNames` |
| **Single-row headers:** no **multi-cell** merges in the header row (e.g. B1:C1) | `rejectMultiCellMergesInSingleRowHeader` |
| **Multi-row headers:** `header_name` = row labels joined with ` > `; horizontal group-row merges propagate on non-leaf rows | `resolvedMultiRowHeaderName` |
| **Multi-row headers:** no **multi-column** merge may intersect `header_end_row` (leaf row) | `rejectMultiColumnMergesIntersectingHeaderEndRow` |
| **Data band:** no **multi-column** merge may intersect `data_start_row`–`data_end_row` (finest column per row) | `rejectMultiColumnMergesIntersectingDataBand` |
| `IsRowGrouped` | `true` when the data band has a **single-column** merge spanning multiple rows in that column | `dataColumnIsRowGrouped` |
| Merges overlapping a read range must lie **entirely** inside it | `readRawVerticals` (see `raw_vertical/`) |
| Cell semantics (formula, types, whitespace) | `raw_cells/`, not schema fixtures |

### Deferred schema cases (later or error-only)

When in doubt during initial development, prefer explicit schema failure rather than guesswork. Add when real `.xlsx` files require them:

- Single-row header merge with data **only** on anchor column B (pass) vs **only** on C (fail) — policy may stay “no 1-row header merges” instead
- Header skip / merge collapse with partial data columns
- `bound-*` style merges clipped by range (already errors at `readRawVerticals`)
- Hidden rows/columns, formulas, whitespace in headers (usually `raw_cells` / `raw_vertical`)

### Schema track regenerate commands

```bash
# header (testdata/schema_header/{1row,2row,3row}/)
go test ./excel -update -run TestCreateSchemaHeaderFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternalHeader
go test ./excel -update -run TestReadRowsSchemaHeader

# data band
go test ./excel -update -run TestCreateSchemaDataFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternalData
go test ./excel -update -run TestReadRawDataSchemaData
go test ./excel -update -run TestReadRowsSchemaData

# row goldens (all schema tracks)
go test ./excel -update -run 'TestReadRowsSchema'

# coupling
go test ./excel -update -run TestCreateSchemaCouplingFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternalCoupling
```
