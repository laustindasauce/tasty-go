package models

import (
	"time"

	"github.com/austinbspencer/tasty-go/constants"
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
	FuturesMarginRateMultiplier              StringToFloat32 `json:"futures-margin-rate-multiplier"`
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
	SmallNotionalFuturesMarginRateMultiplier StringToFloat32 `json:"small-notional-futures-margin-rate-multiplier"`
	CMTAOverride                             int             `json:"cmta-override"`
	IsEquityOfferingEnabled                  bool            `json:"is-equity-offering-enabled"`
	IsEquityOfferingClosingOnly              bool            `json:"is-equity-offering-closing-only"`
	EnhancedFraudSafeguardsEnabledAt         time.Time       `json:"enhanced-fraud-safeguards-enabled-at"`
	UpdatedAt                                time.Time       `json:"updated-at"`
}

type AccountBalance struct {
	AccountNumber                      string          `json:"account-number"`
	CashBalance                        StringToFloat32 `json:"cash-balance"`
	LongEquityValue                    StringToFloat32 `json:"long-equity-value"`
	ShortEquityValue                   StringToFloat32 `json:"short-equity-value"`
	LongDerivativeValue                StringToFloat32 `json:"long-derivative-value"`
	ShortDerivativeValue               StringToFloat32 `json:"short-derivative-value"`
	LongFuturesValue                   StringToFloat32 `json:"long-futures-value"`
	ShortFuturesValue                  StringToFloat32 `json:"short-futures-value"`
	LongFuturesDerivativeValue         StringToFloat32 `json:"long-futures-derivative-value"`
	ShortFuturesDerivativeValue        StringToFloat32 `json:"short-futures-derivative-value"`
	LongMargineableValue               StringToFloat32 `json:"long-margineable-value"`
	ShortMargineableValue              StringToFloat32 `json:"short-margineable-value"`
	MarginEquity                       StringToFloat32 `json:"margin-equity"`
	EquityBuyingPower                  StringToFloat32 `json:"equity-buying-power"`
	DerivativeBuyingPower              StringToFloat32 `json:"derivative-buying-power"`
	DayTradingBuyingPower              StringToFloat32 `json:"day-trading-buying-power"`
	FuturesMarginRequirement           StringToFloat32 `json:"futures-margin-requirement"`
	AvailableTradingFunds              StringToFloat32 `json:"available-trading-funds"`
	MaintenanceRequirement             StringToFloat32 `json:"maintenance-requirement"`
	MaintenanceCallValue               StringToFloat32 `json:"maintenance-call-value"`
	RegTCallValue                      StringToFloat32 `json:"reg-t-call-value"`
	DayTradingCallValue                StringToFloat32 `json:"day-trading-call-value"`
	DayEquityCallValue                 StringToFloat32 `json:"day-equity-call-value"`
	NetLiquidatingValue                StringToFloat32 `json:"net-liquidating-value"`
	CashAvailableToWithdraw            StringToFloat32 `json:"cash-available-to-withdraw"`
	DayTradeExcess                     StringToFloat32 `json:"day-trade-excess"`
	PendingCash                        StringToFloat32 `json:"pending-cash"`
	PendingCashEffect                  string          `json:"pending-cash-effect"`
	LongCryptocurrencyValue            StringToFloat32 `json:"long-cryptocurrency-value"`
	ShortCryptocurrencyValue           StringToFloat32 `json:"short-cryptocurrency-value"`
	CryptocurrencyMarginRequirement    StringToFloat32 `json:"cryptocurrency-margin-requirement"`
	UnsettledCryptocurrencyFiatAmount  StringToFloat32 `json:"unsettled-cryptocurrency-fiat-amount"`
	UnsettledCryptocurrencyFiatEffect  string          `json:"unsettled-cryptocurrency-fiat-effect"`
	ClosedLoopAvailableBalance         StringToFloat32 `json:"closed-loop-available-balance"`
	EquityOfferingMarginRequirement    StringToFloat32 `json:"equity-offering-margin-requirement"`
	LongBondValue                      StringToFloat32 `json:"long-bond-value"`
	BondMarginRequirement              StringToFloat32 `json:"bond-margin-requirement"`
	SnapshotDate                       string          `json:"snapshot-date"`
	TimeOfDay                          string          `json:"time-of-day"`
	RegTMarginRequirement              StringToFloat32 `json:"reg-t-margin-requirement"`
	FuturesOvernightMarginRequirement  StringToFloat32 `json:"futures-overnight-margin-requirement"`
	FuturesIntradayMarginRequirement   StringToFloat32 `json:"futures-intraday-margin-requirement"`
	MaintenanceExcess                  StringToFloat32 `json:"maintenance-excess"`
	PendingMarginInterest              StringToFloat32 `json:"pending-margin-interest"`
	ApexStartingDayMarginEquity        StringToFloat32 `json:"apex-starting-day-margin-equity"`
	BuyingPowerAdjustment              StringToFloat32 `json:"buying-power-adjustment"`
	BuyingPowerAdjustmentEffect        string          `json:"buying-power-adjustment-effect"`
	EffectiveCryptocurrencyBuyingPower StringToFloat32 `json:"effective-cryptocurrency-buying-power"`
	UpdatedAt                          time.Time       `json:"updated-at"`
}

