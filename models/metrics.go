package models

type DividendInfo struct {
	OccurredDate string          `json:"occurred-date"`
	Amount       StringToFloat32 `json:"amount"`
}

type EarningsInfo struct {
	OccurredDate string          `json:"occurred-date"`
	Eps          StringToFloat32 `json:"eps"`
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
