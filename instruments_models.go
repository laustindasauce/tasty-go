package tasty

import (
	"time"
)

type DestinationVenueSymbol struct {
	ID                   int    `json:"id"`
	Symbol               string `json:"symbol"`
	DestinationVenue     string `json:"destination-venue"`
	MaxQuantityPrecision int    `json:"max-quantity-precision"`
	MaxPricePrecision    int    `json:"max-price-precision"`
	Routable             bool   `json:"routable"`
}

type CryptocurrencyInfo struct {
	ID                      int                      `json:"id"`
	Symbol                  string                   `json:"symbol"`
	InstrumentType          InstrumentType           `json:"instrument-type"`
	ShortDescription        string                   `json:"short-description"`
	Description             string                   `json:"description"`
	IsClosingOnly           bool                     `json:"is-closing-only"`
	Active                  bool                     `json:"active"`
	TickSize                StringToFloat32          `json:"tick-size"`
	StreamerSymbol          string                   `json:"streamer-symbol"`
	DestinationVenueSymbols []DestinationVenueSymbol `json:"destination-venue-symbols"`
}

type TickSize struct {
	Value     StringToFloat32 `json:"value"`
	Threshold StringToFloat32 `json:"threshold"`
	Symbol    string          `json:"symbol"`
}

type Equity struct {
	ID                             int             `json:"id"`
	Symbol                         string          `json:"symbol"`
	InstrumentType                 InstrumentType  `json:"instrument-type"`
	Cusip                          string          `json:"cusip"`
	ShortDescription               string          `json:"short-description"`
	IsIndex                        bool            `json:"is-index"`
	ListedMarket                   string          `json:"listed-market"`
	Description                    string          `json:"description"`
	Lendability                    string          `json:"lendability"`
	BorrowRate                     StringToFloat32 `json:"borrow-rate"`
	HaltedAt                       string          `json:"halted-at"`
	StopsTradingAt                 time.Time       `json:"stops-trading-at"`
	MarketTimeInstrumentCollection string          `json:"market-time-instrument-collection"`
	IsClosingOnly                  bool            `json:"is-closing-only"`
	IsOptionsClosingOnly           bool            `json:"is-options-closing-only"`
	Active                         bool            `json:"active"`
	IsFractionalQuantityEligible   bool            `json:"is-fractional-quantity-eligible"`
	IsIlliquid                     bool            `json:"is-illiquid"`
	IsEtf                          bool            `json:"is-etf"`
	StreamerSymbol                 string          `json:"streamer-symbol"`
	TickSizes                      []TickSize      `json:"tick-sizes"`
	OptionTickSizes                []TickSize      `json:"option-tick-sizes"`
}

type EquityOption struct {
	Symbol                         string          `json:"symbol"`
	InstrumentType                 InstrumentType  `json:"instrument-type"`
	Active                         bool            `json:"active"`
	ListedMarket                   string          `json:"listed-market"`
	StrikePrice                    StringToFloat32 `json:"strike-price"`
	RootSymbol                     string          `json:"root-symbol"`
	UnderlyingSymbol               string          `json:"underlying-symbol"`
	ExpirationDate                 string          `json:"expiration-date"`
	ExerciseStyle                  string          `json:"exercise-style"`
	SharesPerContract              int             `json:"shares-per-contract"`
	OptionType                     OptionType      `json:"option-type"`
	OptionChainType                string          `json:"option-chain-type"`
	ExpirationType                 string          `json:"expiration-type"`
	SettlementType                 string          `json:"settlement-type"`
	HaltedAt                       string          `json:"halted-at"`
	StopsTradingAt                 time.Time       `json:"stops-trading-at"`
	MarketTimeInstrumentCollection string          `json:"market-time-instrument-collection"`
	DaysToExpiration               int             `json:"days-to-expiration"`
	ExpiresAt                      time.Time       `json:"expires-at"`
	IsClosingOnly                  bool            `json:"is-closing-only"`
	OldSecurityNumber              string          `json:"old-security-number"`
	StreamerSymbol                 string          `json:"streamer-symbol"`
}

type FutureETFEquivalent struct {
	Symbol        string `json:"symbol"`
	ShareQuantity int    `json:"share-quantity"`
}

type Roll struct {
	Name               string `json:"name"`
	ActiveCount        int    `json:"active-count"`
	CashSettled        bool   `json:"cash-settled"`
	BusinessDaysOffset int    `json:"business-days-offset"`
	FirstNotice        bool   `json:"first-notice"`
}

