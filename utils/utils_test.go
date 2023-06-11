package utils

import (
	"testing"
	"time"

	"github.com/austinbspencer/tasty-go/constants"
	"github.com/stretchr/testify/require"
)

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

func TestGetEquity(t *testing.T) {
	sym := OCCSymbology{
		Symbol:     "AAPL",
		Strike:     185,
		OptionType: constants.Call,
		Expiration: time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC),
	}
	occSymbol := sym.GetOCCSymbology()

	require.Equal(t, "AAPL  230616C00185000", occSymbol)
}
