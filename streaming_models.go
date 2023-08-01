package tasty

// Response from the API quote streamer request.
type QuoteStreamerTokenAuthResult struct {
	// API quote token unique to the customer identified by the session
	// Quote streamer tokens are valid for 24 hours.
	Token     string `json:"token"`
	DXLinkURL string `json:"dxlink-url"`
	Level     string `json:"level"`
}
