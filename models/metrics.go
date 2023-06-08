package models

import "time"

type DividendInfo struct {
	OccurredDate string          `json:"occurred-date"`
	Amount       StringToFloat32 `json:"amount"`
}

type EarningsInfo struct {
	OccurredDate string          `json:"occurred-date"`
	Eps          StringToFloat32 `json:"eps"`
}

type Earnings struct {
	Visible            bool            `json:"visible"`
	ExpectedReportDate string          `json:"expected-report-date"`
	Estimated          bool            `json:"estimated"`
	LateFlag           int             `json:"late-flag"`
	QuarterEndDate     string          `json:"quarter-end-date"`
	ActualEPS          StringToFloat32 `json:"actual-eps"`
	ConsensusEstimate  StringToFloat32 `json:"consensus-estimate"`
	UpdatedAt          time.Time       `json:"updated-at"`
}

type OptionExpirationImpliedVolatility struct {
	ExpirationDate    string          `json:"expiration-date"`
	SettlementType    string          `json:"settlement-type"`
	OptionChainType   string          `json:"option-chain-type"`
	ImpliedVolatility StringToFloat32 `json:"implied-volatility"`
}

type MarketMetricInfo struct {
	Symbol                              string                              `json:"symbol"`
	ImpliedVolatilityIndex              StringToFloat32                     `json:"implied-volatility-index"`
	ImpliedVolatilityIndex5DayChange    StringToFloat32                     `json:"implied-volatility-index-5-day-change"`
	ImpliedVolatilityRank               StringToFloat32                     `json:"implied-volatility-rank"`
	ImpliedVolatilityPercentile         StringToFloat32                     `json:"implied-volatility-percentile"`
	Liquidity                           StringToFloat32                     `json:"liquidity"`
	LiquidityRank                       StringToFloat32                     `json:"liquidity-rank"`
	LiquidityRating                     int                                 `json:"liquidity-rating"`
	OptionExpirationImpliedVolatilities []OptionExpirationImpliedVolatility `json:"option-expiration-implied-volatilities"`
}

type LiquidityRunningState struct {
	Sum       StringToFloat32 `json:"sum"`
	Count     int             `json:"count"`
	StartedAt time.Time       `json:"started-at"`
	UpdatedAt time.Time       `json:"updated-at"`
}

type MarketMetricVolatility struct {
	Symbol                                 string                              `json:"symbol"`
	ImpliedVolatilityIndex                 StringToFloat32                     `json:"implied-volatility-index"`
	ImpliedVolatilityIndex5DayChange       StringToFloat32                     `json:"implied-volatility-index-5-day-change"`
	ImpliedVolatilityIndexRank             StringToFloat32                     `json:"implied-volatility-index-rank"`
	TosImpliedVolatilityIndexRank          StringToFloat32                     `json:"tos-implied-volatility-index-rank"`
	TwImpliedVolatilityIndexRank           StringToFloat32                     `json:"tw-implied-volatility-index-rank"`
	TosImpliedVolatilityIndexRankUpdatedAt time.Time                           `json:"tos-implied-volatility-index-rank-updated-at"`
	ImpliedVolatilityIndexRankSource       string                              `json:"implied-volatility-index-rank-source"`
	ImpliedVolatilityPercentile            StringToFloat32                     `json:"implied-volatility-percentile"`
	ImpliedVolatilityUpdatedAt             time.Time                           `json:"implied-volatility-updated-at"`
	LiquidityValue                         StringToFloat32                     `json:"liquidity-value"`
	LiquidityRank                          StringToFloat32                     `json:"liquidity-rank"`
	LiquidityRating                        int                                 `json:"liquidity-rating"`
	CreatedAt                              string                              `json:"created-at"`
	UpdatedAt                              time.Time                           `json:"updated-at"`
	OptionExpirationImpliedVolatilities    []OptionExpirationImpliedVolatility `json:"option-expiration-implied-volatilities"`
	LiquidityRunningState                  LiquidityRunningState               `json:"liquidity-running-state"`
	Beta                                   StringToFloat32                     `json:"beta"`
	BetaUpdatedAt                          time.Time                           `json:"beta-updated-at"`
	CorrSpy3month                          StringToFloat32                     `json:"corr-spy-3month"`
	DividendRatePerShare                   StringToFloat32                     `json:"dividend-rate-per-share"`
	AnnualDividendPerShare                 StringToFloat32                     `json:"annual-dividend-per-share"`
	DividendYield                          StringToFloat32                     `json:"dividend-yield"`
	DividendExDate                         string                              `json:"dividend-ex-date"`
	DividendNextDate                       string                              `json:"dividend-next-date"`
	DividendPayDate                        string                              `json:"dividend-pay-date"`
	DividendUpdatedAt                      time.Time                           `json:"dividend-updated-at"`
	Earnings                               Earnings                            `json:"earnings"`
	ListedMarket                           string                              `json:"listed-market"`
	Lendability                            string                              `json:"lendability"`
	BorrowRate                             StringToFloat32                     `json:"borrow-rate"`
	MarketCap                              int                                 `json:"market-cap"`
	ImpliedVolatility30Day                 StringToFloat32                     `json:"implied-volatility-30-day"`
	HistoricalVolatility30Day              StringToFloat32                     `json:"historical-volatility-30-day"`
	HistoricalVolatility60Day              StringToFloat32                     `json:"historical-volatility-60-day"`
	HistoricalVolatility90Day              StringToFloat32                     `json:"historical-volatility-90-day"`
	IvHv30DayDifference                    StringToFloat32                     `json:"iv-hv-30-day-difference"`
	PriceEarningsRatio                     StringToFloat32                     `json:"price-earnings-ratio"`
	EarningsPerShare                       StringToFloat32                     `json:"earnings-per-share"`
}

type MarketMetricsQuery struct {
	// Symbols is the list of symbols
	Symbols []string `url:"symbols,comma"`
}
