package jsonutil

import (
	"path/filepath"
	"testing"
)

// Read unmarshals a JSON object from a file on disk.
func TestRead(t *testing.T) {
	path := filepath.Join("testdata", "read_ok", "data.json")

	var got struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	if err := Read(path, &got); err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got.Name != "alpha" || got.Count != 2 {
		t.Fatalf("got %+v", got)
	}
}

// Missing file path returns an error.
func TestRead_file_missing(t *testing.T) {
	path := filepath.Join("testdata", "read_missing", "data.json")

	var v map[string]string
	if err := Read(path, &v); err == nil {
		t.Fatal("expected error")
	}
}

// Invalid JSON in the file returns an error.
func TestRead_invalid_json(t *testing.T) {
	path := filepath.Join("testdata", "read_invalid_json", "data.json")

	var v map[string]string
	if err := Read(path, &v); err == nil {
		t.Fatal("expected error")
	}
}
