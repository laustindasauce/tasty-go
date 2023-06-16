package tasty

type InstrumentType string
type TimeOfDay string
type TimeInForce string
type OrderType string
type PriceEffect string
type OrderAction string
type Direction string
type Indicator string
type Comparator string
type OrderRuleAction string
type TimeBack string
type Cryptocurrency string
type Lendability string
type OptionType string
type MonthCode string
type Exchange string
type SortOrder string

// The normal flow for a filled order would be Received -> Routed -> In Flight -> Live -> Filled.
// Order status updates come in real-time to websocket clients that have sent the account-subscribe message.
type OrderStatus string

const (
	// InstrumentType.
	Bond           InstrumentType = "Bond"
	Crypto         InstrumentType = "Cryptocurrency"
	CurrencyPair   InstrumentType = "Currency Pair"
	EquityIT       InstrumentType = "Equity"
	EquityOffering InstrumentType = "Equity Offering"
	EquityOptionIT InstrumentType = "Equity Option"
	FutureIT       InstrumentType = "Future"
	FutureOptionIT InstrumentType = "Future Option"
	Index          InstrumentType = "Index"
	Unknown        InstrumentType = "Unknown"
	WarrantIT      InstrumentType = "Warrant"
	// TimeOfDay.
	EndOfDay       TimeOfDay = "EOD"
	BeginningOfDay TimeOfDay = "BOD"
	// TimeInForce.
	Day    TimeInForce = "Day"
	GTC    TimeInForce = "GTC"
	GTD    TimeInForce = "GTD"
	Ext    TimeInForce = "Ext"
	GTCExt TimeInForce = "GTC Ext"
	IOC    TimeInForce = "IOC"
	// OrderType.
	Limit           OrderType = "Limit"
	Market          OrderType = "Market"
	MarketableLimit OrderType = "Marketable Limit"
	Stop            OrderType = "Stop"
	StopLimit       OrderType = "Stop Limit"
	NotionalMarket  OrderType = "Notional Market"
	// PriceEffect.
	Credit PriceEffect = "Credit"
	Debit  PriceEffect = "Debit"
	None   PriceEffect = "None"
	// OrderAction.
	STO  OrderAction = "Sell to Open"
	STC  OrderAction = "Sell to Close"
	BTO  OrderAction = "Buy to Open"
	BTC  OrderAction = "Buy to Close"
	Sell OrderAction = "Sell"
	Buy  OrderAction = "Buy"
	// Direction.
	Long  Direction = "Long"
	Short Direction = "Short"
	// Indicator.
	Last Indicator = "last"
	// Comparator.
	GTE Comparator = "gte"
	LTE Comparator = "lte"
	// OrderRuleAction.
	Route  OrderRuleAction = "route"
	Cancel OrderRuleAction = "cancel"
	// TimeBack.
	OneDay      TimeBack = "1d"
	OneWeek     TimeBack = "1w"
	OneMonth    TimeBack = "1m"
	ThreeMonths TimeBack = "3m"
	SixMonths   TimeBack = "6m"
	OneYear     TimeBack = "1y"
	All         TimeBack = "all"
	// Cryptocurrency.
	Cardano     Cryptocurrency = "ADA"
	Bitcoin     Cryptocurrency = "BTC"
	BitcoinCash Cryptocurrency = "BCH"
	Polkadot    Cryptocurrency = "DOT"
	Dogecoin    Cryptocurrency = "DOGE"
	Ethereum    Cryptocurrency = "ETH"
	Litecoin    Cryptocurrency = "LTC"
	Solana      Cryptocurrency = "SOL"
	// Lendability.
	EasyToBorrow   Lendability = "Easy To Borrow"
	LocateRequired Lendability = "Locate Required"
	Preborrow      Lendability = "Preborrow"
	// OptionType.
	Call OptionType = "C"
	Put  OptionType = "P"
	// MonthCode.
	January   MonthCode = "F"
	February  MonthCode = "G"
	March     MonthCode = "H"
	April     MonthCode = "J"
	May       MonthCode = "K"
	June      MonthCode = "M"
	July      MonthCode = "N"
	August    MonthCode = "Q"
	September MonthCode = "U"
	October   MonthCode = "V"
	November  MonthCode = "X"
	December  MonthCode = "Z"
	// Exchange.
	CME    Exchange = "CME"
	SMALLS Exchange = "SMALLS"
	CFE    Exchange = "CFE"
	CBOED  Exchange = "CBOED"
	// OrderStatus.

	// Initial order state.
	Received OrderStatus = "Received"
	// Order is on its way out of Tastytrade's system.
	Routed OrderStatus = "Routed"
	// Order is en route to the exchange.
	InFlight OrderStatus = "In Flight"
	// Order is live at the exchange.
	Live OrderStatus = "Live"
	// Customer has requested to cancel the order.
	// Awaiting a 'cancelled' message from the exchange.
	CancelRequested OrderStatus = "Cancel Requested"
	// Customer has submitted a replacement order.
	// This order is awaiting a 'cancelled' message from the exchange.
	ReplaceRequested OrderStatus = "Replace Requested"
	// This pertains to replacement orders. It means the replacement order
	// is awaiting a 'cancelled' message for the order it is replacing.
	Contingent OrderStatus = "Contingent"
	// Order has been fully filled.
	Filled OrderStatus = "Filled"
	// Order is cancelled.
	Cancelled OrderStatus = "Cancelled"
	// Order has expired. Usually applies to an option order.
	Expired OrderStatus = "Expired"
	// Order has been rejected by either Tastytrade or the exchange.
	Rejected OrderStatus = "Rejected"
	// Administrator has manually removed this order from customer account.
	Removed OrderStatus = "Removed"
	// Administrator has manually removed part of this order from customer account.
	PartiallyRemoved OrderStatus = "Partially Removed"
	// SortOrder.
	Asc  SortOrder = "Asc"
	Desc SortOrder = "Desc"
)
