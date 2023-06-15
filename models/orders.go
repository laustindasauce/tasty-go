package models

import (
	"time"

	"github.com/austinbspencer/tasty-go/constants"
)

type Order struct {
	ID                       int                      `json:"id"`
	AccountNumber            string                   `json:"account-number"`
	TimeInForce              constants.TimeInForce    `json:"time-in-force"`
	GtcDate                  string                   `json:"gtc-date"`
	OrderType                constants.OrderType      `json:"order-type"`
	Size                     int                      `json:"size"`
	UnderlyingSymbol         string                   `json:"underlying-symbol"`
	UnderlyingInstrumentType constants.InstrumentType `json:"underlying-instrument-type"`
	Price                    StringToFloat32          `json:"price"`
	PriceEffect              constants.PriceEffect    `json:"price-effect"`
	Value                    StringToFloat32          `json:"value"`
	ValueEffect              constants.PriceEffect    `json:"value-effect"`
	StopTrigger              StringToFloat32          `json:"stop-trigger"`
	Status                   constants.OrderStatus    `json:"status"`
	ContingentStatus         string                   `json:"contingent-status"`
	ConfirmationStatus       string                   `json:"confirmation-status"`
	Cancellable              bool                     `json:"cancellable"`
	CancelledAt              time.Time                `json:"cancelled-at"`
	CancelUserID             string                   `json:"cancel-user-id"`
	CancelUsername           string                   `json:"cancel-username"`
	Editable                 bool                     `json:"editable"`
	Edited                   bool                     `json:"edited"`
	ExtExchangeOrderNumber   string                   `json:"ext-exchange-order-number"`
	ExtClientOrderID         string                   `json:"ext-client-order-id"`
	ExtGlobalOrderNumber     int                      `json:"ext-global-order-number"`
	ReplacingOrderID         string                   `json:"replacing-order-id"`
	ReplacesOrderID          string                   `json:"replaces-order-id"`
	ReceivedAt               time.Time                `json:"received-at"`
	UpdatedAt                int                      `json:"updated-at"`
	InFlightAt               string                   `json:"in-flight-at"`
	LiveAt                   string                   `json:"live-at"`
	RejectReason             string                   `json:"reject-reason"`
	UserID                   string                   `json:"user-id"`
	Username                 string                   `json:"username"`
	TerminalAt               time.Time                `json:"terminal-at"`
	ComplexOrderID           int                      `json:"complex-order-id"`
	ComplexOrderTag          string                   `json:"complex-order-tag"`
	PreflightID              string                   `json:"preflight-id"`
	Legs                     []OrderLeg               `json:"legs"`
	Rules                    OrderRules               `json:"rules"`
}

type ComplexOrder struct {
	ID                                   int             `json:"id"`
	AccountNumber                        string          `json:"account-number"`
	Type                                 string          `json:"type"`
	TerminalAt                           string          `json:"terminal-at"`
	RatioPriceThreshold                  StringToFloat32 `json:"ratio-price-threshold"`
	RatioPriceComparator                 string          `json:"ratio-price-comparator"`
	RatioPriceIsThresholdBasedOnNotional bool            `json:"ratio-price-is-threshold-based-on-notional"`
	// RelatedOrders Non-current orders. This includes replaced orders, unfilled orders, and terminal orders.
	RelatedOrders []RelatedOrder `json:"related-orders"`
	// Orders with complex-order-tag: '::order'. For example, 'OTO::order' for OTO complex orders.
	Orders []Order `json:"orders"`
	// TriggerOrder Order with complex-order-tag: '::trigger-order'. For example, 'OTO::trigger-order for OTO complex orders.
	TriggerOrder Order `json:"trigger-order"`
}

type OrderInfo struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	PreflightID string `json:"preflight-id"`
}

type RelatedOrder struct {
	ID               int    `json:"id"`
	ComplexOrderID   int    `json:"complex-order-id"`
	ComplexOrderTag  string `json:"complex-order-tag"`
	ReplacesOrderID  string `json:"replaces-order-id"`
	ReplacingOrderID string `json:"replacing-order-id"`
	Status           string `json:"status"`
}

type OrderLeg struct {
	InstrumentType    constants.InstrumentType `json:"instrument-type"`
	Symbol            string                   `json:"symbol"`
	Quantity          int                      `json:"quantity"`
	RemainingQuantity int                      `json:"remaining-quantity"`
	Action            constants.OrderAction    `json:"action"`
	Fills             []OrderFill              `json:"fills"`
}

