package excel

import (
	"fmt"
	"strings"

	"github.com/richardimaoka/file-map-experiment/jsonutil"
	"github.com/xuri/excelize/v2"
)

// SchemaPrep is needed before reading Excel table schema
type SchemaPrep struct {
	SchemaName     string `json:"schema_name"`
	SheetName      string `json:"sheet_name"`
	Comment        string `json:"comment"`
	IsReady        bool   `json:"is_ready"`
	HeaderStartRow int    `json:"header_start_row"`
	HeaderEndRow   int    `json:"header_end_row"`
	DataStartRow   int    `json:"data_start_row"`
	DataEndRow     int    `json:"data_end_row"`
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

func (p SchemaPrep) DataRange() CellRange {
	return NewCellRange(
		p.StartColumn,
		p.DataStartRow,
		p.EndColumn,
		p.DataEndRow,
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

	if p.DataStartRow != 0 || p.DataEndRow != 0 {
		if p.DataStartRow < 1 {
			return fmt.Errorf("SchemaPrep.Validate: data_start_row must be >= 1")
		}
		if p.DataEndRow < 1 {
			return fmt.Errorf("SchemaPrep.Validate: data_end_row must be >= 1")
		}
		if p.DataStartRow > p.DataEndRow {
			return fmt.Errorf("SchemaPrep.Validate: data_start_row must be <= data_end_row")
		}
	}

	return nil
}

type Config struct {
	ExcelFileName string       `json:"excel_file_name"`
	SchemaPreps   []SchemaPrep `json:"schema_preps"` // manually updated in config.json
}

func (c Config) GetSchemaPrep(schemaName string) (SchemaPrep, error) {
	for _, prep := range c.SchemaPreps {
		if prep.SchemaName == schemaName {
			return prep, nil
		}
	}

	return SchemaPrep{}, fmt.Errorf("Config.GetSchemaPrep: schema %q not found in config", schemaName)
}

// Validate checks the config contents, but doesn't validate if the directory holding config.json
// is in a good shape (e.g. if the Excel file exists).
func (c Config) Validate() error {
	if strings.TrimSpace(c.ExcelFileName) == "" {
		return fmt.Errorf("Config.Validate: excel_file_name is not set")
	}
	if !strings.HasSuffix(c.ExcelFileName, ".xlsx") && !strings.HasSuffix(c.ExcelFileName, ".xlsm") {
		return fmt.Errorf("Config.Validate: excel_file_name must end with .xlsx or .xlsm")
	}

	return nil
}

// The excelFileName argument is needed, assuming you already know
// which Excel file to read upon creating config.json
func InitConfig(excelFileName string) Config {
	return Config{
		ExcelFileName: excelFileName,
	}
}

// enrichConfig sets cfg.SchemaPreps from mergeSchemaPrepsFromWorkbook.
func enrichConfig(cfg Config, excelPath string) (Config, error) {
	preps, err := mergeSchemaPrepsFromWorkbook(cfg.SchemaPreps, excelPath)
	if err != nil {
		return Config{}, err
	}
	cfg.SchemaPreps = preps
	return cfg, nil
}

func schemaPrepIsLocked(p SchemaPrep) bool {
	return p.IsReady || p.Comment != ""
}

// mergeSchemaPrepsFromWorkbook lists one prep per workbook tab. Locked preps (is_ready or comment) stay as-is if sheet_name exists; missing sheet is error. Unlocked existing preps are dropped.
func mergeSchemaPrepsFromWorkbook(existing []SchemaPrep, excelPath string) ([]SchemaPrep, error) {
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, fmt.Errorf("open workbook: %w", err)
	}
	defer f.Close()

	sheetNames := f.GetSheetList()
	sheets := make(map[string]struct{}, len(sheetNames))
	for _, name := range sheetNames {
		if name == "" {
			continue
		}
		sheets[name] = struct{}{}
	}

	lockedBySheet := make(map[string]SchemaPrep)
	for _, prep := range existing {
		if !schemaPrepIsLocked(prep) {
			continue
		}
		if prep.SheetName == "" {
			return nil, fmt.Errorf(
				"schema_prep %q has is_ready or comment set but sheet_name is empty; set sheet_name to the Excel tab name",
				prep.SchemaName,
			)
		}
		if _, ok := sheets[prep.SheetName]; !ok {
			return nil, fmt.Errorf(
				"schema_prep %q (sheet %q): sheet not in workbook after copy; if the sheet was renamed, update sheet_name in config.json to the new tab name; if the sheet was removed, clear is_ready and comment on this prep so BuildDir can drop it, then adjust schema_preps if needed",
				prep.SchemaName,
				prep.SheetName,
			)
		}
		if _, dup := lockedBySheet[prep.SheetName]; dup {
			return nil, fmt.Errorf(
				"multiple schema_preps use sheet %q; keep one prep per sheet",
				prep.SheetName,
			)
		}
		lockedBySheet[prep.SheetName] = prep
	}

	merged := make([]SchemaPrep, 0, len(sheetNames))
	for _, sheetName := range sheetNames {
		if sheetName == "" {
			continue
		}
		if prep, ok := lockedBySheet[sheetName]; ok {
			merged = append(merged, prep)
			continue
		}
		merged = append(merged, defaultSchemaPrepForSheet(sheetName))
	}
	return merged, nil
}

func defaultSchemaPrepForSheet(sheetName string) SchemaPrep {
	return SchemaPrep{
		SchemaName:     sheetName,
		SheetName:      sheetName,
		Comment:        "",
		IsReady:        false,
		HeaderStartRow: 1,
		HeaderEndRow:   1,
		StartColumn:    "A",
		EndColumn:      "A",
	}
}

func ReadConfig(jsonPath string) (Config, error) {
	var config Config
	err := jsonutil.Read(jsonPath, &config)
	if err != nil {
		return Config{}, fmt.Errorf("excel.ReadConfig: %w", err)
	}

	if err := config.Validate(); err != nil {
		return Config{}, fmt.Errorf("excel.ReadConfig: %w", err)
	}

	return config, nil
}
