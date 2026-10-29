#!/usr/bin/env bash
# Regenerate Excel workbook fixtures under excel/testdata (go test -update).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export APP_BASE_DIR="$ROOT"

exec go test ./excel -count=1 -update -run 'TestCreate.*FixturesInTestdata'
