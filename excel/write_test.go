package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

// Fixture variants under excel/testdata/ (same fictional ad data, different sheet layout):
//   - simple: single-row headers, no grouped rows, no summary rows
//   - multi_row_headers
//   - grouped_rows: 広告アカウント / 広告キャンペーン / 広告グループ merged for repeated values
//   - summary_rows
//   - origin_d5: simple table with upper-left at D5
//   - complex: multi-row headers, grouped rows, summary rows, origin at D5

func TestCreateAdsReportFixturesInTestdata(t *testing.T) {
	dir := filepath.Join("testdata")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name    string
		file    string
		writeFn func(f *excelize.File, sheet string) error
	}{
		{"simple", "ads_report_simple.xlsx", writeAdsReportSimple},
		{"multi_row_headers", "ads_report_multi_row_headers.xlsx", writeAdsReportMultiRowHeaders},
		{"grouped_rows", "ads_report_grouped_rows.xlsx", writeAdsReportGroupedRows},
		{"summary_rows", "ads_report_summary_rows.xlsx", writeAdsReportSummaryRows},
		{"complex", "ads_report_complex.xlsx", writeAdsReportComplex},
		{"origin_d5", "ads_report_origin_d5.xlsx", writeAdsReportOriginD5},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(dir, tc.file)
			f := excelize.NewFile()
			defer f.Close()

			sheet := f.GetSheetName(0)
			if err := tc.writeFn(f, sheet); err != nil {
				t.Fatal(err)
			}
			if err := f.SaveAs(path); err != nil {
				t.Fatal(err)
			}
			info, err := os.Stat(path)
			if err != nil {
				t.Fatal(err)
			}
			if info.Size() == 0 {
				t.Fatalf("%s was not written", tc.file)
			}
		})
	}
}

type adsReportSpec struct {
	account  string
	campaign string
	group    string
	ad       string
	creative string
}

func adsReportLeafHeaders() []interface{} {
	return []interface{}{
		"広告アカウント", "広告キャンペーン", "広告グループ", "広告", "クリエイティブ",
		"インプレッション", "費用", "Conversion",
	}
}

// Fictional 株式会社サクラブレンド group: coffee, energy drink, and ice dessert brands.
func adsReportSpecs() []adsReportSpec {
	return []adsReportSpec{
		{"サクラブレンド・コーヒー", "2026春_新焙煎豆", "関東_コーヒー愛好家", "モーニングセット告知", "sakura_morning_300x250.png"},
		{"サクラブレンド・コーヒー", "2026春_新焙煎豆", "関東_コーヒー愛好家", "店頭POP連動", "sakura_pop_728x90.jpg"},
		{"サクラブレンド・コーヒー", "2026春_新焙煎豆", "関東_コーヒー愛好家", "アプリ初回クーポン", "sakura_app_coupon_320x50.png"},
		{"サクラブレンド・コーヒー", "2026春_新焙煎豆", "全国_20代女性", "ラテアート動画", "sakura_latte_video_1080x1920.jpg"},
		{"サクラブレンド・コーヒー", "2026春_新焙煎豆", "全国_20代女性", "春カフェハッシュタグ", "sakura_hashtag_1200x628.png"},
		{"サクラブレンド・コーヒー", "母の日_ギフト箱", "全国_ギフト検討層", "限定ラッピング", "sakura_gift_wrap_300x250.png"},
		{"サクラブレンド・コーヒー", "母の日_ギフト箱", "全国_ギフト検討層", "早割10日間", "sakura_early_bird_728x90.jpg"},

		{"蒼天エナジー", "夏の覚醒キャンペーン", "全国_ゲーム好き", "夜更かし応援バナー", "sora_gamer_night_300x250.png"},
		{"蒼天エナジー", "夏の覚醒キャンペーン", "全国_ゲーム好き", "eスポーツ大会協賛", "sora_esports_728x90.jpg"},
		{"蒼天エナジー", "夏の覚醒キャンペーン", "全国_ゲーム好き", "限定缶デザイン", "sora_limited_can_320x50.png"},
		{"蒼天エナジー", "夏の覚醒キャンペーン", "全国_ゲーム好き", "配信者タイアップ", "sora_streamer_1080x1920.jpg"},
		{"蒼天エナジー", "夏の覚醒キャンペーン", "関西_夜間配信視聴", "午前2時の相棒", "sora_kansai_late_300x250.png"},
		{"蒼天エナジー", "夏の覚醒キャンペーン", "関西_夜間配信視聴", "コンビニ2本セット", "sora_kansai_combo_728x90.jpg"},
		{"蒼天エナジー", "夏の覚醒キャンペーン", "関西_夜間配信視聴", "新フレーバー試飲", "sora_kansai_trial_320x50.png"},
		{"蒼天エナジー", "ゼロ糖ハイボール新味", "全国_健康志向", "糖質オフ訴求", "sora_zero_sugar_300x250.png"},
		{"蒼天エナジー", "ゼロ糖ハイボール新味", "全国_健康志向", "カロリー表示強調", "sora_zero_cal_728x90.jpg"},
		{"蒼天エナジー", "ゼロ糖ハイボール新味", "全国_健康志向", "ジム併設自販機", "sora_zero_gym_320x50.png"},

		{"ミルキーフロスト・アイス", "今夏の新作アイス", "家族層_週末", "週末4個入り箱", "frost_family_box_300x250.png"},
		{"ミルキーフロスト・アイス", "今夏の新作アイス", "家族層_週末", "子ども向けキャラ缶", "frost_kids_char_728x90.jpg"},
		{"ミルキーフロスト・アイス", "今夏の新作アイス", "家族層_週末", "アウトドアクーラー", "frost_outdoor_320x50.png"},
		{"ミルキーフロスト・アイス", "今夏の新作アイス", "家族層_週末", "スーパー店頭試食", "frost_instore_taste_1080x1920.jpg"},
		{"ミルキーフロスト・アイス", "今夏の新作アイス", "SNSインフルエンサー", "抹茶チョコ溶けない", "frost_matcha_300x250.png"},
		{"ミルキーフロスト・アイス", "今夏の新作アイス", "SNSインフルエンサー", "溶けない動画30秒", "frost_meltproof_728x90.jpg"},
		{"ミルキーフロスト・アイス", "今夏の新作アイス", "SNSインフルエンサー", "投稿キャンペーン", "frost_ugc_1200x628.png"},
		{"ミルキーフロスト・アイス", "今夏の新作アイス", "SNSインフルエンサー", "限定パッケージ", "frost_limited_pack_320x50.png"},
	}
}

