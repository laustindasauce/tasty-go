package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContainsInt(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	require.False(t, ContainsInt(nums, 10))
	require.True(t, ContainsInt(nums, 3))
}
