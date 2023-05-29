package models

type PositionEntry struct {
	InstrumentSymbol string `json:"instrument-symbol"`
	InstrumentType   string `json:"instrument-type"`
	Quantity         string `json:"quantity"`
	ClosePrice       string `json:"close-price"`
	FixingPrice      string `json:"fixing-price"`
}

type Holding struct {
	Description                  string          `json:"description"`
	MarginRequirement            string          `json:"margin-requirement"`
	MarginRequirementEffect      string          `json:"margin-requirement-effect"`
	InitialRequirement           string          `json:"initial-requirement"`
	InitialRequirementEffect     string          `json:"initial-requirement-effect"`
	MaintenanceRequirement       string          `json:"maintenance-requirement"`
	MaintenanceRequirementEffect string          `json:"maintenance-requirement-effect"`
	IncludesWorkingOrder         bool            `json:"includes-working-order"`
	BuyingPower                  string          `json:"buying-power"`
	BuyingPowerEffect            string          `json:"buying-power-effect"`
	PositionEntries              []PositionEntry `json:"position-entries"`
}

type MarginGroup struct {
	Description                   string    `json:"description"`
	Code                          string    `json:"code"`
	UnderlyingSymbol              string    `json:"underlying-symbol"`
	UnderlyingType                string    `json:"underlying-type"`
	ExpectedPriceRangeUpPercent   string    `json:"expected-price-range-up-percent"`
	ExpectedPriceRangeDownPercent string    `json:"expected-price-range-down-percent"`
	PointOfNoReturnPercent        string    `json:"point-of-no-return-percent"`
	MarginCalculationType         string    `json:"margin-calculation-type"`
	MarginRequirement             string    `json:"margin-requirement"`
	MarginRequirementEffect       string    `json:"margin-requirement-effect"`
	InitialRequirement            string    `json:"initial-requirement"`
	InitialRequirementEffect      string    `json:"initial-requirement-effect"`
	MaintenanceRequirement        string    `json:"maintenance-requirement"`
	MaintenanceRequirementEffect  string    `json:"maintenance-requirement-effect"`
	BuyingPower                   string    `json:"buying-power"`
	BuyingPowerEffect             string    `json:"buying-power-effect"`
	Holdings                      []Holding `json:"groups"`
	PriceIncreasePercent          string    `json:"price-increase-percent"`
	PriceDecreasePercent          string    `json:"price-decrease-percent"`
}

type MarginRequirements struct {
	AccountNumber                string        `json:"account-number"`
	Description                  string        `json:"description"`
	MarginCalculationType        string        `json:"margin-calculation-type"`
	OptionLevel                  string        `json:"option-level"`
	MarginRequirement            string        `json:"margin-requirement"`
	MarginRequirementEffect      string        `json:"margin-requirement-effect"`
	InitialRequirement           string        `json:"initial-requirement"`
	InitialRequirementEffect     string        `json:"initial-requirement-effect"`
	MaintenanceRequirement       string        `json:"maintenance-requirement"`
	MaintenanceRequirementEffect string        `json:"maintenance-requirement-effect"`
	MarginEquity                 string        `json:"margin-equity"`
	MarginEquityEffect           string        `json:"margin-equity-effect"`
	OptionBuyingPower            string        `json:"option-buying-power"`
	OptionBuyingPowerEffect      string        `json:"option-buying-power-effect"`
	RegTMarginRequirement        string        `json:"reg-t-margin-requirement"`
	RegTMarginRequirementEffect  string        `json:"reg-t-margin-requirement-effect"`
	RegTOptionBuyingPower        string        `json:"reg-t-option-buying-power"`
	RegTOptionBuyingPowerEffect  string        `json:"reg-t-option-buying-power-effect"`
	MaintenanceExcess            string        `json:"maintenance-excess"`
	MaintenanceExcessEffect      string        `json:"maintenance-excess-effect"`
	MarginGroups                 []MarginGroup `json:"groups"`
	LastStateTimestamp           int           `json:"last-state-timestamp"`
}
