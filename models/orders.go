package models

import "github.com/austinbspencer/tasty-go/constants"

type Order struct {
	ID                       string          `json:"id"`
	AccountNumber            string          `json:"account-number"`
	TimeInForce              string          `json:"time-in-force"`
	GtcDate                  string          `json:"gtc-date"`
	OrderType                string          `json:"order-type"`
	Size                     string          `json:"size"`
	UnderlyingSymbol         string          `json:"underlying-symbol"`
	UnderlyingInstrumentType string          `json:"underlying-instrument-type"`
	Price                    StringToFloat32 `json:"price"`
	PriceEffect              string          `json:"price-effect"`
	Value                    StringToFloat32 `json:"value"`
	ValueEffect              string          `json:"value-effect"`
	StopTrigger              string          `json:"stop-trigger"`
	Status                   string          `json:"status"`
	ContingentStatus         string          `json:"contingent-status"`
	ConfirmationStatus       string          `json:"confirmation-status"`
	Cancellable              bool            `json:"cancellable"`
	CancelledAt              string          `json:"cancelled-at"`
	CancelUserID             string          `json:"cancel-user-id"`
	CancelUsername           string          `json:"cancel-username"`
	Editable                 bool            `json:"editable"`
	Edited                   bool            `json:"edited"`
	ReplacingOrderID         string          `json:"replacing-order-id"`
	ReplacesOrderID          string          `json:"replaces-order-id"`
	ReceivedAt               string          `json:"received-at"`
	UpdatedAt                string          `json:"updated-at"`
	InFlightAt               string          `json:"in-flight-at"`
	LiveAt                   string          `json:"live-at"`
	RejectReason             string          `json:"reject-reason"`
	UserID                   string          `json:"user-id"`
	Username                 string          `json:"username"`
	TerminalAt               string          `json:"terminal-at"`
	ComplexOrderID           string          `json:"complex-order-id"`
	ComplexOrderTag          string          `json:"complex-order-tag"`
	PreflightID              string          `json:"preflight-id"`
	Legs                     []OrderLeg      `json:"legs"`
	OrderRule                OrderRule       `json:"order-rule"`
}

type ComplexOrder struct {
	ID                                   string          `json:"id"`
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
	ID               string `json:"id"`
	ComplexOrderID   string `json:"complex-order-id"`
	ComplexOrderTag  string `json:"complex-order-tag"`
	ReplacesOrderID  string `json:"replaces-order-id"`
	ReplacingOrderID string `json:"replacing-order-id"`
	Status           string `json:"status"`
}

type OrderLeg struct {
	InstrumentType    string      `json:"instrument-type"`
	Symbol            string      `json:"symbol"`
	Quantity          string      `json:"quantity"`
	RemainingQuantity string      `json:"remaining-quantity"`
	Action            string      `json:"action"`
	Fills             []OrderFill `json:"fills"`
}

type OrderFill struct {
	ExtGroupFillID   string          `json:"ext-group-fill-id"`
	ExtExecID        string          `json:"ext-exec-id"`
	FillID           string          `json:"fill-id"`
	Quantity         string          `json:"quantity"`
	FillPrice        StringToFloat32 `json:"fill-price"`
	FilledAt         string          `json:"filled-at"`
	DestinationVenue string          `json:"destination-venue"`
}

type OrderRule struct {
	RouteAfter      string           `json:"route-after"`
	RoutedAt        string           `json:"routed-at"`
	CancelAt        string           `json:"cancel-at"`
	CancelledAt     string           `json:"cancelled-at"`
	OrderConditions []OrderCondition `json:"order-conditions"`
}

type OrderCondition struct {
	ID                         string                `json:"id"`
	Action                     string                `json:"action"`
	Symbol                     string                `json:"symbol"`
	InstrumentType             string                `json:"instrument-type"`
	Indicator                  string                `json:"indicator"`
	Comparator                 string                `json:"comparator"`
	Threshold                  StringToFloat32       `json:"threshold"`
	IsThresholdBasedOnNotional bool                  `json:"is-threshold-based-on-notional"`
	TriggeredAt                string                `json:"triggered-at"`
	TriggeredValue             StringToFloat32       `json:"triggered-value"`
	PriceComponents            []OrderPriceComponent `json:"price-components"`
}

type OrderPriceComponent struct {
	Symbol            string                   `json:"symbol"`
	InstrumentType    constants.InstrumentType `json:"instrument-type"`
	Quantity          StringToFloat32          `json:"quantity"`
	QuantityDirection constants.Direction      `json:"quantity-direction"`
}

