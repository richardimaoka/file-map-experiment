package instruction

import (
	"fmt"
	"path/filepath"
	"strings"
)

func formatSourceWorkbookPath(resolver TopLevelResolver, name string) (string, error) {
	_ = resolver
	_ = name
	return "", fmt.Errorf("formatSourceWorkbookPath: not implemented yet")
}

func formatSourceExcelPaths(resolver TopLevelResolver) ([]string, error) {
	_ = resolver
	return nil, fmt.Errorf("formatSourceExcelPaths: not implemented yet")
}

func buildFormatSchema(name string) error {
	top := NewTopLevelResolver()

	sourcePath, err := formatSourceWorkbookPath(top, name)
	if err != nil {
		return fmt.Errorf("buildFormatSchema: %w", err)
	}

	relDirPath := top.FormatSchemaRelDir(name)
	if err := BuildDir(sourcePath, relDirPath); err != nil {
		return fmt.Errorf("buildFormatSchema: %w", err)
	}
	return nil
}

func formatName(excelPath string) (string, error) {
	excelFileName := filepath.Base(excelPath)
	_ = excelFileName
	return "", fmt.Errorf("dirName not yet implemented")
}

// BuildDirs bootstraps every format_schemas/<name>/ tree from workbooks in format_sources/.
// It does not modify format_sources/. Paths come from [TopLevelResolver].
func BuildDirs() error {
	top := NewTopLevelResolver()

	excelPaths, err := formatSourceExcelPaths(top)
	if err != nil {
		return fmt.Errorf("BuildDirs: %w", err)
	}

	for _, excelPath := range excelPaths {
		name, err := formatName(excelPath)
		if err != nil {
			return fmt.Errorf("BuildDirs: %w", err)
		}
		if err := buildFormatSchema(name); err != nil {
			return fmt.Errorf("BuildDirs: %w", err)
		}
	}

	return nil
}

// BuildFormatSchema bootstraps format_schemas/<name>/ from the canonical workbook in
// format_sources/ (same base name as <name>, .xlsx or .xlsm).
func BuildFormatSchema(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("BuildFormatSchema: name is not set")
	}

	if err := buildFormatSchema(name); err != nil {
		return fmt.Errorf("BuildFormatSchema: %w", err)
	}
	return nil
}
