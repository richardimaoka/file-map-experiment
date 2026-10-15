package excel

import (
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/env"
)

type PathResolver struct {
	dirPath       string
	excelFileName string
}

func NewPathResolver(dirPath, excelFileName string) *PathResolver {
	baseDir := env.AppBaseDir()
	fullPath := filepath.Join(baseDir, dirPath)

	return &PathResolver{
		dirPath:       fullPath,
		excelFileName: excelFileName,
	}
}

func (p *PathResolver) DirPath() string {
	return p.dirPath
}

func (p *PathResolver) ExcelPath() string {
	path := filepath.Join(p.dirPath, p.excelFileName)
	return path
}

func (p *PathResolver) ConfigPath() string {
	path := filepath.Join(p.dirPath, "config.json")
	return path
}

func (p *PathResolver) SchemaDirPath() string {
	path := filepath.Join(p.dirPath, "schemas")
	return path
}

func (p *PathResolver) SchemaPath(schemaName string) string {
	path := filepath.Join(p.dirPath, "schemas", schemaName+".json")
	return path
}
