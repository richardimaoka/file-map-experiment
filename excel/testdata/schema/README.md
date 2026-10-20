# Schema fixtures (`readSchemaInternal` → `TableSchema`)

Golden tests call `readSchemaInternal` with `schema_prep.json` and compare to `table_schema.json` (passing cases only). Each case directory has `workbook.xlsx`, optional `raw_headers.json` for reference, and data cells in the band defined by `data_start_row` / `data_end_row` (values like `r2cA` unless noted).

Row reading later uses `RowData` (`[]map[string]string`): each row is a flat JSON object keyed by **header name**. Schema rules are chosen so every non-empty data column has a stable key.

## Design rules (current)

| Rule | Where enforced |
|------|----------------|
| Any Excel column with **non-empty data** must appear in `TableSchema.Columns` with a **non-empty** `header_name` | `validateDataColumnsHaveRowDataKeys` |
| **Single-row headers:** no **multi-cell** merges in the header row (e.g. B1:C1) | `rejectMultiCellMergesInSingleRowHeader` |
| **Multi-row headers:** `header_name` = row labels joined with ` > `; horizontal group-row merges propagate on non-leaf rows | `resolvedMultiRowHeaderName` |
| **Multi-row headers:** no **multi-column** merge may intersect `header_end_row` (leaf row) | `rejectMultiColumnMergesIntersectingHeaderEndRow` |
| `IsRowGrouped` | Hardcoded `false` until data-band row-merge grouping is defined |
| Merges overlapping a read range must lie **entirely** inside it | `readRawVerticals` (see `testdata/raw_vertical`) |
| Cell semantics (formula, types, whitespace) | `testdata/raw_cells`, not schema fixtures |

**Clean cases first:** supported layouts only; ambiguous Excel (empty header + data, merged 1-row header, data under non-anchor merge column) are separate **error** fixtures or deferred until real files require them.

## Fixture layout (planned trees)

| Directory | Status | Role |
|-----------|--------|------|
| [`schema_header_1row/`](schema_header_1row/README.md) | **In progress** | Header track: one header row |
| [`schema_header_2row/`](schema_header_2row/README.md) | **In progress** | Header track: two header rows |
| `schema_header_3row/` | Planned | Header track: three header rows |
| [`schema_data/`](schema_data/README.md) | **In progress** | Data track: fixed simple header, vary data band |
| `schema_coupling/` (name TBD) | Planned | Few cases where header + data rules interact |

Legacy: [`raw_header/`](../raw_header/README.md) may be removed once schema trees subsume it.

---

## Current: `schema_header_1row` (~5 cases)

See [schema_header_1row/README.md](schema_header_1row/README.md).

| Case | Result |
|------|--------|
| `1-row-simple` | Pass |
| `1-row-line-break` | Pass |
| `1-row-empty-header-valid` | Pass (empty header C, no data in C) |
| `1-row-empty-header` | Error (data in C, no RowData key) |
| `1-row-horizontal-merge` | Error (merged single-row header disallowed) |

---

## Current: `schema_header_2row` (9 cases)

See [schema_header_2row/README.md](schema_header_2row/README.md). Data band: `r{row}c{col}` on rows 3–5 (column span matches each case’s `end_column`).

| Case | Result |
|------|--------|
| `2-row-simple` | Pass |
| `2-row-empty-cell` | Pass |
| `2-row-horizontal-merge` | Pass |
| `2-row-vertical-merge` | Pass |
| `2-row-bidirectional-merge` | Error (multi-column merge on leaf row) |
| `2-row-horizontal-merge-repeated` | Pass |
| `2-row-vertical-merge-repeated` | Pass |
| `2-row-bidirectional-merge-rep` | Error (multi-column merge on leaf row) |
| `2-row-mixed-merge` | Error (multi-column merge on leaf row) |

---

## Planned: `schema_header_3row` (~2–3 cases)

| Suggested case | Focus |
|----------------|--------|
| `3-row-simple` | Three header rows |
| `3-row-*-merge` | One “nasty” merge case (not a full matrix) |

---

## Current: `schema_data` (4 cases; more planned)

See [schema_data/README.md](schema_data/README.md).

| Case | Result |
|------|--------|
| `data-all-string` | Pass |
| `data-empty-column` | Pass (header C, no data in C) |
| `data-empty-cell` | Pass |
| `data-mixed-types` | Pass |

**Same header everywhere** (`1-row-simple`). Vary only the data band.

| Planned case | Focus |
|--------------|--------|
| `data-vertical-merge` | Row grouping in **data** → drives `IsRowGrouped` when implemented |
| `data-horizontal-merge` | Only if schema / row read must see data merges |
| `data-formula` / `data-error` | Only if schema exposes them (else `raw_cells`) |

---

## Planned: coupling (~2–3 cases)

Only where implementation **joins** header and data (not the full cross-product of header × data cases). Examples:

- Grouped multi-row header + specific data merge pattern
- Column in prep range with header omitted but data present (explicit error case, if not already covered)

---

## Deferred / edge bucket (later or error-only)

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

# 2-row header fixtures
go test ./excel -run TestCreateSchemaHeader2RowFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternal2Row
```

Future packages will mirror this pattern (`TestCreateSchemaHeader2RowFixturesInTestdata`, etc.).

## Related read-layer fixtures

| Tree | Tests |
|------|--------|
| `raw_cells` | `readCellValue` |
| `raw_vertical` | `readRawVerticals` |
| `raw_header` | `readRawHeaders` (legacy overlap with schema headers) |
