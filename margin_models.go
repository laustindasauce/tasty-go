package tasty

type PositionEntry struct {
	InstrumentSymbol    string          `json:"instrument-symbol"`
	InstrumentType      InstrumentType  `json:"instrument-type"`
	Quantity            StringToFloat32 `json:"quantity"`
	AverageOpenPrice    StringToFloat32 `json:"average-open-price"`
	ClosePrice          StringToFloat32 `json:"close-price"`
	FixingPrice         StringToFloat32 `json:"fixing-price"`
	StrikePrice         StringToFloat32 `json:"strike-price,omitempty"`
	OptionType          OptionType      `json:"option-type,omitempty"`
	DeliverableQuantity StringToFloat32 `json:"deliverable-quantity,omitempty"`
	ExpirationDate      string          `json:"expiration-date,omitempty"`
}

type Holding struct {
	Description                  string          `json:"description"`
	MarginRequirement            StringToFloat32 `json:"margin-requirement"`
	MarginRequirementEffect      PriceEffect     `json:"margin-requirement-effect"`
	InitialRequirement           StringToFloat32 `json:"initial-requirement"`
	InitialRequirementEffect     PriceEffect     `json:"initial-requirement-effect"`
	MaintenanceRequirement       StringToFloat32 `json:"maintenance-requirement"`
	MaintenanceRequirementEffect PriceEffect     `json:"maintenance-requirement-effect"`
	IncludesWorkingOrder         bool            `json:"includes-working-order"`
	BuyingPower                  StringToFloat32 `json:"buying-power"`
	BuyingPowerEffect            PriceEffect     `json:"buying-power-effect"`
	PositionEntries              []PositionEntry `json:"position-entries"`
}

type MarginGroup struct {
	Description                   string          `json:"description"`
	Code                          string          `json:"code"`
	UnderlyingSymbol              string          `json:"underlying-symbol"`
	UnderlyingType                string          `json:"underlying-type"`
	ExpectedPriceRangeUpPercent   StringToFloat32 `json:"expected-price-range-up-percent"`
	ExpectedPriceRangeDownPercent StringToFloat32 `json:"expected-price-range-down-percent"`
	PointOfNoReturnPercent        StringToFloat32 `json:"point-of-no-return-percent"`
	MarginCalculationType         string          `json:"margin-calculation-type"`
	MarginRequirement             StringToFloat32 `json:"margin-requirement"`
	MarginRequirementEffect       PriceEffect     `json:"margin-requirement-effect"`
	InitialRequirement            StringToFloat32 `json:"initial-requirement"`
	InitialRequirementEffect      PriceEffect     `json:"initial-requirement-effect"`
	MaintenanceRequirement        StringToFloat32 `json:"maintenance-requirement"`
	MaintenanceRequirementEffect  PriceEffect     `json:"maintenance-requirement-effect"`
	BuyingPower                   StringToFloat32 `json:"buying-power"`
	BuyingPowerEffect             PriceEffect     `json:"buying-power-effect"`
	Holdings                      []Holding       `json:"groups"`
	PriceIncreasePercent          StringToFloat32 `json:"price-increase-percent"`
	PriceDecreasePercent          StringToFloat32 `json:"price-decrease-percent"`
}

type MarginRequirements struct {
	AccountNumber                string          `json:"account-number"`
	Description                  string          `json:"description"`
	MarginCalculationType        string          `json:"margin-calculation-type"`
	OptionLevel                  string          `json:"option-level"`
	MarginRequirement            StringToFloat32 `json:"margin-requirement"`
	MarginRequirementEffect      PriceEffect     `json:"margin-requirement-effect"`
	InitialRequirement           StringToFloat32 `json:"initial-requirement"`
	InitialRequirementEffect     PriceEffect     `json:"initial-requirement-effect"`
	MaintenanceRequirement       StringToFloat32 `json:"maintenance-requirement"`
	MaintenanceRequirementEffect PriceEffect     `json:"maintenance-requirement-effect"`
	MarginEquity                 StringToFloat32 `json:"margin-equity"`
	MarginEquityEffect           PriceEffect     `json:"margin-equity-effect"`
	OptionBuyingPower            StringToFloat32 `json:"option-buying-power"`
	OptionBuyingPowerEffect      PriceEffect     `json:"option-buying-power-effect"`
	RegTMarginRequirement        StringToFloat32 `json:"reg-t-margin-requirement"`
	RegTMarginRequirementEffect  PriceEffect     `json:"reg-t-margin-requirement-effect"`
	RegTOptionBuyingPower        StringToFloat32 `json:"reg-t-option-buying-power"`
	RegTOptionBuyingPowerEffect  PriceEffect     `json:"reg-t-option-buying-power-effect"`
	MaintenanceExcess            StringToFloat32 `json:"maintenance-excess"`
	MaintenanceExcessEffect      PriceEffect     `json:"maintenance-excess-effect"`
	MarginGroups                 []MarginGroup   `json:"groups"`
	LastStateTimestamp           int             `json:"last-state-timestamp"`
}

type EffectiveMarginRequirements struct {
	UnderlyingSymbol       string          `json:"underlying-symbol"`
	LongEquityInitial      StringToFloat32 `json:"long-equity-initial"`
	ShortEquityInitial     StringToFloat32 `json:"short-equity-initial"`
	LongEquityMaintenance  StringToFloat32 `json:"long-equity-maintenance"`
	ShortEquityMaintenance StringToFloat32 `json:"short-equity-maintenance"`
	NakedOptionStandard    StringToFloat32 `json:"naked-option-standard"`
	NakedOptionMinimum     StringToFloat32 `json:"naked-option-minimum"`
	NakedOptionFloor       StringToFloat32 `json:"naked-option-floor"`
	ClearingIDentifier     string          `json:"clearing-identifier"`
	IsDeleted              bool            `json:"is-deleted"`
}

type MarginType struct {
	Name     string `json:"name"`
	IsMargin bool   `json:"is-margin"`
}
