package tasty

import (
	"time"
)

type Lots struct {
	ID                string          `json:"id"`
	TransactionID     int             `json:"transaction-id"`
	Quantity          StringToFloat32 `json:"quantity"`
	Price             StringToFloat32 `json:"price"`
	QuantityDirection string          `json:"quantity-direction"`
	ExecutedAt        string          `json:"executed-at"`
	TransactionDate   string          `json:"transaction-date"`
}

type Transaction struct {
	ID                               int             `json:"id"`
	AccountNumber                    string          `json:"account-number"`
	Symbol                           string          `json:"symbol"`
	InstrumentType                   InstrumentType  `json:"instrument-type"`
	UnderlyingSymbol                 string          `json:"underlying-symbol"`
	TransactionType                  string          `json:"transaction-type"`
	TransactionSubType               OrderAction     `json:"transaction-sub-type"`
	Description                      string          `json:"description"`
	Action                           OrderAction     `json:"action"`
	Quantity                         StringToFloat32 `json:"quantity"`
	Price                            StringToFloat32 `json:"price"`
	ExecutedAt                       time.Time       `json:"executed-at"`
	TransactionDate                  string          `json:"transaction-date"`
	Value                            StringToFloat32 `json:"value"`
	ValueEffect                      PriceEffect     `json:"value-effect"`
	RegulatoryFees                   StringToFloat32 `json:"regulatory-fees"`
	RegulatoryFeesEffect             PriceEffect     `json:"regulatory-fees-effect"`
	ClearingFees                     StringToFloat32 `json:"clearing-fees"`
	ClearingFeesEffect               PriceEffect     `json:"clearing-fees-effect"`
	OtherCharge                      StringToFloat32 `json:"other-charge"`
	OtherChargeEffect                PriceEffect     `json:"other-charge-effect"`
	OtherChargeDescription           string          `json:"other-charge-description"`
	NetValue                         StringToFloat32 `json:"net-value"`
	NetValueEffect                   PriceEffect     `json:"net-value-effect"`
	Commission                       StringToFloat32 `json:"commission"`
	CommissionEffect                 PriceEffect     `json:"commission-effect"`
	ProprietaryIndexOptionFees       StringToFloat32 `json:"proprietary-index-option-fees"`
	ProprietaryIndexOptionFeesEffect PriceEffect     `json:"proprietary-index-option-fees-effect"`
	IsEstimatedFee                   bool            `json:"is-estimated-fee"`
	ExtExchangeOrderNumber           string          `json:"ext-exchange-order-number"`
	ExtGlobalOrderNumber             int             `json:"ext-global-order-number"`
	ExtGroupID                       string          `json:"ext-group-id"`
	ExtGroupFillID                   string          `json:"ext-group-fill-id"`
	ExtExecID                        string          `json:"ext-exec-id"`
	ExecID                           string          `json:"exec-id"`
	Exchange                         string          `json:"exchange"`
	OrderID                          int             `json:"order-id"`
	ReversesID                       int             `json:"reverses-id"`
	ExchangeAffiliationIDentifier    string          `json:"exchange-affiliation-identifier"`
	CostBasisReconciliationDate      string          `json:"cost-basis-reconciliation-date"`
	Lots                             Lots            `json:"lots"`
	LegCount                         int             `json:"leg-count"`
	DestinationVenue                 string          `json:"destination-venue"`
	AgencyPrice                      StringToFloat32 `json:"agency-price"`
	PrincipalPrice                   StringToFloat32 `json:"principal-price"`
}

type TransactionFees struct {
	TotalFees       StringToFloat32 `json:"total-fees"`
	TotalFeesEffect PriceEffect     `json:"total-fees-effect"`
}

// The query for account transactions.
type TransactionsQuery struct {
	// Default value 250
	PerPage int `url:"per-page,omitempty"`
	// Default value 0
	PageOffset int `url:"page-offset,omitempty"`
	// The order to sort results in. Defaults to Desc, Accepts Desc or Asc.
	Sort SortOrder `url:"sort,omitempty"`
	// Filter based on transaction_type
	Type string `url:"type,omitempty"`
	// Allows filtering on multiple transaction_types
	Types []string `url:"types[],omitempty"`
	// Filter based on transaction_sub_type
	SubTypes []string `url:"sub-type[],omitempty"`
	// The start date of transactions to query.
	StartDate time.Time `layout:"2006-01-02" url:"start-date,omitempty"`
	// The end date of transactions to query. Defaults to now.
	EndDate time.Time `layout:"2006-01-02" url:"end-date,omitempty"`
	// The type of instrument i.e. InstrumentType
	InstrumentType InstrumentType `url:"instrument-type,omitempty"`
	// The Stock Ticker Symbol AAPL, OCC Option Symbol AAPL 191004P00275000,
	// TW Future Symbol /ESZ9, or TW Future Option Symbol ./ESZ9 EW4U9 190927P2975
	Symbol string `url:"symbol,omitempty"`
	// The Underlying Symbol. The Ticker Symbol FB or
	// TW Future Symbol with out date component /M6E or
	// the full TW Future Symbol /ESU9
	UnderlyingSymbol string `url:"underlying-symbol,omitempty"`
	// The action of the transaction. i.e. OrderAction
	Action OrderAction `url:"action,omitempty"`
	// Account partition key
	PartitionKey string `url:"partition-key,omitempty"`
	// The full TW Future Symbol /ESZ9 or
	// /NGZ19 if two year digit are appropriate
	FuturesSymbol string `url:"futures-symbol"`
	// DateTime start range for filtering transactions in full date-time
	StartAt time.Time `layout:"2006-01-02T15:04:05Z" url:"start-at,omitempty"`
	// DateTime end range for filtering transactions in full date-time
	EndAt time.Time `layout:"2006-01-02T15:04:05Z" url:"end-at,omitempty"`
}