type Future struct {
	Symbol                       string              `json:"symbol"`
	ProductCode                  string              `json:"product-code"`
	ContractSize                 StringToFloat32     `json:"contract-size"`
	TickSize                     StringToFloat32     `json:"tick-size"`
	NotionalMultiplier           StringToFloat32     `json:"notional-multiplier"`
	MainFraction                 StringToFloat32     `json:"main-fraction"`
	SubFraction                  StringToFloat32     `json:"sub-fraction"`
	DisplayFactor                StringToFloat32     `json:"display-factor"`
	LastTradeDate                string              `json:"last-trade-date"`
	ExpirationDate               string              `json:"expiration-date"`
	ClosingOnlyDate              string              `json:"closing-only-date"`
	Active                       bool                `json:"active"`
	ActiveMonth                  bool                `json:"active-month"`
	NextActiveMonth              bool                `json:"next-active-month"`
	IsClosingOnly                bool                `json:"is-closing-only"`
	FirstNoticeDate              string              `json:"first-notice-date"`
	StopsTradingAt               time.Time           `json:"stops-trading-at"`
	ExpiresAt                    time.Time           `json:"expires-at"`
	ProductGroup                 string              `json:"product-group"`
	Exchange                     string              `json:"exchange"`
	RollTargetSymbol             string              `json:"roll-target-symbol"`
	StreamerExchangeCode         string              `json:"streamer-exchange-code"`
	StreamerSymbol               string              `json:"streamer-symbol"`
	BackMonthFirstCalendarSymbol bool                `json:"back-month-first-calendar-symbol"`
	IsTradeable                  bool                `json:"is-tradeable"`
	TrueUnderlyingSymbol         string              `json:"true-underlying-symbol"`
	FutureETFEquivalent          FutureETFEquivalent `json:"future-etf-equivalent"`
	FutureProduct                FutureProduct       `json:"future-product"`
	TickSizes                    []TickSize          `json:"tick-sizes"`
	OptionTickSizes              []TickSize          `json:"option-tick-sizes"`
	SpreadTickSizes              []TickSize          `json:"spread-tick-sizes"`
}

type FutureOptionProduct struct {
	RootSymbol              string          `json:"root-symbol"`
	CashSettled             bool            `json:"cash-settled"`
	Code                    string          `json:"code"`
	LegacyCode              string          `json:"legacy-code"`
	ClearportCode           string          `json:"clearport-code"`
	ClearingCode            string          `json:"clearing-code"`
	ClearingExchangeCode    string          `json:"clearing-exchange-code"`
	ClearingPriceMultiplier StringToFloat32 `json:"clearing-price-multiplier"`
	DisplayFactor           StringToFloat32 `json:"display-factor"`
	Exchange                string          `json:"exchange"`
	ProductType             string          `json:"product-type"`
	ExpirationType          string          `json:"expiration-type"`
	SettlementDelayDays     int             `json:"settlement-delay-days"`
	IsRollover              bool            `json:"is-rollover"`
	ProductSubtype          string          `json:"product-subtype"`
	MarketSector            string          `json:"market-sector"`
	FutureProduct           FutureProduct   `json:"future-product"`
}

type FutureProduct struct {
	RootSymbol                   string                `json:"root-symbol"`
	Code                         string                `json:"code"`
	Description                  string                `json:"description"`
	ClearingCode                 string                `json:"clearing-code"`
	ClearingExchangeCode         string                `json:"clearing-exchange-code"`
	ClearportCode                string                `json:"clearport-code"`
	LegacyCode                   string                `json:"legacy-code"`
	Exchange                     string                `json:"exchange"`
	LegacyExchangeCode           string                `json:"legacy-exchange-code"`
	ProductType                  string                `json:"product-type"`
	ListedMonths                 []string              `json:"listed-months"`
	ActiveMonths                 []string              `json:"active-months"`
	NotionalMultiplier           StringToFloat32       `json:"notional-multiplier"`
	TickSize                     StringToFloat32       `json:"tick-size"`
	DisplayFactor                StringToFloat32       `json:"display-factor"`
	BaseTick                     int                   `json:"base-tick"`
	SubTick                      int                   `json:"sub-tick"`
	StreamerExchangeCode         string                `json:"streamer-exchange-code"`
	SmallNotional                bool                  `json:"small-notional"`
	BackMonthFirstCalendarSymbol bool                  `json:"back-month-first-calendar-symbol"`
	FirstNotice                  bool                  `json:"first-notice"`
	CashSettled                  bool                  `json:"cash-settled"`
	ContractLimit                int                   `json:"contract-limit"`
	SecurityGroup                string                `json:"security-group"`
	ProductSubtype               string                `json:"product-subtype"`
	TrueUnderlyingCode           string                `json:"true-underlying-code"`
	MarketSector                 string                `json:"market-sector"`
	OptionProducts               []FutureOptionProduct `json:"option-products"`
	Roll                         Roll                  `json:"roll"`
}

