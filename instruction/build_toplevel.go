package instruction

import (
	"fmt"
	"log/slog"
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
	name := filepath.Base(excelPath)
	name = strings.Replace(name, ".xlsx", "", 1)
	name = strings.Replace(name, ".xlsm", "", 1)
	name = strings.Replace(name, "【公式】", "", 1)
	name = strings.Replace(name, "依頼FMT", "", 1)
	name = strings.Replace(name, "※DDA依頼不可※", "", 1)
	name = strings.Replace(name, "01_更新者", "", 1)
	name = strings.Replace(name, "_yyyymm_案件名_", "", 1)
	name = strings.Replace(name, "_yymmdd_案件名_", "", 1)
	name = strings.Replace(name, "（アプリ領域）", "アプリ領域 ", 1)
	name = strings.Replace(name, "入稿指示書", "", 1)
	name = strings.Replace(name, "_入稿シート_", "", 1)

	// 半角カッコに囲まれた部分
	firstOpen := strings.Index(name, "(")
	lastOpen := strings.LastIndex(name, "(")
	firstClose := strings.Index(name, ")")
	lastClose := strings.LastIndex(name, ")")
	// 半角カッコがそれぞれちょうどひとつだけ存在し、かつ"("が")"より前にある場合
	if firstOpen != -1 && firstOpen == lastOpen &&
		firstClose != -1 && firstClose == lastClose &&
		firstOpen < firstClose {
		name = strings.Replace(name, "(", "", 1)
		name = strings.Replace(name, ")", "", 1)
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

	slog.Info("instruction.BuildDirs: starting", "workbook_count", len(excelPaths))

	for _, excelPath := range excelPaths {
		formatName := formatName(excelPath)
		relDir := top.FormatSchemaRelDir(formatName)

		if err := BuildDir(excelPath, relDir); err != nil {
			return fmt.Errorf("BuildDirs: %w", err)
		}
	}

	return nil
}
