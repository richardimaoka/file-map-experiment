package shinoken

import (
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/env"
	"github.com/richardimaoka/file-map-experiment/instruction"
)

// Forward instructions case layout under $APP_BASE_DIR/<relDirPath>/:
//
//	instructions/params.json
//	instructions/filtered_params.json
//	instructions/text.json
//	instructions/specific_default_instruction.json
//	instructions/variants.json
//	instructions/instructions.json
//	config.json

type ForwardResolver struct {
	appBaseDir string
	relDirPath string

	Instruction instruction.Resolver
}

func NewForwardResolver(relDirPath string) ForwardResolver {
	baseDir := env.AppBaseDir()

	return ForwardResolver{
		appBaseDir: baseDir,
		relDirPath: relDirPath,
	}
}

func (r ForwardResolver) FullDirPath() string {
	return filepath.Join(r.appBaseDir, r.relDirPath)
}

func (r ForwardResolver) ConfigPath() string {
	return filepath.Join(r.FullDirPath(), "config.json")
}

func (r ForwardResolver) ParamsJsonPath() string {
	return filepath.Join(r.FullDirPath(), "params.json")
}

func (r ForwardResolver) FilteredParamsPath() string {
	return filepath.Join(r.FullDirPath(), "filtered_params.json")
}

func (r ForwardResolver) InstructionParamsPath() string {
	return filepath.Join(r.FullDirPath(), "instruction_params.json")
}

func (r ForwardResolver) TextJsonPath() string {
	return filepath.Join(r.FullDirPath(), "text.json")
}

func (r ForwardResolver) SpecificDefaultInstructionPath() string {
	return filepath.Join(r.FullDirPath(), "specific_default_instruction.json")
}

func (r ForwardResolver) VariantsJsonPath() string {
	return filepath.Join(r.FullDirPath(), "variants.json")
}

func (r ForwardResolver) InstructionsJsonPath() string {
	return filepath.Join(r.FullDirPath(), "instructions.json")
}
