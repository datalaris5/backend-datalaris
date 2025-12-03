package utils

import (
	"errors"
	"strconv"
)

func AnyToUintPtr(v any) (*uint, error) {
	switch val := v.(type) {
	case float64:
		u := uint(val)
		return &u, nil
	case int:
		u := uint(val)
		return &u, nil
	case int64:
		u := uint(val)
		return &u, nil
	case string:
		n, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, err
		}
		u := uint(n)
		return &u, nil
	default:
		return nil, errors.New("unsupported type for uint pointer conversion")
	}
}

func UintPtr(v uint) *uint {
	return &v
}
func StringPtr(v string) *string {
	return &v
}
func BoolPtr(v bool) *bool {
	return &v
}

func ParseUintParam(idStr string) uint {
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {

	}
	id := uint(idUint64)
	return id
}
