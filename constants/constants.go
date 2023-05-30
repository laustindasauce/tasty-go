package constants

type InstrumentType string

const (
	Bond           InstrumentType = "Bond"
	Crypto         InstrumentType = "Cryptocurrency"
	CurrencyPair   InstrumentType = "Currency Pair"
	Equity         InstrumentType = "Equity"
	EquityOffering InstrumentType = "Equity Offering"
	EquityOption   InstrumentType = "Equity Option"
	Future         InstrumentType = "Future"
	FutureOption   InstrumentType = "Future Option"
	Index          InstrumentType = "Index"
	Unknown        InstrumentType = "Unknown"
	Warrant        InstrumentType = "Warrant"
)

type TimeOfDay *string

var (
	eod                      = "EOD"
	bod                      = "BOD"
	EndOfDay       TimeOfDay = &eod
	BeginningOfDay TimeOfDay = &bod
)

type TimeInForce *string

var (
	day                = "Day"
	gtc                = "GTC"
	gtd                = "GTD"
	ext                = "Ext"
	gtcExt             = "GTC Ext"
	ioc                = "IOC"
	Day    TimeInForce = &day
	GTC    TimeInForce = &gtc
	GTD    TimeInForce = &gtd
	Ext    TimeInForce = &ext
	GTCExt TimeInForce = &gtcExt
	IOC    TimeInForce = &ioc
)

type OrderType *string

var (
	limit                     = "Limit"
	market                    = "Market"
	marketableLimit           = "Marketable Limit"
	stop                      = "Stop"
	stopLimit                 = "Stop Limit"
	notionalMarket            = "Notional Market"
	Limit           OrderType = &limit
	Market          OrderType = &market
	MarketableLimit OrderType = &marketableLimit
	Stop            OrderType = &stop
	StopLimit       OrderType = &stopLimit
	NotionalMarket  OrderType = &notionalMarket
)

type PriceEffect *string

var (
	credit             = "Credit"
	debit              = "Debit"
	Credit PriceEffect = &credit
	Debit  PriceEffect = &debit
)

type OrderAction *string

var (
	sto              = "Sell to Open"
	stc              = "Sell to Close"
	bto              = "Buy to Open"
	btc              = "Buy to Close"
	sell             = "Sell"
	buy              = "Buy"
	STO  OrderAction = &sto
	STC  OrderAction = &stc
	BTO  OrderAction = &bto
	BTC  OrderAction = &btc
	Sell OrderAction = &sell
	Buy  OrderAction = &buy
)

type Direction *string

var (
	long            = "long"
	short           = "short"
	Long  Direction = &long
	Short Direction = &short
)

type Indicator string

const (
	Last Indicator = "last"
)

type Comparator *string

var (
	gte            = "gte"
	lte            = "lte"
	GTE Comparator = &gte
	LTE Comparator = &lte
)

type OrderRuleAction *string

var (
	route                  = "route"
	cancel                 = "cancel"
	Route  OrderRuleAction = &route
	Cancel OrderRuleAction = &cancel
)
