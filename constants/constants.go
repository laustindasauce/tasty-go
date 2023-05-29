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
