package excel

import (
	"fmt"
	"strings"

	"github.com/richardimaoka/file-map-experiment/jsonfile"
	"github.com/xuri/excelize/v2"
)

// SchemaPrep is needed before reading Excel table schema
type SchemaPrep struct {
	SchemaName     string `json:"schema_name"`
	SheetName      string `json:"sheet_name"`
	IsReady        bool   `json:"is_ready"`
	HeaderStartRow int    `json:"header_start_row"`
	HeaderEndRow   int    `json:"header_end_row"`
	DataStartRow   int    `json:"data_start_row"`
	StartColumn    string `json:"start_column"`
	EndColumn      string `json:"end_column"`
}

func (p SchemaPrep) HeaderRange() CellRange {
	return NewCellRange(
		p.StartColumn,
		p.HeaderStartRow,
		p.EndColumn,
		p.HeaderEndRow,
	)
}

func (p SchemaPrep) Validate() error {
	if strings.TrimSpace(p.SchemaName) == "" {
		return fmt.Errorf("SchemaPrep.Validate: schema_name is not set")
	}
	if strings.TrimSpace(p.SheetName) == "" {
		return fmt.Errorf("SchemaPrep.Validate: sheet_name is not set")
	}

	startCol, err := excelize.ColumnNameToNumber(p.StartColumn)
	if err != nil {
		return fmt.Errorf("SchemaPrep.Validate: start_column: %w", err)
	}
	endCol, err := excelize.ColumnNameToNumber(p.EndColumn)
	if err != nil {
		return fmt.Errorf("SchemaPrep.Validate: end_column: %w", err)
	}
	if startCol > endCol {
		return fmt.Errorf("SchemaPrep.Validate: start_column must be <= end_column")
	}
	if p.HeaderStartRow < 1 {
		return fmt.Errorf("SchemaPrep.Validate: header_start_row must be >= 1")
	}
	if p.HeaderEndRow < 1 {
		return fmt.Errorf("SchemaPrep.Validate: header_end_row must be >= 1")
	}
	if p.HeaderStartRow > p.HeaderEndRow {
		return fmt.Errorf("SchemaPrep.Validate: header_start_row must be <= header_end_row")
	}

	return nil
}

type ExcelConfig struct {
	ExcelFileName string       `json:"excel_file_name"`
	SchemaPreps   []SchemaPrep `json:"schema_preps"` // manually updated in config.json
}

func (c ExcelConfig) GetSchemaPrep(schemaName string) (SchemaPrep, error) {
	for _, prep := range c.SchemaPreps {
		if prep.SchemaName == schemaName {
			return prep, nil
		}
	}

	return SchemaPrep{}, fmt.Errorf("ExcelConfig.GetSchemaPrep: schema %q not found in config", schemaName)
}

// The excelFileName argument is needed, assuming you already know
// which Excel file to read upon creating config.json
func InitConfig(excelFileName string) ExcelConfig {
	return ExcelConfig{
		ExcelFileName: excelFileName,
	}
}

// Validate checks the config contents, but doesn't validate if the directory holding config.json
// is in a good shape (e.g. if the Excel file exists).
func (c ExcelConfig) Validate() error {
	if strings.TrimSpace(c.ExcelFileName) == "" {
		return fmt.Errorf("ExcelConfig.Validate: excel_file_name is not set")
	}
	if !strings.HasSuffix(c.ExcelFileName, ".xlsx") && !strings.HasSuffix(c.ExcelFileName, ".xlsm") {
		return fmt.Errorf("ExcelConfig.Validate: excel_file_name must end with .xlsx or .xlsm")
	}

	return nil
}

func ReadExcelConfig(jsonPath string) (ExcelConfig, error) {
	var config ExcelConfig
	err := jsonfile.Read(jsonPath, &config)
	if err != nil {
		return ExcelConfig{}, fmt.Errorf("ReadExcelConfig: %w", err)
	}

	if err := config.Validate(); err != nil {
		return ExcelConfig{}, fmt.Errorf("ReadExcelConfig: %w", err)
	}

	return config, nil
}
