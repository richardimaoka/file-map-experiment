package shinoken

import (
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/excel"
)

// Parameters case layout under $APP_BASE_DIR/<relDirPath>/:
//
//	<excelFileName>
//	data/<schemaName>.json
//	schemas/<schemaName>.json
//	config.json

type ParametersResolver struct {
	excel.Resolver
}

func NewParametersResolver(relDirPath string) ParametersResolver {
	return ParametersResolver{Resolver: excel.NewResolver(relDirPath)}
}

func (r ParametersResolver) DataDir() string {
	return filepath.Join(r.FullDirPath(), "data")
}

func (r ParametersResolver) SchemasDir() string {
	return filepath.Join(r.FullDirPath(), "schemas")
}

func (r ParametersResolver) DataPath(schemaName string) string {
	return filepath.Join(r.DataDir(), schemaName+".json")
}

func (r ParametersResolver) SchemaPath(schemaName string) string {
	return filepath.Join(r.SchemasDir(), schemaName+".json")
}
