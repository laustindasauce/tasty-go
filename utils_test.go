package tasty //nolint:testpackage // testing private field

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFutureSymbology(t *testing.T) {
	future := FutureSymbology{ProductCode: "ES",
		MonthCode: December, YearDigit: 9}

	fcc := FutureOptionsSymbology{
		OptionContractCode: "EW4U9",
		FutureContractCode: future.Build(),
		OptionType:         Put,
		Strike:             2975,
		Expiration:         time.Date(2019, 9, 27, 0, 0, 0, 0, time.Local),
	}

	require.Equal(t, "/ESZ9", future.Build())
	require.Equal(t, "./ESZ9 EW4U9 190927P2975", fcc.Build())
}

func TestYearDigitFromFS(t *testing.T) {
	num, err := yearDigitFromFS("/ESZ9")
	require.NoError(t, err)

	require.Equal(t, 9, num)

	num, err = yearDigitFromFS("/VXX22")
	require.NoError(t, err)

	require.Equal(t, 22, num)

	_, err = yearDigitFromFS("/VXX//")
	require.Error(t, err)

	require.Equal(t, "missing valid year", err.Error())
}

func TestNewFSFromString(t *testing.T) {
	sym, err := NewFSFromString("/ESZ9")
	require.NoError(t, err)

	require.Equal(t, 9, sym.YearDigit)
	require.Equal(t, December, sym.MonthCode)
	require.Equal(t, "ES", sym.ProductCode)

	_, err = NewFSFromString("/ESA9")
	require.Error(t, err)

	require.Equal(t, "invalid month code: A", err.Error())

	_, err = NewFSFromString("EESA9")
	require.Error(t, err)

	require.Equal(t, "future symbol must start with '/'", err.Error())

	_, err = NewFSFromString("/VXXYY")
	require.Error(t, err)

	require.Equal(t, "missing valid year", err.Error())
}

func TestNewFOSFromString(t *testing.T) {
	sym, err := NewFOSFromString("./CLZ2 LO1X2 221104C91")
	require.NoError(t, err)

	require.Equal(t, "LO1X2", sym.OptionContractCode)
	require.Equal(t, "/CLZ2", sym.FutureContractCode)
	require.Equal(t, Call, sym.OptionType)
	require.Equal(t, 91, sym.Strike)
	require.Equal(t, time.Date(2022, time.November, 4, 0, 0, 0, 0, time.UTC).Format(time.RFC1123),
		sym.Expiration.Format(time.RFC1123))

	_, err = NewFOSFromString("/ESA9")
	require.Error(t, err)

	require.Equal(t, "invalid future options symbol structure", err.Error())

	_, err = NewFOSFromString("./CLZ2 LO1X2 22110")
	require.Error(t, err)

	require.Equal(t, "invalid future options symbol structure: strike info too short", err.Error())

	_, err = NewFOSFromString("/CLZ2 LO1X2 221104C91")
	require.Error(t, err)

	require.Equal(t, "future options symbol must start with './'", err.Error())

	_, err = NewFOSFromString("./CLZZ LO1X2 221104C91")
	require.Error(t, err)

	require.Equal(t, "invalid contract code: missing valid year", err.Error())

	_, err = NewFOSFromString("./CLZZ2 LO1X2 2f1104C91")
	require.Error(t, err)

	require.Equal(t, "parsing time \"2f1104\" as \"060102\": cannot parse \"2f1104\" as \"06\"", err.Error())
}

func TestContainsInt(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	require.False(t, containsInt(nums, 10))
	require.True(t, containsInt(nums, 3))
}

func TestContainsString(t *testing.T) {
	vals := []string{"1", "2", "3", "4", "5"}

	require.False(t, containsString(vals, "10"))
	require.True(t, containsString(vals, "3"))
}

func TestGetSymbolWithPadding(t *testing.T) {
	require.Equal(t, "AAPL  ", getSymbolWithPadding("AAPL"))
	require.Equal(t, "AMD   ", getSymbolWithPadding("AMD"))
	require.Equal(t, "MU    ", getSymbolWithPadding("MU"))
	require.Equal(t, "X     ", getSymbolWithPadding("X"))
}

func TestGetStrikeWithPadding(t *testing.T) {
	require.Equal(t, "00645500", getStrikeWithPadding(645.5))
	require.Equal(t, "00185000", getStrikeWithPadding(185))
	require.Equal(t, "00015500", getStrikeWithPadding(15.5))
	require.Equal(t, "00012000", getStrikeWithPadding(12))
	require.Equal(t, "00005000", getStrikeWithPadding(5))
}

func TestGetEquitySymbol(t *testing.T) {
	sym := EquityOptionsSymbology{
		Symbol:     "AAPL",
		Strike:     185,
		OptionType: Call,
		Expiration: time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
	}
	occSymbol := sym.Build()

	require.Equal(t, "AAPL  230616C00185000", occSymbol)
}

func TestGetEquitySymbolFromSymbol(t *testing.T) {
	sym := EquityOptionsSymbology{
		Symbol:     "AAPL",
		Strike:     185,
		OptionType: Call,
		Expiration: time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
	}
	occSymbol, err := NewOCCFromString(sym.Build())
	require.NoError(t, err)

	require.Equal(t, sym.Symbol, occSymbol.Symbol)
	require.Equal(t, sym.Strike, occSymbol.Strike)
	require.Equal(t, sym.OptionType, occSymbol.OptionType)
	require.Equal(t,
		sym.Expiration.Format(time.RFC1123),
		occSymbol.Expiration.Format(time.RFC1123))
}

func TestGetEquitySymbolFromSymbolErrors(t *testing.T) {
	_, err := NewOCCFromString("AAPL  230616C0018500")
	require.Error(t, err)
	require.Equal(t, "invalid occ symbol", err.Error())

	_, err = NewOCCFromString("AAPL  s30616C00185000")
	require.Error(t, err)
	require.Equal(t, "parsing time \"s30616\" as \"060102\": cannot parse \"s30616\" as \"06\"", err.Error())

	_, err = NewOCCFromString("AAPL  230616M00185000")
	require.Error(t, err)
	require.Equal(t, "unknown option type: M", err.Error())

	_, err = NewOCCFromString("AAPL  230616C0018f000")
	require.Error(t, err)
	require.Equal(t, "invalid option strike: 0018f000", err.Error())
}
