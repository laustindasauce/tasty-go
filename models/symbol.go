package models

// SymbolData model
type SymbolData struct {
	// Symbol is the stock ticker symbol
	Symbol string `json:"symbol"`
	// Description is the company name
	Description string `json:"description"`
}
