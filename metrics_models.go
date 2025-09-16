package tasty

import (
	"time"

	"github.com/shopspring/decimal"
)

type DividendInfo struct {
	OccurredDate string          `json:"occurred-date"`
	Amount       decimal.Decimal `json:"amount"`
}

type EarningsInfo struct {
	OccurredDate string          `json:"occurred-date"`
	Eps          decimal.Decimal `json:"eps"`
}

type Earnings struct {
	Visible            bool            `json:"visible"`
	ExpectedReportDate string          `json:"expected-report-date"`
	Estimated          bool            `json:"estimated"`
	LateFlag           int             `json:"late-flag"`
	QuarterEndDate     string          `json:"quarter-end-date"`
	ActualEPS          decimal.Decimal `json:"actual-eps"`
	ConsensusEstimate  decimal.Decimal `json:"consensus-estimate"`
	UpdatedAt          time.Time       `json:"updated-at"`
}

type OptionExpirationImpliedVolatility struct {
	ExpirationDate    string          `json:"expiration-date"`
	SettlementType    string          `json:"settlement-type"`
	OptionChainType   string          `json:"option-chain-type"`
	ImpliedVolatility decimal.Decimal `json:"implied-volatility"`
}

type MarketMetricInfo struct {
	Symbol                              string                              `json:"symbol"`
	ImpliedVolatilityIndex              decimal.Decimal                     `json:"implied-volatility-index"`
	ImpliedVolatilityIndex5DayChange    decimal.Decimal                     `json:"implied-volatility-index-5-day-change"`
	ImpliedVolatilityRank               decimal.Decimal                     `json:"implied-volatility-rank"`
	ImpliedVolatilityPercentile         decimal.Decimal                     `json:"implied-volatility-percentile"`
	Liquidity                           decimal.Decimal                     `json:"liquidity"`
	LiquidityRank                       decimal.Decimal                     `json:"liquidity-rank"`
	LiquidityRating                     int                                 `json:"liquidity-rating"`
	OptionExpirationImpliedVolatilities []OptionExpirationImpliedVolatility `json:"option-expiration-implied-volatilities"`
}

type LiquidityRunningState struct {
	Sum       decimal.Decimal `json:"sum"`
	Count     int             `json:"count"`
	StartedAt time.Time       `json:"started-at"`
	UpdatedAt time.Time       `json:"updated-at"`
}

type MarketMetricVolatility struct {
	Symbol                                 string                              `json:"symbol"`
	ImpliedVolatilityIndex                 decimal.Decimal                     `json:"implied-volatility-index"`
	ImpliedVolatilityIndex5DayChange       decimal.Decimal                     `json:"implied-volatility-index-5-day-change"`
	ImpliedVolatilityIndexRank             decimal.Decimal                     `json:"implied-volatility-index-rank"`
	TosImpliedVolatilityIndexRank          decimal.Decimal                     `json:"tos-implied-volatility-index-rank"`
	TwImpliedVolatilityIndexRank           decimal.Decimal                     `json:"tw-implied-volatility-index-rank"`
	TosImpliedVolatilityIndexRankUpdatedAt time.Time                           `json:"tos-implied-volatility-index-rank-updated-at"`
	ImpliedVolatilityIndexRankSource       string                              `json:"implied-volatility-index-rank-source"`
	ImpliedVolatilityPercentile            decimal.Decimal                     `json:"implied-volatility-percentile"`
	ImpliedVolatilityUpdatedAt             time.Time                           `json:"implied-volatility-updated-at"`
	LiquidityValue                         decimal.Decimal                     `json:"liquidity-value"`
	LiquidityRank                          decimal.Decimal                     `json:"liquidity-rank"`
	LiquidityRating                        int                                 `json:"liquidity-rating"`
	CreatedAt                              string                              `json:"created-at"`
	UpdatedAt                              time.Time                           `json:"updated-at"`
	OptionExpirationImpliedVolatilities    []OptionExpirationImpliedVolatility `json:"option-expiration-implied-volatilities"`
	LiquidityRunningState                  LiquidityRunningState               `json:"liquidity-running-state"`
	Beta                                   decimal.Decimal                     `json:"beta"`
	BetaUpdatedAt                          time.Time                           `json:"beta-updated-at"`
	CorrSpy3month                          decimal.Decimal                     `json:"corr-spy-3month"`
	DividendRatePerShare                   decimal.Decimal                     `json:"dividend-rate-per-share"`
	AnnualDividendPerShare                 decimal.Decimal                     `json:"annual-dividend-per-share"`
	DividendYield                          decimal.Decimal                     `json:"dividend-yield"`
	DividendExDate                         string                              `json:"dividend-ex-date"`
	DividendNextDate                       string                              `json:"dividend-next-date"`
	DividendPayDate                        string                              `json:"dividend-pay-date"`
	DividendUpdatedAt                      time.Time                           `json:"dividend-updated-at"`
	Earnings                               Earnings                            `json:"earnings"`
	ListedMarket                           string                              `json:"listed-market"`
	Lendability                            string                              `json:"lendability"`
	BorrowRate                             decimal.Decimal                     `json:"borrow-rate"`
	MarketCap                              int                                 `json:"market-cap"`
	ImpliedVolatility30Day                 decimal.Decimal                     `json:"implied-volatility-30-day"`
	HistoricalVolatility30Day              decimal.Decimal                     `json:"historical-volatility-30-day"`
	HistoricalVolatility60Day              decimal.Decimal                     `json:"historical-volatility-60-day"`
	HistoricalVolatility90Day              decimal.Decimal                     `json:"historical-volatility-90-day"`
	IvHv30DayDifference                    decimal.Decimal                     `json:"iv-hv-30-day-difference"`
	PriceEarningsRatio                     decimal.Decimal                     `json:"price-earnings-ratio"`
	EarningsPerShare                       decimal.Decimal                     `json:"earnings-per-share"`
}

