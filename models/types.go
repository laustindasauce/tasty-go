package models

import (
	"strconv"
	"strings"
)

type StringToFloat32 float32

func (foe *StringToFloat32) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		if foe != nil {
			*foe = 0
		}
		return nil
	}
	num := strings.ReplaceAll(string(data), "\"", "")
	n, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return err
	}
	*foe = StringToFloat32(n)
	return nil
}
