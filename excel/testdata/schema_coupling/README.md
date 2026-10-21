# `schema_coupling`

Cases for `readSchemaInternal` where **multi-row header** rules and **data-band** rules interact (not covered by exercising each track alone).

Each case directory has:

- `workbook.xlsx`
- `schema_prep.json`
- `table_schema.json` — golden for **passing** cases only
- `raw_headers.json` — optional reference from `raw_header` when the header layout matches

| Case | Result |
|------|--------|
| `coupling-2row-data-vmerge` | Pass — `2-row-simple` header, vertical merge **A3:A4** → column A `is_row_grouped` |
| `coupling-2row-hdr-data-vmerge` | Pass — `2-row-vertical-merge` header + same data merge on A |
| `coupling-2row-empty-header-data` | Error — column C has no resolved `header_name` but data in C |

## Regenerate

```bash
go test ./excel -run TestCreateSchemaCouplingFixturesInTestdata
go test ./excel -update -run TestReadSchemaInternalCoupling
```

See [../schema/README.md](../schema/README.md) for global schema fixture rules.
