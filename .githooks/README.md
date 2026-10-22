# Git hooks

## Pre-commit: excel testdata workbooks

Before each commit, `pre-commit` runs `scripts/regenerate-excel-testdata.sh` (`go test -update`, Create tests only) and stages any changed `excel/testdata/**/*.xlsx`.

**One-time setup** (per clone):

```bash
./scripts/install-git-hooks.sh
```

**All fixtures** (`.xlsx` + JSON goldens), same `-update` flag:

```bash
./scripts/regenerate-testdata.sh
```

Workbooks only:

```bash
./scripts/regenerate-excel-testdata.sh
```

Plain `go test ./excel` does not rewrite testdata; use `-update` (via the scripts or directly).
