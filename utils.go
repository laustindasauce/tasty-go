package tasty

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var validMonthCodes = []string{string(January), string(February), string(March),
	string(April), string(May), string(June), string(July), string(August),
	string(September), string(October), string(November), string(December)}

func yearDigitFromFS(symbol string) (int, error) {
	var num string
	rns := []rune(symbol)
	for i := len(rns) - 1; i >= 0; i-- {
		if unicode.IsNumber(rns[i]) {
			num = string(symbol[i]) + num
		} else {
			// break after first non number
			break
		}
	}

	if num == "" {
		return 0, errors.New("missing valid year")
	}

	res, _ := strconv.Atoi(num)

	return res, nil
}

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

// Parse the future symbol into FutureSymbology struct.
func NewFSFromString(symbol string) (FutureSymbology, error) {
	var sym FutureSymbology

	if len(symbol) < 5 || len(symbol) > 6 {
		return sym, errors.New("invalid futures symbol")
	}

	if strings.Index(symbol, "/") != 0 {
		return sym, errors.New("future symbol must start with '/'")
	}

	year, err := yearDigitFromFS(symbol)
	if err != nil {
		return sym, err
	}

	sym.YearDigit = year

	idx := len(symbol) - len(fmt.Sprint(year))

	sym.MonthCode = MonthCode(symbol[idx-1 : idx])

	if !containsString(validMonthCodes, string(sym.MonthCode)) {
		return sym, fmt.Errorf("invalid month code: %s", sym.MonthCode)
	}

	sym.ProductCode = symbol[1 : idx-1]

	return sym, nil
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

// Parse the future options symbol into FutureOptionsSymbology struct.
func NewFOSFromString(symbol string) (FutureOptionsSymbology, error) {
	var sym FutureOptionsSymbology

	components := strings.Split(symbol, " ")

	if len(components) != 3 {
		return sym, errors.New("invalid future options symbol structure")
	}

	if len(components[2]) < 8 {
		return sym, errors.New("invalid future options symbol structure: strike info too short")
	}

	if strings.Index(symbol, "./") != 0 {
		return sym, errors.New("future options symbol must start with './'")
	}

	// pass first component without period
	futureContractCode, err := NewFSFromString(components[0][1:])
	if err != nil {
		return sym, fmt.Errorf("invalid contract code: %s", err.Error())
	}
	sym.FutureContractCode = futureContractCode.Build()

	// This includes option product code and expiry info
	sym.OptionContractCode = components[1]

	expiry, optionType, strike, err := getExpiryTypeStrike(components[2])
	if err != nil {
		return sym, err
	}

	sym.Expiration = expiry
	sym.OptionType = optionType
	sym.Strike = int(strike)

	return sym, nil
}

// containsInt returns whether or not the int exists in the slice.
func containsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// containsString returns whether or not the string exists in the slice.
func containsString(s []string, e string) bool {
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

// Parse occ symbol into EquityOptionsSymbology struct.
func NewOCCFromString(occSymbol string) (EquityOptionsSymbology, error) {
	var sym EquityOptionsSymbology
	if len(occSymbol) != 21 {
		return sym, errors.New("invalid occ symbol")
	}

	// Assuming valid symbol
	sym.Symbol = strings.ToUpper(strings.TrimSpace(occSymbol[:6]))

	expiry, optionType, strike, err := getExpiryTypeStrike(occSymbol[6:])
	if err != nil {
		return sym, err
	}

	sym.Expiration = expiry
	sym.OptionType = optionType
	sym.Strike = strike / float32(1000)

	return sym, nil
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

func getExpiryTypeStrike(symbol string) (time.Time, OptionType, float32, error) {
	var expiry time.Time
	var optionType OptionType
	var strike float32

	expiry, err := time.Parse("060102", symbol[:6])
	if err != nil {
		return expiry, optionType, strike, err
	}

	optionType = OptionType(symbol[6:7])
	if optionType != Call && optionType != Put {
		return expiry, optionType, strike, fmt.Errorf("unknown option type: %s", optionType)
	}

	strike64, err := strconv.ParseFloat(symbol[7:], 32)
	if err != nil {
		return expiry, optionType, strike, fmt.Errorf("invalid option strike: %s", symbol[7:])
	}

	strike = float32(strike64)

	return expiry, optionType, strike, nil
}
