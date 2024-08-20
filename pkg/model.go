package pkg

type TTS struct {
	Text           string
	File           string
	VoiceArg       string
	Rate           string
	Volume         string
	Pitch          string
	WordsInCue     int
	WriteMedia     string
	WriteSubtitles string
	Proxy          string
	//Texts          []string
}

type VoiceTag struct {
	ContentCategories  []string `json:"contentCategories"`
	VoicePersonalities []string `json:"voicePersonalities"`
}

type Voice struct {
	Name           string   `json:"name"`
	ShortName      string   `json:"shortName"`
	Gender         string   `json:"gender"`
	Locale         string   `json:"locale"`
	SuggestedCodec string   `json:"suggestedCodec"`
	FriendlyName   string   `json:"friendlyName"`
	Status         string   `json:"status"`
	VoiceTag       VoiceTag `json:"voiceTag"`
}