type FutureOption struct {
	Symbol               string              `json:"symbol"`
	UnderlyingSymbol     string              `json:"underlying-symbol"`
	ProductCode          string              `json:"product-code"`
	ExpirationDate       string              `json:"expiration-date"`
	RootSymbol           string              `json:"root-symbol"`
	OptionRootSymbol     string              `json:"option-root-symbol"`
	StrikePrice          StringToFloat32     `json:"strike-price"`
	Exchange             string              `json:"exchange"`
	ExchangeSymbol       string              `json:"exchange-symbol"`
	StreamerSymbol       string              `json:"streamer-symbol"`
	OptionType           OptionType          `json:"option-type"`
	ExerciseStyle        string              `json:"exercise-style"`
	IsVanilla            bool                `json:"is-vanilla"`
	IsPrimaryDeliverable bool                `json:"is-primary-deliverable"`
	FuturePriceRatio     StringToFloat32     `json:"future-price-ratio"`
	Multiplier           StringToFloat32     `json:"multiplier"`
	UnderlyingCount      StringToFloat32     `json:"underlying-count"`
	IsConfirmed          bool                `json:"is-confirmed"`
	NotionalValue        StringToFloat32     `json:"notional-value"`
	DisplayFactor        StringToFloat32     `json:"display-factor"`
	SecurityExchange     string              `json:"security-exchange"`
	SxID                 string              `json:"sx-id"`
	SettlementType       string              `json:"settlement-type"`
	StrikeFactor         StringToFloat32     `json:"strike-factor"`
	MaturityDate         string              `json:"maturity-date"`
	IsExercisableWeekly  bool                `json:"is-exercisable-weekly"`
	LastTradeTime        string              `json:"last-trade-time"`
	DaysToExpiration     int                 `json:"days-to-expiration"`
	IsClosingOnly        bool                `json:"is-closing-only"`
	Active               bool                `json:"active"`
	StopsTradingAt       time.Time           `json:"stops-trading-at"`
	ExpiresAt            time.Time           `json:"expires-at"`
	FutureOptionProduct  FutureOptionProduct `json:"future-option-product"`
}

type QuantityDecimalPrecision struct {
	InstrumentType            InstrumentType `json:"instrument-type"`
	Symbol                    string         `json:"symbol"`
	Value                     int            `json:"value"`
	MinimumIncrementPrecision int            `json:"minimum-increment-precision"`
}

type Warrant struct {
	Symbol         string         `json:"symbol"`
	InstrumentType InstrumentType `json:"instrument-type"`
	Cusip          string         `json:"cusip"`
	ListedMarket   string         `json:"listed-market"`
	Description    string         `json:"description"`
	IsClosingOnly  bool           `json:"is-closing-only"`
	Active         bool           `json:"active"`
}

type Strike struct {
	StrikePrice        StringToFloat32 `json:"strike-price"`
	Call               string          `json:"call"`
	CallStreamerSymbol string          `json:"call-streamer-symbol"`
	Put                string          `json:"put"`
	PutStreamerSymbol  string          `json:"put-streamer-symbol"`
}

type FuturesExpiration struct {
	UnderlyingSymbol     string          `json:"underlying-symbol"`
	RootSymbol           string          `json:"root-symbol"`
	OptionRootSymbol     string          `json:"option-root-symbol"`
	OptionContractSymbol string          `json:"option-contract-symbol"`
	Asset                string          `json:"asset"`
	ExpirationDate       string          `json:"expiration-date"`
	DaysToExpiration     int             `json:"days-to-expiration"`
	ExpirationType       string          `json:"expiration-type"`
	SettlementType       string          `json:"settlement-type"`
	NotionalValue        StringToFloat32 `json:"notional-value"`
	DisplayFactor        StringToFloat32 `json:"display-factor"`
	StrikeFactor         StringToFloat32 `json:"strike-factor"`
	StopsTradingAt       time.Time       `json:"stops-trading-at"`
	ExpiresAt            time.Time       `json:"expires-at"`
	TickSizes            []TickSize      `json:"tick-sizes"`
	Strikes              []Strike        `json:"strikes"`
}

type Expiration struct {
	ExpirationType   string   `json:"expiration-type"`
	ExpirationDate   string   `json:"expiration-date"`
	DaysToExpiration int      `json:"days-to-expiration"`
	SettlementType   string   `json:"settlement-type"`
	Strikes          []Strike `json:"strikes"`
}

