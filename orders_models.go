package tasty

import (
	"time"
)

type Order struct {
	ID                       int             `json:"id"`
	AccountNumber            string          `json:"account-number"`
	TimeInForce              TimeInForce     `json:"time-in-force"`
	GtcDate                  string          `json:"gtc-date"`
	OrderType                OrderType       `json:"order-type"`
	Size                     int             `json:"size"`
	UnderlyingSymbol         string          `json:"underlying-symbol"`
	UnderlyingInstrumentType InstrumentType  `json:"underlying-instrument-type"`
	Price                    StringToFloat32 `json:"price"`
	PriceEffect              PriceEffect     `json:"price-effect"`
	Value                    StringToFloat32 `json:"value"`
	ValueEffect              PriceEffect     `json:"value-effect"`
	StopTrigger              StringToFloat32 `json:"stop-trigger"`
	Status                   OrderStatus     `json:"status"`
	ContingentStatus         string          `json:"contingent-status"`
	ConfirmationStatus       string          `json:"confirmation-status"`
	Cancellable              bool            `json:"cancellable"`
	CancelledAt              time.Time       `json:"cancelled-at"`
	CancelUserID             string          `json:"cancel-user-id"`
	CancelUsername           string          `json:"cancel-username"`
	Editable                 bool            `json:"editable"`
	Edited                   bool            `json:"edited"`
	ExtExchangeOrderNumber   string          `json:"ext-exchange-order-number"`
	ExtClientOrderID         string          `json:"ext-client-order-id"`
	ExtGlobalOrderNumber     int             `json:"ext-global-order-number"`
	ReplacingOrderID         string          `json:"replacing-order-id"`
	ReplacesOrderID          string          `json:"replaces-order-id"`
	ReceivedAt               time.Time       `json:"received-at"`
	UpdatedAt                int             `json:"updated-at"`
	InFlightAt               string          `json:"in-flight-at"`
	LiveAt                   string          `json:"live-at"`
	RejectReason             string          `json:"reject-reason"`
	UserID                   string          `json:"user-id"`
	Username                 string          `json:"username"`
	TerminalAt               time.Time       `json:"terminal-at"`
	ComplexOrderID           int             `json:"complex-order-id"`
	ComplexOrderTag          string          `json:"complex-order-tag"`
	PreflightID              string          `json:"preflight-id"`
	Legs                     []OrderLeg      `json:"legs"`
	Rules                    OrderRules      `json:"rules"`
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
	InstrumentType    InstrumentType `json:"instrument-type"`
	Symbol            string         `json:"symbol"`
	Quantity          int            `json:"quantity"`
	RemainingQuantity int            `json:"remaining-quantity"`
	Action            OrderAction    `json:"action"`
	Fills             []OrderFill    `json:"fills"`
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
	ID                         int                   `json:"id"`
	Action                     OrderRuleAction       `json:"action"`
	Symbol                     string                `json:"symbol"`
	InstrumentType             InstrumentType        `json:"instrument-type"`
	Indicator                  Indicator             `json:"indicator"`
	Comparator                 Comparator            `json:"comparator"`
	Threshold                  StringToFloat32       `json:"threshold"`
	IsThresholdBasedOnNotional bool                  `json:"is-threshold-based-on-notional"`
	TriggeredAt                string                `json:"triggered-at"`
	TriggeredValue             StringToFloat32       `json:"triggered-value"`
	PriceComponents            []OrderPriceComponent `json:"price-components"`
}

type OrderPriceComponent struct {
	Symbol            string         `json:"symbol"`
	InstrumentType    InstrumentType `json:"instrument-type"`
	Quantity          int            `json:"quantity"`
	QuantityDirection Direction      `json:"quantity-direction"`
}

type NewOrderECR struct {
	// (Required) The length in time before the order expires.
	TimeInForce TimeInForce `json:"time-in-force"`
	GtcDate     string      `json:"gtc-date,omitempty"`
	// (Required) The type of order in regards to the price.
	OrderType   OrderType `json:"order-type"`
	StopTrigger float32   `json:"stop-trigger,omitempty"`
	Price       float32   `json:"price,omitempty"`
	// (Required) If pay or receive payment for placing the order.
	PriceEffect PriceEffect `json:"price-effect"`
	Value       float32     `json:"value,omitempty"`
	// If pay or receive payment for placing the notional market order.
	// i.e. Credit or Debit
	ValueEffect  PriceEffect   `json:"value-effect"`
	Source       string        `json:"source,omitempty"`
	PartitionKey string        `json:"partition-key,omitempty"`
	PreflightID  string        `json:"preflight-id,omitempty"`
	Legs         []NewOrderLeg `json:"legs,omitempty"`
}

type NewOrderLeg struct {
	InstrumentType InstrumentType `json:"instrument-type"`
	Symbol         string         `json:"symbol"`
	Quantity       int            `json:"quantity,omitempty"`
	Action         OrderAction    `json:"action"`
}

type FeeCalculation struct {
	RegulatoryFees                   StringToFloat32 `json:"regulatory-fees"`
	RegulatoryFeesEffect             PriceEffect     `json:"regulatory-fees-effect"`
	ClearingFees                     StringToFloat32 `json:"clearing-fees"`
	ClearingFeesEffect               PriceEffect     `json:"clearing-fees-effect"`
	Commission                       StringToFloat32 `json:"commission"`
	CommissionEffect                 PriceEffect     `json:"commission-effect"`
	ProprietaryIndexOptionFees       StringToFloat32 `json:"proprietary-index-option-fees"`
	ProprietaryIndexOptionFeesEffect PriceEffect     `json:"proprietary-index-option-fees-effect"`
	TotalFees                        StringToFloat32 `json:"total-fees"`
	TotalFeesEffect                  PriceEffect     `json:"total-fees-effect"`
}

