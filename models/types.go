package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type StringToFloat32 float32

// UnmarshalJSON is the custom unmarshaler interface
func (foe *StringToFloat32) UnmarshalJSON(data []byte) error {
	if string(data) == "\"\"" {
		if foe != nil {
			*foe = 0
		}
		return nil
	}

	num := strings.ReplaceAll(string(data), "\"", "")

	if num == "NaN" {
		if foe != nil {
			*foe = 0
		}
		return nil
	}

	n, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return err
	}
	*foe = StringToFloat32(n)
	return nil
}

// MarshalJSON is the custom marshaler interface
func (foe StringToFloat32) MarshalJSON() ([]byte, error) {
	return json.Marshal(float32(foe))
}

type SimpleDate time.Time

// UnmarshalJSON is the custom unmarshaler interface
func (j *SimpleDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = SimpleDate(t)
	return nil
}

// MarshalJSON is the custom marshaler interface
func (j SimpleDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

type Pagination struct {
	PerPage            int     `json:"per-page"`
	PageOffset         int     `json:"page-offset"`
	ItemOffset         int     `json:"item-offset"`
	TotalItems         int     `json:"total-items"`
	TotalPages         int     `json:"total-pages"`
	CurrentItemCount   int     `json:"current-item-count"`
	PreviousLink       *string `json:"previous-link"`
	NextLink           *string `json:"next-link"`
	PagingLinkTemplate *string `json:"paging-link-template"`
}
