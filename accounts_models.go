package tasty

import (
	"time"

	"github.com/shopspring/decimal"
)

type Account struct {
	AccountNumber         string    `json:"account-number"`
	ExternalID            string    `json:"external-id"`
	OpenedAt              time.Time `json:"opened-at"`
	Nickname              string    `json:"nickname"`
	AccountTypeName       string    `json:"account-type-name"`
	DayTraderStatus       bool      `json:"day-trader-status"`
	ClosedAt              time.Time `json:"closed-at"`
	IsClosed              bool      `json:"is-closed"`
	IsFirmError           bool      `json:"is-firm-error"`
	IsFirmProprietary     bool      `json:"is-firm-proprietary"`
	IsFuturesApproved     bool      `json:"is-futures-approved"`
	IsTestDrive           bool      `json:"is-test-drive"`
	MarginOrCash          string    `json:"margin-or-cash"`
	IsForeign             bool      `json:"is-foreign"`
	FundingDate           string    `json:"funding-date"`
	InvestmentObjective   string    `json:"investment-objective"`
	LiquidityNeeds        string    `json:"liquidity-needs"`
	RiskTolerance         string    `json:"risk-tolerance"`
	InvestmentTimeHorizon string    `json:"investment-time-horizon"`
	FuturesAccountPurpose string    `json:"futures-account-purpose"`
	ExternalFDID          string    `json:"external-fdid"`
	SuitableOptionsLevel  string    `json:"suitable-options-level"`
	CreatedAt             time.Time `json:"created-at"`
	SubmittingUserID      string    `json:"submitting-user-id"`
	AuthorityLevel        string    `json:"authority-level"`
}

type AccountTradingStatus struct {
	AccountNumber                            string          `json:"account-number"`
	AutotradeAccountType                     string          `json:"autotrade-account-type"`
	ClearingAccountNumber                    string          `json:"clearing-account-number"`
	ClearingAggregationIdentifier            string          `json:"clearing-aggregation-identifier"`
	DayTradeCount                            int             `json:"day-trade-count"`
	EquitiesMarginCalculationType            string          `json:"equities-margin-calculation-type"`
	FeeScheduleName                          string          `json:"fee-schedule-name"`
	FuturesMarginRateMultiplier              decimal.Decimal `json:"futures-margin-rate-multiplier"`
	HasIntradayEquitiesMargin                bool            `json:"has-intraday-equities-margin"`
	ID                                       int             `json:"id"`
	IsAggregatedAtClearing                   bool            `json:"is-aggregated-at-clearing"`
	IsClosed                                 bool            `json:"is-closed"`
	IsClosingOnly                            bool            `json:"is-closing-only"`
	IsCryptocurrencyClosingOnly              bool            `json:"is-cryptocurrency-closing-only"`
	IsCryptocurrencyEnabled                  bool            `json:"is-cryptocurrency-enabled"`
	IsFrozen                                 bool            `json:"is-frozen"`
	IsFullEquityMarginRequired               bool            `json:"is-full-equity-margin-required"`
	IsFuturesClosingOnly                     bool            `json:"is-futures-closing-only"`
	IsFuturesIntraDayEnabled                 bool            `json:"is-futures-intra-day-enabled"`
	IsFuturesEnabled                         bool            `json:"is-futures-enabled"`
	IsInDayTradeEquityMaintenanceCall        bool            `json:"is-in-day-trade-equity-maintenance-call"`
	IsInMarginCall                           bool            `json:"is-in-margin-call"`
	IsPatternDayTrader                       bool            `json:"is-pattern-day-trader"`
	IsPortfolioMarginEnabled                 bool            `json:"is-portfolio-margin-enabled"`
	IsRiskReducingOnly                       bool            `json:"is-risk-reducing-only"`
	IsSmallNotionalFuturesIntraDayEnabled    bool            `json:"is-small-notional-futures-intra-day-enabled"`
	IsRollTheDayForwardEnabled               bool            `json:"is-roll-the-day-forward-enabled"`
	AreFarOtmNetOptionsRestricted            bool            `json:"are-far-otm-net-options-restricted"`
	OptionsLevel                             string          `json:"options-level"`
	PdtResetOn                               string          `json:"pdt-reset-on"`
	ShortCallsEnabled                        bool            `json:"short-calls-enabled"`
	SmallNotionalFuturesMarginRateMultiplier decimal.Decimal `json:"small-notional-futures-margin-rate-multiplier"`
	CMTAOverride                             int             `json:"cmta-override"`
	IsEquityOfferingEnabled                  bool            `json:"is-equity-offering-enabled"`
	IsEquityOfferingClosingOnly              bool            `json:"is-equity-offering-closing-only"`
	EnhancedFraudSafeguardsEnabledAt         time.Time       `json:"enhanced-fraud-safeguards-enabled-at"`
	UpdatedAt                                time.Time       `json:"updated-at"`
}

