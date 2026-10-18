package excel

import (
	"encoding/json"
	"flag"
	"reflect"
	"strings"
	"testing"
)

var updateGoldensFlag = flag.Bool("update", false, "rewrite JSON golden files")

func updateGoldens() bool {
	return *updateGoldensFlag
}

func fatalIfJSONGoldenMismatch(t *testing.T, got, want interface{}) {
	t.Helper()
	if reflect.DeepEqual(got, want) {
		return
	}
	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		t.Fatalf("marshal got: %v", err)
	}
	wantJSON, err := json.MarshalIndent(want, "", "  ")
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
			b.WriteString("- " + w + "\n")
		}
		if i < len(gotLines) {
			b.WriteString("+ " + g + "\n")
		}
	}
	return b.String()
}
