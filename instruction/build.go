package instruction

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/richardimaoka/file-map-experiment/excel"
	"github.com/richardimaoka/file-map-experiment/fs"
	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

func updateDefaultDataFromConfig(resolver Resolver) error {
	if err := os.MkdirAll(resolver.DefaultsDir(), 0o755); err != nil {
		return fmt.Errorf("updateDefaultDataFromConfig: %w", err)
	}

	schemaPaths, err := fs.FilesWithExtensions(resolver.SchemasDir(), ".json")
	if err != nil {
		return fmt.Errorf("updateDefaultDataFromConfig: %w", err)
	}

	for _, schemaPath := range schemaPaths {
		fileName := filepath.Base(schemaPath)
		schemaName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		rows, err := excel.ReadRowsFromExcel(resolver.Resolver, schemaName)
		if err != nil {
			return fmt.Errorf("updateDefaultDataFromConfig: %w", err)
		}
		if err := jsonutil.Write(resolver.DefaultDataPath(schemaName), rows); err != nil {
			return fmt.Errorf("updateDefaultDataFromConfig: %w", err)
		}
	}
	return nil
}

// BuildDir bootstraps an instruction case at relDirPath:
// runs [excel.BuildDir], then rewrites defaults/<schemaName>.json from the copied workbook
// for each schemas/<schemaName>.json present after bootstrap. For the layout, see [Resolver].
func BuildDir(sourceExcelPath string, relDirPath string) error {
	if err := excel.BuildDir(sourceExcelPath, relDirPath); err != nil {
		return fmt.Errorf("instruction.BuildDir: %w", err)
	}

	resolver := NewResolver(relDirPath)

	if err := updateDefaultDataFromConfig(resolver); err != nil {
		return fmt.Errorf("instruction.BuildDir: %w", err)
	}

	return nil
}
