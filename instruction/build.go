package instruction

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/richardimaoka/file-map-experiment/excel"
)

func writeDefaultData(resolver Resolver, schemaName string) error {
	_ = resolver
	_ = schemaName
	return fmt.Errorf("writeDefaultData: not implemented yet")
}

func ensureDefaultDataExists(resolver Resolver) error {
	_ = resolver
	return fmt.Errorf("ensureDefaultDataExists: not implemented yet")
}

// BuildDir bootstraps an instruction case at relDirPath:
// runs [excel.BuildDir], then writes defaults/<schemaName>.json sidecars.
// It does not create schemas/ entries beyond what [excel.BuildDir] already skips. For the layout, see [Resolver].
func BuildDir(sourceExcelPath string, relDirPath string) error {
	excelName := filepath.Base(sourceExcelPath)
	if excelName == "." || excelName == "/" || excelName == "" {
		return fmt.Errorf("BuildDir: invalid excel path %q", sourceExcelPath)
	}

	if err := excel.BuildDir(sourceExcelPath, relDirPath); err != nil {
		return fmt.Errorf("BuildDir: %w", err)
	}

	resolver := NewResolver(relDirPath)

	// ensure defaults/<schemaName>.json exist
	if err := ensureDefaultDataExists(resolver); err != nil {
		return fmt.Errorf("BuildDir: %w", err)
	}

	return nil
}

// WriteDefaultData persists defaults/<schemaName>.json for the case at relDirPath.
func WriteDefaultData(relDirPath string, schemaName string) error {
	if strings.TrimSpace(schemaName) == "" {
		return fmt.Errorf("WriteDefaultData: schemaName is not set")
	}

	resolver := NewResolver(relDirPath)
	if err := writeDefaultData(resolver, schemaName); err != nil {
		return fmt.Errorf("WriteDefaultData: %w", err)
	}
	return nil
}
