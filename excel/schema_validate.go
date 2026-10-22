package excel

import (
	"fmt"
	"os"
	"reflect"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
)

// validateExcelAgainstSchemaDetailedDiff controls whether a schema mismatch warning
// includes a JSON line diff (persisted schema vs workbook-derived schema).
const validateExcelAgainstSchemaDetailedDiff = false

// ValidateExcelAgainstSchema loads schemas/<name>.json, derives the schema from the
// instance workbook, and warns on mismatch. It never returns an error; on success it
// returns the workbook-derived schema.
func ValidateExcelAgainstSchema(resolver PathResolver, schemaName string) TableSchema {
	persisted, err := ReadSchemaFromJson(resolver.SchemaPath(schemaName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ValidateExcelAgainstSchema: warning: could not load persisted schema %q: %v\n", schemaName, err)
		return TableSchema{}
	}

	fromExcel, err := ReadSchemaFromExcel(resolver, schemaName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ValidateExcelAgainstSchema: warning: could not read schema from workbook for %q: %v\n", schemaName, err)
		return persisted
	}

	if reflect.DeepEqual(persisted, fromExcel) {
		return fromExcel
	}

	warnExcelSchemaMismatch(schemaName, persisted, fromExcel)
	return fromExcel
}

func warnExcelSchemaMismatch(schemaName string, persisted, fromExcel TableSchema) {
	if !validateExcelAgainstSchemaDetailedDiff {
		fmt.Fprintf(os.Stderr, "ValidateExcelAgainstSchema: warning: workbook schema does not match persisted schema %q\n", schemaName)
		return
	}

	persistedJSON, err := jsonfile.MarshalIndent(persisted)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ValidateExcelAgainstSchema: warning: workbook schema does not match persisted schema %q (marshal persisted: %v)\n", schemaName, err)
		return
	}
	fromExcelJSON, err := jsonfile.MarshalIndent(fromExcel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ValidateExcelAgainstSchema: warning: workbook schema does not match persisted schema %q (marshal workbook: %v)\n", schemaName, err)
		return
	}

	fmt.Fprintf(os.Stderr,
		"ValidateExcelAgainstSchema: warning: workbook schema does not match persisted schema %q (-persisted +workbook):\n%s",
		schemaName,
		jsonLineDiff(string(persistedJSON), string(fromExcelJSON)),
	)
}
