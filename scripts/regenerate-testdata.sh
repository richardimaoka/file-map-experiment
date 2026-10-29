#!/usr/bin/env bash
# Regenerate excel testdata: Excel workbook fixtures and JSON golden files (go test -update).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export APP_BASE_DIR="$ROOT"

RUN='TestCreate.*FixturesInTestdata|TestReadSchemaInternal|TestReadRawDataSchemaData|TestReadRowsSchema|TestReadRawVerticals|TestLoadMergedRanges|TestReadCellValue|TestReadSchemaFromExcel|TestReadRowsFromExcel'

exec go test ./excel -count=1 -update -run "$RUN"
