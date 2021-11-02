package util

import (
	"fmt"
	"strconv"
)

// ConvertToInt converts any type of number interface to a regular old integer.
// Runes will be treated as int32 and bytes as uint8. Strings will be converted
// using strings.Atoi. An error will be returned if the passed parameter cannot
// be converted to an int
func ConvertToInt(i interface{}) (int, error) {
	switch v := i.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	}
	return 0, fmt.Errorf("%v is not an int", i)
}
