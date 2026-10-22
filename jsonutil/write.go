package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// MarshalIndent encodes v as indented JSON (2 spaces). Strings are not HTML-escaped.
func MarshalIndent(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("jsonutil.MarshalIndent: %w", err)
	}
	return buf.Bytes(), nil
}

// Write marshals v as indented JSON and writes it to path.
func Write(path string, v any) error {
	data, err := MarshalIndent(v)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("jsonutil.Write: %w", err)
	}
	return nil
}