type BuyingPowerEffect struct {
	ChangeInMarginRequirement            StringToFloat32 `json:"change-in-margin-requirement"`
	ChangeInMarginRequirementEffect      PriceEffect     `json:"change-in-margin-requirement-effect"`
	ChangeInBuyingPower                  StringToFloat32 `json:"change-in-buying-power"`
	ChangeInBuyingPowerEffect            PriceEffect     `json:"change-in-buying-power-effect"`
	CurrentBuyingPower                   StringToFloat32 `json:"current-buying-power"`
	CurrentBuyingPowerEffect             PriceEffect     `json:"current-buying-power-effect"`
	NewBuyingPower                       StringToFloat32 `json:"new-buying-power"`
	NewBuyingPowerEffect                 PriceEffect     `json:"new-buying-power-effect"`
	IsolatedOrderMarginRequirement       StringToFloat32 `json:"isolated-order-margin-requirement"`
	IsolatedOrderMarginRequirementEffect PriceEffect     `json:"isolated-order-margin-requirement-effect"`
	IsSpread                             bool            `json:"is-spread"`
	Impact                               StringToFloat32 `json:"impact"`
	Effect                               PriceEffect     `json:"effect"`
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
	TimeInForce TimeInForce `json:"time-in-force"`
	GtcDate     string      `json:"gtc-date"`
	// OrderType The type of order in regards to the price.
	OrderType   OrderType       `json:"order-type"`
	StopTrigger StringToFloat32 `json:"stop-trigger,omitempty"`
	Price       StringToFloat32 `json:"price,omitempty"`
	// PriceEffect If pay or receive payment for placing the order.
	PriceEffect PriceEffect     `json:"price-effect,omitempty"`
	Value       StringToFloat32 `json:"value,omitempty"`
	// ValueEffect If pay or receive payment for placing the notional market order.
	ValueEffect  PriceEffect `json:"value-effect,omitempty"`
	Source       string      `json:"source,omitempty"`
	PartitionKey string      `json:"partition-key,omitempty"`
	PreflightID  string      `json:"preflight-id,omitempty"`
}

type NewOrder struct {
	TimeInForce  TimeInForce   `json:"time-in-force"`
	GtcDate      string        `json:"gtc-date"`
	OrderType    OrderType     `json:"order-type"`
	StopTrigger  float32       `json:"stop-trigger,omitempty"`
	Price        float32       `json:"price,omitempty"`
	PriceEffect  PriceEffect   `json:"price-effect,omitempty"`
	Value        float32       `json:"value,omitempty"`
	ValueEffect  PriceEffect   `json:"value-effect,omitempty"`
	Source       string        `json:"source,omitempty"`
	PartitionKey string        `json:"partition-key,omitempty"`
	PreflightID  string        `json:"preflight-id,omitempty"`
	Legs         []NewOrderLeg `json:"legs"`
	Rules        NewOrderRules `json:"rules,omitempty"`
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
	InstrumentType InstrumentType `json:"instrument-type"`
	// The Ratio quantity in relation to the symbol
	Quantity int `json:"quantity"`
	// The quantity direction(ie Long or Short) in relation to the symbol
	QuantityDirection Direction `json:"quantity-direction"`
}

type NewOrderCondition struct {
	// The action in which the trigger is enacted. i.e. route and cancel
	Action OrderRuleAction `json:"action"`
	// The symbol to apply the condition to.
	Symbol string `json:"symbol"`
	// The instrument's type in relation to the condition.
	InstrumentType string `json:"instrument-type"`
	// The indicator for the trigger, currently only supports last
	Indicator Indicator `json:"indicator"`
	// How to compare against the threshold.
	Comparator Comparator `json:"comparator"`
	// The price at which the condition triggers.
	Threshold       float32                  `json:"threshold"`
	PriceComponents []NewOrderPriceComponent `json:"price-components"`
}

// The query for account orders.
type OrdersQuery struct {
	// Default value 10
	PerPage int `url:"per-page,omitempty"`
	// Default value 0
	PageOffset int `url:"page-offset,omitempty"`
	// The start date of orders to query.
	StartDate time.Time `layout:"2006-01-02" url:"start-date,omitempty"`
	// The end date of orders to query.
	EndDate time.Time `layout:"2006-01-02" url:"end-date,omitempty"`
	// The Underlying Symbol. The Ticker Symbol FB or
	// TW Future Symbol with out date component /M6E or
	// the full TW Future Symbol /ESU9
	UnderlyingSymbol string `url:"underlying-symbol,omitempty"`
	// Status of the order
	Status []OrderStatus `url:"status[],omitempty"`
	// The full TW Future Symbol /ESZ9 or
	// /NGZ19 if two year digit are appropriate
	FuturesSymbol string `url:"futures-symbol"`
	// Underlying instrument type i.e. InstrumentType
	UnderlyingInstrumentType InstrumentType `url:"underlying-instrument-type,omitempty"`
	// The order to sort results in. Defaults to Desc, Accepts Desc or Asc.
	Sort SortOrder `url:"sort,omitempty"`
	// DateTime start range for filtering transactions in full date-time
	StartAt time.Time `layout:"2006-01-02T15:04:05Z" url:"start-at,omitempty"`
	// DateTime end range for filtering transactions in full date-time
	EndAt time.Time `layout:"2006-01-02T15:04:05Z" url:"end-at,omitempty"`
}
