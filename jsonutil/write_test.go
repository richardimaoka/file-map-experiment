package jsonutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWrite_roundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")

	in := struct {
		Label string `json:"label"`
		N     int    `json:"n"`
	}{Label: "a>b", N: 1}

	if err := Write(path, in); err != nil {
		t.Fatalf("Write: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(data) != "{\n  \"label\": \"a>b\",\n  \"n\": 1\n}\n" {
		t.Fatalf("unexpected file contents:\n%s", data)
	}

	var out struct {
		Label string `json:"label"`
		N     int    `json:"n"`
	}
	if err := Read(path, &out); err != nil {
		t.Fatalf("Read: %v", err)
	}
	if out != in {
		t.Fatalf("got %+v want %+v", out, in)
	}
}
