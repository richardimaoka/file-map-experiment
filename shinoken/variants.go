package shinoken

type MetaStoriesData struct {
	CR         string `json:"cr(画像・素材)"`
	WebSiteURL string `json:"ウェブサイトのURL"`
}

type MetaReelData struct {
	CR         string `json:"cr(画像・素材)"`
	WebSiteURL string `json:"ウェブサイトのURL"`
	Text       string `json:"テキスト"`
}

type Variant struct {
	// 配信面カスタマイズ
	FBStories        MetaStoriesData `json:"Facebook ストーリーズ"`
	FBReels          MetaStoriesData `json:"Facebook リール"`
	IGStories        MetaStoriesData `json:"Instagram ストーリーズ"`
	IGReels          MetaStoriesData `json:"Instagram リール"`
	MessengerStories MetaStoriesData `json:"Messener ストーリーズ"`

	// 基本情報・アカウント
	CampaignName string `json:"キャンペーン名"`
	AdSetName    string `json:"広告セット名"`
	AdName       string `json:"広告名"`

	// 素材・クリエイティブ設定
	CRFileName  string `json:"cr(画像・動画ファイル名)"`
	AspectRatio string `json:"素材アスペクト比"`

	// テキスト情報
	Text        string `json:"テキスト（tx）"`
	Headline    string `json:"見出し（ds）"`
	Description string `json:"説明（ex）"`
	WebsiteURL  string `json:"ウェブサイトのURL"`
}
