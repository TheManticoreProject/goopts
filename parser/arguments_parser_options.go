package parser

type ArgumentsParserOptions struct {
	ShowBannerOnHelp bool

	ShowBannerOnRun bool
}

// SetOptShowBannerOnHelp sets the option to show the banner on help.
//
// Parameters:
// - showBannerOnHelp: A boolean indicating whether to show the banner on help.
func (ap *ArgumentsParser) SetOptShowBannerOnHelp(showBannerOnHelp bool) {
	ap.Options.ShowBannerOnHelp = showBannerOnHelp
}

// SetOptShowBannerOnRun sets the option to show the banner on run.
//
// Parameters:
// - showBannerOnRun: A boolean indicating whether to show the banner on run.
func (ap *ArgumentsParser) SetOptShowBannerOnRun(showBannerOnRun bool) {
	ap.Options.ShowBannerOnRun = showBannerOnRun
}
