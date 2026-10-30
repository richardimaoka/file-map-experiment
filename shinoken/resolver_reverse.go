package shinoken

import "path/filepath"

// Reverse (past) instructions case layout under $APP_BASE_DIR/<relDirPath>/:
//
//	instructions/past_instructions.json
//	instructions/past_specificinstructions.json
//	instructions/past_specific_common_instruction.json
//	instructions/past_specific_default_instruction.json
//	specific_default_instruction.json
//	instructions/variant_instructions.json
//	config.json

type ReverseResolver struct {
	ForwardResolver
}

func NewReverseResolver(relDirPath string) ReverseResolver {
	return ReverseResolver{ForwardResolver: NewForwardResolver(relDirPath)}
}

func (r ReverseResolver) PastInstructionsPath() string {
	return filepath.Join(r.InstructionsDir(), "past_instructions.json")
}

func (r ReverseResolver) PastSpecificInstructionsPath() string {
	return filepath.Join(r.InstructionsDir(), "past_specificinstructions.json")
}

func (r ReverseResolver) PastSpecificCommonInstructionPath() string {
	return filepath.Join(r.InstructionsDir(), "past_specific_common_instruction.json")
}

func (r ReverseResolver) PastSpecificDefaultInstructionPath() string {
	return filepath.Join(r.InstructionsDir(), "past_specific_default_instruction.json")
}

func (r ReverseResolver) SpecificDefaultInstructionPath() string {
	return filepath.Join(r.FullDirPath(), "specific_default_instruction.json")
}

func (r ReverseResolver) VariantInstructionsPath() string {
	return filepath.Join(r.InstructionsDir(), "variant_instructions.json")
}
