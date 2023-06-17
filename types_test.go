package tasty_test

import (
	"encoding/json"
	"testing"

	"github.com/austinbspencer/tasty-go"
	"github.com/stretchr/testify/require"
)

func TestStringToFloat32(t *testing.T) {
	type test struct {
		Key tasty.StringToFloat32
	}

	testMap := map[string]string{
		"Key": "155.55",
	}

	test1JSON, err := json.Marshal(testMap)
	require.NoError(t, err)

	test1Res := new(test)

	err = json.Unmarshal(test1JSON, test1Res)
	require.NoError(t, err)

	require.Equal(t, float32(155.55), float32(test1Res.Key))

	res, err := test1Res.Key.MarshalJSON()
	require.NoError(t, err)

	require.Equal(t, "155.55", string(res))

	// Testing non numerical string
	testMap["Key"] = "."

	test2JSON, err := json.Marshal(testMap)
	require.NoError(t, err)

	test2Res := new(test)
	err = json.Unmarshal(test2JSON, test2Res)
	require.Error(t, err)

	// Testing empty string
	testMap["Key"] = ""

	test3JSON, err := json.Marshal(testMap)
	require.NoError(t, err)

	test3Res := new(test)
	err = json.Unmarshal(test3JSON, test3Res)
	require.NoError(t, err)

	res, err = test3Res.Key.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, "0", string(res))
}
