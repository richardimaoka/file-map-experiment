# Raw merged-range fixtures

Each case is `excel/testdata/raw_merged/<case>/` with:

- `workbook.xlsx` — single sheet named like the case directory
- `merged_ranges.json` — golden `[]MergedRange` from `loadMergedRanges`

All cases are also combined in `test_cases.xlsx` at this directory (one sheet per case, same names).

Regenerate workbooks:

```bash
go test ./excel -run TestCreateRawMergedFixturesInTestdata
```

Regenerate goldens:

```bash
UPDATE_RAW_MERGED_GOLDEN=1 go test ./excel -run TestLoadMergedRanges
```

| Directory        | Description                          |
|------------------|--------------------------------------|
| `none`           | No merged cells                      |
| `horizontal`     | One row, B1:C1                       |
| `vertical`       | One column, A1:A2                    |
| `block`          | Rectangle B1:C2                      |
| `two-horizontal` | Two separate row merges (A1:B1, D1:E1) |
