# Raw cell fixtures

Each case is `excel/testdata/raw_cells/<case>/` with:

- `workbook.xlsx` — single sheet named like the case directory
- `merged_ranges.json` — golden `[]MergedRange` from `loadMergedRanges`
- `<cell>.json` — golden `CellValue` from `readCellValue` (e.g. `B1.json`)

All cases are also combined in `test_cases.xlsx` at this directory (one sheet per case, same names).

Regenerate workbooks:

```bash
go test ./excel -run TestCreateRawCellsFixturesInTestdata
```

Regenerate `merged_ranges.json` and cell goldens:

```bash
go test ./excel -update -run 'TestLoadMergedRanges|TestReadCellValue'
```

| Directory        | Description                          |
|------------------|--------------------------------------|
| `none`           | No merged cells                      |
| `horizontal`     | One row, B1:C1                       |
| `vertical`       | One column, A1:A2                    |
| `block`          | Rectangle B1:C2                      |
| `two-horizontal` | Two separate row merges (A1:B1, D1:E1) |
