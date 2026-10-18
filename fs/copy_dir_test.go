package fs

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir_copies_nested_files(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")

	mustMkdir(t, filepath.Join(src, "nested"))
	mustWrite(t, filepath.Join(src, "top.txt"), []byte("top"))
	mustWrite(t, filepath.Join(src, "nested", "inner.txt"), []byte("inner"))

	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir: %v", err)
	}

	assertFileBytes(t, filepath.Join(dst, "top.txt"), []byte("top"))
	assertFileBytes(t, filepath.Join(dst, "nested", "inner.txt"), []byte("inner"))
}

func TestCopyDir_leaves_extra_destination_files(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")

	mustWrite(t, filepath.Join(src, "from_src.txt"), []byte("src"))
	mustMkdir(t, dst)
	mustWrite(t, filepath.Join(dst, "only_in_dst.txt"), []byte("keep"))

	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir: %v", err)
	}

	assertFileBytes(t, filepath.Join(dst, "from_src.txt"), []byte("src"))
	assertFileBytes(t, filepath.Join(dst, "only_in_dst.txt"), []byte("keep"))
}

func TestCopyDir_leaves_extra_nested_destination_files(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")

	mustWrite(t, filepath.Join(src, "shared", "updated.txt"), []byte("from-src"))
	mustMkdir(t, filepath.Join(dst, "shared"))
	mustWrite(t, filepath.Join(dst, "shared", "updated.txt"), []byte("old"))
	mustMkdir(t, filepath.Join(dst, "extra"))
	mustWrite(t, filepath.Join(dst, "extra", "nested", "only_dst.txt"), []byte("keep"))

	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir: %v", err)
	}

	assertFileBytes(t, filepath.Join(dst, "shared", "updated.txt"), []byte("from-src"))
	assertFileBytes(t, filepath.Join(dst, "extra", "nested", "only_dst.txt"), []byte("keep"))
}

func TestCopyDir_overwrites_existing_files(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")

	mustWrite(t, filepath.Join(src, "file.txt"), []byte("new"))
	mustMkdir(t, dst)
	mustWrite(t, filepath.Join(dst, "file.txt"), []byte("old"))

	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir: %v", err)
	}

	assertFileBytes(t, filepath.Join(dst, "file.txt"), []byte("new"))
}

func TestCopyDir_empty_source(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")

	mustMkdir(t, src)

	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir: %v", err)
	}

	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatal("expected destination directory")
	}
}

func TestCopyDir_source_missing(t *testing.T) {
	root := t.TempDir()
	err := CopyDir(filepath.Join(root, "missing"), filepath.Join(root, "dst"))
	if err == nil {
		t.Fatal("expected error when source is missing")
	}
}

func TestCopyDir_source_not_directory(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "file.txt")
	mustWrite(t, src, []byte("x"))

	err := CopyDir(src, filepath.Join(root, "dst"))
	if err == nil {
		t.Fatal("expected error when source is not a directory")
	}
}

func mustMkdir(t *testing.T, path string) {
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func mustWrite(t *testing.T, path string, content []byte) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertFileBytes(t *testing.T, path string, want []byte) {
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !bytes.Equal(want, got) {
		t.Fatalf("%s: got %q want %q", path, got, want)
	}
}
