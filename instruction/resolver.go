package instruction

import (
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/excel"
)

// Resolver extends [excel.Resolver] with instruction case paths under the same case directory.
//
// Expected layout (excel case tree plus instruction sidecars):
//
//	$APP_BASE_DIR/<relDirPath>/
//	  <excelFileName>   # .xlsx/.xlsm; name from config.json
//	  config.json
//	  schemas/
//	    <schemaName>.json
//	  defaults/
//	    <schemaName>.json
//
// Example: NewResolver("cases/client1/instruction-sheet")
type Resolver struct {
	excel.Resolver
}

func NewResolver(relDirPath string) Resolver {
	return Resolver{Resolver: excel.NewResolver(relDirPath)}
}

func (r Resolver) DefaultDataPath(schemaName string) string {
	return filepath.Join(r.FullDirPath(), "defaults", schemaName+".json")
}
