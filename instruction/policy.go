package instruction

import (
	"fmt"
	"strings"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

// PolicyFileName is the optional sidecar next to schema_prep.json / workbook fixtures.
const PolicyFileName = "sheet_policy.json"

// ColumnType classifies the data band for writers and validators (declared in SheetPolicy, not inferred).
const (
	ColumnTypeValue       = "value"
	ColumnTypeFormula     = "formula"
	ColumnTypeMixed       = "mixed"       // value and formula cells in the same column
	ColumnTypeIncremental = "incremental" // first row literal 1, following rows increment via formula
)

// SheetPolicy is write/validation metadata for one sheet in an instruction workbook.
// Kept separate from excel.TableSchema (structural).
type SheetPolicy struct {
	SchemaName           string            `json:"schema_name"`
	PreserveFormulaCells bool              `json:"preserve_formula_cells,omitempty"`
	ColumnTypes          map[string]string `json:"column_types,omitempty"` // excel column letter → ColumnType*
}

func ReadSheetPolicyFromJSON(jsonPath string) (SheetPolicy, error) {
	var policy SheetPolicy
	if err := jsonfile.Read(jsonPath, &policy); err != nil {
		return SheetPolicy{}, fmt.Errorf("ReadSheetPolicyFromJSON: %w", err)
	}
	if err := policy.Validate(); err != nil {
		return SheetPolicy{}, fmt.Errorf("ReadSheetPolicyFromJSON: %w", err)
	}
	return policy, nil
}

func (p SheetPolicy) Validate() error {
	if strings.TrimSpace(p.SchemaName) == "" {
		return fmt.Errorf("SheetPolicy.Validate: schema_name is not set")
	}
	for col, typ := range p.ColumnTypes {
		if !isKnownColumnType(typ) {
			return fmt.Errorf("SheetPolicy.Validate: column_types[%q]: unknown column_type %q", col, typ)
		}
	}
	return nil
}

func (p SheetPolicy) ColumnType(excelColumn string) string {
	if p.ColumnTypes == nil {
		return ""
	}
	return p.ColumnTypes[excelColumn]
}

func isKnownColumnType(typ string) bool {
	switch typ {
	case ColumnTypeValue, ColumnTypeFormula, ColumnTypeMixed, ColumnTypeIncremental:
		return true
	default:
		return false
	}
}