type NestedFuture struct {
	Symbol           string    `json:"symbol"`
	RootSymbol       string    `json:"root-symbol"`
	ExpirationDate   string    `json:"expiration-date"`
	DaysToExpiration int       `json:"days-to-expiration"`
	ActiveMonth      bool      `json:"active-month"`
	NextActiveMonth  bool      `json:"next-active-month"`
	StopsTradingAt   time.Time `json:"stops-trading-at"`
	ExpiresAt        time.Time `json:"expires-at"`
}

type OptionChains struct {
	UnderlyingSymbol string              `json:"underlying-symbol"`
	RootSymbol       string              `json:"root-symbol"`
	ExerciseStyle    string              `json:"exercise-style"`
	Expirations      []FuturesExpiration `json:"expirations"`
}

type NestedFuturesOptionChains struct {
	Futures      []NestedFuture `json:"futures"`
	OptionChains []OptionChains `json:"option-chains"`
}

type Deliverable struct {
	ID              int             `json:"id"`
	RootSymbol      string          `json:"root-symbol"`
	DeliverableType string          `json:"deliverable-type"`
	Description     string          `json:"description"`
	Amount          StringToFloat32 `json:"amount"`
	Symbol          string          `json:"symbol"`
	InstrumentType  InstrumentType  `json:"instrument-type"`
	Percent         StringToFloat32 `json:"percent"`
}

type NestedOptionChains struct {
	UnderlyingSymbol  string        `json:"underlying-symbol"`
	RootSymbol        string        `json:"root-symbol"`
	OptionChainType   string        `json:"option-chain-type"`
	SharesPerContract int           `json:"shares-per-contract"`
	TickSizes         []TickSize    `json:"tick-sizes"`
	Deliverables      []Deliverable `json:"deliverables"`
	Expirations       []Expiration  `json:"expirations"`
}

type CompactOptionChains struct {
	UnderlyingSymbol  string        `json:"underlying-symbol"`
	RootSymbol        string        `json:"root-symbol"`
	OptionChainType   string        `json:"option-chain-type"`
	SettlementType    string        `json:"settlement-type"`
	SharesPerContract int           `json:"shares-per-contract"`
	ExpirationType    string        `json:"expiration-type"`
	Deliverables      []Deliverable `json:"deliverables"`
	Symbols           []string      `json:"symbols"`
	StreamerSymbols   []string      `json:"streamer-symbols"`
}

type ActiveEquitiesQuery struct {
	// PerPage equities to return per page. Default: 1000
	PerPage int `url:"per-page"`
	// PageOffset defaults to 0
	PageOffset int `url:"page-offset"`
	// Lendability Available values : Easy To Borrow, Locate Required, Preborrow
	Lendability Lendability `url:"lendability"`
}

type EquitiesQuery struct {
	// The symbols of the equity(s), i.e AAPL
	Symbols []string `url:"symbol[]"`
	// Available values : Easy To Borrow, Locate Required, Preborrow
	Lendability Lendability `url:"lendability"`
	// Flag indicating if equity is an index instrument
	IsIndex bool `url:"is-index"`
	// Flag indicating if equity is an etf instrument
	IsETF bool `url:"is-etf"`
}

type EquityOptionsQuery struct {
	// The symbol(s) of the equity option(s) using OCC Symbology, i.e. [FB 180629C00200000]
	Symbols []string `url:"symbol[]"`
	// Whether an option is available for trading with the broker.
	// Terminology is somewhat misleading as this is generally used
	// to filter non-standard / flex options out.
	Active bool `url:"active"`
	// Include expired options
	WithExpired bool `url:"with-expired"`
}

type FuturesQuery struct {
	// The symbol(s) of the future(s), i.e. symbol[]=ESZ9. Leading forward slash is not required.
	Symbols []string `url:"symbol[],omitempty"`
	// The product code of the future(s), i.e. ES or 6A
	// Ignored if Symbols parameter is given
	ProductCode []string `url:"product-code[],omitempty"`
}

type FutureOptionsQuery struct {
	// The symbol(s) of the future(s), i.e. symbol[]=ESZ9. Leading forward slash is not required.
	Symbols []string `url:"symbol[]"`
	// Future option root, i.e. EW3 or SO
	OptionRootSymbol string `url:"option-root-symbol,omitempty"`
	// Expiration date
	ExpirationDate time.Time `layout:"2006-01-02" url:"expiration-date,omitempty"`
	// P(ut) or C(all)
	OptionType OptionType `url:"option-type,omitempty"`
	// Strike price using display factor
	StrikePrice float32 `url:"strike-price,omitempty"`
}
