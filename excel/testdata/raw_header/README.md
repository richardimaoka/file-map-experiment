# Raw header layout fixtures

Each case lives in its own directory with `workbook.xlsx` (single sheet named like the directory, e.g. `1-row-simple`; Excel allows at most 31 characters). All cases are also combined in `test_cases.xlsx` at this directory (one sheet per case, same names). Layouts are defined in `raw_header_test.go`; regenerate fixtures with:

```bash
go test ./excel -run TestCreateRawHeaderFixturesInTestdata
```

Only header regions are filled; data rows below the header are intentionally empty.

Each case directory also has:

- `schema_prep.json` — `SchemaPrep` for the header region (`schema_name` and `sheet_name` match the directory name, `is_ready`: true). `start_column` / `end_column` bound the header area; wholly empty columns inside that range (e.g. column C in `1-row-empty-header`) are omitted from `readRawColumns` output.
- `raw_columns.json` — golden output of `readRawColumns` (JSON array of `RawColumn`).

Regenerate `raw_columns.json` after changing layouts or reader behavior:

```bash
UPDATE_RAW_HEADER_GOLDEN=1 go test ./excel -run TestReadRawColumns
```

| Directory                            | Description                          |
|--------------------------------------|--------------------------------------|
| `1-row-simple`                       | 1-row, simple                        |
| `1-row-empty-header`                 | 1-row, empty header                  |
| `1-row-line-break`                   | 1-row, line break                    |
| `1-row-horizontal-merge`             | 1-row, horizontal merge (B–C)        |
| `2-row-simple`                       | 2-row, simple                        |
| `2-row-empty-cell`                   | 2-row, empty cell                    |
| `2-row-horizontal-merge`             | 2-row, horizontal merge              |
| `2-row-horizontal-merge-repeated`    | 2-row, horizontal merge repeated     |
| `2-row-vertical-merge`               | 2-row, vertical merge                |
| `2-row-vertical-merge-repeated`      | 2-row, vertical merge repeated       |
| `2-row-bidirectional-merge`          | 2-row, bidirectional merge           |
| `2-row-bidirectional-merge-rep`      | 2-row, bidirectional merge repeated |
| `2-row-mixed-merge`                  | 2-row, mixed merge                   |
