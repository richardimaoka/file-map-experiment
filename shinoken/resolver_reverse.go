package shinoken

import (
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/env"
)

// Reverse (past) instructions case layout under $APP_BASE_DIR/<relDirPath>/:
//
//	instructions/past_instructions.json
//	instructions/past_specificinstructions.json
//	instructions/past_specific_common_instruction.json
//	specific_default_instruction.json
//	instructions/variant_instructions.json
//	config.json

type ReverseResolver struct {
	appBaseDir string
	relDirPath string
}

func NewReverseResolver(relDirPath string) ReverseResolver {
	baseDir := env.AppBaseDir()

	return ReverseResolver{
		appBaseDir: baseDir,
		relDirPath: relDirPath,
	}
}

func (r ReverseResolver) FullDirPath() string {
	return filepath.Join(r.appBaseDir, r.relDirPath)
}

func (r ReverseResolver) ConfigPath() string {
	return filepath.Join(r.FullDirPath(), "config.json")
}

func (r ReverseResolver) PastInstructionsPath() string {
	return filepath.Join(r.FullDirPath(), "past_instructions.json")
}

func (r ReverseResolver) PastSpecificInstructionsPath() string {
	return filepath.Join(r.FullDirPath(), "past_specific_instructions.json")
}

func (r ReverseResolver) PastSpecificCommonInstructionPath() string {
	return filepath.Join(r.FullDirPath(), "past_specific_common_instruction.json")
}

func (r ReverseResolver) PastSpecificDefaultInstructionPath() string {
	return filepath.Join(r.FullDirPath(), "past_specific_default_instruction.json")
}

func (r ReverseResolver) SpecificDefaultInstructionPath() string {
	return filepath.Join(r.FullDirPath(), "specific_default_instruction.json")
}

func (r ReverseResolver) VariantInstructionsPath() string {
	return filepath.Join(r.FullDirPath(), "variant_instructions.json")
}