type adsMetrics struct {
	impressions int64
	cost        float64
	conversions int64
}

func adsMetricsForIndex(i int) adsMetrics {
	impressions := int64(12000 + i*1750)
	return adsMetrics{
		impressions: impressions,
		cost:        float64(impressions) * 0.11,
		conversions: impressions/2500 + int64(i%4),
	}
}

func specToRow(s adsReportSpec, m adsMetrics) []interface{} {
	return []interface{}{
		s.account, s.campaign, s.group, s.ad, s.creative,
		m.impressions, m.cost, m.conversions,
	}
}

func writeAdsReportSimple(f *excelize.File, sheet string) error {
	specs := adsReportSpecs()
	if len(specs) > 25 {
		return fmt.Errorf("simple fixture: at most 25 data rows, got %d", len(specs))
	}
	if err := setSheetRowAt(f, sheet, 1, "A", adsReportLeafHeaders()); err != nil {
		return err
	}
	for i, s := range specs {
		if err := setSheetRowAt(f, sheet, i+2, "A", specToRow(s, adsMetricsForIndex(i))); err != nil {
			return err
		}
	}
	return applyAdsReportSheetFormat(f, sheet, sheetFormatOpts{
		headerFrom: "A1", headerTo: "H1",
		freezeYSplit: 1, topLeftCell: "A2",
	})
}

func writeAdsReportMultiRowHeaders(f *excelize.File, sheet string) error {
	groupRow := []interface{}{
		"配信設定", "配信設定", "配信設定", "配信設定", "クリエイティブ",
		"成果", "成果", "成果",
	}
	if err := setSheetRowAt(f, sheet, 1, "A", groupRow); err != nil {
		return err
	}
	if err := setSheetRowAt(f, sheet, 2, "A", adsReportLeafHeaders()); err != nil {
		return err
	}
	specs := adsReportSpecs()
	for i, s := range specs {
		if err := setSheetRowAt(f, sheet, i+3, "A", specToRow(s, adsMetricsForIndex(i))); err != nil {
			return err
		}
	}
	return applyAdsReportSheetFormat(f, sheet, sheetFormatOpts{
		headerFrom: "A1", headerTo: "H2",
		freezeYSplit: 2, topLeftCell: "A3",
	})
}