type AccountBalance struct {
	AccountNumber                      string          `json:"account-number"`
	CashBalance                        decimal.Decimal `json:"cash-balance"`
	LongEquityValue                    decimal.Decimal `json:"long-equity-value"`
	ShortEquityValue                   decimal.Decimal `json:"short-equity-value"`
	LongDerivativeValue                decimal.Decimal `json:"long-derivative-value"`
	ShortDerivativeValue               decimal.Decimal `json:"short-derivative-value"`
	LongFuturesValue                   decimal.Decimal `json:"long-futures-value"`
	ShortFuturesValue                  decimal.Decimal `json:"short-futures-value"`
	LongFuturesDerivativeValue         decimal.Decimal `json:"long-futures-derivative-value"`
	ShortFuturesDerivativeValue        decimal.Decimal `json:"short-futures-derivative-value"`
	LongMargineableValue               decimal.Decimal `json:"long-margineable-value"`
	ShortMargineableValue              decimal.Decimal `json:"short-margineable-value"`
	MarginEquity                       decimal.Decimal `json:"margin-equity"`
	EquityBuyingPower                  decimal.Decimal `json:"equity-buying-power"`
	DerivativeBuyingPower              decimal.Decimal `json:"derivative-buying-power"`
	DayTradingBuyingPower              decimal.Decimal `json:"day-trading-buying-power"`
	FuturesMarginRequirement           decimal.Decimal `json:"futures-margin-requirement"`
	AvailableTradingFunds              decimal.Decimal `json:"available-trading-funds"`
	MaintenanceRequirement             decimal.Decimal `json:"maintenance-requirement"`
	MaintenanceCallValue               decimal.Decimal `json:"maintenance-call-value"`
	RegTCallValue                      decimal.Decimal `json:"reg-t-call-value"`
	DayTradingCallValue                decimal.Decimal `json:"day-trading-call-value"`
	DayEquityCallValue                 decimal.Decimal `json:"day-equity-call-value"`
	NetLiquidatingValue                decimal.Decimal `json:"net-liquidating-value"`
	CashAvailableToWithdraw            decimal.Decimal `json:"cash-available-to-withdraw"`
	DayTradeExcess                     decimal.Decimal `json:"day-trade-excess"`
	PendingCash                        decimal.Decimal `json:"pending-cash"`
	PendingCashEffect                  PriceEffect     `json:"pending-cash-effect"`
	LongCryptocurrencyValue            decimal.Decimal `json:"long-cryptocurrency-value"`
	ShortCryptocurrencyValue           decimal.Decimal `json:"short-cryptocurrency-value"`
	CryptocurrencyMarginRequirement    decimal.Decimal `json:"cryptocurrency-margin-requirement"`
	UnsettledCryptocurrencyFiatAmount  decimal.Decimal `json:"unsettled-cryptocurrency-fiat-amount"`
	UnsettledCryptocurrencyFiatEffect  PriceEffect     `json:"unsettled-cryptocurrency-fiat-effect"`
	ClosedLoopAvailableBalance         decimal.Decimal `json:"closed-loop-available-balance"`
	EquityOfferingMarginRequirement    decimal.Decimal `json:"equity-offering-margin-requirement"`
	LongBondValue                      decimal.Decimal `json:"long-bond-value"`
	BondMarginRequirement              decimal.Decimal `json:"bond-margin-requirement"`
	SnapshotDate                       string          `json:"snapshot-date"`
	TimeOfDay                          string          `json:"time-of-day"`
	RegTMarginRequirement              decimal.Decimal `json:"reg-t-margin-requirement"`
	FuturesOvernightMarginRequirement  decimal.Decimal `json:"futures-overnight-margin-requirement"`
	FuturesIntradayMarginRequirement   decimal.Decimal `json:"futures-intraday-margin-requirement"`
	MaintenanceExcess                  decimal.Decimal `json:"maintenance-excess"`
	PendingMarginInterest              decimal.Decimal `json:"pending-margin-interest"`
	ApexStartingDayMarginEquity        decimal.Decimal `json:"apex-starting-day-margin-equity"`
	BuyingPowerAdjustment              decimal.Decimal `json:"buying-power-adjustment"`
	BuyingPowerAdjustmentEffect        PriceEffect     `json:"buying-power-adjustment-effect"`
	EffectiveCryptocurrencyBuyingPower decimal.Decimal `json:"effective-cryptocurrency-buying-power"`
	UpdatedAt                          time.Time       `json:"updated-at"`
}

