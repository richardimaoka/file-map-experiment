# Raw header layout fixtures

Each case lives in its own directory with `workbook.xlsx` (single sheet named `Sheet1`). Layouts are defined in `raw_header_test.go`; regenerate fixtures with:

```bash
go test ./excel -run TestCreateRawHeaderFixturesInTestdata
```

Only header regions are filled; data rows below the header are intentionally empty.

| Directory                            | Description                          |
|--------------------------------------|--------------------------------------|
| `1-row-simple`                       | 1-row, simple                        |
| `1-row-empty-header`                 | 1-row, empty header                  |
| `1-row-line-break`                   | 1-row, line break                    |
| `2-row-simple`                       | 2-row, simple                        |
| `2-row-empty-cell`                   | 2-row, empty cell                    |
| `2-row-horizontal-merge`             | 2-row, horizontal merge              |
| `2-row-horizontal-merge-repeated`    | 2-row, horizontal merge repeated     |
| `2-row-vertical-merge`               | 2-row, vertical merge                |
| `2-row-vertical-merge-repeated`      | 2-row, vertical merge repeated       |
| `2-row-bidirectional-merge`          | 2-row, bidirectional merge           |
| `2-row-bidirectional-merge-repeated` | 2-row, bidirectional merge repeated |
| `2-row-mixed-merge`                  | 2-row, mixed merge                   |
