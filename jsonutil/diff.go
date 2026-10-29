package jsonutil

import "strings"

// LineDiff formats a minimal line-oriented diff between two indented JSON strings.
func LineDiff(want, got string) string {
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
