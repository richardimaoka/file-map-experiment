package jsonutil

import "testing"

func TestLineDiff(t *testing.T) {
	want := "{\n  \"a\": 1\n}\n"
	got := "{\n  \"a\": 2\n}\n"
	diff := LineDiff(want, got)
	if diff != "-   \"a\": 1\n+   \"a\": 2\n" {
		t.Fatalf("unexpected diff:\n%s", diff)
	}
}
