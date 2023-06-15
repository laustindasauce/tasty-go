package utils

import (
	"fmt"
	"time"

	"github.com/austinbspencer/tasty-go/constants"
)

// FutureSymbology is a struct to help build future contract codes
// Futures symbols always start with a slash followed by the contract code. The contract code
// consists of 3 things: the product code (A-Z), the month code, and a 1-2 digit number for the year.
// For a full list of futures products that we support, you can hit GET /instruments/future-products.
// Each item in the response has a code, which is the product code.
type FutureSymbology struct {
	// Don't include the / in the symbol i.e. ES
	ProductCode string
	MonthCode   constants.MonthCode
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
	OptionType         constants.OptionType
	Strike             int
	Expiration         time.Time
}

// Builds the future option into correct symbology.
func (foSym FutureOptionsSymbology) Build() string {
	codes := fmt.Sprintf(".%s %s", foSym.FutureContractCode, foSym.OptionContractCode)
	expiryString := foSym.Expiration.Format("060102")
	return fmt.Sprintf("%s %s%s%d", codes, expiryString, foSym.OptionType, foSym.Strike)
}
