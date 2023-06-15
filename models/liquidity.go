package models

// NetLiqOHLC Represents the open, high, low, close values for a time interval.
// The close value is the value at time. Also includes pending cash values and
// a sum of the two for convenience.
type NetLiqOHLC struct {
	Open             StringToFloat32 `json:"open"`
	High             StringToFloat32 `json:"high"`
	Low              StringToFloat32 `json:"low"`
	Close            StringToFloat32 `json:"close"`
	PendingCashOpen  StringToFloat32 `json:"pending-cash-open"`
	PendingCashHigh  StringToFloat32 `json:"pending-cash-high"`
	PendingCashLow   StringToFloat32 `json:"pending-cash-low"`
	PendingCashClose StringToFloat32 `json:"pending-cash-close"`
	TotalOpen        StringToFloat32 `json:"total-open"`
	TotalHigh        StringToFloat32 `json:"total-high"`
	TotalLow         StringToFloat32 `json:"total-low"`
	TotalClose       StringToFloat32 `json:"total-close"`
	Time             string          `json:"time"`
}
