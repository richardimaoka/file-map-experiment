# Git hooks

## Pre-commit: excel testdata workbooks

Before each commit, `pre-commit` runs `scripts/regenerate-excel-testdata.sh` (`go test -update`, Create tests only). If that leaves an unstaged diff under `excel/testdata`, the commit is **aborted**; review `git diff`, `git add` the workbooks you want, and run `git commit` again.

**One-time setup** (per clone):

```bash
./scripts/install-git-hooks.sh
```

**All testdata** (Excel workbook fixtures + JSON goldens), same `-update` flag:

```bash
./scripts/regenerate-testdata.sh
```

Workbooks only:

```bash
./scripts/regenerate-excel-testdata.sh
```

Plain `go test ./excel` does not rewrite testdata; use `-update` (via the scripts or directly).
