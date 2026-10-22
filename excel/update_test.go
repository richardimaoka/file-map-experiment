package excel

import (
	"flag"
	"reflect"
	"strings"
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
	t.Fatalf("golden JSON mismatch (-want +got):\n%s", jsonLineDiff(string(wantJSON), string(gotJSON)))
}

func jsonLineDiff(want, got string) string {
	wantLines := strings.Split(strings.TrimSuffix(want, "\n"), "\n")
	gotLines := strings.Split(strings.TrimSuffix(got, "\n"), "\n")
	n := len(wantLines)
	if len(gotLines) > n {
		n = len(gotLines)
	}
	var b strings.Builder
	for i := 0; i < n; i++ {
		w := ""
		if i < len(wantLines) {
			w = wantLines[i]
		}
		g := ""
		if i < len(gotLines) {
			g = gotLines[i]
		}
		if w == g {
			continue
		}
		if i < len(wantLines) {
			b.WriteString("- ")
			b.WriteString(w)
			b.WriteByte('\n')
		}
		if i < len(gotLines) {
			b.WriteString("+ ")
			b.WriteString(g)
			b.WriteByte('\n')
		}
	}
	return b.String()
}
