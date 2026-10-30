package shinoken

import (
	"fmt"
	"os"
	"strings"

	"github.com/richardimaoka/file-map-experiment/fs"
	"github.com/richardimaoka/file-map-experiment/jsonutil"
)

func ensureAssetConfigExists(resolver AssetsResolver) error {
	configPath := resolver.ConfigPath()

	// if already exists, return early
	if fs.FileExists(configPath) {
		return nil
	}

	var paramsExcelFile, textExcelFile string

	dirPath := resolver.FullDirPath()
	excelFiles, err := fs.FilesWithExtensions(dirPath, ".xlsx", ".xlsm")
	if err != nil {
		return fmt.Errorf("ensureAssetConfigExists: could not find Excel files, %w", err)
	}

	for _, filename := range excelFiles {
		if strings.Contains(filename, "パラメータ") {
			paramsExcelFile = filename
			continue
		}

		if strings.Contains(filename, "入稿内容確認") {
			textExcelFile = filename
			continue
		}
	}

	if paramsExcelFile == "" {
		return fmt.Errorf("ensureAssetConfigExists: failed to find the params Excel file, %w", err)
	}

	if textExcelFile == "" {
		return fmt.Errorf("ensureAssetConfigExists: failed to find the text Excel file, %w", err)
	}

	config := CreateAssetsConfig(paramsExcelFile, textExcelFile)
	if err := jsonutil.Write(configPath, config); err != nil {
		return fmt.Errorf("ensureAssetConfigExists: %w", err)
	}

	return nil
}

func ensureParamsExcelExists(resolver AssetsResolver) error {
	configPath := resolver.ConfigPath()
	config, err := ReadAssetsConfig(configPath)
	if err != nil {
		return fmt.Errorf("ensureParamsExcelExists: %w", err)
	}

	paramsExcelPath := resolver.ParamsExcelPath(config)
	if err := manualFileExists(paramsExcelPath); err != nil {
		return fmt.Errorf("ensureParamsExcelExists: %s", err)
	}

	return nil
}

func ensureTextExcelExists(resolver AssetsResolver) error {
	configPath := resolver.ConfigPath()
	config, err := ReadAssetsConfig(configPath)
	if err != nil {
		return fmt.Errorf("ensureTextExcelExists: %w", err)
	}

	textExcelPath := resolver.TextExcelPath(config)
	if err := manualFileExists(textExcelPath); err != nil {
		return fmt.Errorf("ensureTextExcelExists: %s", err)
	}

	return nil
}

// BuildAssetsDir bootstraps an assets case directory for [AssetsResolver]:
// - create the directory if it does not exist
// - ensure config.json exists (write it from workbook filenames in the case dir when missing)
// - ensure the params workbook named in config exists
// - ensure the text workbook named in config exists
//
// For the full on-disk layout, see [AssetsResolver].
func BuildAssetsDir(resolver AssetsResolver) error {
	dirPath := resolver.FullDirPath()
	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		return fmt.Errorf("BuildAssetsDir: failed to create dir, %w", err)
	}

	if err := ensureAssetConfigExists(resolver); err != nil {
		return fmt.Errorf("BuildAssetsDir: %w", err)
	}

	if err := ensureParamsExcelExists(resolver); err != nil {
		return fmt.Errorf("BuildAssetsDir: %w", err)
	}

	if err := ensureTextExcelExists(resolver); err != nil {
		return fmt.Errorf("BuildAssetsDir: %w", err)
	}

	return nil
}
