# Schema fixtures (`readSchemaInternal` → `TableSchema`)

Golden tests call `readSchemaInternal` with `schema_prep.json` and compare to `table_schema.json` (passing cases only). Each case directory has `workbook.xlsx` and data cells in the band defined by `data_start_row` / `data_end_row` (values like `r2cA` unless noted).

Row reading later uses `RowData` (`[]map[string]string`): each row is a flat JSON object keyed by **header name**. Schema rules are chosen so every non-empty data column has a stable key.

## Design rules (current)

| Rule | Where enforced |
|------|----------------|
| Any Excel column with **non-empty data** must appear in `TableSchema.Columns` with a **non-empty** `header_name` | `validateDataColumnsHaveRowDataKeys` |
| **Non-empty `header_name` values must be unique** across columns (RowData keys) | `rejectDuplicateHeaderNames` |
| **Single-row headers:** no **multi-cell** merges in the header row (e.g. B1:C1) | `rejectMultiCellMergesInSingleRowHeader` |
| **Multi-row headers:** `header_name` = row labels joined with ` > `; horizontal group-row merges propagate on non-leaf rows | `resolvedMultiRowHeaderName` |
| **Multi-row headers:** no **multi-column** merge may intersect `header_end_row` (leaf row) | `rejectMultiColumnMergesIntersectingHeaderEndRow` |
| **Data band:** no **multi-column** merge may intersect `data_start_row`–`data_end_row` (finest column per row) | `rejectMultiColumnMergesIntersectingDataBand` |
| `IsRowGrouped` | `true` when the data band has a **single-column** merge spanning multiple rows in that column | `dataColumnIsRowGrouped` |
| Merges overlapping a read range must lie **entirely** inside it | `readRawVerticals` (see `testdata/raw_vertical`) |
| Cell semantics (formula, types, whitespace) | `testdata/raw_cells`, not schema fixtures |

**Clean cases first:** supported layouts only; ambiguous Excel (empty header + data, merged 1-row header, data under non-anchor merge column) are separate **error** fixtures or deferred until real files require them.

## Fixture layout

| Status | Meaning |
|--------|---------|
| **Complete (initial pass)** | All cases listed for that track are implemented, with tests and goldens; no further cases planned for this pass (deferred/edge items stay in the bucket below). |
| **In progress** | Some cases exist; more are listed under that track or here as planned. |
| **Planned** | Tree not started yet. |

| Directory | Status | Role |
|-----------|--------|------|
| [`schema_header/`](schema_header/README.md) | **Complete (initial pass)** | Header tracks: `1row/` (6), `2row/` (10), `3row/` (3) |
| [`schema_data/`](schema_data/README.md) | **Complete (initial pass)** | Data track: fixed simple header, vary data band (9 cases) |
| [`schema_coupling/`](schema_coupling/README.md) | **Complete (initial pass)** | Header + data interaction (3 cases) |

All tracks above are implemented for the initial `readSchemaInternal` fixture pass; further cases belong in the [deferred / edge](#deferred--edge-bucket-later-or-error-only) bucket unless real files require otherwise.

Header layouts are defined in Go (`excel/fixture_header_layouts_test.go`) and written by `TestCreateSchemaHeaderFixturesInTestdata`; cases add a filled data band and `table_schema.json` (**passing cases only**; error fixtures omit it).

## Tracks (case tables)

Case names, layouts, and pass/error notes:

- [schema_header/README.md](../schema_header/README.md)
- [schema_data/README.md](../schema_data/README.md)
- [schema_coupling/README.md](../schema_coupling/README.md)

---

## Deferred / edge bucket (later or error-only)

For the initial development phase, when in doubt, prefer explicit schema creation failure rather than guesswork to make it up.
Not part of the initial “clean” pass; add when real `.xlsx` files require them or to document errors:

- Single-row header merge with data **only** on anchor column B (pass) vs **only** on C (fail) — policy may stay “no 1-row header merges” instead
- Header skip / merge collapse with partial data columns
- `bound-*` style merges clipped by range (already errors at `readRawVerticals`)
- Hidden rows/columns, formulas, whitespace in headers (usually `raw_cells` / `raw_vertical`)

---

## Test commands

All fixtures (Excel workbooks + JSON goldens): `./scripts/regenerate-testdata.sh` from repo root (`APP_BASE_DIR` set in script).

Per-track examples (`-update` rewrites workbooks via `TestCreate*` and JSON via read tests):

```bash
# header fixtures (1row / 2row / 3row under testdata/schema_header/)
go test ./excel -update -run TestCreateSchemaHeaderFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternalHeader
go test ./excel -update -run TestReadRowsSchemaHeader

# data band fixtures
go test ./excel -update -run TestCreateSchemaDataFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternalData
go test ./excel -update -run TestReadRawDataSchemaData
go test ./excel -update -run TestReadRowsSchemaData

# row goldens (all schema tracks)
go test ./excel -update -run 'TestReadRowsSchema'

# coupling fixtures
go test ./excel -update -run TestCreateSchemaCouplingFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternalCoupling
```

## Related read-layer fixtures

| Tree | Tests |
|------|--------|
| `raw_cells` | `readCellValue` |
| `raw_vertical` | `readRawVerticals` |