type OrderDryRunECR struct {
	// TimeInForce The length in time before the order expires.
	TimeInForce constants.TimeInForce `json:"time-in-force"`
	GtcDate     string                `json:"gtc-date"`
	// OrderType The type of order in regards to the price.
	OrderType   constants.OrderType `json:"order-type"`
	StopTrigger StringToFloat32     `json:"stop-trigger"`
	Price       StringToFloat32     `json:"price"`
	// PriceEffect If pay or receive payment for placing the order.
	PriceEffect constants.PriceEffect `json:"price-effect"`
	Value       StringToFloat32       `json:"value"`
	// ValueEffect
	ValueEffect  constants.PriceEffect `json:"value-effect"`
	Source       string                `json:"source"`
	PartitionKey string                `json:"partition-key"`
	PreflightID  string                `json:"preflight-id"`
	Legs         []OrderDryRunLeg      `json:"legs"`
}

type OrderDryRunLeg struct {
	InstrumentType constants.InstrumentType `json:"instrument-type"`
	Symbol         *string                  `json:"symbol"`
	Quantity       StringToFloat32          `json:"quantity"`
	Action         constants.OrderAction    `json:"action"`
}

type PlacedOrderResponse struct {
	Order             Order        `json:"order"`
	ComplexOrder      ComplexOrder `json:"complex-order"`
	Warnings          []OrderInfo  `json:"warnings"`
	Errors            []OrderInfo  `json:"errors"`
	BuyingPowerEffect string       `json:"buying-power-effect"`
	FeeCalculation    string       `json:"fee-calculation"`
}

// OrderReplacement Replaces a live order with a new one. Subsequent fills
// of the original order will abort the replacement.
type OrderReplacement struct {
	// TimeInForce The length in time before the order expires.
	TimeInForce constants.TimeInForce `json:"time-in-force"`
	GtcDate     string                `json:"gtc-date"`
	// OrderType The type of order in regards to the price.
	OrderType   constants.OrderType `json:"order-type"`
	StopTrigger StringToFloat32     `json:"stop-trigger"`
	Price       StringToFloat32     `json:"price"`
	// PriceEffect If pay or receive payment for placing the order.
	PriceEffect constants.PriceEffect `json:"price-effect"`
	Value       StringToFloat32       `json:"value"`
	// ValueEffect If pay or receive payment for placing the notional market order.
	ValueEffect  constants.PriceEffect `json:"value-effect"`
	Source       string                `json:"source"`
	PartitionKey string                `json:"partition-key"`
	PreflightID  string                `json:"preflight-id"`
}

type OrderDryRun struct {
	TimeInForce  constants.TimeInForce `json:"time-in-force"`
	GtcDate      string                `json:"gtc-date"`
	OrderType    constants.OrderType   `json:"order-type"`
	StopTrigger  StringToFloat32       `json:"stop-trigger"`
	Price        StringToFloat32       `json:"price"`
	PriceEffect  constants.PriceEffect `json:"price-effect"`
	Value        StringToFloat32       `json:"value"`
	ValueEffect  constants.PriceEffect `json:"value-effect"`
	Source       string                `json:"source"`
	PartitionKey string                `json:"partition-key"`
	PreflightID  string                `json:"preflight-id"`
	Legs         OrderDryRunLeg        `json:"legs"`
	Rules        OrderDryRunRules      `json:"rules"`
}

type OrderDryRunRules struct {
	// RouteAfter Earliest time an order should route at
	RouteAfter string `json:"route-after"`
	// CancelAt Latest time an order should be canceled at
	CancelAt   string                 `json:"cancel-at"`
	Conditions []OrderDryRunCondition `json:"conditions"`
}

type OrderDryRunPriceComponent struct {
	// The symbol to apply the condition to.
	Symbol string `json:"symbol"`
	// The instrument's type in relation to the symbol.
	InstrumentType constants.InstrumentType `json:"instrument-type"`
	// The Ratio quantity in relation to the symbol
	Quantity float32 `json:"quantity"`
	// The quantity direction(ie Long or Short) in relation to the symbol
	QuantityDirection constants.Direction `json:"quantity-direction"`
}

type OrderDryRunCondition struct {
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
	Threshold       float32                     `json:"threshold"`
	PriceComponents []OrderDryRunPriceComponent `json:"price-components"`
}
