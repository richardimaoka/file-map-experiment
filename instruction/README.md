# `instruction`

**Instruction** workbooks (heavy multi-sheet Excel with formulas) add policy metadata on top of generic `excel` schema reading.

Dependency: **`instruction` → `excel` only** (`excel` does not import or document this package).

| Piece | Package | Typical file |
|-------|---------|----------------|
| Bounds | `excel` | `schema_prep.json` |
| Structure | `excel` | `table_schema.json` + workbook |
| Write / validation policy | `instruction` | `sheet_policy.json` |

`LoadSheet(caseDir, excelPath)` returns `Sheet` with `Prep`, `Schema`, and optional `Policy`.

Example fixtures: [`testdata/data-formula/`](testdata/data-formula/) (workbook reused from `excel/testdata/schema_data/data-formula/workbook.xlsx` in tests).

```bash
go test ./instruction/...
```
