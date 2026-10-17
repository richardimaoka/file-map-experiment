package fs

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")
	content := []byte("copy-me")

	if err := os.WriteFile(src, content, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(content, got) {
		t.Fatalf("dst content: got %q want %q", got, content)
	}
}

func TestCopyFile_overwrites_existing_destination(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.txt")
	dst := filepath.Join(dir, "dst.txt")

	if err := os.WriteFile(dst, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(src, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("CopyFile: %v", err)
	}

	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal([]byte("new"), got) {
		t.Fatalf("dst content: got %q want %q", got, "new")
	}
}

func TestCopyFile_source_missing(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "dst.txt")

	err := CopyFile(filepath.Join(dir, "missing.txt"), dst)
	if err == nil {
		t.Fatal("expected error when source is missing")
	}
	if _, statErr := os.Stat(dst); statErr == nil {
		t.Fatal("destination should not be created when copy fails")
	}
}
