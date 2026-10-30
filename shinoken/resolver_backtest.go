package shinoken

import (
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/env"
	"github.com/richardimaoka/file-map-experiment/instruction"
)

// Backtest case layout under $APP_BASE_DIR/<relDirPath>/:
//
//	forward/
//	reverse/
//	config.json

type BacktestResolver struct {
	appBaseDir string
	relDirPath string

	Forward      ForwardResolver
	Reverse      ReverseResolver
	Instructions instruction.Resolver
}

func NewBacktestResolver(relDirPath string) BacktestResolver {
	baseDir := env.AppBaseDir()

	return BacktestResolver{
		appBaseDir: baseDir,
		relDirPath: relDirPath,
	}
}

func (r BacktestResolver) FullDirPath() string {
	return filepath.Join(r.appBaseDir, r.relDirPath)
}
