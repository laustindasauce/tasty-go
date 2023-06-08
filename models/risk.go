package models

// MarginRequirement model
type MarginRequirement struct {
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