type OrderFill struct {
	ExtGroupFillID   string          `json:"ext-group-fill-id"`
	ExtExecID        string          `json:"ext-exec-id"`
	FillID           string          `json:"fill-id"`
	Quantity         int             `json:"quantity"`
	FillPrice        StringToFloat32 `json:"fill-price"`
	FilledAt         string          `json:"filled-at"`
	DestinationVenue string          `json:"destination-venue"`
}

type OrderRules struct {
	RouteAfter  string           `json:"route-after"`
	RoutedAt    string           `json:"routed-at"`
	CancelAt    string           `json:"cancel-at"`
	CancelledAt string           `json:"cancelled-at"`
	Conditions  []OrderCondition `json:"conditions"`
}

type OrderCondition struct {
	ID                         int                       `json:"id"`
	Action                     constants.OrderRuleAction `json:"action"`
	Symbol                     string                    `json:"symbol"`
	InstrumentType             constants.InstrumentType  `json:"instrument-type"`
	Indicator                  constants.Indicator       `json:"indicator"`
	Comparator                 constants.Comparator      `json:"comparator"`
	Threshold                  StringToFloat32           `json:"threshold"`
	IsThresholdBasedOnNotional bool                      `json:"is-threshold-based-on-notional"`
	TriggeredAt                string                    `json:"triggered-at"`
	TriggeredValue             StringToFloat32           `json:"triggered-value"`
	PriceComponents            []OrderPriceComponent     `json:"price-components"`
}

type OrderPriceComponent struct {
	Symbol            string                   `json:"symbol"`
	InstrumentType    constants.InstrumentType `json:"instrument-type"`
	Quantity          int                      `json:"quantity"`
	QuantityDirection constants.Direction      `json:"quantity-direction"`
}

type NewOrderECR struct {
	// (Required) The length in time before the order expires.
	TimeInForce constants.TimeInForce `json:"time-in-force"`
	GtcDate     string                `json:"gtc-date,omitempty"`
	// (Required) The type of order in regards to the price.
	OrderType   constants.OrderType `json:"order-type"`
	StopTrigger float32             `json:"stop-trigger,omitempty"`
	Price       float32             `json:"price,omitempty"`
	// (Required) If pay or receive payment for placing the order.
	PriceEffect constants.PriceEffect `json:"price-effect"`
	Value       float32               `json:"value,omitempty"`
	// If pay or receive payment for placing the notional market order.
	// i.e. Credit or Debit
	ValueEffect  constants.PriceEffect `json:"value-effect"`
	Source       string                `json:"source,omitempty"`
	PartitionKey string                `json:"partition-key,omitempty"`
	PreflightID  string                `json:"preflight-id,omitempty"`
	Legs         []NewOrderLeg         `json:"legs,omitempty"`
}

type NewOrderLeg struct {
	InstrumentType constants.InstrumentType `json:"instrument-type"`
	Symbol         string                   `json:"symbol"`
	Quantity       int                      `json:"quantity,omitempty"`
	Action         constants.OrderAction    `json:"action"`
}

type FeeCalculation struct {
	RegulatoryFees                   StringToFloat32       `json:"regulatory-fees"`
	RegulatoryFeesEffect             constants.PriceEffect `json:"regulatory-fees-effect"`
	ClearingFees                     StringToFloat32       `json:"clearing-fees"`
	ClearingFeesEffect               constants.PriceEffect `json:"clearing-fees-effect"`
	Commission                       StringToFloat32       `json:"commission"`
	CommissionEffect                 constants.PriceEffect `json:"commission-effect"`
	ProprietaryIndexOptionFees       StringToFloat32       `json:"proprietary-index-option-fees"`
	ProprietaryIndexOptionFeesEffect constants.PriceEffect `json:"proprietary-index-option-fees-effect"`
	TotalFees                        StringToFloat32       `json:"total-fees"`
	TotalFeesEffect                  constants.PriceEffect `json:"total-fees-effect"`
}

