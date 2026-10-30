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

type ParamsResolver struct {
	excel.Resolver
}

func NewParametersResolver(relDirPath string) ParamsResolver {
	return ParamsResolver{Resolver: excel.NewResolver(relDirPath)}
}

func (r ParamsResolver) DataDir() string {
	return filepath.Join(r.FullDirPath(), "data")
}

func (r ParamsResolver) ParamsSchemaPath(config ParamsConfig) string {
	return filepath.Join(r.SchemasDir(), config.SchemaName+".json")
}

func (r ParamsResolver) ParamsAllPath(config ParamsConfig) string {
	return filepath.Join(r.DataDir(), config.SchemaName+".json")
}
