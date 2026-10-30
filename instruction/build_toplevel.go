package instruction

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/richardimaoka/file-map-experiment/fs"
)

func formatSourceExcelPaths(resolver TopLevelResolver) ([]string, error) {
	parent := resolver.FormatSourcesDir()

	excelFiles, err := fs.FilesWithExtensions(parent, ".xlsx", ".xlsm")
	if err != nil {
		return nil, fmt.Errorf("formatSourceExcelPath: %s", err)
	}

	excelPaths := []string{}
	for _, filename := range excelFiles {
		excelPaths = append(excelPaths, filepath.Join(parent, filename))
	}

	return excelPaths, nil
}

func formatName(excelPath string) string {
	excelFileName := filepath.Base(excelPath)

	var name string
	name = strings.Replace(excelFileName, ".xlsx", "", 1)
	name = strings.Replace(excelFileName, "01_更新者_yyyymmdd_案件名", "", 1)
	name = strings.Replace(excelFileName, "yyyymmdd_案件名_", "", 1)
	name = strings.Replace(excelFileName, "【公式】", "", 1)
	name = strings.Replace(excelFileName, "依頼FMT", "", 1)
	name = strings.Replace(excelFileName, "（アプリ領域）", "", 1)
	name = strings.Replace(excelFileName, "入稿指示書", "", 1)
	name = strings.Replace(excelFileName, "_入稿シート_", "", 1)

	// 半角カッコに囲まれた部分
	firstOpen := strings.Index(excelFileName, "(")
	lastOpen := strings.LastIndex(excelFileName, "(")
	firstClose := strings.Index(excelFileName, ")")
	lastClose := strings.LastIndex(excelFileName, ")")
	// 半角カッコがそれぞれちょうどひとつだけ存在し、かつ"("が")"より前にある場合
	if firstOpen != -1 && firstOpen == lastOpen &&
		firstClose != -1 && firstClose == lastClose &&
		firstOpen < firstClose {
		name = strings.Replace(excelFileName, "(", "", 1)
		name = strings.Replace(excelFileName, ")", "", 1)
	}

	return name

}

// BuildDirs bootstraps every format_schemas/<name>/ tree from workbooks in format_sources/.
// It does not modify format_sources/. Paths come from [TopLevelResolver].
func BuildDirs() error {
	top := NewTopLevelResolver()

	excelPaths, err := formatSourceExcelPaths(top)
	if err != nil {
		return fmt.Errorf("BuildDirs: %w", err)
	}

	for _, excelPath := range excelPaths {
		formatName := formatName(excelPath)
		formatSchemaDir := top.FormatSchemaDir(formatName)

		if err := BuildDir(excelPath, formatSchemaDir); err != nil {
			return fmt.Errorf("BuildDirs: %w", err)
		}
	}

	return nil
}
