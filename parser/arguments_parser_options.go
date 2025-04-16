package parser

type ArgumentsParserOptions struct {
	ShowBannerOnHelp bool

	ShowBannerOnRun bool
}

func (ap *ArgumentsParser) SetOptShowBannerOnHelp(showBannerOnHelp bool) {
	ap.Options.ShowBannerOnHelp = showBannerOnHelp
}

func (ap *ArgumentsParser) SetOptShowBannerOnRun(showBannerOnRun bool) {
	ap.Options.ShowBannerOnRun = showBannerOnRun
}
