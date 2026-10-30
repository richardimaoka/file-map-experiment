package shinoken

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Parameter is one row from the parameters sheet (params.json).
type Parameter struct {
	GA                string `json:"GA"`
	LPName            string `json:"LP名"`
	No                string `json:"No"`
	URL               string `json:"URL"`
	Campaign          string `json:"キャンペーン"`
	CampaignSeq       string `json:"キャンペーン連番"`
	URLWithParam      string `json:"パラメータ付きURL　入稿※に必要なURL"`
	Media             string `json:"メディア"`
	AdGroup           string `json:"広告グループ"`
	AdGroupSeq        string `json:"広告グループ連番"`
	AdName            string `json:"広告名"`
	AdNameSeq         string `json:"広告名連番"`
	AdManagementTool  string `json:"広告管理ツール"`
	SquareImageName   string `json:"画像名(正方形)"`
	VerticalImageName string `json:"画像名(縦長)"`
}

// UnmarshalJSONは json.Unmarshaler インターフェースの実装です。
// json.Unmarshal 実行時に自動的に呼び出され、日本語空白文字の揺らぎを吸収してマッピングします。
func (p *Parameter) UnmarshalJSON(data []byte) error {
	// 一度 map[string]interface{}　としてデコード
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Helper: map空文字列を安全に取得する関数
	getStr := func(key string) string {
		if v, ok := raw[key]; ok && v != nil {
			if str, ok := v.(string); ok {
				return str
			}
		}
		return ""
	}

	// 通常のキーをマッピング
	p.GA = getStr("GA")
	p.LPName = getStr("LP名")
	p.No = getStr("No")
	p.URL = getStr("URL")
	p.Campaign = getStr("キャンペーン")
	p.CampaignSeq = getStr("キャンペーン連番")
	p.Media = getStr("メディア")
	p.AdGroup = getStr("広告グループ")
	p.AdGroupSeq = getStr("広告グループ連番")
	p.AdName = getStr("広告名")
	p.AdNameSeq = getStr("広告名連番")
	p.AdManagementTool = getStr("広告管理ツール")
	p.SquareImageName = getStr("画像名(正方形)")
	p.VerticalImageName = getStr("画像名(縦長)")

	// 全角半角スペースなどの表記揺れがあるキーの揺らぎ吸収
	for k, v := range raw {
		if strings.HasPrefix(k, "パラメータ付きURL") {
			if str, ok := v.(string); ok {
				p.URLWithParam = str
				break
			}
		}
	}

	return nil
}

// MarshalJSONは json.Marshaler インターフェースの実装です。
// 入稿用の固定キー名で JSON を出力します（UnmarshalJSON で吸収する表記揺れは正規キーに統一）。
func (p *Parameter) MarshalJSON() ([]byte, error) {
	raw := map[string]string{
		"GA":       p.GA,
		"LP名":      p.LPName,
		"No":       p.No,
		"URL":      p.URL,
		"キャンペーン":   p.Campaign,
		"キャンペーン連番": p.CampaignSeq,
		"パラメータ付きURL　入稿※に必要なURL": p.URLWithParam, // 元の全角スペース付きキー名で出力
		"メディア":     p.Media,
		"広告グループ":   p.AdGroup,
		"広告グループ連番": p.AdGroupSeq,
		"広告名":      p.AdName,
		"広告名連番":    p.AdNameSeq,
		"広告管理ツール":  p.AdManagementTool,
		"画像名(正方形)": p.SquareImageName,
		"画像名(縦長)":  p.VerticalImageName,
	}
	return json.Marshal(raw)
}

func ReadParametersFromJson(jsonPath string) ([]Parameter, error) {
	jsonRaw, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("ReadParametersFromJson: %v", err)
	}

	var params []Parameter
	// ここで json.Unmarshal を呼ぶと 内部で各要素の p.UnmarsalJSON() が自動実行される
	if err := json.Unmarshal(jsonRaw, &params); err != nil {
		return nil, fmt.Errorf("ReadParametersFromJson: failed to unmarshal JSON from %s, %v", jsonPath, err)
	}

	return params, nil
}
