package utils

import (
	"fmt"
	"time"

	"github.com/austinbspencer/tasty-go/constants"
)

// EquityOptionsSymbology is a struct to help build option symbol in correct OCC Symbology
// Root symbol of the underlying stock or ETF, padded with spaces to 6 characters.
// Expiration date, 6 digits in the format yymmdd. Option type, either P or C, for
// put or call.
type EquityOptionsSymbology struct {
	Symbol     string
	OptionType constants.OptionType
	Strike     float32
	Expiration time.Time
}

// Builds the equity option into correct symbology.
func (sym EquityOptionsSymbology) Build() string {
	expiryString := sym.Expiration.Format("060102")
	strikeString := getStrikeWithPadding(sym.Strike)
	symbol := getSymbolWithPadding(sym.Symbol)
	return fmt.Sprintf("%s%s%s%s", symbol, expiryString, sym.OptionType, strikeString)
}

// convert the strike into a string with correct padding.
func getStrikeWithPadding(strike float32) string {
	strikeString := fmt.Sprintf("%d", int(strike*1000))
	for len(strikeString) < 8 {
		strikeString = "0" + strikeString
	}
	return strikeString
}

// convert the symbol into a string with correct padding.
func getSymbolWithPadding(symbol string) string {
	for len(symbol) < 6 {
		symbol += " "
	}

	return symbol
}
