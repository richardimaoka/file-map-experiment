# Schema fixtures (`readSchemaInternal` → `TableSchema`)

Golden tests call `readSchemaInternal` with `schema_prep.json` and compare to `table_schema.json` (passing cases only). Each case directory has `workbook.xlsx`, optional `raw_headers.json` for reference, and data cells in the band defined by `data_start_row` / `data_end_row` (values like `r2cA` unless noted).

Row reading later uses `RowData` (`[]map[string]string`): each row is a flat JSON object keyed by **header name**. Schema rules are chosen so every non-empty data column has a stable key.

## Design rules (current)

| Rule | Where enforced |
|------|----------------|
| Any Excel column with **non-empty data** must appear in `TableSchema.Columns` with a **non-empty** `header_name` | `validateDataColumnsHaveRowDataKeys` |
| **Single-row headers:** no **multi-cell** merges in the header row (e.g. B1:C1) | `rejectMultiCellMergesInSingleRowHeader` |
| **Multi-row headers:** merge logic not implemented yet (`columnSchemasFromMultiRowHeader`) | — |
| `IsRowGrouped` | Hardcoded `false` until data-band row-merge grouping is defined |
| Merges overlapping a read range must lie **entirely** inside it | `readRawVerticals` (see `testdata/raw_vertical`) |
| Cell semantics (formula, types, whitespace) | `testdata/raw_cells`, not schema fixtures |

**Clean cases first:** supported layouts only; ambiguous Excel (empty header + data, merged 1-row header, data under non-anchor merge column) are separate **error** fixtures or deferred until real files require them.

## Fixture layout (planned trees)

| Directory | Status | Role |
|-----------|--------|------|
| [`schema_header_1row/`](schema_header_1row/README.md) | **In progress** | Header track: one header row |
| `schema_header_2row/` | Planned | Header track: two header rows |
| `schema_header_3row/` | Planned | Header track: three header rows |
| `schema_data/` | Planned | Data track: fixed simple header, vary data band |
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

## Planned: `schema_header_2row` (~5–6 cases)

Boring data band (empty or uniform `r{row}c{col}`). Header-focused.

| Suggested case | Focus |
|----------------|--------|
| `2-row-simple` | Two header rows, plain names |
| `2-row-empty-cell` | Empty cell in header grid |
| `2-row-horizontal-merge` | Group row merge across columns |
| `2-row-vertical-merge` | Vertical merge in header (one column) |
| `2-row-bidirectional-merge` | One 2D header merge |
| `2-row-horizontal-merge-repeated` | Optional: many columns / repeated pattern |

Reuse header layouts from `testdata/raw_header` where names match.

---

## Planned: `schema_header_3row` (~2–3 cases)

| Suggested case | Focus |
|----------------|--------|
| `3-row-simple` | Three header rows |
| `3-row-*-merge` | One “nasty” merge case (not a full matrix) |

---

## Planned: `schema_data` (~5–8 cases)

**Same header everywhere** (e.g. copy `1-row-simple` header + `schema_prep` bounds). Vary only the data band.

| Suggested case | Focus |
|----------------|--------|
| `data-all-string` | Baseline string values |
| `data-mixed-types` | Number vs text per column (display strings) |
| `data-empty-cell` | Sparse values inside the grid |
| `data-empty-column` | One column blank for all data rows |
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
# 1-row header fixtures (today)
go test ./excel -run TestCreateSchemaHeader1RowFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternal1Row
```

Future packages will mirror this pattern (`TestCreateSchemaHeader2RowFixturesInTestdata`, etc.).

## Related read-layer fixtures

| Tree | Tests |
|------|--------|
| `raw_cells` | `readCellValue` |
| `raw_vertical` | `readRawVerticals` |
| `raw_header` | `readRawHeaders` (legacy overlap with schema headers) |