type AccountPosition struct {
	AccountNumber                 string          `json:"account-number"`
	Symbol                        string          `json:"symbol"`
	InstrumentType                InstrumentType  `json:"instrument-type"`
	UnderlyingSymbol              string          `json:"underlying-symbol"`
	Quantity                      int             `json:"quantity"`
	QuantityDirection             Direction       `json:"quantity-direction"`
	ClosePrice                    decimal.Decimal `json:"close-price"`
	AverageOpenPrice              decimal.Decimal `json:"average-open-price"`
	AverageYearlyMarketClosePrice decimal.Decimal `json:"average-yearly-market-close-price"`
	AverageDailyMarketClosePrice  decimal.Decimal `json:"average-daily-market-close-price"`
	Mark                          decimal.Decimal `json:"mark"`
	MarkPrice                     decimal.Decimal `json:"mark-price"`
	Multiplier                    int             `json:"multiplier"`
	CostEffect                    PriceEffect     `json:"cost-effect"`
	IsSuppressed                  bool            `json:"is-suppressed"`
	IsFrozen                      bool            `json:"is-frozen"`
	RestrictedQuantity            int             `json:"restricted-quantity"`
	ExpiresAt                     time.Time       `json:"expires-at"`
	FixingPrice                   decimal.Decimal `json:"fixing-price"`
	DeliverableType               string          `json:"deliverable-type"`
	RealizedDayGain               decimal.Decimal `json:"realized-day-gain"`
	RealizedDayGainEffect         PriceEffect     `json:"realized-day-gain-effect"`
	RealizedDayGainDate           string          `json:"realized-day-gain-date"`
	RealizedToday                 decimal.Decimal `json:"realized-today"`
	RealizedTodayEffect           PriceEffect     `json:"realized-today-effect"`
	RealizedTodayDate             string          `json:"realized-today-date"`
	CreatedAt                     time.Time       `json:"created-at"`
	UpdatedAt                     time.Time       `json:"updated-at"`
}

type AccountBalanceSnapshots struct {
	AccountNumber                      string          `json:"account-number"`
	CashBalance                        decimal.Decimal `json:"cash-balance"`
	LongEquityValue                    decimal.Decimal `json:"long-equity-value"`
	ShortEquityValue                   decimal.Decimal `json:"short-equity-value"`
	LongDerivativeValue                decimal.Decimal `json:"long-derivative-value"`
	ShortDerivativeValue               decimal.Decimal `json:"short-derivative-value"`
	LongFuturesValue                   decimal.Decimal `json:"long-futures-value"`
	ShortFuturesValue                  decimal.Decimal `json:"short-futures-value"`
	LongFuturesDerivativeValue         decimal.Decimal `json:"long-futures-derivative-value"`
	ShortFuturesDerivativeValue        decimal.Decimal `json:"short-futures-derivative-value"`
	LongMargineableValue               decimal.Decimal `json:"long-margineable-value"`
	ShortMargineableValue              decimal.Decimal `json:"short-margineable-value"`
	MarginEquity                       decimal.Decimal `json:"margin-equity"`
	EquityBuyingPower                  decimal.Decimal `json:"equity-buying-power"`
	DerivativeBuyingPower              decimal.Decimal `json:"derivative-buying-power"`
	DayTradingBuyingPower              decimal.Decimal `json:"day-trading-buying-power"`
	FuturesMarginRequirement           decimal.Decimal `json:"futures-margin-requirement"`
	AvailableTradingFunds              decimal.Decimal `json:"available-trading-funds"`
	MaintenanceRequirement             decimal.Decimal `json:"maintenance-requirement"`
	MaintenanceCallValue               decimal.Decimal `json:"maintenance-call-value"`
	RegTCallValue                      decimal.Decimal `json:"reg-t-call-value"`
	DayTradingCallValue                decimal.Decimal `json:"day-trading-call-value"`
	DayEquityCallValue                 decimal.Decimal `json:"day-equity-call-value"`
	NetLiquidatingValue                decimal.Decimal `json:"net-liquidating-value"`
	CashAvailableToWithdraw            decimal.Decimal `json:"cash-available-to-withdraw"`
	DayTradeExcess                     decimal.Decimal `json:"day-trade-excess"`
	PendingCash                        decimal.Decimal `json:"pending-cash"`
	PendingCashEffect                  PriceEffect     `json:"pending-cash-effect"`
	LongCryptocurrencyValue            decimal.Decimal `json:"long-cryptocurrency-value"`
	ShortCryptocurrencyValue           decimal.Decimal `json:"short-cryptocurrency-value"`
	CryptocurrencyMarginRequirement    decimal.Decimal `json:"cryptocurrency-margin-requirement"`
	UnsettledCryptocurrencyFiatAmount  decimal.Decimal `json:"unsettled-cryptocurrency-fiat-amount"`
	UnsettledCryptocurrencyFiatEffect  PriceEffect     `json:"unsettled-cryptocurrency-fiat-effect"`
	ClosedLoopAvailableBalance         decimal.Decimal `json:"closed-loop-available-balance"`
	EquityOfferingMarginRequirement    decimal.Decimal `json:"equity-offering-margin-requirement"`
	LongBondValue                      decimal.Decimal `json:"long-bond-value"`
	BondMarginRequirement              decimal.Decimal `json:"bond-margin-requirement"`
	SnapshotDate                       string          `json:"snapshot-date"`
	RegTMarginRequirement              decimal.Decimal `json:"reg-t-margin-requirement"`
	FuturesOvernightMarginRequirement  decimal.Decimal `json:"futures-overnight-margin-requirement"`
	FuturesIntradayMarginRequirement   decimal.Decimal `json:"futures-intraday-margin-requirement"`
	MaintenanceExcess                  decimal.Decimal `json:"maintenance-excess"`
	PendingMarginInterest              decimal.Decimal `json:"pending-margin-interest"`
	EffectiveCryptocurrencyBuyingPower decimal.Decimal `json:"effective-cryptocurrency-buying-power"`
	UpdatedAt                          time.Time       `json:"updated-at"`
}

