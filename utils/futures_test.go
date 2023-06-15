package utils_test

import (
	"testing"
	"time"

	"github.com/austinbspencer/tasty-go/constants"
	"github.com/austinbspencer/tasty-go/utils"
	"github.com/stretchr/testify/require"
)

func TestFutureSymbology(t *testing.T) {
	future := utils.FutureSymbology{ProductCode: "ES", MonthCode: constants.December, YearDigit: 9}

	fcc := utils.FutureOptionsSymbology{
		OptionContractCode: "EW4U9",
		FutureContractCode: future.Build(),
		OptionType:         constants.Put,
		Strike:             2975,
		Expiration:         time.Date(2019, 9, 27, 0, 0, 0, 0, time.Local),
	}

	require.Equal(t, "/ESZ9", future.Build())
	require.Equal(t, "./ESZ9 EW4U9 190927P2975", fcc.Build())
}
