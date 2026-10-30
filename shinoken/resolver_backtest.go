package shinoken

import "path/filepath"

// Backtest case layout under $APP_BASE_DIR/<relDirPath>/:
//
//	forward/
//	reverse/
//	config.json

type BacktestResolver struct {
	caseDir
}

func NewBacktestResolver(relDirPath string) BacktestResolver {
	return BacktestResolver{caseDir: newCaseDir(relDirPath)}
}

func (r BacktestResolver) ForwardDir() string {
	return filepath.Join(r.FullDirPath(), "forward")
}

func (r BacktestResolver) ReverseDir() string {
	return filepath.Join(r.FullDirPath(), "reverse")
}
