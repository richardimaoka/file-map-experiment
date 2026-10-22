package excel

import (
	"flag"
	"reflect"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

var updateFlag = flag.Bool("update", false, "rewrite checked-in testdata: JSON golden files and Excel workbook fixtures")

// updateEnabled reports whether -update was passed (see scripts/regenerate-testdata.sh).
func updateEnabled() bool {
	return *updateFlag
}

// TestCreate*FixturesInTestdata rewrite Excel workbook fixtures only when -update is set.
func skipUnlessUpdateEnabled(t *testing.T) {
	if !updateEnabled() {
		t.Skip("skipping workbook fixture regeneration (use go test -update)")
	}
}

func marshalJSONGoldenIndent(v interface{}) ([]byte, error) {
	return jsonfile.MarshalIndent(v)
}

func fatalIfJSONGoldenMismatch(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		return
	}
	gotJSON, err := marshalJSONGoldenIndent(got)
	if err != nil {
		t.Fatalf("marshal got: %v", err)
	}
	wantJSON, err := marshalJSONGoldenIndent(want)
	if err != nil {
		t.Fatalf("marshal want: %v", err)
	}
	t.Fatalf("golden JSON mismatch (-want +got):\n%s", jsonfile.LineDiff(string(wantJSON), string(gotJSON)))
}
