package excel

import "testing"

// TestCreate*FixturesInTestdata rewrite excel/testdata workbooks. They run only with go test -update
// (same flag as JSON goldens). See scripts/regenerate-testdata.sh.
func skipUnlessUpdateGoldens(t *testing.T) {
	if !updateGoldens() {
		t.Skip("skipping workbook regeneration (use go test -update)")
	}
}
