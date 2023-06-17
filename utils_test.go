package tasty //nolint:testpackage // testing private field

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFutureSymbology(t *testing.T) {
	future := FutureSymbology{ProductCode: "ES", MonthCode: December, YearDigit: 9}

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

func TestContainsInt(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	require.False(t, containsInt(nums, 10))
	require.True(t, containsInt(nums, 3))
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
