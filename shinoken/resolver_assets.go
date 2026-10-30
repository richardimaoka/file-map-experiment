package shinoken

import (
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/env"
)

// Assets case layout under $APP_BASE_DIR/<relDirPath>/:
//
//	<parameters-excel>.xlsx
//	<text-excel>.xlsx
//	specific_default_instruction.json
//	config.json

type AssetsResolver struct {
	appBaseDir string
	relDirPath string
}

func NewAssetsResolver(relDirPath string) AssetsResolver {
	baseDir := env.AppBaseDir()

	return AssetsResolver{
		appBaseDir: baseDir,
		relDirPath: relDirPath,
	}
}

func (r AssetsResolver) FullDirPath() string {
	return filepath.Join(r.appBaseDir, r.relDirPath)
}

func (r AssetsResolver) ConfigPath() string {
	return filepath.Join(r.FullDirPath(), "config.json")
}

func (r AssetsResolver) SpecificDefaultInstructionPath(fileName string) string {
	return filepath.Join(r.FullDirPath(), "specific_default_instruction.json")
}

func (r AssetsResolver) ParamsExcelPath(config AssetsConfig) string {
	return filepath.Join(r.FullDirPath(), config.ParamsExcelFile)
}

func (r AssetsResolver) TextExcelPath(config AssetsConfig) string {
	return filepath.Join(r.FullDirPath(), config.TextExcelFile)
}