type BuyingPowerEffect struct {
	ChangeInMarginRequirement            StringToFloat32       `json:"change-in-margin-requirement"`
	ChangeInMarginRequirementEffect      constants.PriceEffect `json:"change-in-margin-requirement-effect"`
	ChangeInBuyingPower                  StringToFloat32       `json:"change-in-buying-power"`
	ChangeInBuyingPowerEffect            constants.PriceEffect `json:"change-in-buying-power-effect"`
	CurrentBuyingPower                   StringToFloat32       `json:"current-buying-power"`
	CurrentBuyingPowerEffect             constants.PriceEffect `json:"current-buying-power-effect"`
	NewBuyingPower                       StringToFloat32       `json:"new-buying-power"`
	NewBuyingPowerEffect                 constants.PriceEffect `json:"new-buying-power-effect"`
	IsolatedOrderMarginRequirement       StringToFloat32       `json:"isolated-order-margin-requirement"`
	IsolatedOrderMarginRequirementEffect constants.PriceEffect `json:"isolated-order-margin-requirement-effect"`
	IsSpread                             bool                  `json:"is-spread"`
	Impact                               StringToFloat32       `json:"impact"`
	Effect                               constants.PriceEffect `json:"effect"`
}

type OrderResponse struct {
	Order             Order             `json:"order"`
	ComplexOrder      ComplexOrder      `json:"complex-order"`
	Warnings          []OrderInfo       `json:"warnings"`
	Errors            []OrderInfo       `json:"errors"`
	BuyingPowerEffect BuyingPowerEffect `json:"buying-power-effect"`
	FeeCalculation    FeeCalculation    `json:"fee-calculation"`
}

type OrderErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Errors  []OrderInfo `json:"errors"`
}

// OrderReplacement Replaces a live order with a new one. Subsequent fills
// of the original order will abort the replacement.
type OrderReplacement struct {
	// TimeInForce The length in time before the order expires.
	TimeInForce constants.TimeInForce `json:"time-in-force"`
	GtcDate     string                `json:"gtc-date"`
	// OrderType The type of order in regards to the price.
	OrderType   constants.OrderType `json:"order-type"`
	StopTrigger StringToFloat32     `json:"stop-trigger,omitempty"`
	Price       StringToFloat32     `json:"price,omitempty"`
	// PriceEffect If pay or receive payment for placing the order.
	PriceEffect constants.PriceEffect `json:"price-effect,omitempty"`
	Value       StringToFloat32       `json:"value,omitempty"`
	// ValueEffect If pay or receive payment for placing the notional market order.
	ValueEffect  constants.PriceEffect `json:"value-effect,omitempty"`
	Source       string                `json:"source,omitempty"`
	PartitionKey string                `json:"partition-key,omitempty"`
	PreflightID  string                `json:"preflight-id,omitempty"`
}

type NewOrder struct {
	TimeInForce  constants.TimeInForce `json:"time-in-force"`
	GtcDate      string                `json:"gtc-date"`
	OrderType    constants.OrderType   `json:"order-type"`
	StopTrigger  float32               `json:"stop-trigger,omitempty"`
	Price        float32               `json:"price,omitempty"`
	PriceEffect  constants.PriceEffect `json:"price-effect,omitempty"`
	Value        float32               `json:"value,omitempty"`
	ValueEffect  constants.PriceEffect `json:"value-effect,omitempty"`
	Source       string                `json:"source,omitempty"`
	PartitionKey string                `json:"partition-key,omitempty"`
	PreflightID  string                `json:"preflight-id,omitempty"`
	Legs         []NewOrderLeg         `json:"legs"`
	Rules        NewOrderRules         `json:"rules,omitempty"`
}

type NewOrderRules struct {
	// RouteAfter Earliest time an order should route at
	RouteAfter string `json:"route-after"`
	// CancelAt Latest time an order should be canceled at
	CancelAt   string              `json:"cancel-at"`
	Conditions []NewOrderCondition `json:"conditions"`
}

type NewOrderPriceComponent struct {
	// The symbol to apply the condition to.
	Symbol string `json:"symbol"`
	// The instrument's type in relation to the symbol.
	InstrumentType constants.InstrumentType `json:"instrument-type"`
	// The Ratio quantity in relation to the symbol
	Quantity int `json:"quantity"`
	// The quantity direction(ie Long or Short) in relation to the symbol
	QuantityDirection constants.Direction `json:"quantity-direction"`
}

type NewOrderCondition struct {
	// The action in which the trigger is enacted. i.e. route and cancel
	Action constants.OrderRuleAction `json:"action"`
	// The symbol to apply the condition to.
	Symbol string `json:"symbol"`
	// The instrument's type in relation to the condition.
	InstrumentType string `json:"instrument-type"`
	// The indicator for the trigger, currently only supports last
	Indicator constants.Indicator `json:"indicator"`
	// How to compare against the threshold.
	Comparator constants.Comparator `json:"comparator"`
	// The price at which the condition triggers.
	Threshold       float32                  `json:"threshold"`
	PriceComponents []NewOrderPriceComponent `json:"price-components"`
}
