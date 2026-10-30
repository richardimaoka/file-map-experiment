package shinoken

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

func TestReadAssetsConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	want := CreateAssetsConfig("params.xlsx", "text.xlsx")
	if err := jsonutil.Write(path, want); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := ReadAssetsConfig(path)
	if err != nil {
		t.Fatalf("ReadAssetsConfig: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestReadParamsConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	want := CreateParamsConfig("params.xlsx", "更新用シート", "更新用シート")
	if err := jsonutil.Write(path, want); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := ReadParamsConfig(path)
	if err != nil {
		t.Fatalf("ReadParamsConfig: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestReadForwardConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	want := CreateForwardConfig("fmt-a", "schema-a", "Sheet1")
	if err := jsonutil.Write(path, want); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := ReadForwardConfig(path)
	if err != nil {
		t.Fatalf("ReadForwardConfig: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestAssetsConfigValidate_requiresFields(t *testing.T) {
	cfg := AssetsConfig{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestForwardConfigValidate_requiresFields(t *testing.T) {
	cfg := ForwardConfig{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}
