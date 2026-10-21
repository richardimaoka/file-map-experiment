# Excel package test fixtures

Golden and workbook fixtures for `excel` live under this directory. Each tree has its own `README.md` with case tables and regenerate commands.

Requires `APP_BASE_DIR` (see module `.env` / `excel/testmain_test.go`).

## Fixture trees

| Directory | Tests | Status |
|-----------|--------|--------|
| [`raw_cells/`](raw_cells/README.md) | `readCellValue`, merged ranges | Stable |
| [`raw_vertical/`](raw_vertical/README.md) | `readRawVerticals` | Stable |
| [`raw_header/`](raw_header/README.md) | `readRawHeaders` | Legacy overlap with schema header tracks |
| [`schema/`](schema/README.md) | Index for `readSchemaInternal` → `TableSchema` | — |
| [`schema_header_1row/`](schema_header_1row/README.md) | Header track (1 row) | **Complete** (initial pass, 6 cases) |
| [`schema_header_2row/`](schema_header_2row/README.md) | Header track (2 rows) | **Complete** (initial pass, 10 cases) |
| [`schema_header_3row/`](schema_header_3row/README.md) | Header track (3 rows) | **Complete** (initial pass, 3 cases) |
| [`schema_data/`](schema_data/README.md) | Data band track | **Complete** (initial pass, 9 cases) |
| [`schema_coupling/`](schema_coupling/README.md) | Header + data interaction | **Complete** (initial pass, 3 cases) |
| [`samples/`](samples/) | Realistic `.xlsx` samples (not golden-driven) | Ad hoc |

Schema design rules and the full case matrix: [schema/README.md](schema/README.md).
