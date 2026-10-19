# Raw vertical range fixtures (`readRawRange`)

Each case lives in `excel/testdata/raw_vertical/<case>/` with:

- `workbook.xlsx` — single sheet named like the case directory (≤31 characters)
- `read_range.json` — `sheet_name` and `range` (`CellRange`) passed to `readRawRange`
- `raw_verticals.json` — golden `[]RawVertical`

Case directory names use a **prefix** that matches one **combined workbook** (several sheets, one per case in that group). Combined files sit in this directory:

| Combined workbook | Prefix | Cases |
|-------------------|--------|--------|
| `combined_extent.xlsx` | `ext-` | `ext-1-row-simple`, `ext-3-row-simple`, `ext-1x1-single-cell` |
| `combined_sparse.xlsx` | `sparse-` | `sparse-3-row-empty-cell`, `sparse-3-row-empty-column`, `sparse-1-row-gap-column` |
| `combined_merge.xlsx` | `merge-` | `merge-1-row-horizontal`, `merge-2-row-vertical`, `merge-2-row-block`, `merge-2-row-bidirectional`, `merge-1-row-two-horizontal` |
| `combined_merge_boundary.xlsx` | `bound-` | `bound-merge-above`, `bound-merge-below`, `bound-merge-left`, `bound-merge-right` |
| `combined_hidden.xlsx` | `hidden-` | `hidden-3-row-hidden-row`, `hidden-3-row-hidden-column` |
| `combined_column.xlsx` | `col-` | `col-1-row-wide`, `col-1-row-column-aa` |
| `combined_range.xlsx` | `range-` | `range-3-row-all-empty`, `range-3-row-below-used` |
| `combined_stress.xlsx` | `stress-` | `stress-2-row-many-merge` |

Regenerate per-case workbooks, `read_range.json`, and combined workbooks:

```bash
go test ./excel -run TestCreateRawVerticalFixturesInTestdata
```

Regenerate `raw_verticals.json`:

```bash
go test ./excel -update -run TestReadRawRange
```

Cell value semantics (formula, types, whitespace, etc.) are covered under `testdata/raw_cells`.
