package models

import (
	"time"

	"github.com/austinbspencer/tasty-go/constants"
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
	ID                               int                      `json:"id"`
	AccountNumber                    string                   `json:"account-number"`
	Symbol                           string                   `json:"symbol"`
	InstrumentType                   constants.InstrumentType `json:"instrument-type"`
	UnderlyingSymbol                 string                   `json:"underlying-symbol"`
	TransactionType                  string                   `json:"transaction-type"`
	TransactionSubType               string                   `json:"transaction-sub-type"`
	Description                      string                   `json:"description"`
	Action                           constants.OrderAction    `json:"action"`
	Quantity                         StringToFloat32          `json:"quantity"`
	Price                            StringToFloat32          `json:"price"`
	ExecutedAt                       time.Time                `json:"executed-at"`
	TransactionDate                  string                   `json:"transaction-date"`
	Value                            StringToFloat32          `json:"value"`
	ValueEffect                      string                   `json:"value-effect"`
	RegulatoryFees                   StringToFloat32          `json:"regulatory-fees"`
	RegulatoryFeesEffect             string                   `json:"regulatory-fees-effect"`
	ClearingFees                     StringToFloat32          `json:"clearing-fees"`
	ClearingFeesEffect               string                   `json:"clearing-fees-effect"`
	OtherCharge                      StringToFloat32          `json:"other-charge"`
	OtherChargeEffect                string                   `json:"other-charge-effect"`
	OtherChargeDescription           string                   `json:"other-charge-description"`
	NetValue                         StringToFloat32          `json:"net-value"`
	NetValueEffect                   string                   `json:"net-value-effect"`
	Commission                       StringToFloat32          `json:"commission"`
	CommissionEffect                 string                   `json:"commission-effect"`
	ProprietaryIndexOptionFees       StringToFloat32          `json:"proprietary-index-option-fees"`
	ProprietaryIndexOptionFeesEffect string                   `json:"proprietary-index-option-fees-effect"`
	IsEstimatedFee                   bool                     `json:"is-estimated-fee"`
	ExtExchangeOrderNumber           string                   `json:"ext-exchange-order-number"`
	ExtGlobalOrderNumber             int                      `json:"ext-global-order-number"`
	ExtGroupID                       string                   `json:"ext-group-id"`
	ExtGroupFillID                   string                   `json:"ext-group-fill-id"`
	ExtExecID                        string                   `json:"ext-exec-id"`
	ExecID                           string                   `json:"exec-id"`
	Exchange                         string                   `json:"exchange"`
	OrderID                          int                      `json:"order-id"`
	ReversesID                       int                      `json:"reverses-id"`
	ExchangeAffiliationIDentifier    string                   `json:"exchange-affiliation-identifier"`
	CostBasisReconciliationDate      string                   `json:"cost-basis-reconciliation-date"`
	Lots                             Lots                     `json:"lots"`
	LegCount                         int                      `json:"leg-count"`
	DestinationVenue                 string                   `json:"destination-venue"`
	AgencyPrice                      StringToFloat32          `json:"agency-price"`
	PrincipalPrice                   StringToFloat32          `json:"principal-price"`
}

type TransactionFees struct {
	TotalFees       StringToFloat32 `json:"total-fees"`
	TotalFeesEffect string          `json:"total-fees-effect"`
}
