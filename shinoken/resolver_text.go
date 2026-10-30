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

func (r TextResolver) SchemasDir() string {
	return filepath.Join(r.FullDirPath(), "schemas")
}

func (r TextResolver) DataPath(schemaName string) string {
	return filepath.Join(r.DataDir(), schemaName+".json")
}

func (r TextResolver) SchemaPath(schemaName string) string {
	return filepath.Join(r.SchemasDir(), schemaName+".json")
}
