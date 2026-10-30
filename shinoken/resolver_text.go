package shinoken

import (
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/excel"
)

// Text case layout under $APP_BASE_DIR/<relDirPath>/:
//
//	<excelFileName>
//	data/<schemaName>.json
//	schemas/<schemaName>.json
//	config.json

type TextResolver struct {
	excel.Resolver
}

func NewTextResolver(relDirPath string) TextResolver {
	return TextResolver{Resolver: excel.NewResolver(relDirPath)}
}

func (r TextResolver) DataDir() string {
	return filepath.Join(r.FullDirPath(), "data")
}

func (r TextResolver) TextSchemaPath(config TextConfig) string {
	return filepath.Join(r.SchemasDir(), config.SchemaName+".json")
}

func (r TextResolver) TextDataPath(config TextConfig) string {
	return filepath.Join(r.DataDir(), config.SchemaName+".json")
}
