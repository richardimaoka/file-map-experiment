package shinoken

import "path/filepath"

// Forward instructions case layout under $APP_BASE_DIR/<relDirPath>/:
//
//	instructions/params.json
//	instructions/filtered_params.json
//	instructions/instruction_params.json
//	instructions/text.json
//	instructions/default_instruction.json
//	instructions/variants.json
//	instructions/instructions.json
//	config.json

type ForwardResolver struct {
	caseDir
}

func NewForwardResolver(relDirPath string) ForwardResolver {
	return ForwardResolver{caseDir: newCaseDir(relDirPath)}
}

func (r ForwardResolver) InstructionsDir() string {
	return filepath.Join(r.FullDirPath(), "instructions")
}

func (r ForwardResolver) ParamsPath() string {
	return filepath.Join(r.InstructionsDir(), "params.json")
}

func (r ForwardResolver) FilteredParamsPath() string {
	return filepath.Join(r.InstructionsDir(), "filtered_params.json")
}

func (r ForwardResolver) InstructionParamsPath() string {
	return filepath.Join(r.InstructionsDir(), "instruction_params.json")
}

func (r ForwardResolver) TextPath() string {
	return filepath.Join(r.InstructionsDir(), "text.json")
}

func (r ForwardResolver) DefaultInstructionPath() string {
	return filepath.Join(r.InstructionsDir(), "default_instruction.json")
}

func (r ForwardResolver) VariantsPath() string {
	return filepath.Join(r.InstructionsDir(), "variants.json")
}

func (r ForwardResolver) InstructionsPath() string {
	return filepath.Join(r.InstructionsDir(), "instructions.json")
}
