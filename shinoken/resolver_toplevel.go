package shinoken

import (
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/env"
)

const dataShinokenRel = "data/shinoken"

// TopLevelResolver resolves shared shinoken paths under APP_BASE_DIR (not a single case relDirPath).
//
// *RelDir methods are relative to APP_BASE_DIR; *Dir methods are absolute.
//
// Expected layout:
//
//	$APP_BASE_DIR/
//	  data/shinoken/
//	    assets/
//	    backtest/
//	    new_instructions/
//	    past_instructions/
//	    params/
//	    text/
type TopLevelResolver struct {
	appBaseDir string
}

func NewTopLevelResolver() TopLevelResolver {
	return TopLevelResolver{appBaseDir: env.AppBaseDir()}
}

func (t TopLevelResolver) RootRelDir() string {
	return dataShinokenRel
}

func (t TopLevelResolver) RootDir() string {
	return filepath.Join(t.appBaseDir, dataShinokenRel)
}

func (t TopLevelResolver) AssetsRelDir() string {
	return filepath.Join(dataShinokenRel, "assets")
}

func (t TopLevelResolver) AssetsDir() string {
	return filepath.Join(t.RootDir(), "assets")
}

func (t TopLevelResolver) BacktestRelDir() string {
	return filepath.Join(dataShinokenRel, "backtest")
}

func (t TopLevelResolver) BacktestDir() string {
	return filepath.Join(t.RootDir(), "backtest")
}

func (t TopLevelResolver) NewInstructionsRelDir() string {
	return filepath.Join(dataShinokenRel, "new_instructions")
}

func (t TopLevelResolver) NewInstructionsDir() string {
	return filepath.Join(t.RootDir(), "new_instructions")
}

func (t TopLevelResolver) PastInstructionsRelDir() string {
	return filepath.Join(dataShinokenRel, "past_instructions")
}

func (t TopLevelResolver) PastInstructionsDir() string {
	return filepath.Join(t.RootDir(), "past_instructions")
}

func (t TopLevelResolver) ParamsRelDir() string {
	return filepath.Join(dataShinokenRel, "params")
}

func (t TopLevelResolver) ParamsDir() string {
	return filepath.Join(t.RootDir(), "params")
}

func (t TopLevelResolver) TextRelDir() string {
	return filepath.Join(dataShinokenRel, "text")
}

func (t TopLevelResolver) TextDir() string {
	return filepath.Join(t.RootDir(), "text")
}

type caseDir struct {
	appBaseDir string
	relDirPath string
}

func newCaseDir(relDirPath string) caseDir {
	return caseDir{
		appBaseDir: env.AppBaseDir(),
		relDirPath: relDirPath,
	}
}

func (c caseDir) RelDirPath() string {
	return c.relDirPath
}

func (c caseDir) FullDirPath() string {
	return filepath.Join(c.appBaseDir, c.relDirPath)
}

func (c caseDir) ConfigPath() string {
	return filepath.Join(c.FullDirPath(), "config.json")
}
