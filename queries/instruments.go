package queries

import (
	"time"

	"github.com/austinbspencer/tasty-go/constants"
)

type ActiveEquities struct {
	// PerPage equities to return per page. Default: 1000
	PerPage int `url:"per-page"`
	// PageOffset defaults to 0
	PageOffset int `url:"page-offset"`
	// Lendability Available values : Easy To Borrow, Locate Required, Preborrow
	Lendability constants.Lendability `url:"lendability"`
}

type Equities struct {
	// The symbols of the equity(s), i.e AAPL
	Symbols []string `url:"symbol[]"`
	// Available values : Easy To Borrow, Locate Required, Preborrow
	Lendability constants.Lendability `url:"lendability"`
	// Flag indicating if equity is an index instrument
	IsIndex bool `url:"is-index"`
	// Flag indicating if equity is an etf instrument
	IsETF bool `url:"is-etf"`
}

type EquityOptions struct {
	// The symbol(s) of the equity option(s) using OCC Symbology, i.e. [FB 180629C00200000]
	Symbols []string `url:"symbol[]"`
	// Whether an option is available for trading with the broker.
	// Terminology is somewhat misleading as this is generally used
	// to filter non-standard / flex options out.
	Active bool `url:"active"`
	// Include expired options
	WithExpired bool `url:"with-expired"`
}

type Futures struct {
	// The symbol(s) of the future(s), i.e. symbol[]=ESZ9. Leading forward slash is not required.
	Symbols []string `url:"symbol[],omitempty"`
	// The product code of the future(s), i.e. ES or 6A
	// Ignored if Symbols parameter is given
	ProductCode []string `url:"product-code[],omitempty"`
}

type FutureOptions struct {
	// The symbol(s) of the future(s), i.e. symbol[]=ESZ9. Leading forward slash is not required.
	Symbols []string `url:"symbol[]"`
	// Future option root, i.e. EW3 or SO
	OptionRootSymbol string `url:"option-root-symbol,omitempty"`
	// Expiration date
	ExpirationDate time.Time `layout:"2006-01-02" url:"expiration-date,omitempty"`
	// P(ut) or C(all)
	OptionType constants.OptionType `url:"option-type,omitempty"`
	// Strike price using display factor
	StrikePrice float32 `url:"strike-price,omitempty"`
}
