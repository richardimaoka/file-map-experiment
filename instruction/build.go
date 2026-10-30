package instruction

import (
	"fmt"

	"github.com/richardimaoka/file-map-experiment/excel"
)

func ensureDefaultDataExists(resolver Resolver) error {
	_ = resolver
	return fmt.Errorf("ensureDefaultDataExists: not implemented yet")
}

// BuildDir bootstraps an instruction case at relDirPath:
// runs [excel.BuildDir], then writes defaults/<schemaName>.json sidecars.
// schemas/ is rewritten by [excel.BuildDir] when schema_preps are ready. For the layout, see [Resolver].
func BuildDir(sourceExcelPath string, relDirPath string) error {
	if err := excel.BuildDir(sourceExcelPath, relDirPath); err != nil {
		return fmt.Errorf("instruction.BuildDir: %w", err)
	}

	resolver := NewResolver(relDirPath)

	// ensure defaults/<schemaName>.json exist
	if err := ensureDefaultDataExists(resolver); err != nil {
		return fmt.Errorf("instruction.BuildDir: %w", err)
	}

	return nil
}
