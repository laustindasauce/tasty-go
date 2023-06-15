package utils_test

import (
	"testing"

	"github.com/austinbspencer/tasty-go/utils"
	"github.com/stretchr/testify/require"
)

func TestContainsInt(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	require.False(t, utils.ContainsInt(nums, 10))
	require.True(t, utils.ContainsInt(nums, 3))
}
