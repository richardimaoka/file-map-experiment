# Schema fixtures (`readSchemaInternal` → `TableSchema`)

Golden tests call `readSchemaInternal` with `schema_prep.json` and compare to `table_schema.json` (passing cases only). Each case directory has `workbook.xlsx`, optional `raw_headers.json` for reference, and data cells in the band defined by `data_start_row` / `data_end_row` (values like `r2cA` unless noted).

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
| [`schema_header_1row/`](schema_header_1row/README.md) | **Complete (initial pass)** | Header track: one header row (6 cases) |
| [`schema_header_2row/`](schema_header_2row/README.md) | **Complete (initial pass)** | Header track: two header rows (10 cases) |
| [`schema_data/`](schema_data/README.md) | — | Data track: fixed simple header, vary data band (9 cases) |
| [`schema_header_3row/`](schema_header_3row/README.md) | **Complete (initial pass)** | Header track: three header rows (3 cases) |
| [`schema_coupling/`](schema_coupling/README.md) | **Complete (initial pass)** | Header + data interaction (3 cases) |

Header layouts for `schema_header_*` match same-named cases under [`raw_header/`](../raw_header/README.md) where applicable; schema cases add a filled data band and `table_schema.json` (passing cases only). `raw_header` remains for `readRawHeaders` goldens and may shrink once only schema trees are needed.

---

## `schema_header_1row` (6 cases) — initial pass complete

See [schema_header_1row/README.md](schema_header_1row/README.md).

| Case | Result |
|------|--------|
| `1-row-simple` | Pass |
| `1-row-line-break` | Pass |
| `1-row-empty-header-valid` | Pass (empty header C, no data in C) |
| `1-row-empty-header` | Error (data in C, no RowData key) |
| `1-row-horizontal-merge` | Error (merged single-row header disallowed) |
| `1-row-duplicate-header` | Error (duplicate header name) |

---

## `schema_header_2row` (10 cases) — initial pass complete

See [schema_header_2row/README.md](schema_header_2row/README.md). Data band: `r{row}c{col}` on rows 3–5 (column span matches each case’s `end_column`). Error cases have no `table_schema.json`.

| Case | Result |
|------|--------|
| `2-row-simple` | Pass |
| `2-row-duplicate-header` | Error (duplicate header name) |
| `2-row-empty-cell` | Pass |
| `2-row-horizontal-merge` | Pass |
| `2-row-vertical-merge` | Pass |
| `2-row-bidirectional-merge` | Error (multi-column merge on leaf row) |
| `2-row-horizontal-merge-repeated` | Pass |
| `2-row-vertical-merge-repeated` | Pass |
| `2-row-bidirectional-merge-rep` | Error (multi-column merge on leaf row) |
| `2-row-mixed-merge` | Error (multi-column merge on leaf row) |

---

## `schema_header_3row` (3 cases) — initial pass complete

See [schema_header_3row/README.md](schema_header_3row/README.md). Data band: `r{row}c{col}` on rows 4–6 (column span matches each case’s `end_column`). Error cases have no `table_schema.json`.

| Case | Result |
|------|--------|
| `3-row-simple` | Pass |
| `3-row-bidirectional-merge` | Error (multi-column merge on leaf row) |
| `3-row-mixed-merge` | Pass (horizontal group/subgroup + vertical leaf merges; prep and data A–H) |

---

## `schema_data` (9 cases)

See [schema_data/README.md](schema_data/README.md).

| Case | Result |
|------|--------|
| `data-all-string` | Pass |
| `data-empty-column` | Pass (header C, no data in C) |
| `data-empty-cell` | Pass |
| `data-mixed-types` | Pass |
| `data-vertical-merge` | Pass (column A `is_row_grouped` true) |
| `data-vmerge-a6-b3` | Pass (A 6×2 + B 3×4 row groups, rows 2–13) |
| `data-vmerge-b6-c3` | Pass (B 6×2 + C 3×4 row groups, rows 2–13) |
| `data-horizontal-merge` | Error (multi-column merge in data band) |
| `data-formula` | Pass (`table_schema` unchanged; formula in `raw_data.json`) |

**Same header everywhere** (`1-row-simple`). Vary only the data band. Formula **errors** (`#DIV/0!`, etc.) stay under `raw_cells`.

---

## `schema_coupling` (3 cases) — initial pass complete

See [schema_coupling/README.md](schema_coupling/README.md). Data band: `r{row}c{col}` on rows 3–5 (2-row header prep). Error case has no `table_schema.json`.

| Case | Result |
|------|--------|
| `coupling-2row-data-vmerge` | Pass (`is_row_grouped` on A) |
| `coupling-2row-hdr-data-vmerge` | Pass (vertical header + vertical data merge on A) |
| `coupling-2row-empty-header-data` | Error (data in C, no RowData key) |

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

```bash
# 1-row header fixtures
go test ./excel -run TestCreateSchemaHeader1RowFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternal1Row

# data band fixtures
go test ./excel -run TestCreateSchemaDataFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternalData
go test ./excel -update -run TestReadRawDataSchemaData

# 2-row header fixtures
go test ./excel -run TestCreateSchemaHeader2RowFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternal2Row

# 3-row header fixtures
go test ./excel -run TestCreateSchemaHeader3RowFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternal3Row

# coupling fixtures
go test ./excel -run TestCreateSchemaCouplingFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternalCoupling
```

## Related read-layer fixtures

| Tree | Tests |
|------|--------|
| `raw_cells` | `readCellValue` |
| `raw_vertical` | `readRawVerticals` |
| `raw_header` | `readRawHeaders` (legacy overlap with schema headers) |
