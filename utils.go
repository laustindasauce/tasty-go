package tasty

import (
	"fmt"
	"time"
)

// FutureSymbology is a struct to help build future contract codes
// Futures symbols always start with a slash followed by the contract code. The contract code
// consists of 3 things: the product code (A-Z), the month code, and a 1-2 digit number for the year.
// For a full list of futures products that we support, you can hit GET /instruments/future-products.
// Each item in the response has a code, which is the product code.
type FutureSymbology struct {
	// Don't include the / in the symbol i.e. ES
	ProductCode string
	MonthCode   MonthCode
	YearDigit   int
}

// Builds the future symbol into correct symbology.
func (fcc FutureSymbology) Build() string {
	return fmt.Sprintf("/%s%s%d", fcc.ProductCode, fcc.MonthCode, fcc.YearDigit)
}

// FutureOptionsSymbology is a struct to help build option symbol in correct Future Options Symbology
// Both the future product code and the future option product code will be present in a future option
// symbol. The future option symbol will always start with a ./, followed by the future contract code,
// the option contract code, the expiration date, option type (C/P), and strike price.
type FutureOptionsSymbology struct {
	OptionContractCode string
	// Should start with / (You can use the FutureSymbology struct's Build method)
	FutureContractCode string
	OptionType         OptionType
	Strike             int
	Expiration         time.Time
}

// Builds the future option into correct symbology.
func (foSym FutureOptionsSymbology) Build() string {
	codes := fmt.Sprintf(".%s %s", foSym.FutureContractCode, foSym.OptionContractCode)
	expiryString := foSym.Expiration.Format("060102")
	return fmt.Sprintf("%s %s%s%d", codes, expiryString, foSym.OptionType, foSym.Strike)
}

// ContainsInt returns whether or not the int exists in the slice.
func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// EquityOptionsSymbology is a struct to help build option symbol in correct OCC Symbology
// Root symbol of the underlying stock or ETF, padded with spaces to 6 characters.
// Expiration date, 6 digits in the format yymmdd. Option type, either P or C, for
// put or call.
type EquityOptionsSymbology struct {
	Symbol     string
	OptionType OptionType
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
