#!/usr/bin/env bash
# Point this repo at tracked hooks under .githooks/ (run once per clone).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

chmod +x "$ROOT/.githooks/pre-commit"
chmod +x "$ROOT/scripts/regenerate-excel-testdata.sh"
chmod +x "$ROOT/scripts/regenerate-testdata.sh"

git config core.hooksPath .githooks
echo "Installed git hooks: core.hooksPath=.githooks"
