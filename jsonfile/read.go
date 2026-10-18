package jsonfile

import (
	"encoding/json"
	"fmt"
	"os"
)

func Read(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("jsonfile.Read: %w", err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("jsonfile.Read: %w", err)
	}
	return nil
}