// PositionLimit model.
type PositionLimit struct {
	ID                          int    `json:"id"`
	AccountNumber               string `json:"account-number"`
	EquityOrderSize             int    `json:"equity-order-size"`
	EquityOptionOrderSize       int    `json:"equity-option-order-size"`
	FutureOrderSize             int    `json:"future-order-size"`
	FutureOptionOrderSize       int    `json:"future-option-order-size"`
	UnderlyingOpeningOrderLimit int    `json:"underlying-opening-order-limit"`
	EquityPositionSize          int    `json:"equity-position-size"`
	EquityOptionPositionSize    int    `json:"equity-option-position-size"`
	FuturePositionSize          int    `json:"future-position-size"`
	FutureOptionPositionSize    int    `json:"future-option-position-size"`
}

type AccountAuthorityDecorator struct {
	Account        Account `json:"account"`
	AuthorityLevel string  `json:"authority-level"`
}

type AccountType struct {
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	IsTaxAdvantaged     bool         `json:"is-tax-advantaged"`
	HasMultipleOwners   bool         `json:"has-multiple-owners"`
	IsPubliclyAvailable bool         `json:"is-publicly-available"`
	MarginTypes         []MarginType `json:"margin-types"`
}

// Optional filtering available for position endpoint.
type AccountPositionQuery struct {
	// UnderlyingSymbol An array of Underlying symbol(s) for positions
	UnderlyingSymbol []string `url:"underlying-symbol[],omitempty"`
	// Symbol A single symbol. Stock Ticker Symbol AAPL,
	// OCC Option Symbol AAPL 191004P00275000,
	// TW Future Symbol /ESZ9, or
	// TW Future Option Symbol ./ESZ9 EW4U9 190927P2975
	Symbol string `url:"symbol,omitempty"`
	// InstrumentType The type of Instrument
	InstrumentType InstrumentType `url:"instrument-type,omitempty"`
	// IncludeClosedPositions If closed positions should be included in the query
	IncludeClosedPositions bool `url:"include-closed-positions,omitempty"`
	// UnderlyingProductCode The underlying Future's Product code. i.e ES
	UnderlyingProductCode string `url:"underlying-product-code,omitempty"`
	// PartitionKeys Account partition keys
	PartitionKeys []string `url:"partition-keys[],omitempty"`
	// NetPositions Returns net positions grouped by instrument type and symbol
	NetPositions bool `url:"net-positions,omitempty"`
	// IncludeMarks Include current quote mark (note: can decrease performance)
	IncludeMarks bool `url:"include-marks,omitempty"`
}

// Filtering for the account balance snapshots endpoint.
type AccountBalanceSnapshotsQuery struct {
	// SnapshotDate The day of the balance snapshot to retrieve
	SnapshotDate time.Time `layout:"2006-01-02" url:"snapshot-date,omitempty"`
	// TimeOfDay The abbreviation for the time of day. Default value: EOD
	// Available values: EOD, BOD.
	TimeOfDay TimeOfDay `url:"time-of-day"`
}
