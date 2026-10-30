package shinoken

import "path/filepath"

// Assets case layout under $APP_BASE_DIR/<relDirPath>/:
//
//	<parameters-excel>.xlsx
//	<text-excel>.xlsx
//	specific_default_instruction.json
//	config.json

type AssetsResolver struct {
	caseDir
}

func NewAssetsResolver(relDirPath string) AssetsResolver {
	return AssetsResolver{caseDir: newCaseDir(relDirPath)}
}

func (r AssetsResolver) SpecificDefaultInstructionPath() string {
	return filepath.Join(r.FullDirPath(), "specific_default_instruction.json")
}

func (r AssetsResolver) ParamsExcelPath(config AssetsConfig) string {
	return filepath.Join(r.FullDirPath(), config.ParamsExcelFile)
}

func (r AssetsResolver) TextExcelPath(config AssetsConfig) string {
	return filepath.Join(r.FullDirPath(), config.TextExcelFile)
}
