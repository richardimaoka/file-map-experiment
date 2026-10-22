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
| [`schema/`](schema/README.md) | Index for `readSchemaInternal` → `TableSchema` | — |
| [`schema_header/`](schema_header/README.md) | Header tracks (`1row/`, `2row/`, `3row/`) | **Complete** (initial pass, 19 cases) |
| [`schema_data/`](schema_data/README.md) | Data band track | **Complete** (initial pass, 9 cases) |
| [`schema_coupling/`](schema_coupling/README.md) | Header + data interaction | **Complete** (initial pass, 3 cases) |
| [`read_api/`](read_api/README.md) | `ReadSchemaFromExcel`, `ReadRowsFromExcel` (resolver layout) | Smoke (`data_all_string`, `ads_report_simple`, `ads_report_multi_row_headers`) |
| [`samples/`](samples/) | Realistic `.xlsx` samples (not golden-driven) | Ad hoc |

Schema design rules, track index, and regenerate commands: [schema/README.md](schema/README.md).
