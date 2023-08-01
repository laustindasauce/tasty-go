package tasty

import "github.com/shopspring/decimal"

type PositionEntry struct {
	InstrumentSymbol    string          `json:"instrument-symbol"`
	InstrumentType      InstrumentType  `json:"instrument-type"`
	Quantity            decimal.Decimal `json:"quantity"`
	AverageOpenPrice    StringToFloat32 `json:"average-open-price"`
	ClosePrice          decimal.Decimal `json:"close-price"`
	FixingPrice         StringToFloat32 `json:"fixing-price"`
	StrikePrice         decimal.Decimal `json:"strike-price,omitempty"`
	OptionType          OptionType      `json:"option-type,omitempty"`
	DeliverableQuantity decimal.Decimal `json:"deliverable-quantity,omitempty"`
	ExpirationDate      string          `json:"expiration-date,omitempty"`
}

type Holding struct {
	Description                  string          `json:"description"`
	MarginRequirement            decimal.Decimal `json:"margin-requirement"`
	MarginRequirementEffect      PriceEffect     `json:"margin-requirement-effect"`
	InitialRequirement           decimal.Decimal `json:"initial-requirement"`
	InitialRequirementEffect     PriceEffect     `json:"initial-requirement-effect"`
	MaintenanceRequirement       decimal.Decimal `json:"maintenance-requirement"`
	MaintenanceRequirementEffect PriceEffect     `json:"maintenance-requirement-effect"`
	IncludesWorkingOrder         bool            `json:"includes-working-order"`
	BuyingPower                  decimal.Decimal `json:"buying-power"`
	BuyingPowerEffect            PriceEffect     `json:"buying-power-effect"`
	PositionEntries              []PositionEntry `json:"position-entries"`
}

type MarginGroup struct {
	Description                   string          `json:"description"`
	Code                          string          `json:"code"`
	UnderlyingSymbol              string          `json:"underlying-symbol"`
	UnderlyingType                string          `json:"underlying-type"`
	ExpectedPriceRangeUpPercent   decimal.Decimal `json:"expected-price-range-up-percent"`
	ExpectedPriceRangeDownPercent decimal.Decimal `json:"expected-price-range-down-percent"`
	PointOfNoReturnPercent        decimal.Decimal `json:"point-of-no-return-percent"`
	MarginCalculationType         string          `json:"margin-calculation-type"`
	MarginRequirement             decimal.Decimal `json:"margin-requirement"`
	MarginRequirementEffect       PriceEffect     `json:"margin-requirement-effect"`
	InitialRequirement            decimal.Decimal `json:"initial-requirement"`
	InitialRequirementEffect      PriceEffect     `json:"initial-requirement-effect"`
	MaintenanceRequirement        decimal.Decimal `json:"maintenance-requirement"`
	MaintenanceRequirementEffect  PriceEffect     `json:"maintenance-requirement-effect"`
	BuyingPower                   decimal.Decimal `json:"buying-power"`
	BuyingPowerEffect             PriceEffect     `json:"buying-power-effect"`
	Holdings                      []Holding       `json:"groups"`
	PriceIncreasePercent          decimal.Decimal `json:"price-increase-percent"`
	PriceDecreasePercent          decimal.Decimal `json:"price-decrease-percent"`
}

type MarginRequirements struct {
	AccountNumber                string          `json:"account-number"`
	Description                  string          `json:"description"`
	MarginCalculationType        string          `json:"margin-calculation-type"`
	OptionLevel                  string          `json:"option-level"`
	MarginRequirement            decimal.Decimal `json:"margin-requirement"`
	MarginRequirementEffect      PriceEffect     `json:"margin-requirement-effect"`
	InitialRequirement           decimal.Decimal `json:"initial-requirement"`
	InitialRequirementEffect     PriceEffect     `json:"initial-requirement-effect"`
	MaintenanceRequirement       decimal.Decimal `json:"maintenance-requirement"`
	MaintenanceRequirementEffect PriceEffect     `json:"maintenance-requirement-effect"`
	MarginEquity                 decimal.Decimal `json:"margin-equity"`
	MarginEquityEffect           PriceEffect     `json:"margin-equity-effect"`
	OptionBuyingPower            decimal.Decimal `json:"option-buying-power"`
	OptionBuyingPowerEffect      PriceEffect     `json:"option-buying-power-effect"`
	RegTMarginRequirement        decimal.Decimal `json:"reg-t-margin-requirement"`
	RegTMarginRequirementEffect  PriceEffect     `json:"reg-t-margin-requirement-effect"`
	RegTOptionBuyingPower        decimal.Decimal `json:"reg-t-option-buying-power"`
	RegTOptionBuyingPowerEffect  PriceEffect     `json:"reg-t-option-buying-power-effect"`
	MaintenanceExcess            decimal.Decimal `json:"maintenance-excess"`
	MaintenanceExcessEffect      PriceEffect     `json:"maintenance-excess-effect"`
	MarginGroups                 []MarginGroup   `json:"groups"`
	LastStateTimestamp           int             `json:"last-state-timestamp"`
}

type EffectiveMarginRequirements struct {
	UnderlyingSymbol       string          `json:"underlying-symbol"`
	LongEquityInitial      decimal.Decimal `json:"long-equity-initial"`
	ShortEquityInitial     decimal.Decimal `json:"short-equity-initial"`
	LongEquityMaintenance  decimal.Decimal `json:"long-equity-maintenance"`
	ShortEquityMaintenance decimal.Decimal `json:"short-equity-maintenance"`
	NakedOptionStandard    decimal.Decimal `json:"naked-option-standard"`
	NakedOptionMinimum     decimal.Decimal `json:"naked-option-minimum"`
	NakedOptionFloor       decimal.Decimal `json:"naked-option-floor"`
	ClearingIDentifier     string          `json:"clearing-identifier"`
	IsDeleted              bool            `json:"is-deleted"`
}

type MarginType struct {
	Name     string `json:"name"`
	IsMargin bool   `json:"is-margin"`
}

// Response object for the margin requirements public configuration request.
type MarginRequirementsGlobalConfiguration struct {
	RiskFreeRate decimal.Decimal `json:"risk-free-rate"`
}
