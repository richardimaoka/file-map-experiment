package instruction

import (
	"fmt"
	"slices"
	"strings"

	"github.com/richardimaoka/file-map-experiment/excel"
	"github.com/xuri/excelize/v2"
)

type Config struct {
	excel.Config
}

func ReadConfig(jsonPath string) (Config, error) {
	cfg, err := excel.ReadConfig(jsonPath)
	if err != nil {
		return Config{}, fmt.Errorf("instruciton.ReadConfig: %v", err)
	}
	return Config{Config: cfg}, nil
}

const prepCommentNeedsReview = "要確認"

var managementSheets = []string{
	"表紙",
	"変更履歴",
	"入力規則",
	"入稿FMTの説明",
	"原稿規定",
	"マイグレーション履歴",
	"シート定義",
	"必須チェック",
	"対応表_転記",
	"デフォルト設定",
	"自動入稿時確認シート",
}

// enrichConfig adjusts cfg.SchemaPreps for instruction workbooks: drops management sheets,
// flags other sheets without 入力規則 text or data validation with Comment 要確認.
func enrichConfig(cfg Config, excelPath string) (Config, error) {
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		return Config{}, fmt.Errorf("enrichConfig: %w", err)
	}
	defer f.Close()

	filtered := make([]excel.SchemaPrep, 0, len(cfg.SchemaPreps))
	for _, prep := range cfg.SchemaPreps {
		sheetName := prep.SheetName
		if sheetName == "" {
			return Config{}, fmt.Errorf("enrichConfig: schema %q: sheet_name is not set", prep.SchemaName)
		}

		idx, err := f.GetSheetIndex(sheetName)
		if err != nil {
			return Config{}, fmt.Errorf("enrichConfig: sheet %q: %w", sheetName, err)
		}
		if idx < 0 {
			return Config{}, fmt.Errorf("enrichConfig: sheet %q not found in workbook", sheetName)
		}

		if isManagementSheet(sheetName) {
			continue
		}

		hasRulesText, err := sheetContainsText(f, sheetName, "入力規則")
		if err != nil {
			return Config{}, fmt.Errorf("enrichConfig: %w", err)
		}
		hasDataValidation, err := sheetHasDataValidation(f, sheetName)
		if err != nil {
			return Config{}, fmt.Errorf("enrichConfig: %w", err)
		}

		if !hasRulesText && !hasDataValidation {
			prep.Comment = prepCommentNeedsReview
		}
		filtered = append(filtered, prep)
	}

	cfg.SchemaPreps = filtered
	return cfg, nil
}

func isManagementSheet(name string) bool {
	return slices.Contains(managementSheets, name)
}

func sheetContainsText(f *excelize.File, sheetName, needle string) (bool, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return false, fmt.Errorf("read sheet %q: %w", sheetName, err)
	}
	for _, row := range rows {
		for _, cell := range row {
			if strings.Contains(cell, needle) {
				return true, nil
			}
		}
	}
	return false, nil
}

func sheetHasDataValidation(f *excelize.File, sheetName string) (bool, error) {
	dvs, err := f.GetDataValidations(sheetName)
	if err != nil {
		return false, fmt.Errorf("data validations on sheet %q: %w", sheetName, err)
	}
	return len(dvs) > 0, nil
}
