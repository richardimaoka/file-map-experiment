package instruction

import (
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/env"
)

// TopLevelResolver resolves shared instruction asset paths under APP_BASE_DIR
// (not under a per-case relDirPath; see [Resolver] for case directories).
//
// Methods named *RelDir return paths relative to APP_BASE_DIR (for [BuildDir] and [NewResolver]).
// All other *Dir and *Path methods return absolute paths for file I/O.
//
// Expected layout:
//
//	$APP_BASE_DIR/
//	  assets/入稿指示書/
//	    format_sources/   # canonical workbooks (format embedded in sheets)
//	    format_schemas/
//	      <name>/         # per-format case tree (bootstrap via [BuildFormatSchema] or [BuildDirs])
//
// Example: NewTopLevelResolver().FormatSchemaRelDir("client-a") for BuildDir;
// NewTopLevelResolver().FormatSchemaDir("client-a") for reading files.
type TopLevelResolver struct {
	appBaseDir string
}

func NewTopLevelResolver() TopLevelResolver {
	return TopLevelResolver{appBaseDir: env.AppBaseDir()}
}

func (t TopLevelResolver) FormatSourcesDir() string {
	return filepath.Join(t.appBaseDir, "assets", "入稿指示書", "format_sources")
}

func (t TopLevelResolver) FormatSchemasRoot() string {
	return filepath.Join(t.appBaseDir, "assets", "入稿指示書", "format_schemas")
}

func (t TopLevelResolver) FormatSchemaRelDir(formatName string) string {
	return filepath.Join("assets", "入稿指示書", "format_schemas", formatName)
}

func (t TopLevelResolver) FormatSchemaDir(formatName string) string {
	return filepath.Join(t.FormatSchemasRoot(), formatName)
}
