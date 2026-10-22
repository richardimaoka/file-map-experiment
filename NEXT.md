# Next work (memo)

_Last updated: 2026-10-23._

## Suggested next step

**`ValidateExcelAgainstSchema`** — Load persisted `schemas/<name>.json`, read headers (and bounds from prep) on a **instance** workbook, and confirm it matches the schema (column set, grouped headers, row-grouping flags, etc.). Stub exists in `excel/schema.go`; pairs with `ReadRowsFromExcel` for the “same schema, many weekly files” workflow.

## Optional next

- **Rename helpers** — `writeRawHeader*` / `writeRawHeaderWorkbook` → `writeFixture*` (mechanical).
- **Toolchain** — Go **1.26.6+** for **GO-2026-6088** (`encoding/xml` via excelize open/read).

## Defer for now

- Excel **write** path / formula preservation  
- Expanding **`instruction/`** (workbook-level policy, production layout for ~30 books)

## Commands

```bash
export APP_BASE_DIR="$(pwd)"   # repo root; see .env

# Regenerate Excel workbook fixtures + JSON goldens
./scripts/regenerate-testdata.sh

# Public read API smoke
go test ./excel -run 'TestReadSchemaFrom|TestReadRowsFromExcel'

# Git hooks (workbook regen on commit; abort if unstaged excel/testdata diff)
./scripts/install-git-hooks.sh

# Vulnerability scan
govulncheck ./...
```

Fixture layout and per-track commands: [excel/testdata/README.md](excel/testdata/README.md), [excel/testdata/schema/README.md](excel/testdata/schema/README.md).
