package excel

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

// Writes excel/testdata/ads_report.xlsx with a flat ad-performance grid for read/build tests.
func TestCreateAdsReportExcelInTestdata(t *testing.T) {
	dir := filepath.Join("testdata")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, "ads_report.xlsx")
	rows := adsReportRows()
	dataRows := len(rows) - 1
	if dataRows < 1 || dataRows > 25 {
		t.Fatalf("expected 1–25 data rows, got %d", dataRows)
	}

	f := excelize.NewFile()
	defer f.Close()

	sheet := f.GetSheetName(0)
	for i, row := range rows {
		cell := fmt.Sprintf("A%d", i+1)
		if err := f.SetSheetRow(sheet, cell, &row); err != nil {
			t.Fatal(err)
		}
	}

	if err := applyAdsReportSheetFormat(f, sheet); err != nil {
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
		t.Fatal("ads_report.xlsx was not written")
	}
}

// Fictional 株式会社サクラブレンド group: coffee, energy drink, and ice dessert brands.
func adsReportRows() [][]interface{} {
	headers := []interface{}{
		"広告アカウント", "広告キャンペーン", "広告グループ", "広告", "クリエイティブ",
		"インプレッション", "費用", "Conversion",
	}
	rows := [][]interface{}{headers}

	specs := []struct {
		account  string
		campaign string
		group    string
		ad       string
		creative string
	}{
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

	for i, s := range specs {
		impressions := int64(12000 + i*1750)
		cost := float64(impressions) * 0.11
		conversions := impressions/2500 + int64(i%4)

		rows = append(rows, []interface{}{
			s.account, s.campaign, s.group, s.ad, s.creative,
			impressions, cost, conversions,
		})
	}

	return rows
}

func applyAdsReportSheetFormat(f *excelize.File, sheet string) error {
	// Excel column widths (character units) for ads_report.xlsx content.
	colWidths := []struct {
		col   string
		width float64
	}{
		{"A", 26}, // 広告アカウント (e.g. ミルキーフロスト・アイス)
		{"B", 22}, // 広告キャンペーン
		{"C", 21}, // 広告グループ
		{"D", 22}, // 広告
		{"E", 35}, // クリエイティブ file names
		{"F", 18}, // インプレッション
		{"G", 12}, // 費用
		{"H", 12}, // Conversion
	}
	for _, cw := range colWidths {
		if err := f.SetColWidth(sheet, cw.col, cw.col, cw.width); err != nil {
			return err
		}
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
	if err := f.SetCellStyle(sheet, "A1", "H1", headerStyle); err != nil {
		return err
	}

	return f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})
}