func writeAdsReportGroupedRows(f *excelize.File, sheet string) error {
	specs := adsReportSpecs()
	if err := setSheetRowAt(f, sheet, 1, "A", adsReportLeafHeaders()); err != nil {
		return err
	}
	for i, s := range specs {
		if err := setSheetRowAt(f, sheet, i+2, "A", specToRow(s, adsMetricsForIndex(i))); err != nil {
			return err
		}
	}
	dataStart := 2
	dataEnd := len(specs) + 1
	if err := mergeAdsReportGroupedRows(f, sheet, "A", "B", "C", dataStart, dataEnd, specs); err != nil {
		return err
	}
	return applyAdsReportSheetFormat(f, sheet, sheetFormatOpts{
		headerFrom: "A1", headerTo: "H1",
		freezeYSplit: 1, topLeftCell: "A2",
	})
}

func mergeAdsReportGroupedRows(
	f *excelize.File, sheet string,
	accountCol, campaignCol, groupCol string,
	startRow, endRow int, specs []adsReportSpec,
) error {
	if endRow < startRow {
		return nil
	}
	rowCount := endRow - startRow + 1
	accountKey := func(i int) string { return specs[i].account }
	campaignKey := func(i int) string { return specs[i].account + "\x00" + specs[i].campaign }
	groupKey := func(i int) string {
		return specs[i].account + "\x00" + specs[i].campaign + "\x00" + specs[i].group
	}
	if err := mergeSheetColumnByKey(f, sheet, accountCol, startRow, rowCount, accountKey); err != nil {
		return err
	}
	if err := mergeSheetColumnByKey(f, sheet, campaignCol, startRow, rowCount, campaignKey); err != nil {
		return err
	}
	return mergeSheetColumnByKey(f, sheet, groupCol, startRow, rowCount, groupKey)
}

func mergeSheetColumnByKey(f *excelize.File, sheet, col string, startRow, rowCount int, key func(int) string) error {
	if rowCount == 0 {
		return nil
	}
	groupStart := startRow
	currentKey := key(0)
	for i := 1; i <= rowCount; i++ {
		if i == rowCount || key(i) != currentKey {
			groupEnd := startRow + i - 1
			if groupEnd > groupStart {
				from, err := excelize.JoinCellName(col, groupStart)
				if err != nil {
					return err
				}
				to, err := excelize.JoinCellName(col, groupEnd)
				if err != nil {
					return err
				}
				if err := f.MergeCell(sheet, from, to); err != nil {
					return err
				}
			}
			if i < rowCount {
				groupStart = startRow + i
				currentKey = key(i)
			}
		}
	}
	return nil
}

func writeAdsReportSummaryRows(f *excelize.File, sheet string) error {
	specs := adsReportSpecs()
	if err := setSheetRowAt(f, sheet, 1, "A", adsReportLeafHeaders()); err != nil {
		return err
	}
	var sumImp int64
	var sumCost float64
	var sumConv int64
	for i, s := range specs {
		m := adsMetricsForIndex(i)
		sumImp += m.impressions
		sumCost += m.cost
		sumConv += m.conversions
		if err := setSheetRowAt(f, sheet, i+2, "A", specToRow(s, m)); err != nil {
			return err
		}
	}
	summaryRow := len(specs) + 2
	if err := setSheetRowAt(f, sheet, summaryRow, "A", []interface{}{
		"合計", "", "", "", "", sumImp, sumCost, sumConv,
	}); err != nil {
		return err
	}
	if err := styleSummaryRow(f, sheet, summaryRow, "A", "H"); err != nil {
		return err
	}
	return applyAdsReportSheetFormat(f, sheet, sheetFormatOpts{
		headerFrom: "A1", headerTo: "H1",
		freezeYSplit: 1, topLeftCell: "A2",
	})
}

