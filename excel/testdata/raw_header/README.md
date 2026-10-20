# Raw header layout fixtures

Layouts shared with [schema_header_1row](../schema_header_1row/README.md) and [schema_header_2row](../schema_header_2row/README.md) (same directory names; schema cases add a data band and `readSchemaInternal` goldens). See [../schema/README.md](../schema/README.md) for the full schema matrix.

Each case lives in its own directory with `workbook.xlsx` (single sheet named like the directory, e.g. `1-row-simple`; Excel allows at most 31 characters). All cases are also combined in `raw_header_tests.xlsx` at this directory (one sheet per case, same names). Layouts are defined in `raw_header_test.go`; regenerate fixtures with:

```bash
go test ./excel -run TestCreateRawHeaderFixturesInTestdata
```

Only header regions are filled; data rows below the header are intentionally empty.

Each case directory also has:

- `schema_prep.json` — full `SchemaPrep` (`schema_name` and `sheet_name` match the directory name, `is_ready`: true). Header rows use `header_start_row` / `header_end_row`; `start_column` / `end_column` bound both header and data. Data rows use `data_start_row` / `data_end_row` (first row after the header through a short empty band below it). Every column in range appears in `readRawHeaders` output (including empty cells).
- `raw_headers.json` — golden output of `readRawHeaders` (JSON array of `RawHeader`).

Regenerate `raw_headers.json` after changing layouts or reader behavior:

```bash
go test ./excel -update -run TestReadRawHeaders
```

| Directory                            | Description                          |
|--------------------------------------|--------------------------------------|
| `1-row-simple`                       | 1-row, simple                        |
| `1-row-empty-header`                 | 1-row, empty header                  |
| `1-row-line-break`                   | 1-row, line break                    |
| `1-row-horizontal-merge`             | 1-row, horizontal merge (B–C)        |
| `1-row-duplicate-header`             | 1-row, duplicate header text (B, C)  |
| `2-row-simple`                       | 2-row, simple                        |
| `2-row-duplicate-header`             | 2-row, duplicate resolved header (B, C) |
| `2-row-empty-cell`                   | 2-row, empty cell                    |
| `2-row-horizontal-merge`             | 2-row, horizontal merge              |
| `2-row-horizontal-merge-repeated`    | 2-row, horizontal merge repeated     |
| `2-row-vertical-merge`               | 2-row, vertical merge                |
| `2-row-vertical-merge-repeated`      | 2-row, vertical merge repeated       |
| `2-row-bidirectional-merge`          | 2-row, bidirectional merge           |
| `2-row-bidirectional-merge-rep`      | 2-row, bidirectional merge repeated |
| `2-row-mixed-merge`                  | 2-row, mixed merge                   |
