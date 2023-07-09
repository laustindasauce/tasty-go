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