type MarketDataQuery struct {
	// Ticker(s) of indices
	Index []string `url:"index[]"`
	// The symbol(s) of the equities
	Equity []string `url:"equity[]"`
	// The symbol(s) of the equity option(s) using OCC Symbology, i.e. [FB 180629C00200000]
	EquityOption []string `url:"equity-option[]"`
	// Symbol(s) for future, i.e. [ES]
	Future []string `url:"future[]"`
	// Symbol(s) for future option(s) using proper symbology, i.e. [./ESZ5 EW4V5 251024C5850]
	FutureOption []string `url:"future-option[]"`
	// The symbol(s) of the cryptocurrency, i.e. [BTC]
	Cryptocurrency []string `url:"cryptocurrency[]"`
}

type MarketData struct {
	Symbol             string           `json:"symbol"`
	InstrumentType     InstrumentType   `json:"instrument-type"`
	UpdatedAt          time.Time        `json:"updated-at"`
	Bid                *decimal.Decimal `json:"bid,omitempty"`
	BidSize            *decimal.Decimal `json:"bid-size,omitempty"`
	Ask                *decimal.Decimal `json:"ask,omitempty"`
	AskSize            *decimal.Decimal `json:"ask-size,omitempty"`
	Mid                *decimal.Decimal `json:"mid,omitempty"`
	Mark               *decimal.Decimal `json:"mark,omitempty"`
	Last               *decimal.Decimal `json:"last,omitempty"`
	LastMkt            *decimal.Decimal `json:"last-mkt,omitempty"`
	Open               *decimal.Decimal `json:"open,omitempty"`
	DayHighPrice       *decimal.Decimal `json:"day-high-price,omitempty"`
	DayLowPrice        *decimal.Decimal `json:"day-low-price,omitempty"`
	Close              *decimal.Decimal `json:"close,omitempty"`
	ClosePriceType     string           `json:"close-price-type"`
	PrevClose          *decimal.Decimal `json:"prev-close,omitempty"`
	PrevClosePriceType string           `json:"prev-close-price-type"`
	SummaryDate        string           `json:"summary-date"`
	PrevCloseDate      string           `json:"prev-close-date"`
	IsTradingHalted    bool             `json:"is-trading-halted"`
	HaltStartTime      int              `json:"halt-start-time"`
	HaltEndTime        int              `json:"halt-end-time"`
	Volume             *decimal.Decimal `json:"volume,omitempty"`
	Volatility         *decimal.Decimal `json:"volatility,omitempty"`
	Delta              *decimal.Decimal `json:"delta,omitempty"`
	Gamma              *decimal.Decimal `json:"gamma,omitempty"`
	Theta              *decimal.Decimal `json:"theta,omitempty"`
	Rho                *decimal.Decimal `json:"rho,omitempty"`
	Vega               *decimal.Decimal `json:"vega,omitempty"`
	TheoPrice          *decimal.Decimal `json:"theo-price,omitempty"`
	DxMark             *decimal.Decimal `json:"dx-mark,omitempty"`
	TickSize           *decimal.Decimal `json:"tick-size,omitempty"`
	OpenInterest       int              `json:"open-interest"`
	Beta               *decimal.Decimal `json:"beta,omitempty"`
	DividendAmount     *decimal.Decimal `json:"dividend-amount,omitempty"`
	DividendFrequency  *decimal.Decimal `json:"dividend-frequency,omitempty"`
	LowLimitPrice      *decimal.Decimal `json:"low-limit-price,omitempty"`
	HighLimitPrice     *decimal.Decimal `json:"high-limit-price,omitempty"`
	YearLowPrice       *decimal.Decimal `json:"year-low-price,omitempty"`
	YearHighPrice      *decimal.Decimal `json:"year-high-price,omitempty"`
}