func writeAdsReportComplex(f *excelize.File, sheet string) error {
	const startCol = "D"
	const headerTopRow = 5
	const leafHeaderRow = 6
	const dataStartRow = 7

	if err := f.SetCellValue(sheet, "D5", "配信設定"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "D5", "G5"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "H5", "クリエイティブ"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "H5", "H6"); err != nil {
		return err
	}
	if err := f.SetCellValue(sheet, "I5", "成果"); err != nil {
		return err
	}
	if err := f.MergeCell(sheet, "I5", "K5"); err != nil {
		return err
	}
	if err := setSheetRowAt(f, sheet, leafHeaderRow, startCol, adsReportLeafHeaders()); err != nil {
		return err
	}

	specs := adsReportSpecs()
	var sumImp int64
	var sumCost float64
	var sumConv int64
	for i, s := range specs {
		m := adsMetricsForIndex(i)
		sumImp += m.impressions
		sumCost += m.cost
		sumConv += m.conversions
		if err := setSheetRowAt(f, sheet, dataStartRow+i, startCol, specToRow(s, m)); err != nil {
			return err
		}
	}
	dataEnd := dataStartRow + len(specs) - 1
	if err := mergeAdsReportGroupedRows(f, sheet, "D", "E", "F", dataStartRow, dataEnd, specs); err != nil {
		return err
	}

	summaryRow := dataEnd + 1
	if err := setSheetRowAt(f, sheet, summaryRow, startCol, []interface{}{
		"合計", "", "", "", "", sumImp, sumCost, sumConv,
	}); err != nil {
		return err
	}
	if err := styleSummaryRow(f, sheet, summaryRow, "D", "K"); err != nil {
		return err
	}
	return applyAdsReportSheetFormat(f, sheet, sheetFormatOpts{
		tableStartCol: "D",
		headerFrom:    "D5", headerTo: "K6",
		freezeXSplit: 3, freezeYSplit: leafHeaderRow,
		topLeftCell: "D7",
	})
}

func writeAdsReportOriginD5(f *excelize.File, sheet string) error {
	const startCol = "D"
	const headerRow = 5

	if err := setSheetRowAt(f, sheet, headerRow, startCol, adsReportLeafHeaders()); err != nil {
		return err
	}
	specs := adsReportSpecs()
	for i, s := range specs {
		if err := setSheetRowAt(f, sheet, headerRow+1+i, startCol, specToRow(s, adsMetricsForIndex(i))); err != nil {
			return err
		}
	}
	return applyAdsReportSheetFormat(f, sheet, sheetFormatOpts{
		tableStartCol: "D",
		headerFrom:    "D5", headerTo: "K5",
		freezeXSplit: 3, freezeYSplit: headerRow,
		topLeftCell: "D6",
	})
}

func setSheetRowAt(f *excelize.File, sheet string, row int, startCol string, values []interface{}) error {
	cell, err := excelize.JoinCellName(startCol, row)
	if err != nil {
		return err
	}
	return f.SetSheetRow(sheet, cell, &values)
}

type sheetFormatOpts struct {
	tableStartCol string // first column of the 8-column table (A or D)
	headerFrom    string
	headerTo      string
	freezeXSplit  int
	freezeYSplit  int
	topLeftCell   string
}

// Widths for 広告アカウント … Conversion (Excel character units).
var adsReportTableColWidths = []float64{26, 22, 21, 22, 35, 18, 12, 12}

func setAdsReportTableColumnWidths(f *excelize.File, sheet, startCol string) error {
	start, err := excelize.ColumnNameToNumber(startCol)
	if err != nil {
		return err
	}
	for i, width := range adsReportTableColWidths {
		col, err := excelize.ColumnNumberToName(start + i)
		if err != nil {
			return err
		}
		if err := f.SetColWidth(sheet, col, col, width); err != nil {
			return err
		}
	}
	return nil
}

func applyAdsReportSheetFormat(f *excelize.File, sheet string, opts sheetFormatOpts) error {
	startCol := opts.tableStartCol
	if startCol == "" {
		startCol = "A"
	}
	if err := setAdsReportTableColumnWidths(f, sheet, startCol); err != nil {
		return err
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#D9E1F2"},
		},
		Font: &excelize.Font{Bold: true},
	})
	if err != nil {
		return err
	}
	if err := f.SetCellStyle(sheet, opts.headerFrom, opts.headerTo, headerStyle); err != nil {
		return err
	}

	return f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		XSplit:      opts.freezeXSplit,
		YSplit:      opts.freezeYSplit,
		TopLeftCell: opts.topLeftCell,
		ActivePane:  "bottomLeft",
	})
}

func styleSummaryRow(f *excelize.File, sheet string, row int, fromCol, toCol string) error {
	style, err := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#FFF2CC"},
		},
		Font: &excelize.Font{Bold: true},
	})
	if err != nil {
		return err
	}
	from, err := excelize.JoinCellName(fromCol, row)
	if err != nil {
		return err
	}
	to, err := excelize.JoinCellName(toCol, row)
	if err != nil {
		return err
	}
	return f.SetCellStyle(sheet, from, to, style)
}
