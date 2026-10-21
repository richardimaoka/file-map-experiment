package instruction

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/richardimaoka/file-map-experiment/excel"
	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

// Sheet is one logical table: excel bounds (prep), structural schema, optional policy.
type Sheet struct {
	Prep   excel.SchemaPrep
	Schema excel.TableSchema
	Policy *SheetPolicy
}

// LoadSheet reads schema_prep.json, builds TableSchema from the workbook, and optional sheet_policy.json.
func LoadSheet(caseDir, excelPath string) (Sheet, error) {
	prepPath := filepath.Join(caseDir, "schema_prep.json")
	var prep excel.SchemaPrep
	if err := jsonfile.Read(prepPath, &prep); err != nil {
		return Sheet{}, fmt.Errorf("LoadSheet: %w", err)
	}

	schema, err := excel.ReadTableSchema(excelPath, prep)
	if err != nil {
		return Sheet{}, fmt.Errorf("LoadSheet: %w", err)
	}

	policyPath := filepath.Join(caseDir, PolicyFileName)
	var policy *SheetPolicy
	if _, err := os.Stat(policyPath); err == nil {
		p, err := ReadSheetPolicyFromJSON(policyPath)
		if err != nil {
			return Sheet{}, fmt.Errorf("LoadSheet: %w", err)
		}
		if p.SchemaName != prep.SchemaName {
			return Sheet{}, fmt.Errorf(
				"LoadSheet: %s schema_name %q does not match prep %q",
				PolicyFileName, p.SchemaName, prep.SchemaName,
			)
		}
		policy = &p
	} else if !os.IsNotExist(err) {
		return Sheet{}, fmt.Errorf("LoadSheet: %w", err)
	}

	return Sheet{Prep: prep, Schema: schema, Policy: policy}, nil
}