type AccountPosition struct {
	AccountNumber                 string                   `json:"account-number"`
	Symbol                        string                   `json:"symbol"`
	InstrumentType                constants.InstrumentType `json:"instrument-type"`
	UnderlyingSymbol              string                   `json:"underlying-symbol"`
	Quantity                      int                      `json:"quantity"`
	QuantityDirection             constants.Direction      `json:"quantity-direction"`
	ClosePrice                    StringToFloat32          `json:"close-price"`
	AverageOpenPrice              StringToFloat32          `json:"average-open-price"`
	AverageYearlyMarketClosePrice StringToFloat32          `json:"average-yearly-market-close-price"`
	AverageDailyMarketClosePrice  StringToFloat32          `json:"average-daily-market-close-price"`
	Mark                          StringToFloat32          `json:"mark"`
	MarkPrice                     StringToFloat32          `json:"mark-price"`
	Multiplier                    int                      `json:"multiplier"`
	CostEffect                    string                   `json:"cost-effect"`
	IsSuppressed                  bool                     `json:"is-suppressed"`
	IsFrozen                      bool                     `json:"is-frozen"`
	RestrictedQuantity            int                      `json:"restricted-quantity"`
	ExpiresAt                     time.Time                `json:"expires-at"`
	FixingPrice                   StringToFloat32          `json:"fixing-price"`
	DeliverableType               string                   `json:"deliverable-type"`
	RealizedDayGain               StringToFloat32          `json:"realized-day-gain"`
	RealizedDayGainEffect         string                   `json:"realized-day-gain-effect"`
	RealizedDayGainDate           string                   `json:"realized-day-gain-date"`
	RealizedToday                 StringToFloat32          `json:"realized-today"`
	RealizedTodayEffect           string                   `json:"realized-today-effect"`
	RealizedTodayDate             string                   `json:"realized-today-date"`
	CreatedAt                     time.Time                `json:"created-at"`
	UpdatedAt                     time.Time                `json:"updated-at"`
}

type AccountBalanceSnapshots struct {
	AccountNumber                      string          `json:"account-number"`
	CashBalance                        StringToFloat32 `json:"cash-balance"`
	LongEquityValue                    StringToFloat32 `json:"long-equity-value"`
	ShortEquityValue                   StringToFloat32 `json:"short-equity-value"`
	LongDerivativeValue                StringToFloat32 `json:"long-derivative-value"`
	ShortDerivativeValue               StringToFloat32 `json:"short-derivative-value"`
	LongFuturesValue                   StringToFloat32 `json:"long-futures-value"`
	ShortFuturesValue                  StringToFloat32 `json:"short-futures-value"`
	LongFuturesDerivativeValue         StringToFloat32 `json:"long-futures-derivative-value"`
	ShortFuturesDerivativeValue        StringToFloat32 `json:"short-futures-derivative-value"`
	LongMargineableValue               StringToFloat32 `json:"long-margineable-value"`
	ShortMargineableValue              StringToFloat32 `json:"short-margineable-value"`
	MarginEquity                       StringToFloat32 `json:"margin-equity"`
	EquityBuyingPower                  StringToFloat32 `json:"equity-buying-power"`
	DerivativeBuyingPower              StringToFloat32 `json:"derivative-buying-power"`
	DayTradingBuyingPower              StringToFloat32 `json:"day-trading-buying-power"`
	FuturesMarginRequirement           StringToFloat32 `json:"futures-margin-requirement"`
	AvailableTradingFunds              StringToFloat32 `json:"available-trading-funds"`
	MaintenanceRequirement             StringToFloat32 `json:"maintenance-requirement"`
	MaintenanceCallValue               StringToFloat32 `json:"maintenance-call-value"`
	RegTCallValue                      StringToFloat32 `json:"reg-t-call-value"`
	DayTradingCallValue                StringToFloat32 `json:"day-trading-call-value"`
	DayEquityCallValue                 StringToFloat32 `json:"day-equity-call-value"`
	NetLiquidatingValue                StringToFloat32 `json:"net-liquidating-value"`
	CashAvailableToWithdraw            StringToFloat32 `json:"cash-available-to-withdraw"`
	DayTradeExcess                     StringToFloat32 `json:"day-trade-excess"`
	PendingCash                        StringToFloat32 `json:"pending-cash"`
	PendingCashEffect                  string          `json:"pending-cash-effect"`
	LongCryptocurrencyValue            StringToFloat32 `json:"long-cryptocurrency-value"`
	ShortCryptocurrencyValue           StringToFloat32 `json:"short-cryptocurrency-value"`
	CryptocurrencyMarginRequirement    StringToFloat32 `json:"cryptocurrency-margin-requirement"`
	UnsettledCryptocurrencyFiatAmount  StringToFloat32 `json:"unsettled-cryptocurrency-fiat-amount"`
	UnsettledCryptocurrencyFiatEffect  string          `json:"unsettled-cryptocurrency-fiat-effect"`
	ClosedLoopAvailableBalance         StringToFloat32 `json:"closed-loop-available-balance"`
	EquityOfferingMarginRequirement    StringToFloat32 `json:"equity-offering-margin-requirement"`
	LongBondValue                      StringToFloat32 `json:"long-bond-value"`
	BondMarginRequirement              StringToFloat32 `json:"bond-margin-requirement"`
	SnapshotDate                       string          `json:"snapshot-date"`
	RegTMarginRequirement              StringToFloat32 `json:"reg-t-margin-requirement"`
	FuturesOvernightMarginRequirement  StringToFloat32 `json:"futures-overnight-margin-requirement"`
	FuturesIntradayMarginRequirement   StringToFloat32 `json:"futures-intraday-margin-requirement"`
	MaintenanceExcess                  StringToFloat32 `json:"maintenance-excess"`
	PendingMarginInterest              StringToFloat32 `json:"pending-margin-interest"`
	EffectiveCryptocurrencyBuyingPower StringToFloat32 `json:"effective-cryptocurrency-buying-power"`
	UpdatedAt                          time.Time       `json:"updated-at"`
}

// PositionLimit model
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
