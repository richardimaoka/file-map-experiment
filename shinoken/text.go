package shinoken

import (
	"encoding/json"
	"fmt"
	"os"
)

type TextRequirement struct {
	ClickDefinition     string `json:"クリック定義"`
	Segment             string `json:"セグメント"`
	BodyTesxt           string `json:"テキスト (本文)"`
	Device              string `json:"デバイス"`
	MenuAndDeliverables string `json:"メニュー・入稿物"`
	MediaName           string `json:"媒体名"`
	AdType              string `json:"広告タイプ"`
	AdMaterial          string `json:"広告素材"`
	CharacterCount      string `json:"文字数"`
	CharacterCountLimit string `json:"文字数\n (規定)"`
	Material            string `json:"素材"`
	EndDate             string `json:"終了"`
	Headeline           string `json:"見出し"`
	ViewDefinition      string `json:"視聴定義"`
	Description         string `json:"説明"`
	BillingType         string `json:"課金形態"`
	LandingPageURLs     string `json:"遷移先"`
	DeliveryDesign      string `json:"配信設計"`
	LongHeadline        string `json:"長い見出し"`
	StartDate           string `json:"開始"`
}

func ReadTextReqsFromJson(jsonPath string) ([]TextRequirement, error) {
	jsonRaw, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("ReadTextReqsFromJson: %v", err)
	}

	var reqs []TextRequirement
	if err := json.Unmarshal(jsonRaw, &reqs); err != nil {
		return nil, fmt.Errorf("ReadTextReqsFromJson: failed to unmarshal JSON from %s, %v", jsonPath, err)
	}

	return reqs, nil
}
